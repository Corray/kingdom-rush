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
	BossKills   int          // V4 Phase 5: boss 击杀累计 (渲染层观察增量触发顿帧, 同 goldFlash 模式)
	MeteorCD    float64      // V5 Phase 5: 陨石雨剩余冷却秒 (0 = 可用)
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
	g.MeteorCD = 0 // V5 Phase 5: 每关开局技能就绪
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
	// V5 Phase 5: 陨石雨冷却衰减
	if g.MeteorCD > 0 {
		g.MeteorCD -= dt
		if g.MeteorCD < 0 {
			g.MeteorCD = 0
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
	// move (V5 Phase 4: 减速状态生效 + 计时衰减)
	for _, e := range g.Enemies {
		if e.Dead || e.Escaped {
			continue
		}
		e.PathIdx += e.EffectiveSpeed() * dt
		if e.SlowTimer > 0 {
			e.SlowTimer -= dt
		}
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
		// V5 Phase 2: 目标选择抽 pickTarget 纯函数 (targeting.go),
		// 按塔的 Target 策略选 (默认 First = 原"最前优先"行为)
		target := pickTarget(t, g.Enemies, g.Path)
		if target != nil {
			t.cooldown = lvl.Cooldown
			// V2.6: push 视觉特效  (V3 Phase 3b: shoot 带 tower kind 决定 bullet sprite)
			g.Effects = append(g.Effects,
				makeShootEffect(t.Pos, target.Pos(g.Path), towerSpec.Color, t.Kind),
				makeHitEffect(target.Pos(g.Path)),
			)
			g.pushSound(shootSound(t.Kind)) // V4: 射击音按塔型分 (伤害即时结算, 不另设 hit 音)
			g.damageEnemy(target, lvl.Damage)
			// V5 Phase 4: Frost 命中减速 (击杀后施加无意义但无害)
			if lvl.Slow > 0 {
				target.ApplySlow(lvl.Slow)
			}
			// V5 Phase 3: Cannon 溅射 — 主目标位置 Splash 半径内
			// 其他敌人吃 splashFactor 折扣伤害 (飞行过滤与主目标一致)
			if lvl.Splash > 0 {
				tp := target.Pos(g.Path)
				splashDmg := int(math.Round(float64(lvl.Damage) * splashFactor))
				for _, e := range g.Enemies {
					if e == target || e.Dead || e.Escaped {
						continue
					}
					if enemySpecs[e.Kind].Flying && !towerSpec.HitsFlying {
						continue
					}
					ep := e.Pos(g.Path)
					ddx := float64(ep.X - tp.X)
					ddy := float64(ep.Y - tp.Y)
					if math.Sqrt(ddx*ddx+ddy*ddy) <= lvl.Splash {
						g.damageEnemy(e, splashDmg)
					}
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

// AdjustVolume: V4 Phase 2 — 音量档 ±delta (0-10), 立即持久化。
// 存档失败不阻断 (音量是偏好设置, 丢失可重调)。
func (g *Game) AdjustVolume(delta int) {
	g.Save.SetVolumeLevel(g.Save.VolumeLevel() + delta)
	_ = StoreSave(g.Save)
	g.Msg = fmt.Sprintf("Volume %d/%d", g.Save.VolumeLevel(), maxVolume)
}

// splashFactor: V5 Phase 3 — 溅射伤害折扣 (主目标全额, 周围打折)。
const splashFactor = 0.5

// damageEnemy: V5 Phase 3 前置重构 — 伤害结算唯一入口 (主目标与
// 溅射目标同路径): HP 扣减 + 伤害飘字 + 致死转 killEnemy。
func (g *Game) damageEnemy(e *Enemy, dmg int) {
	e.HP -= dmg
	// V4 Phase 3/4: 插值坐标 (对齐渲染层平滑位置), 飘字 + 死亡动画共用
	fx, fy := pathLerp(g.Path, e.PathIdx)
	g.Effects = append(g.Effects, makeDamageText(fx, fy, dmg))
	if e.HP <= 0 && !e.Dead {
		g.killEnemy(e, fx, fy)
	}
}

// killEnemy: 致死处理唯一出处 (音效/Boss 计数/死亡动画/赏金/Spawner
// 召唤)。fix-pattern-scan 家族约束: 任何新增伤害来源 (溅射/技能/...)
// 必须经 damageEnemy → 本函数, 禁止内联复制击杀逻辑。
func (g *Game) killEnemy(e *Enemy, fx, fy float64) {
	e.Dead = true
	g.pushSound(SndEnemyDeath)
	if e.Kind == EBoss {
		g.BossKills++ // V4 Phase 5: 渲染层观察此计数触发顿帧
	}
	reward := enemySpecs[e.Kind].Reward
	// 赏金飘字右偏 0.3 cell, 与伤害字错开
	g.Effects = append(g.Effects,
		makeDeathEffect(fx, fy, e.Kind),
		makeGoldText(fx+0.3, fy, reward))
	g.Gold += reward
	// V3.6: Spawner 死时 spawn 2 个 ENormal 在同 PathIdx
	if e.Kind == ESpawner {
		normSpec := enemySpecs[ENormal]
		g.Enemies = append(g.Enemies,
			&Enemy{Kind: ENormal, HP: normSpec.HP, MaxHP: normSpec.HP, PathIdx: e.PathIdx},
			&Enemy{Kind: ENormal, HP: normSpec.HP, MaxHP: normSpec.HP, PathIdx: e.PathIdx},
		)
	}
}

// V5 Phase 5: 陨石雨主动技能参数。
const (
	meteorCooldownS = 25.0
	meteorRadius    = 2.0 // cell
	meteorDamage    = 60
)

// CastMeteor: V5 Phase 5 — 在 at 处释放陨石雨: 半径内全敌 AoE
// (含飞行 — 天上砸下来的) 经 damageEnemy 统一路径结算, 触发冷却。
// 冷却中 / 非战斗 phase 返回 false。
func (g *Game) CastMeteor(at Point) bool {
	if g.MeteorCD > 0 || g.Phase != PhasePlaying {
		return false
	}
	g.MeteorCD = meteorCooldownS
	g.pushSound(SndMeteor)
	hits := 0
	for _, e := range g.Enemies {
		if e.Dead || e.Escaped {
			continue
		}
		ep := e.Pos(g.Path)
		dx := float64(ep.X - at.X)
		dy := float64(ep.Y - at.Y)
		if math.Sqrt(dx*dx+dy*dy) <= meteorRadius {
			g.damageEnemy(e, meteorDamage)
			hits++
		}
	}
	// 火焰特效: 半径内每个 cell 一团火 (复用 EHit 渲染)
	for ddy := -int(meteorRadius); ddy <= int(meteorRadius); ddy++ {
		for ddx := -int(meteorRadius); ddx <= int(meteorRadius); ddx++ {
			if math.Sqrt(float64(ddx*ddx+ddy*ddy)) > meteorRadius {
				continue
			}
			g.Effects = append(g.Effects, makeHitEffect(Point{X: at.X + ddx, Y: at.Y + ddy}))
		}
	}
	g.Msg = fmt.Sprintf("Meteor strike! %d hit", hits)
	return true
}

// sellRefundRate: V5 Phase 1 — 卖塔退款比例 (对累计投入)。
const sellRefundRate = 0.7

// sellRefund: 退款额单一出处 (SellTower + HUD 提示共用)。
// math.Round 防 IEEE754 截断 (90 × 0.7 = 62.999… → int 截成 62)。
func sellRefund(kind TowerKind, level int) int {
	return int(math.Round(float64(towerInvested(kind, level)) * sellRefundRate))
}

// SellTower: V5 Phase 1 — 卖光标位置的塔。退款入账 + 金币飘字
// (复用 V4 Phase 4) + SFX; 光标不在塔上仅提示。
func (g *Game) SellTower() {
	for i, t := range g.Towers {
		if t.Pos != g.Cursor {
			continue
		}
		refund := sellRefund(t.Kind, t.Level)
		g.Gold += refund
		g.Towers = append(g.Towers[:i], g.Towers[i+1:]...)
		g.pushSound(SndSell)
		g.Effects = append(g.Effects,
			makeGoldText(float64(t.Pos.X), float64(t.Pos.Y), refund))
		g.Msg = fmt.Sprintf("Sold %s +%dg", towerSpecs[t.Kind].Name, refund)
		return
	}
	g.Msg = "No tower here to sell"
}

// CycleTargeting: V5 Phase 2 — 光标处塔的 targeting 策略循环切换
// (First → Last → Strong → First)。塔级运行时状态, 不入存档。
func (g *Game) CycleTargeting() {
	for _, t := range g.Towers {
		if t.Pos == g.Cursor {
			t.Target = t.Target.Next()
			g.Msg = fmt.Sprintf("%s targeting: %s", towerSpecs[t.Kind].Name, t.Target.Name())
			return
		}
	}
	g.Msg = "No tower here (T = targeting)"
}

// ToggleJuice: V4 Phase 5 — 屏幕反馈特效 (shake / 顿帧) 一键开关, 持久化。
func (g *Game) ToggleJuice() {
	g.Save.JuiceOff = !g.Save.JuiceOff
	_ = StoreSave(g.Save)
	if g.Save.JuiceOff {
		g.Msg = "Screen effects: OFF"
	} else {
		g.Msg = "Screen effects: ON"
	}
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
