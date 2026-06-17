// V12 测试: 多入口 path — Paths() / spawn 均摊 / per-path 取位 / 汇流 / 并集避让。
package main

import "testing"

// newDualPathGame: 双路 fixture — 上路 (y=5) 与下路 (y=9) 末尾汇合到 (5,5)。
// 直接设 Path/Path2 字段绕过 yaml (同 newTestGame 模式)。
func newDualPathGame() *Game {
	levels := []Level{{
		ID: 1, Name: "Dual",
		StartGold: 200, StartLives: 50,
		Path:  []Point{{0, 5}, {1, 5}, {2, 5}, {3, 5}, {4, 5}, {5, 5}},
		Path2: []Point{{0, 9}, {1, 9}, {2, 9}, {3, 9}, {4, 9}, {5, 5}}, // 汇合到 (5,5)
		Waves: []WaveSpec{{Enemies: []EnemyKind{ENormal, ENormal, ENormal, ENormal}}},
	}}
	g := NewGame(levels, NewSave())
	g.StartLevel(0)
	return g
}

// TestLevelPathsSingleVsDual: Paths() — 单路返回 1 条, 双路返回 2 条。
func TestLevelPathsSingleVsDual(t *testing.T) {
	single := Level{Path: []Point{{0, 0}, {1, 0}}}
	if got := single.Paths(); len(got) != 1 {
		t.Errorf("single-path Paths() len = %d, want 1", len(got))
	}
	dual := Level{Path: []Point{{0, 0}, {1, 0}}, Path2: []Point{{0, 2}, {1, 2}}}
	if got := dual.Paths(); len(got) != 2 {
		t.Errorf("dual-path Paths() len = %d, want 2", len(got))
	}
}

// TestBeginRunDualPaths: 双路关 beginRun 后 g.Paths 含两条; pathLookup 是并集。
func TestBeginRunDualPaths(t *testing.T) {
	g := newDualPathGame()
	if len(g.Paths) != 2 {
		t.Fatalf("dual game should have 2 paths, got %d", len(g.Paths))
	}
	// 并集: 两路所有 cell 都不可建塔
	for _, p := range []Point{{2, 5}, {3, 9}, {5, 5}} {
		if !g.pathContains(p) {
			t.Errorf("cell %v should be in pathLookup union", p)
		}
	}
	// 非 path cell 可建
	if g.pathContains(Point{7, 7}) {
		t.Error("(7,7) should not be a path cell")
	}
}

// TestSpawnAlternatesPaths: 双路 spawn 按计数交替 [0,1,0,1] (均摊)。
func TestSpawnAlternatesPaths(t *testing.T) {
	g := newDualPathGame()
	for i := 0; i < 400 && len(g.Enemies) < 4; i++ {
		g.Update(0.1)
	}
	if len(g.Enemies) < 4 {
		t.Fatalf("expected 4 enemies spawned, got %d", len(g.Enemies))
	}
	want := []int{0, 1, 0, 1}
	for i := 0; i < 4; i++ {
		if g.Enemies[i].PathID != want[i] {
			t.Errorf("enemy %d PathID = %d, want %d (均摊交替)", i, g.Enemies[i].PathID, want[i])
		}
	}
}

// TestSinglePathSpawnAllZero: 单路 spawn 全 PathID 0 (零回归)。
func TestSinglePathSpawnAllZero(t *testing.T) {
	g := newTestGame() // 单路
	for i := 0; i < 200 && len(g.Enemies) < 1; i++ {
		g.Update(0.1)
	}
	if len(g.Enemies) == 0 {
		t.Fatal("expected at least 1 enemy")
	}
	for i, e := range g.Enemies {
		if e.PathID != 0 {
			t.Errorf("single-path enemy %d PathID = %d, want 0", i, e.PathID)
		}
	}
}

// TestEnemyPosPerPath: 双路敌按各自 PathID 取位 — 同 idx 不同 path 坐标不同。
func TestEnemyPosPerPath(t *testing.T) {
	g := newDualPathGame()
	e0 := &Enemy{Kind: ENormal, PathID: 0, PathIdx: 0}
	e1 := &Enemy{Kind: ENormal, PathID: 1, PathIdx: 0}
	p0 := e0.Pos(g.Paths[e0.PathID])
	p1 := e1.Pos(g.Paths[e1.PathID])
	if p0 != (Point{0, 5}) {
		t.Errorf("path0 idx0 = %v, want (0,5)", p0)
	}
	if p1 != (Point{0, 9}) {
		t.Errorf("path1 idx0 = %v, want (0,9)", p1)
	}
}

// TestConvergentTailMerge: 汇流 — 两路末尾 cell 相同; path1 敌走到 len-1 escaped。
func TestConvergentTailMerge(t *testing.T) {
	g := newDualPathGame()
	p0, p1 := g.Paths[0], g.Paths[1]
	if p0[len(p0)-1] != p1[len(p1)-1] {
		t.Fatalf("paths must converge: tail0=%v tail1=%v", p0[len(p0)-1], p1[len(p1)-1])
	}
	// path1 敌推进到末尾 → escaped (用各自 path 的 len 判定)
	livesBefore := g.Lives
	g.Enemies = []*Enemy{{Kind: ENormal, HP: 999, MaxHP: 999, PathID: 1,
		PathIdx: float64(len(p1) - 1)}}
	g.Update(0.1)
	if !g.Enemies[0].Escaped {
		t.Error("path1 enemy at tail should escape")
	}
	if g.Lives != livesBefore-1 {
		t.Errorf("escape should cost 1 life: %d → %d", livesBefore, g.Lives)
	}
}
