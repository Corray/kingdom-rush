//go:build !term

// Kenney TD top-down tilesheet 加载 + sprite atlas (V3 接入)。
//
// Pack: kenney_tower-defense-top-down (CC0, public domain)
//       https://kenney.nl/assets/tower-defense-top-down
//
// Tilesheet: 1472×832 px, 23 cols × 13 rows × 64 px tile
// 编号: tile001 在 (col=0, row=0), tileN 在 (col=(N-1)%23, row=(N-1)/23)
package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/sprites/towerDefense_tilesheet.png
var tilesheetBytes []byte

const (
	sheetCols    = 23
	sheetTileSrc = 64 // 源 tile 尺寸 px
)

// 关键 sprite ID 映射 (按 visual id, 详细 catalog 见 ADR-004)
const (
	spriteGrass       = 24
	spriteDirtPath    = 200
	spriteTowerArcher = 249
	spriteTowerCannon = 250
	spriteTowerMagic1 = 206 // 1 rocket
	spriteTowerMagic2 = 204 // 2 rockets
	spriteTowerMagic3 = 205 // 4 rockets
	spriteEnemyNormal = 248
	spriteEnemyFast   = 246
	spriteEnemyGlider = 270
	spriteEnemyBoss   = 247
	// V3 Phase 2 新增
	spriteHitFire = 295 // 橙色火焰 (用作 hit explosion)
	spriteDigit1  = 277
	spriteDigit2  = 278
	spriteDigit3  = 279
	spriteBullet  = 275 // 灰白小子弹 (Archer)
	// V3 Phase 3b: 投射物按塔型分
	spriteCannonball  = 252 // 红色大火箭 (Cannon)
	spriteMagicMissile = 251 // 红色小火箭 (Magic)
)

// bulletSpriteID 按塔型返回投射物 sprite
func bulletSpriteID(k TowerKind) int {
	switch k {
	case TCannon:
		return spriteCannonball
	case TMagic:
		return spriteMagicMissile
	default: // TArcher 与默认
		return spriteBullet
	}
}

// digitSpriteID 返回 level 数字对应的 sprite ID (1-3)
func digitSpriteID(level int) int {
	switch level {
	case 1:
		return spriteDigit1
	case 2:
		return spriteDigit2
	default:
		return spriteDigit3
	}
}

// drawTileAt 绘制 tile 在指定位置, 自定义 scale + alpha
// scaleMul: 相对于默认 cellPx/64 的额外缩放 (1.0=cell 大小, 0.4=mini badge, 1.4=Boss)
// alpha: 0..1, 用于特效 fade
func drawTileAt(dst *ebiten.Image, tileID int, x, y float32, scaleMul, alpha float64) {
	sub := tileSubImage(tileID)
	if sub == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	scale := float64(cellPx) / float64(sheetTileSrc) * scaleMul
	op.GeoM.Scale(scale, scale)
	// 居中: scaleMul != 1 时,offset 把 sprite 居中到 cell
	offset := (1.0 - scaleMul) * float64(cellPx) / 2
	op.GeoM.Translate(float64(x)+offset, float64(y)+offset)
	if alpha < 1.0 {
		op.ColorScale.ScaleAlpha(float32(alpha))
	}
	dst.DrawImage(sub, op)
}

var tilesheet *ebiten.Image

func loadTilesheet() error {
	if tilesheet != nil {
		return nil
	}
	img, _, err := image.Decode(bytes.NewReader(tilesheetBytes))
	if err != nil {
		return fmt.Errorf("decode tilesheet: %w", err)
	}
	tilesheet = ebiten.NewImageFromImage(img)
	return nil
}

// tileSubImage 返回指定 tile ID 的 sub-image (64×64 px from sheet)
func tileSubImage(tileID int) *ebiten.Image {
	if tilesheet == nil {
		return nil
	}
	idx := tileID - 1 // tile001 → idx 0
	col := idx % sheetCols
	row := idx / sheetCols
	rect := image.Rect(col*sheetTileSrc, row*sheetTileSrc,
		(col+1)*sheetTileSrc, (row+1)*sheetTileSrc)
	return tilesheet.SubImage(rect).(*ebiten.Image)
}

// drawTile 在指定屏幕位置画 tile, 自动按 cellPx 缩放
func drawTile(dst *ebiten.Image, tileID int, x, y float32) {
	sub := tileSubImage(tileID)
	if sub == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	scale := float64(cellPx) / float64(sheetTileSrc)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	dst.DrawImage(sub, op)
}

// drawTileScaled 自定义缩放 (Boss 用 1.5x)
func drawTileScaled(dst *ebiten.Image, tileID int, x, y float32, extraScale float64) {
	sub := tileSubImage(tileID)
	if sub == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	scale := float64(cellPx) / float64(sheetTileSrc) * extraScale
	op.GeoM.Scale(scale, scale)
	// 中心对齐: extraScale > 1 时, offset 把 sprite 居中到 cell
	offset := (1.0 - extraScale) * float64(cellPx) / 2
	op.GeoM.Translate(float64(x)+offset, float64(y)+offset)
	dst.DrawImage(sub, op)
}

// towerSpriteID 返回塔型 + 等级对应的 sprite ID
func towerSpriteID(kind TowerKind, level int) int {
	switch kind {
	case TArcher:
		return spriteTowerArcher
	case TCannon:
		return spriteTowerCannon
	case TMagic:
		switch level {
		case 1:
			return spriteTowerMagic1
		case 2:
			return spriteTowerMagic2
		default:
			return spriteTowerMagic3
		}
	}
	return spriteTowerArcher
}

// enemySpriteID 返回敌人型对应的 sprite ID
func enemySpriteID(kind EnemyKind) int {
	switch kind {
	case ENormal:
		return spriteEnemyNormal
	case EFast:
		return spriteEnemyFast
	case EGlider:
		return spriteEnemyGlider
	case EBoss:
		return spriteEnemyBoss
	}
	return spriteEnemyNormal
}
