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
//
// V9: 英雄关内成长 (per-run XP/等级, 被动 XP 由 killEnemy 统一给) + cleave
// 技能 (auto-player 就绪即放, 平衡上界)。仿真证: 英雄升到 L3-L5 + 用 cleave
// 仍不破平衡 (Normal 20/20 / Hard 19/20 同难点 / endless 44) — per-run 后段
// 加载自限 (每关从 L1 起, 早期弱, 难点波在英雄弱时发生)。
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
			for _, path := range g.Paths {
				for _, pc := range path {
					dx, dy := float64(pc.X-p.X), float64(pc.Y-p.Y)
					if dx*dx+dy*dy <= r*r {
						n++
					}
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

// simBlueprint: V16 — 参数化塔蓝图 (kinds 按 built 序循环选型; branch
// 建成即设, L1→L2 升级时生效)。legacy = [A,A,A,F] branch 0, 与 V7.4
// 原策略序列逐位一致 (零回归)。
type simBlueprint struct {
	kinds  []TowerKind
	branch int
}

var legacyBlueprint = simBlueprint{kinds: []TowerKind{TArcher, TArcher, TArcher, TFrost}}

func simStrategy(g *Game, spots []Point, built *int) {
	simStrategyBP(g, spots, built, legacyBlueprint)
}

func simStrategyBP(g *Game, spots []Point, built *int, bp simBlueprint) {
	// 铺塔阶段
	if *built < len(spots) && *built < 8 {
		kind := bp.kinds[*built%len(bp.kinds)]
		if g.Gold >= towerSpecs[kind].Levels[0].Cost {
			g.Selected = kind
			g.Cursor = spots[*built]
			before := len(g.Towers)
			g.TryAction()
			if len(g.Towers) > before {
				g.Towers[len(g.Towers)-1].Branch = bp.branch
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
			g.CastMeteor(front.Pos(g.Paths[front.PathID]))
		}
	}
	// V9: 英雄技能就绪即放 (模拟玩家用 cleave, 平衡上界)
	if g.Hero != nil && g.Hero.AbilityReady() {
		g.CastHeroAbility()
	}
}

// autoPlay: 自动打完一关 (或 60000 步上限 = game-time 100 分钟)。
// heroEnabled=false 时 nil 掉英雄, 用于"有/无英雄"平衡对比。
// 默认 Knight; per-class 对比走 autoPlayClass (V10 P4)。
func autoPlay(levelIdx int, d Difficulty, heroEnabled bool) *Game {
	return autoPlayClass(levelIdx, d, heroEnabled, 0)
}

// autoPlayClass: V10 P4 — 指定英雄职业的仿真 (classIdx = heroClasses index)。
func autoPlayClass(levelIdx int, d Difficulty, heroEnabled bool, classIdx int) *Game {
	return autoPlayTree(levelIdx, d, heroEnabled, classIdx, 0)
}

// autoPlayTree: V11 P4 — 指定职业 + 技能树级的仿真 (treeLvl 0 = 无树
// = V10 基线; treeNodesPerClass = 满树上界)。
func autoPlayTree(levelIdx int, d Difficulty, heroEnabled bool, classIdx, treeLvl int) *Game {
	levels, err := LoadLevels()
	if err != nil {
		panic(err)
	}
	g := NewGame(levels, NewSave())
	for i := 1; i <= len(levels); i++ {
		g.Save.MarkCompleted(i) // 绕 unlock 链
	}
	g.Save.Difficulty = d
	g.Save.HeroChoice = classIdx // V10: beginRun 按此生成职业
	if treeLvl > 0 {             // V11: beginRun 按此快照 perk
		g.Save.TreeNodes = map[string]int{heroClasses[classIdx].Name: treeLvl}
	}
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

// autoPlayBlueprint: V16 — 指定塔蓝图的仿真 (默认 Knight / 无树 / 有英雄)。
func autoPlayBlueprint(levelIdx int, d Difficulty, bp simBlueprint) *Game {
	levels, err := LoadLevels()
	if err != nil {
		panic(err)
	}
	g := NewGame(levels, NewSave())
	for i := 1; i <= len(levels); i++ {
		g.Save.MarkCompleted(i)
	}
	g.Save.Difficulty = d
	g.StartLevel(levelIdx)
	spots := rankedTowerSpots(g)
	built := 0
	for step := 0; step < 60000 && g.Phase == PhasePlaying; step++ {
		if step%5 == 0 {
			simStrategyBP(g, spots, &built, bp)
		}
		g.Update(0.1)
		g.DrainSounds()
	}
	return g
}

// TestBalance_NewTowerMatrix: V16 — 新塔/分支蓝图矩阵。每个蓝图跑
// Normal 20 关, 断言不低于门槛 (18/20 — 允许风格弱项, 但不允许坍塌);
// legacy 基线仍 20/20 由 TestBalance_AllLevelsBeatableOnNormal 守。
func TestBalance_NewTowerMatrix(t *testing.T) {
	withTempSavePath(t, func() {
		blueprints := []struct {
			name string
			bp   simBlueprint
			min  int
		}{
			{"tesla-mix", simBlueprint{kinds: []TowerKind{TTesla, TArcher, TTesla, TFrost}}, 18},
			{"sniper-mix", simBlueprint{kinds: []TowerKind{TArcher, TSniper, TArcher, TFrost}}, 18},
			{"tesla-pure", simBlueprint{kinds: []TowerKind{TTesla, TTesla, TTesla, TFrost}}, 18},
			{"branchB-all", simBlueprint{kinds: []TowerKind{TArcher, TArcher, TArcher, TFrost}, branch: 1}, 18},
			{"gatling-mix", simBlueprint{kinds: []TowerKind{TArcher, TCannon, TArcher, TFrost}, branch: 1}, 18},
		}
		for _, tc := range blueprints {
			wins := 0
			var fails []int
			for idx := 0; idx < 20; idx++ {
				g := autoPlayBlueprint(idx, DiffNormal, tc.bp)
				if g.Phase == PhaseWon {
					wins++
				} else {
					fails = append(fails, idx+1)
				}
			}
			if wins < tc.min {
				t.Errorf("[%s] Normal %d/20 < 门槛 %d (失守: %v)", tc.name, wins, tc.min, fails)
			} else {
				t.Logf("[%s] Normal %d/20 (失守: %v)", tc.name, wins, fails)
			}
		}
	})
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

// TestBalance_HeroLevelReport: V9 诊断 — 报告各关英雄成长到几级 (验证升级
// 速度合理: 太慢则成长无感 / cleave@L3 够不着; 太快则中期碾压)。不断言, 供调参。
func TestBalance_HeroLevelReport(t *testing.T) {
	withTempSavePath(t, func() {
		for idx := 0; idx < 20; idx++ {
			g := autoPlay(idx, DiffNormal, true)
			lv := 0
			if g.Hero != nil {
				lv = g.Hero.Level
			}
			t.Logf("Lv%-2d → hero reached L%d", idx+1, lv)
		}
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

// ============================================================
// V10 P4: 多英雄平衡矩阵 — 三职业 × Normal/Hard
// ============================================================

// TestBalance_AllClassesBeatableOnNormal: Archer/Rogue 也保 Normal 全通
// (Knight 由 TestBalance_AllLevelsBeatableOnNormal 守, 不重跑)。
func TestBalance_AllClassesBeatableOnNormal(t *testing.T) {
	withTempSavePath(t, func() {
		for classIdx := 1; classIdx < len(heroClasses); classIdx++ {
			name := heroClasses[classIdx].Name
			for idx := 0; idx < 20; idx++ {
				g := autoPlayClass(idx, DiffNormal, true, classIdx)
				if g.Phase != PhaseWon {
					t.Errorf("[%s] Lv%-2d ✗ 不可通关 (lives=%d wave=%d/%d)",
						name, idx+1, g.Lives, g.WaveIdx+1, len(g.currentLevel().Waves))
				}
			}
			t.Logf("[%s] Normal 20 关跑完", name)
		}
	})
}

// TestBalance_ClassMatrixHard: 三职业 Hard 矩阵 — 每职业通关数 ≥ 无英雄
// 基线 (净非负, 同 V8 HeroNet 不变量推广到 per-class)。职业间差异只记录
// 不断言 (设计允许强弱有别, 但谁都不能比没英雄更糟)。
func TestBalance_ClassMatrixHard(t *testing.T) {
	withTempSavePath(t, func() {
		baseline := 0
		for idx := 0; idx < 20; idx++ {
			if autoPlay(idx, DiffHard, false).Phase == PhaseWon {
				baseline++
			}
		}
		for classIdx := 0; classIdx < len(heroClasses); classIdx++ {
			name := heroClasses[classIdx].Name
			wins := 0
			for idx := 0; idx < 20; idx++ {
				if autoPlayClass(idx, DiffHard, true, classIdx).Phase == PhaseWon {
					wins++
				}
			}
			if wins < baseline {
				t.Errorf("[%s] Hard 通关 %d < 无英雄基线 %d (净负 = 设计失败)",
					name, wins, baseline)
			}
			t.Logf("[%s] Hard %d/20 (基线 %d, 净增益 %+d)", name, wins, baseline, wins-baseline)
		}
	})
}

// ============================================================
// V11 P4: 技能树上下界 — 无树零回归 (下界, 由上方 V10 矩阵守) +
// 满树上界 (决策 D: 允许碾压 Hard, 但不得低于无树)
// ============================================================

// TestBalance_FullTreeMatrixHard: 满树 Hard 矩阵 — 每职业满树通关数 ≥
// 同职业无树 (perk 必须非负贡献); 允许 20/20 (决策 D earned power)。
func TestBalance_FullTreeMatrixHard(t *testing.T) {
	withTempSavePath(t, func() {
		for classIdx := 0; classIdx < len(heroClasses); classIdx++ {
			name := heroClasses[classIdx].Name
			noTree, fullTree := 0, 0
			for idx := 0; idx < 20; idx++ {
				if autoPlayTree(idx, DiffHard, true, classIdx, 0).Phase == PhaseWon {
					noTree++
				}
				if autoPlayTree(idx, DiffHard, true, classIdx, treeNodesPerClass).Phase == PhaseWon {
					fullTree++
				}
			}
			if fullTree < noTree {
				t.Errorf("[%s] 满树 %d < 无树 %d (perk 净负 = 设计失败)", name, fullTree, noTree)
			}
			t.Logf("[%s] Hard 无树 %d/20 → 满树 %d/20 (perk 增益 %+d)",
				name, noTree, fullTree, fullTree-noTree)
		}
	})
}

// TestBalance_FullTreeNormalAllClear: 满树不破 Normal (仍全通, 平衡上界守护)。
func TestBalance_FullTreeNormalAllClear(t *testing.T) {
	withTempSavePath(t, func() {
		for classIdx := 0; classIdx < len(heroClasses); classIdx++ {
			name := heroClasses[classIdx].Name
			for idx := 0; idx < 20; idx++ {
				g := autoPlayTree(idx, DiffNormal, true, classIdx, treeNodesPerClass)
				if g.Phase != PhaseWon {
					t.Errorf("[%s 满树] Lv%-2d ✗ 不可通关 (lives=%d wave=%d/%d)",
						name, idx+1, g.Lives, g.WaveIdx+1, len(g.currentLevel().Waves))
				}
			}
			t.Logf("[%s] 满树 Normal 20 关跑完", name)
		}
	})
}
