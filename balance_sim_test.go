// V7.4 C5: 平衡仿真器 — auto-player 以"中等玩家下限"策略自动打关,
// 验证 20 关 Normal 可通关性。平衡数值改动从此有回归。
//
// 策略 (刻意保守, 结论只作下限参考):
//   - 塔位: 按射程内 path 覆盖 cell 数排序, 贪心取最优
//   - 蓝图: Archer 为主, 每第 4 塔放 Frost (减速普适); 上限 8 塔
//   - 塔满后升级; 陨石 ready 且场上 ≥6 敌时砸最前
//
// V8: 英雄被动计入 — beginRun 自生成在 path 中点, auto-player 不主动
// rally (保持"未用英雄"的保守下限)。autoPlay 加 heroEnabled 参数, 支持
// 有/无英雄对比 (TestBalance_HeroNetNonNegative 守"英雄不让难度上升")。
package main

import "testing"

// rankedTowerSpots: 非 path cell 按"Archer lvl1 射程内 path cell 数"降序。
func rankedTowerSpots(g *Game) []Point {
	const r = 3.5
	type scored struct {
		p Point
		n int
	}
	var cands []scored
	for y := 0; y < mapH; y++ {
		for x := 0; x < mapW; x++ {
			p := Point{X: x, Y: y}
			if g.pathContains(p) {
				continue
			}
			n := 0
			for _, pc := range g.Path {
				dx, dy := float64(pc.X-p.X), float64(pc.Y-p.Y)
				if dx*dx+dy*dy <= r*r {
					n++
				}
			}
			if n > 0 {
				cands = append(cands, scored{p, n})
			}
		}
	}
	// 简单选择排序 (规模 ~400, 测试场景够用)
	for i := 0; i < len(cands); i++ {
		best := i
		for j := i + 1; j < len(cands); j++ {
			if cands[j].n > cands[best].n {
				best = j
			}
		}
		cands[i], cands[best] = cands[best], cands[i]
	}
	out := make([]Point, len(cands))
	for i, c := range cands {
		out[i] = c.p
	}
	return out
}

func simStrategy(g *Game, spots []Point, built *int) {
	// 铺塔阶段
	if *built < len(spots) && *built < 8 {
		kind := TArcher
		if *built%4 == 3 {
			kind = TFrost
		}
		if g.Gold >= towerSpecs[kind].Levels[0].Cost {
			g.Selected = kind
			g.Cursor = spots[*built]
			before := len(g.Towers)
			g.TryAction()
			if len(g.Towers) > before {
				*built++
			}
		}
		return // 攒钱优先建塔
	}
	// 升级阶段
	for _, tw := range g.Towers {
		if cost, ok := tw.NextUpgradeCost(); ok && g.Gold >= cost {
			g.Cursor = tw.Pos
			g.TryAction()
			return
		}
	}
	// 陨石: 压力大时砸最前
	if g.MeteorCD <= 0 && g.CountAliveEnemies() >= 6 {
		var front *Enemy
		for _, e := range g.Enemies {
			if e.Dead || e.Escaped {
				continue
			}
			if front == nil || e.PathIdx > front.PathIdx {
				front = e
			}
		}
		if front != nil {
			g.CastMeteor(front.Pos(g.Path))
		}
	}
}

// autoPlay: 自动打完一关 (或 60000 步上限 = game-time 100 分钟)。
// heroEnabled=false 时 nil 掉英雄, 用于"有/无英雄"平衡对比。
func autoPlay(levelIdx int, d Difficulty, heroEnabled bool) *Game {
	levels, err := LoadLevels()
	if err != nil {
		panic(err)
	}
	g := NewGame(levels, NewSave())
	for i := 1; i <= len(levels); i++ {
		g.Save.MarkCompleted(i) // 绕 unlock 链
	}
	g.Save.Difficulty = d
	g.StartLevel(levelIdx)
	if !heroEnabled {
		g.Hero = nil // V8: 对比基线 — 去英雄 (Update 已 nil-guard)
	}
	spots := rankedTowerSpots(g)
	built := 0
	for step := 0; step < 60000 && g.Phase == PhasePlaying; step++ {
		if step%5 == 0 {
			simStrategy(g, spots, &built)
		}
		g.Update(0.1)
		g.DrainSounds() // 防队列顶满
	}
	return g
}

