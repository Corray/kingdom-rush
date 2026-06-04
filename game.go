// Game state machine + 主 Update logic。
// Phase: LevelSelect -> Playing -> Won/Lost -> (返回 LevelSelect)
package main

import (
	"fmt"
	"math"
)

const (
	mapW      = 30
	mapH      = 15
	fps       = 30
	wavePrepS = 5.0
	spawnGapS = 0.7
	waveBonus = 30
)

type GamePhase int

const (
	PhaseLevelSelect GamePhase = iota
	PhasePlaying
	PhaseWon
	PhaseLost
)

type Game struct {
	Phase    GamePhase
	Levels   []Level
	LevelIdx int
	Save     Save // 存档:已完成关卡 + unlock 状态

	Path        []Point
	pathLookup  map[Point]bool
	Towers      []*Tower
	Enemies     []*Enemy
	Effects     []*Effect    // V2.6: 攻击视觉特效 (ebiten 渲染消费,terminal 忽略)
	SoundEvents []SoundEvent // V4: SFX 事件队列 (sound.go, 渲染层每帧 drain)
	Gold        int
	Lives       int
	StartLives  int
	WaveIdx     int
	Cursor      Point
	Selected    TowerKind
	Msg         string

	prepTimer  float64
	spawned    int
	spawnTimer float64
}

func NewGame(levels []Level, save Save) *Game {
	return &Game{Phase: PhaseLevelSelect, Levels: levels, Save: save}
}

func (g *Game) StartLevel(idx int) {
	if idx < 0 || idx >= len(g.Levels) {
		return
	}
	lv := g.Levels[idx]
	if !g.Save.IsUnlocked(lv.ID) {
		g.Msg = fmt.Sprintf("Level %d locked — clear Level %d first", lv.ID, lv.ID-1)
		return
	}
	g.Phase = PhasePlaying
	g.LevelIdx = idx
	g.Path = lv.Path
	g.pathLookup = make(map[Point]bool, len(lv.Path))
	for _, p := range lv.Path {
		g.pathLookup[p] = true
	}
	g.Towers = nil
	g.Enemies = nil
	g.Effects = nil
	g.Gold = lv.StartGold
	g.Lives = lv.StartLives
	g.StartLives = lv.StartLives
	g.WaveIdx = 0
	g.Cursor = Point{15, 10}
	g.Selected = TArcher
	g.prepTimer = wavePrepS
	g.spawned = 0
	g.spawnTimer = 0
	g.Msg = fmt.Sprintf("Level %d: %s — get ready!", lv.ID, lv.Name)
}

func (g *Game) BackToMenu() {
	g.Phase = PhaseLevelSelect
	g.Msg = ""
}

func (g *Game) pathContains(p Point) bool {
	return g.pathLookup[p]
}

func (g *Game) currentLevel() *Level {
	if g.LevelIdx < 0 || g.LevelIdx >= len(g.Levels) {
		return nil
	}
	return &g.Levels[g.LevelIdx]
}

func (g *Game) currentWave() *WaveSpec {
	lv := g.currentLevel()
	if lv == nil || g.WaveIdx >= len(lv.Waves) {
		return nil
	}
	return &lv.Waves[g.WaveIdx]
}

