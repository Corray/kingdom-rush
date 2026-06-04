//go:build !term

// kingdom-rush V2 desktop entry (ebiten,默认 build)。
//
// Build:
//
//	go build .                  → ebiten desktop binary (default)
//	go build -tags term .       → V1.7 terminal binary (opt-in,见 term_main.go)
//
// 控制:
//
//	Menu:   1-9, 0 选关 | Q/Esc 退出
//	Game:   Arrows 移光标 | 1/2/3 选塔型 | Space 建塔/升塔
//	        M 回菜单 | Q/Esc 退出
package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	levels, err := LoadLevels()
	if err != nil {
		fmt.Fprintln(os.Stderr, "load levels:", err)
		os.Exit(1)
	}
	if len(levels) == 0 {
		fmt.Fprintln(os.Stderr, "no levels loaded")
		os.Exit(1)
	}

	save, err := LoadSave()
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: load save failed:", err, "(starting fresh)")
		save = NewSave()
	}

	g := NewGame(levels, save)

	ebiten.SetWindowSize(windowW, windowH)
	ebiten.SetWindowTitle("Gopher Defense")
	ebiten.SetTPS(60)

	eg := NewEbitenGame(g)
	if err := ebiten.RunGame(eg); err != nil {
		fmt.Fprintln(os.Stderr, "run game:", err)
		os.Exit(1)
	}
}
