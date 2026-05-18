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

	Path       []Point
	pathLookup map[Point]bool
	Towers     []*Tower
	Enemies    []*Enemy
	Gold       int
	Lives      int
	StartLives int
	WaveIdx    int
	Cursor     Point
	Selected   TowerKind
	Msg        string

	prepTimer  float64
	spawned    int
	spawnTimer float64
}

func NewGame(levels []Level) *Game {
	return &Game{Phase: PhaseLevelSelect, Levels: levels}
}

func (g *Game) StartLevel(idx int) {
	if idx < 0 || idx >= len(g.Levels) {
		return
	}
	lv := g.Levels[idx]
	g.Phase = PhasePlaying
	g.LevelIdx = idx
	g.Path = lv.Path
	g.pathLookup = make(map[Point]bool, len(lv.Path))
	for _, p := range lv.Path {
		g.pathLookup[p] = true
	}
	g.Towers = nil
	g.Enemies = nil
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
		if g.prepTimer < 0 {
			g.prepTimer = 0
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
				g.Phase = PhaseLost
				return
			}
		}
	}
	// shoot
	for _, t := range g.Towers {
		spec := t.Spec()
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
			ep := e.Pos(g.Path)
			dx := float64(ep.X - t.Pos.X)
			dy := float64(ep.Y - t.Pos.Y)
			if math.Sqrt(dx*dx+dy*dy) <= spec.Range && e.PathIdx > maxIdx {
				maxIdx = e.PathIdx
				target = e
			}
		}
		if target != nil {
			target.HP -= spec.Damage
			t.cooldown = spec.Cooldown
			if target.HP <= 0 {
				target.Dead = true
				g.Gold += enemySpecs[target.Kind].Reward
			}
		}
	}
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
				g.Phase = PhaseWon
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
