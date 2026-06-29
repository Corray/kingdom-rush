//go:build !term

// Kenney TD top-down tilesheet 加载 + sprite atlas (V3 接入)。
//
// Pack: kenney_tower-defense-top-down (CC0, public domain)
//
//	https://kenney.nl/assets/tower-defense-top-down
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
	spriteGrass        = 25 // V7.1: 24 是泥土+草边过渡 tile, 平铺成竖条纹 (V3 起选错); 25 = 无缝纯草
	spriteDirtPath     = 200
	spriteTowerArcher  = 249
	spriteTowerCannon  = 250
	spriteTowerMagic1  = 206 // 1 rocket
	spriteTowerMagic2  = 204 // 2 rockets
	spriteTowerMagic3  = 205 // 4 rockets
	spriteTowerFrost   = 227 // V5 Phase 4: 灰色炮塔 (调研确认未占用)
	spriteFrostShot    = 276 // V5 Phase 4: 白色小圆弹
	spriteEnemyNormal  = 248
	spriteEnemyFast    = 246
	spriteEnemyGlider  = 270
	spriteEnemyBoss    = 247
	spriteEnemySpawner = 268 // V3.6: 绿色容器, "summoning pod" 形象
	spriteEnemyShield  = 265 // V13: 绿色方块 (装甲车) — drawTileTint 灰色化
	spriteEnemyRegen   = 291 // V13: 绿色药瓶 (回血语义)
	spriteEnemyHealer  = 292 // V13: 黄色药瓶 (治疗者)
	// V14: 英雄 sprite — 使用车辆/飞机造型 + drawTileTint 职业色
	spriteHeroKnight = 266 // 棕色方块 (坦克/战士)
	spriteHeroArcher = 269 // 灰色飞机 (远程射手)
	spriteHeroRogue  = 253 // 红色火箭 (高速突袭)
	// V3 Phase 2 新增
	spriteHitFire = 295 // 橙色火焰 (用作 hit explosion)
	spriteDigit1  = 277
	spriteDigit2  = 278
	spriteDigit3  = 279
	spriteBullet  = 275 // 灰白小子弹 (Archer)
	// V3 Phase 3b: 投射物按塔型分
	spriteCannonball   = 252 // 红色大火箭 (Cannon)
	spriteMagicMissile = 251 // 红色小火箭 (Magic)
	// V3 Phase 4: UI icon
	spriteGold  = 287 // $ symbol (gold counter)
	spriteLives = 289 // + cross (lives counter)
	// V7.2 M3: 地图装饰 (透明底独立 sprite, 调研见 /tmp/rows4-6)
	spriteDecorTree   = 134
	spriteDecorBush   = 131
	spriteDecorBushSm = 132
	spriteDecorRock   = 137
)

// bulletSpriteID 按塔型返回投射物 sprite
func bulletSpriteID(k TowerKind) int {
	switch k {
	case TCannon:
		return spriteCannonball
	case TMagic:
		return spriteMagicMissile
	case TFrost:
		return spriteFrostShot
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

// drawTileRot: V4 Phase 3 — 以 tile 中心为轴旋转 angle (弧度) 绘制。
// 变换顺序: 平移中心到原点 → 旋转 → 缩放 → 平移到 cell 中心。
// scaleMul / alpha 语义同 drawTileAt。
func drawTileRot(dst *ebiten.Image, tileID int, x, y float32, scaleMul, alpha, angle float64) {
	sub := tileSubImage(tileID)
	if sub == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(sheetTileSrc)/2, -float64(sheetTileSrc)/2)
	op.GeoM.Rotate(angle)
	scale := float64(cellPx) / float64(sheetTileSrc) * scaleMul
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x)+float64(cellPx)/2, float64(y)+float64(cellPx)/2)
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

// drawTileRotTint: V14 — 旋转 + 颜色乘数绘制。
func drawTileRotTint(dst *ebiten.Image, tileID int, x, y float32, scaleMul, alpha, angle float64, cr, cg, cb float32) {
	sub := tileSubImage(tileID)
	if sub == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(sheetTileSrc)/2, -float64(sheetTileSrc)/2)
	op.GeoM.Rotate(angle)
	scale := float64(cellPx) / float64(sheetTileSrc) * scaleMul
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x)+float64(cellPx)/2, float64(y)+float64(cellPx)/2)
	if alpha < 1.0 {
		op.ColorScale.ScaleAlpha(float32(alpha))
	}
	op.ColorScale.Scale(cr, cg, cb, 1)
	dst.DrawImage(sub, op)
}

// drawTileTint: V7.3 D1 — 带颜色乘数绘制 (Frost 塔冰蓝 tint 用)。
func drawTileTint(dst *ebiten.Image, tileID int, x, y float32, cr, cg, cb float32) {
	sub := tileSubImage(tileID)
	if sub == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	scale := float64(cellPx) / float64(sheetTileSrc)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.Scale(cr, cg, cb, 1)
	dst.DrawImage(sub, op)
}

// frostTint: Frost 塔的冰蓝色乘数 (灰色原 sprite × 此值 = 偏蓝)。
const frostTintR, frostTintG, frostTintB = 0.7, 0.95, 1.3

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
	case TFrost:
		return spriteTowerFrost
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
	case ESpawner:
		return spriteEnemySpawner
	case EShield:
		return spriteEnemyShield
	case ERegen:
		return spriteEnemyRegen
	case EHealer:
		return spriteEnemyHealer
	}
	return spriteEnemyNormal
}

func heroSpriteID(c *HeroClass) int {
	switch c.Name {
	case "Archer":
		return spriteHeroArcher
	case "Rogue":
		return spriteHeroRogue
	default:
		return spriteHeroKnight
	}
}

// decorSpriteID: V7.2 M3 — 装饰 kind (0-3) → sprite。
func decorSpriteID(kind int) int {
	switch kind {
	case 1:
		return spriteDecorBush
	case 2:
		return spriteDecorBushSm
	case 3:
		return spriteDecorRock
	default:
		return spriteDecorTree
	}
}
