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
	// 10 关存档兼容: 右列点击不应命中 (numLevels=10)
	if _, ok := menuRowAtPixel(windowW-10, startY+5, 10); ok {
		t.Errorf("右列在只有 10 关时不应命中")
	}
}
