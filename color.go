// 解耦 UI 库的颜色抽象。
// entities.go 中 spec table 用 RGB,renderer 各自转换:
//
//	terminal: RGB → tcell.NewRGBColor
//	ebiten:   RGB → color.RGBA
//
// 这样 game / entities / level 层不依赖任何 UI 库。
package main

type RGB struct {
	R, G, B uint8
}

func rgb(r, g, b uint8) RGB { return RGB{r, g, b} }
