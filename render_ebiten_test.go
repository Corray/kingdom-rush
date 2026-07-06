//go:build !term

// 依赖 ebiten-only 符号 (bgm.go / ebiten_renderer.go) 的测试。
//
// 从 game_test.go 剥出 (2026-06-09): 默认 build 的 ebiten 在无窗口服务器的
// 环境会在包 init 阶段 panic (GLFW initializeGLFW → currentMouseLocation nil),
// 导致整个测试二进制无法启动。把这两个 ebiten-only 测试隔到 `!term` 文件后,
// 其余逻辑测试 (game / hero / balance) 可在 `-tags term` (tcell, 不 init ebiten)
// 下 headless 运行。本文件仍随默认 build 测试 (有窗口服务器时) 一起跑。
package main

import "testing"

func TestBGM_TrackForPhase(t *testing.T) {
	cases := []struct {
		phase GamePhase
		want  bgmTrack
	}{
		{PhaseLevelSelect, bgmMenu},
		{PhasePlaying, bgmBattle},
		{PhaseWon, bgmNone},
		{PhaseLost, bgmNone},
	}
	for _, c := range cases {
		if got := bgmFor(c.phase); got != c.want {
			t.Errorf("bgmFor(%v) = %v, want %v", c.phase, got, c.want)
		}
	}
}

func TestMenu_RowAtPixelTwoColumns(t *testing.T) {
	const startY, rowH = 60, 36
	cases := []struct {
		mx, my int
		want   int
		wantOK bool
	}{
		{10, startY + 5, 0, true},                     // 左列首行
		{10, startY + 9*rowH + 5, 9, true},            // 左列末行
		{windowW - 10, startY + 5, 10, true},          // 右列首行 = Lv11
		{windowW - 10, startY + 9*rowH + 5, 19, true}, // 右列末行 = Lv20
		{10, startY - 5, 0, false},                    // 上方界外
		{10, startY + 10*rowH + 5, 0, false},          // 下方界外
	}
	for _, c := range cases {
		got, ok := menuRowAtPixel(c.mx, c.my, 20)
		if ok != c.wantOK || (ok && got != c.want) {
			t.Errorf("menuRowAtPixel(%d,%d) = (%d,%v), want (%d,%v)",
				c.mx, c.my, got, ok, c.want, c.wantOK)
		}
	}
	// V17 动态列: 10 关 = 单列全宽行 (渲染与 hitbox 同式), 右缘点击
	// 命中该行 (V6 两列硬编码时代此处断言"不命中", 语义已随动态列更新)
	if got, ok := menuRowAtPixel(windowW-10, startY+5, 10); !ok || got != 0 {
		t.Errorf("单列全宽: 右缘应命中首行, got (%d,%v)", got, ok)
	}
}

// V17 P3: 三列 hitbox — 30 关时第三列命中 L21-30, 与 drawLevelSelect 同式。
func TestMenu_RowAtPixelThreeColumns(t *testing.T) {
	const startY, rowH, rowX, colGap = 60, 36, 8, 8
	if menuCols(30) != 3 || menuCols(20) != 2 || menuCols(10) != 1 {
		t.Fatalf("menuCols: 30→%d 20→%d 10→%d, want 3/2/1",
			menuCols(30), menuCols(20), menuCols(10))
	}
	colW := (windowW - rowX*2 - colGap*2) / 3
	col2X := rowX + 2*(colW+colGap) + 10 // 第三列内部
	cases := []struct {
		mx, my int
		want   int
	}{
		{col2X, startY + 5, 20},            // 第三列首行 = Lv21
		{col2X, startY + 9*rowH + 5, 29},   // 第三列末行 = Lv30
		{rowX + 10, startY + 5, 0},         // 第一列首行仍 = Lv1
		{rowX + colW + colGap + 10, startY + 5, 10}, // 第二列首行 = Lv11
	}
	for _, c := range cases {
		got, ok := menuRowAtPixel(c.mx, c.my, 30)
		if !ok || got != c.want {
			t.Errorf("menuRowAtPixel(%d,%d,30) = (%d,%v), want (%d,true)",
				c.mx, c.my, got, ok, c.want)
		}
	}
}

// V17 audit #1: 主题分区守护 — L21+ Twilight 与 Lava/Snow/Desert/Forest 互异。
func TestMapTheme_FiveZones(t *testing.T) {
	zones := []int{1, 8, 14, 18, 21} // 每分区代表关
	type theme struct{ r, g, b float32 }
	seen := map[theme]int{}
	for _, id := range zones {
		r, g, b := mapThemeGrass(id)
		th := theme{r, g, b}
		if prev, dup := seen[th]; dup {
			t.Errorf("Lv%d 与 Lv%d 主题 tint 相同 %v (分区应互异)", id, prev, th)
		}
		seen[th] = id
	}
	// L21 与 L30 同区; L20 (Lava) 与 L21 (Twilight) 必须不同
	r21, g21, b21 := mapThemeGrass(21)
	r30, g30, b30 := mapThemeGrass(30)
	if r21 != r30 || g21 != g30 || b21 != b30 {
		t.Error("L21 与 L30 应同为 Twilight 分区")
	}
	r20, g20, b20 := mapThemeGrass(20)
	if r20 == r21 && g20 == g21 && b20 == b21 {
		t.Error("L20 (Lava) 与 L21 (Twilight) 不应同主题")
	}
	p21, _ := mapThemePath(21)
	p20, _ := mapThemePath(20)
	if p21 == p20 {
		t.Error("L21 路径色不应沿用 Lava")
	}
}
