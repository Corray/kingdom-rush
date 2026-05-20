//go:build !term

// Truetype 字体加载 (V3 Phase 5c)。
//
// 用 Go 团队自带的 gomono.TTF (Go Mono, Apache 2.0 license),
// 通过 golang.org/x/image/font/gofont/gomono 子包获得 raw bytes,
// 无需 fetch 外部 ttf 资源。
//
// drawText() wrapper 兼容 ebitenutil.DebugPrintAt 的 top-left anchor
// 坐标系 — 内部加 ascent offset 转 baseline (ttf text.Draw 用 baseline)。
package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	etext "github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
)

var (
	gameFontFace font.Face
	textWhite    = color.RGBA{R: 240, G: 240, B: 240, A: 255}
)

const fontAscent = 10 // 12pt gomono ascent, top-left → baseline 偏移

func loadGameFont() error {
	if gameFontFace != nil {
		return nil
	}
	tt, err := opentype.Parse(gomono.TTF)
	if err != nil {
		return err
	}
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}
	gameFontFace = face
	return nil
}

// drawText: 兼容 ebitenutil.DebugPrintAt (top-left anchor) 坐标系
// 内部 +fontAscent 转为 baseline anchor
func drawText(screen *ebiten.Image, str string, x, y int) {
	if gameFontFace == nil {
		return
	}
	etext.Draw(screen, str, gameFontFace, x, y+fontAscent, textWhite)
}