func (g *Game) Update(dt float64) {
	if g.Phase != PhasePlaying {
		return
	}
	if g.prepTimer > 0 {
		g.prepTimer -= dt
		if g.prepTimer <= 0 {
			g.prepTimer = 0
			g.pushSound(SndWaveStart) // V4: prep 倒计时归零 = wave 开打 (每 wave 恰一次)
		}
	}
	cur := g.currentWave()
	if cur == nil {
		return
	}
	// spawn
	if g.prepTimer == 0 && g.spawned < len(cur.Enemies) {
		g.spawnTimer += dt
		if g.spawnTimer >= spawnGapS {
			g.spawnTimer = 0
			kind := cur.Enemies[g.spawned]
			spec := enemySpecs[kind]
			g.Enemies = append(g.Enemies, &Enemy{Kind: kind, HP: spec.HP, MaxHP: spec.HP})
			g.spawned++
		}
	}
	// move
	for _, e := range g.Enemies {
		if e.Dead || e.Escaped {
			continue
		}
		speed := enemySpecs[e.Kind].Speed
		e.PathIdx += speed * dt
		if e.PathIdx >= float64(len(g.Path)-1) {
			e.Escaped = true
			g.Lives--
			if g.Lives <= 0 {
				g.pushSound(SndLose)
				g.Phase = PhaseLost
				return
			}
		}
	}
	// shoot — 塔目标 path 最远的可击中敌人 (Cannon 不打飞行)
	for _, t := range g.Towers {
		lvl := t.Spec()
		towerSpec := towerSpecs[t.Kind]
		t.cooldown -= dt
		if t.cooldown > 0 {
			continue
		}
		var target *Enemy
		maxIdx := -1.0
		for _, e := range g.Enemies {
			if e.Dead || e.Escaped {
				continue
			}
			eSpec := enemySpecs[e.Kind]
			if eSpec.Flying && !towerSpec.HitsFlying {
				continue
			}
			ep := e.Pos(g.Path)
			dx := float64(ep.X - t.Pos.X)
			dy := float64(ep.Y - t.Pos.Y)
			if math.Sqrt(dx*dx+dy*dy) <= lvl.Range && e.PathIdx > maxIdx {
				maxIdx = e.PathIdx
				target = e
			}
		}
		if target != nil {
			target.HP -= lvl.Damage
			t.cooldown = lvl.Cooldown
			// V2.6: push 视觉特效  (V3 Phase 3b: shoot 带 tower kind 决定 bullet sprite)
			g.Effects = append(g.Effects,
				makeShootEffect(t.Pos, target.Pos(g.Path), towerSpec.Color, t.Kind),
				makeHitEffect(target.Pos(g.Path)),
			)
			g.pushSound(shootSound(t.Kind)) // V4: 射击音按塔型分 (伤害即时结算, 不另设 hit 音)
			if target.HP <= 0 {
				target.Dead = true
				g.pushSound(SndEnemyDeath)
				g.Gold += enemySpecs[target.Kind].Reward
				// V3.6: Spawner 死时 spawn 2 个 ENormal 在同 PathIdx
				if target.Kind == ESpawner {
					normSpec := enemySpecs[ENormal]
					g.Enemies = append(g.Enemies,
						&Enemy{Kind: ENormal, HP: normSpec.HP, MaxHP: normSpec.HP, PathIdx: target.PathIdx},
						&Enemy{Kind: ENormal, HP: normSpec.HP, MaxHP: normSpec.HP, PathIdx: target.PathIdx},
					)
				}
			}
		}
	}
	// V2.6: 衰减视觉特效, 过期清理
	g.Effects = decayEffects(g.Effects, dt)
	// wave done?
	if g.prepTimer == 0 && g.spawned >= len(cur.Enemies) {
		alive := false
		for _, e := range g.Enemies {
			if !e.Dead && !e.Escaped {
				alive = true
				break
			}
		}
		if !alive {
			lv := g.currentLevel()
			if g.WaveIdx+1 < len(lv.Waves) {
				g.WaveIdx++
				g.prepTimer = wavePrepS
				g.spawned = 0
				g.spawnTimer = 0
				g.Enemies = nil
				g.Gold += waveBonus
				g.Msg = fmt.Sprintf("Wave %d cleared! +%dg, prep for next", g.WaveIdx, waveBonus)
			} else {
				g.pushSound(SndWin)
				g.Phase = PhaseWon
				g.Save.MarkCompleted(lv.ID)
				if err := StoreSave(g.Save); err != nil {
					// 不阻断游戏,仅在 msg 显示;玩家可手动重试或忽略
					g.Msg = fmt.Sprintf("Victory! (save failed: %v)", err)
				} else {
					g.Msg = fmt.Sprintf("Victory! Level %d cleared & saved", lv.ID)
				}
			}
		}
	}
}

func (g *Game) TryAction() {
	// 光标在塔上 → 升级;否则 → 建塔
	var atTower *Tower
	for _, t := range g.Towers {
		if t.Pos == g.Cursor {
			atTower = t
			break
		}
	}
	if atTower != nil {
		cost, can := atTower.NextUpgradeCost()
		if !can {
			g.Msg = "Tower at max level"
			return
		}
		if g.Gold < cost {
			g.Msg = fmt.Sprintf("Need %dg to upgrade", cost)
			return
		}
		g.Gold -= cost
		atTower.Level++
		g.pushSound(SndUpgrade)
		g.Msg = fmt.Sprintf("Upgraded to lvl %d", atTower.Level)
		return
	}
	spec := towerSpecs[g.Selected].Levels[0]
	if g.Gold < spec.Cost {
		g.Msg = fmt.Sprintf("Need %dg for %s", spec.Cost, towerSpecs[g.Selected].Name)
		return
	}
	if g.pathContains(g.Cursor) {
		g.Msg = "Cannot build on path"
		return
	}
	g.Towers = append(g.Towers, &Tower{Pos: g.Cursor, Kind: g.Selected, Level: 1})
	g.Gold -= spec.Cost
	g.pushSound(SndBuild)
	g.Msg = fmt.Sprintf("Built %s", towerSpecs[g.Selected].Name)
}

func (g *Game) CountAliveEnemies() int {
	n := 0
	for _, e := range g.Enemies {
		if !e.Dead && !e.Escaped {
			n++
		}
	}
	return n
}

func (g *Game) MoveCursor(dx, dy int) {
	g.Cursor.X += dx
	g.Cursor.Y += dy
	if g.Cursor.X < 0 {
		g.Cursor.X = 0
	}
	if g.Cursor.X >= mapW {
		g.Cursor.X = mapW - 1
	}
	if g.Cursor.Y < 0 {
		g.Cursor.Y = 0
	}
	if g.Cursor.Y >= mapH {
		g.Cursor.Y = mapH - 1
	}
}