func TestBalance_AllLevelsBeatableOnNormal(t *testing.T) {
	withTempSavePath(t, func() {
		for idx := 0; idx < 20; idx++ {
			g := autoPlay(idx, DiffNormal, true)
			if g.Phase != PhaseWon {
				t.Errorf("Lv%-2d ✗ 不可通关 (phase=%v lives=%d wave=%d/%d)",
					idx+1, g.Phase, g.Lives, g.WaveIdx+1, len(g.currentLevel().Waves))
			} else {
				t.Logf("Lv%-2d ✓ %d★ lives %d/%d",
					idx+1, starsFor(g.Lives, g.StartLives), g.Lives, g.StartLives)
			}
		}
	})
}

func TestBalance_EndlessSurvivesBaseline(t *testing.T) {
	// endless 基础可玩性: 同策略至少撑过 5 波 (下限断言)
	withTempSavePath(t, func() {
		levels, _ := LoadLevels()
		g := NewGame(levels, NewSave())
		g.StartEndless(42)
		spots := rankedTowerSpots(g)
		built := 0
		for step := 0; step < 120000 && g.Phase == PhasePlaying; step++ {
			if step%5 == 0 {
				simStrategy(g, spots, &built)
			}
			g.Update(0.1)
			g.DrainSounds()
		}
		if g.Save.BestWave < 5 {
			t.Errorf("endless 基线: 中等策略应至少撑 5 波, got %d", g.Save.BestWave)
		} else {
			t.Logf("endless ✓ 撑到 wave %d", g.Save.BestWave)
		}
	})
}

func TestBalance_HardDifficultyReport(t *testing.T) {
	// Hard 数据报告 (不断言可通关 — Hard 设计上允许中等策略失败,
	// 仅记录数据供平衡决策)
	withTempSavePath(t, func() {
		wins := 0
		for idx := 0; idx < 20; idx++ {
			g := autoPlay(idx, DiffHard, true)
			if g.Phase == PhaseWon {
				wins++
				t.Logf("Lv%-2d ✓ %d★ lives %d/%d", idx+1,
					starsFor(g.Lives, g.StartLives), g.Lives, g.StartLives)
			} else {
				t.Logf("Lv%-2d ✗ wave %d/%d 失守", idx+1,
					g.WaveIdx+1, len(g.currentLevel().Waves))
			}
		}
		t.Logf("Hard 通关率: %d/20", wins)
	})
}

// TestBalance_HeroNetNonNegative: V8 — 英雄是净正/中性贡献, 绝不让难度上升。
// Hard 通关数 (有英雄) 应 ≥ (无英雄)。同时守护"英雄确实在仿真中在场"
// (若英雄被意外 nil/禁用, 两数相等仍通过, 但配合 hero_test 参战单测,
// 此测专守"加英雄不变难"这一不变量)。
func TestBalance_HeroNetNonNegative(t *testing.T) {
	withTempSavePath(t, func() {
		winsWith, winsWithout := 0, 0
		for idx := 0; idx < 20; idx++ {
			if autoPlay(idx, DiffHard, true).Phase == PhaseWon {
				winsWith++
			}
			if autoPlay(idx, DiffHard, false).Phase == PhaseWon {
				winsWithout++
			}
		}
		if winsWith < winsWithout {
			t.Errorf("英雄不应让通关数下降: 有英雄 %d < 无英雄 %d", winsWith, winsWithout)
		}
		t.Logf("Hard 通关数: 有英雄 %d / 无英雄 %d (英雄增益 %d)",
			winsWith, winsWithout, winsWith-winsWithout)
	})
}
