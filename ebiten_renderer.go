//go:build !term

// Ebiten desktop 渲染 (V2 默认 build)。
//
// Window: 840 × 560 像素 (cellPx=28, mapW=30, mapH=15, + 顶底 UI 带)
// 视觉:
//   - path: 棕色填充矩形
//   - tower: 圆 + 字母标签 (颜色随塔型)
//   - enemy: 圆, boss 半径更大, HP bar 在上方
//   - cursor: 黄色描边方框
//   - 起点/终点: 绿/红填充 + S/E 字符
//
// Terminal 模式见 term_main.go / render.go (build tag `term`)。
package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	cellPx     = 32 // V3: 28→32 for sprite (64 px source / 2 = 32 visual)
	gameAreaW  = mapW * cellPx
	gameAreaH  = mapH * cellPx
	topBarH    = 30
	bottomBarH = 110
	windowW    = gameAreaW
	windowH    = topBarH + gameAreaH + bottomBarH
)

func ebitenColor(c RGB) color.RGBA {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: 255}
}

var (
	eColBg     = color.RGBA{R: 30, G: 30, B: 40, A: 255}
	eColPathBg = color.RGBA{R: 139, G: 69, B: 19, A: 255}
	eColStart  = color.RGBA{R: 50, G: 200, B: 80, A: 255}
	eColEnd    = color.RGBA{R: 200, G: 80, B: 80, A: 255}
	eColCursor = color.RGBA{R: 255, G: 220, B: 80, A: 255}
	eColSelHi  = color.RGBA{R: 80, G: 80, B: 140, A: 255}
	eColHpBg   = color.RGBA{R: 80, G: 80, B: 80, A: 255}
	eColHpFg   = color.RGBA{R: 60, G: 220, B: 80, A: 255}
)

// EbitenGame 实现 ebiten.Game interface。
type EbitenGame struct {
	game *Game
	last time.Time
	quit bool
}

func NewEbitenGame(g *Game) *EbitenGame {
	// V3: 加载 sprite tilesheet (失败不致命, fallback 旧 circle 渲染)
	if err := loadTilesheet(); err != nil {
		fmt.Println("warning: tilesheet load failed:", err, "(falling back to primitive render)")
	}
	return &EbitenGame{game: g, last: time.Now()}
}

func (eg *EbitenGame) Update() error {
	if eg.quit {
		return ebiten.Termination
	}
	now := time.Now()
	dt := now.Sub(eg.last).Seconds()
	if dt > 0.1 {
		dt = 0.1 // cap to avoid lag spike at startup
	}
	eg.last = now
	eg.handleInput()
	eg.game.Update(dt)
	return nil
}

func (eg *EbitenGame) Draw(screen *ebiten.Image) {
	screen.Fill(eColBg)
	if eg.game.Phase == PhaseLevelSelect {
		eg.drawLevelSelect(screen)
	} else {
		eg.drawGame(screen)
	}
}

func (eg *EbitenGame) Layout(_, _ int) (int, int) {
	return windowW, windowH
}

// ============================================================
// Input
// ============================================================

func (eg *EbitenGame) handleInput() {
	g := eg.game
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		eg.quit = true
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.BackToMenu()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.Phase == PhasePlaying {
		g.TryAction()
	}
	if g.Phase == PhasePlaying {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			g.MoveCursor(0, -1)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			g.MoveCursor(0, 1)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			g.MoveCursor(-1, 0)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			g.MoveCursor(1, 0)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDigit1) {
			g.Selected = TArcher
			g.Msg = "Selected Archer"
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDigit2) {
			g.Selected = TCannon
			g.Msg = "Selected Cannon"
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDigit3) {
			g.Selected = TMagic
			g.Msg = "Selected Magic"
		}
		// V2.7: 鼠标 hover → cursor 跟随; 左键 → 等同 Space
		mx, my := ebiten.CursorPosition()
		if p, ok := pixelToCell(mx, my); ok {
			g.Cursor = p
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if _, ok := pixelToCell(mx, my); ok {
				g.TryAction()
			}
		}
		return
	}
	if g.Phase == PhaseLevelSelect {
		digitKeys := []ebiten.Key{
			ebiten.KeyDigit1, ebiten.KeyDigit2, ebiten.KeyDigit3, ebiten.KeyDigit4,
			ebiten.KeyDigit5, ebiten.KeyDigit6, ebiten.KeyDigit7, ebiten.KeyDigit8,
			ebiten.KeyDigit9,
		}
		for i, k := range digitKeys {
			if inpututil.IsKeyJustPressed(k) {
				g.StartLevel(i)
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDigit0) {
			g.StartLevel(9)
		}
		// V2.7: menu 阶段鼠标点击 row 启动关卡
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			_, my := ebiten.CursorPosition()
			if i, ok := menuRowAtPixel(my, len(g.Levels)); ok {
				g.StartLevel(i)
			}
		}
	}
}

// ============================================================
// Drawing helpers
// ============================================================

func fillRect(screen *ebiten.Image, x, y, w, h float32, c color.RGBA) {
	vector.DrawFilledRect(screen, x, y, w, h, c, false)
}

func strokeRect(screen *ebiten.Image, x, y, w, h float32, c color.RGBA, thick float32) {
	vector.StrokeRect(screen, x, y, w, h, thick, c, false)
}

func fillCircle(screen *ebiten.Image, cx, cy, r float32, c color.RGBA) {
	vector.DrawFilledCircle(screen, cx, cy, r, c, false)
}

func cellPos(p Point) (float32, float32) {
	return float32(p.X * cellPx), float32(topBarH + p.Y*cellPx)
}

// pixelToCell 把屏幕像素坐标转回 game cell, 返回 (cell, inGameArea)。
// 鼠标在 game area 外(top/bottom UI bar)时返回 false。
func pixelToCell(mx, my int) (Point, bool) {
	if my < topBarH || my >= topBarH+gameAreaH {
		return Point{}, false
	}
	if mx < 0 || mx >= gameAreaW {
		return Point{}, false
	}
	gx := mx / cellPx
	gy := (my - topBarH) / cellPx
	if gx < 0 || gx >= mapW || gy < 0 || gy >= mapH {
		return Point{}, false
	}
	return Point{X: gx, Y: gy}, true
}

// menuRowAtPixel 返回鼠标所在的 level select row index (0..len(levels)-1), false 表示不在行上
func menuRowAtPixel(my, numLevels int) (int, bool) {
	const startY = 50
	const rowH = 22
	i := (my - startY) / rowH
	if i < 0 || i >= numLevels {
		return 0, false
	}
	return i, true
}

// ============================================================
// Level Select
// ============================================================

func (eg *EbitenGame) drawLevelSelect(screen *ebiten.Image) {
	g := eg.game

	ebitenutil.DebugPrintAt(screen, " Kingdom Rush V2  —  Select Level", 8, 4)

	completed := 0
	for _, lv := range g.Levels {
		if g.Save.IsCompleted(lv.ID) {
			completed++
		}
	}
	ebitenutil.DebugPrintAt(screen,
		fmt.Sprintf(" Progress: %d / %d cleared", completed, len(g.Levels)),
		8, 24)

	for i, lv := range g.Levels {
		y := 50 + i*22
		key := i + 1
		keyStr := fmt.Sprintf("%d", key)
		if key == 10 {
			keyStr = "0"
		}
		status := "[    ]"
		if g.Save.IsCompleted(lv.ID) {
			status = "[DONE]"
		} else if !g.Save.IsUnlocked(lv.ID) {
			status = "[LOCK]"
		}
		line := fmt.Sprintf("  [%s] %s  Lv %2d  %-18s  (waves:%d  gold:%d  lives:%d)",
			keyStr, status, lv.ID, lv.Name, len(lv.Waves), lv.StartGold, lv.StartLives)
		ebitenutil.DebugPrintAt(screen, line, 8, y)
	}

	helpY := 50 + len(g.Levels)*22 + 8
	ebitenutil.DebugPrintAt(screen,
		" Press 1-9 or 0 to start | Q/Esc to quit", 8, helpY)
	if g.Msg != "" {
		ebitenutil.DebugPrintAt(screen, " "+g.Msg, 8, helpY+18)
	}
}

// ============================================================
// In-game Drawing
// ============================================================

func (eg *EbitenGame) drawGame(screen *ebiten.Image) {
	g := eg.game
	lv := g.currentLevel()
	if lv == nil {
		return
	}

	// title
	title := fmt.Sprintf(" KR V2 — Lv %d: %s — Wave %d/%d ",
		lv.ID, lv.Name, g.WaveIdx+1, len(lv.Waves))
	ebitenutil.DebugPrintAt(screen, title, 8, 6)

	// V3 Phase 2: 背景 grass tile 填充非 path 区域 (game area 内)
	if tilesheet != nil {
		for y := 0; y < mapH; y++ {
			for x := 0; x < mapW; x++ {
				p := Point{X: x, Y: y}
				if g.pathContains(p) {
					continue
				}
				px, py := cellPos(p)
				drawTile(screen, spriteGrass, px, py)
			}
		}
	}

	// V3: path 用 dirt tile sprite (若 tilesheet load 失败 fallback 棕色 rect)
	for _, p := range g.Path {
		x, y := cellPos(p)
		if tilesheet != nil {
			drawTile(screen, spriteDirtPath, x, y)
		} else {
			fillRect(screen, x, y, float32(cellPx), float32(cellPx), eColPathBg)
		}
	}
	// start / end markers (S/E 文字 overlay,仍叠加在 path tile 之上)
	if len(g.Path) > 0 {
		sx, sy := cellPos(g.Path[0])
		strokeRect(screen, sx, sy, float32(cellPx), float32(cellPx), eColStart, 2)
		ebitenutil.DebugPrintAt(screen, "S", int(sx)+cellPx/2-3, int(sy)+cellPx/2-6)

		end := g.Path[len(g.Path)-1]
		ex, ey := cellPos(end)
		strokeRect(screen, ex, ey, float32(cellPx), float32(cellPx), eColEnd, 2)
		ebitenutil.DebugPrintAt(screen, "E", int(ex)+cellPx/2-3, int(ey)+cellPx/2-6)
	}

	// V3: towers 用 sprite (按 kind + level 切换造型, Magic 升级换 rocket 数)
	// V3 Phase 2: 加 level badge digit overlay (mini sprite 右上角)
	for _, t := range g.Towers {
		x, y := cellPos(t.Pos)
		if tilesheet != nil {
			drawTile(screen, towerSpriteID(t.Kind, t.Level), x, y)
			// level badge: mini digit sprite at top-right
			badgeScale := 0.45
			miniSize := float32(cellPx) * float32(badgeScale)
			drawTileAt(screen, digitSpriteID(t.Level),
				x+float32(cellPx)-miniSize, y, badgeScale, 1.0)
		} else {
			spec := towerSpecs[t.Kind]
			lvl := t.Spec()
			cx := x + float32(cellPx)/2
			cy := y + float32(cellPx)/2
			fillCircle(screen, cx, cy, float32(cellPx)/2-2, ebitenColor(spec.Color))
			label := fmt.Sprintf("%c%c", lvl.Char1, lvl.Char2)
			ebitenutil.DebugPrintAt(screen, label, int(cx)-6, int(cy)-6)
		}
	}

	// V3: enemies 用 sprite (Boss 放大 1.4x)
	for _, e := range g.Enemies {
		if e.Dead || e.Escaped {
			continue
		}
		p := e.Pos(g.Path)
		x, y := cellPos(p)
		if tilesheet != nil {
			if e.Kind == EBoss {
				drawTileScaled(screen, enemySpriteID(e.Kind), x, y, 1.4)
			} else {
				drawTile(screen, enemySpriteID(e.Kind), x, y)
			}
		} else {
			spec := enemySpecs[e.Kind]
			cx := x + float32(cellPx)/2
			cy := y + float32(cellPx)/2
			radius := float32(cellPx) / 3
			if e.Kind == EBoss {
				radius = float32(cellPx)/2 - 1
			}
			col := ebitenColor(spec.Color)
			if e.HP <= e.MaxHP/3 {
				col = color.RGBA{R: col.R / 2, G: col.G / 2, B: col.B / 2, A: 255}
			}
			fillCircle(screen, cx, cy, radius, col)
		}
		// HP bar 仍画 (sprite 之上)
		hpRatio := float32(e.HP) / float32(e.MaxHP)
		if hpRatio < 0 {
			hpRatio = 0
		}
		barW := float32(cellPx) - 4
		fillRect(screen, x+2, y, barW, 3, eColHpBg)
		fillRect(screen, x+2, y, barW*hpRatio, 3, eColHpFg)
	}

	// V2.7: 射程圈预览 (cursor 在 tower 上 → 显示当前级 range,
	//                    cursor 在空地 → 显示 selected lvl 1 range preview)
	var atTowerHover *Tower
	for _, t := range g.Towers {
		if t.Pos == g.Cursor {
			atTowerHover = t
			break
		}
	}
	if atTowerHover != nil {
		spec := towerSpecs[atTowerHover.Kind]
		col := ebitenColor(spec.Color)
		col.A = 140 // 半透明描边
		cxr, cyr := cellPos(g.Cursor)
		centerX := cxr + float32(cellPx)/2
		centerY := cyr + float32(cellPx)/2
		vector.StrokeCircle(screen, centerX, centerY,
			float32(atTowerHover.Spec().Range)*float32(cellPx), 1.5, col, true)
	} else if !g.pathContains(g.Cursor) {
		spec := towerSpecs[g.Selected]
		base := ebitenColor(spec.Color)
		col := color.RGBA{R: base.R, G: base.G, B: base.B, A: 80} // 更淡
		cxr, cyr := cellPos(g.Cursor)
		centerX := cxr + float32(cellPx)/2
		centerY := cyr + float32(cellPx)/2
		vector.StrokeCircle(screen, centerX, centerY,
			float32(spec.Levels[0].Range)*float32(cellPx), 1, col, true)
	}

	// V2.6: 攻击视觉特效 (在 cursor / status 之前画, 在 enemy 之上)
	// V3 Phase 2: EHit 用 fire sprite (alpha-faded) 替换 V2.6 黄圆
	for _, fx := range g.Effects {
		alpha := uint8(fx.Alpha() * 255)
		c := color.RGBA{R: fx.Color.R, G: fx.Color.G, B: fx.Color.B, A: alpha}
		fx1, fy1 := cellPos(fx.From)
		fx2, fy2 := cellPos(fx.To)
		cx1 := fx1 + float32(cellPx)/2
		cy1 := fy1 + float32(cellPx)/2
		cx2 := fx2 + float32(cellPx)/2
		cy2 := fy2 + float32(cellPx)/2
		if fx.Kind == EShoot {
			// V3 Phase 3: bullet sprite 沿 from→to 飞行
			// V3 Phase 3b: bullet 按塔型选 sprite (Archer 小弹 / Cannon 大火箭 / Magic 小火箭)
			if tilesheet != nil {
				progress := float32(1.0 - fx.Alpha())
				lerpX := cx1 + (cx2-cx1)*progress
				lerpY := cy1 + (cy2-cy1)*progress
				bulletX := lerpX - float32(cellPx)/2
				bulletY := lerpY - float32(cellPx)/2
				drawTileAt(screen, bulletSpriteID(fx.Tower), bulletX, bulletY, 0.55, 1.0)
			} else {
				vector.StrokeLine(screen, cx1, cy1, cx2, cy2, 2, c, true)
			}
		} else { // EHit: fire sprite,fade out + 略微缩小
			if tilesheet != nil {
				scaleFactor := 0.6 + 0.4*fx.Alpha()
				drawTileAt(screen, spriteHitFire, fx2, fy2, scaleFactor, fx.Alpha())
			} else {
				r := float32(cellPx)/2*float32(fx.Alpha()) + 2
				fillCircle(screen, cx2, cy2, r, c)
			}
		}
	}

	// cursor
	cxr, cyr := cellPos(g.Cursor)
	strokeRect(screen, cxr, cyr, float32(cellPx), float32(cellPx), eColCursor, 2)

	// status row
	statusY := topBarH + gameAreaH + 4
	leftStatus := fmt.Sprintf(" Gold: %-4d Lives: %d/%d Enemies: %-3d  %s",
		g.Gold, g.Lives, g.StartLives, g.CountAliveEnemies(), g.Msg)
	ebitenutil.DebugPrintAt(screen, leftStatus, 4, statusY)
	if g.prepTimer > 0 && g.Phase == PhasePlaying {
		prepMsg := fmt.Sprintf("PREP: %.1fs", g.prepTimer)
		ebitenutil.DebugPrintAt(screen, prepMsg, windowW-90, statusY)
	}

	// 塔选择栏 + 升级提示
	selY := statusY + 20
	selX := 8
	var atTower *Tower
	for _, t := range g.Towers {
		if t.Pos == g.Cursor {
			atTower = t
			break
		}
	}
	for _, k := range TowerKinds() {
		spec := towerSpecs[k]
		label := fmt.Sprintf("[%d] %s %dg", int(k)+1, spec.Name, spec.Levels[0].Cost)
		if k == g.Selected {
			tw := len(label)*7 + 4
			fillRect(screen, float32(selX)-2, float32(selY)-2, float32(tw), 16, eColSelHi)
		}
		ebitenutil.DebugPrintAt(screen, label, selX, selY)
		selX += len(label)*7 + 12
	}
	if atTower != nil {
		cost, can := atTower.NextUpgradeCost()
		var hint string
		if can {
			hint = fmt.Sprintf("Space=Upgrade (%dg)", cost)
		} else {
			hint = "MAX LEVEL"
		}
		ebitenutil.DebugPrintAt(screen, hint, selX, selY)
	}

	// help row
	helpY := selY + 22
	ebitenutil.DebugPrintAt(screen,
		" Arrows: move  1/2/3: select  Space: build/upgrade  M: menu  Q/Esc: quit",
		4, helpY)

	// banner
	bannerX := windowW/2 - 130
	bannerY := topBarH + gameAreaH/2
	if g.Phase == PhaseWon {
		msg := fmt.Sprintf(" *** VICTORY *** Level %d cleared! ", lv.ID)
		ebitenutil.DebugPrintAt(screen, msg, bannerX, bannerY)
		ebitenutil.DebugPrintAt(screen, " Press M for menu, Q to quit ", bannerX, bannerY+18)
	}
	if g.Phase == PhaseLost {
		ebitenutil.DebugPrintAt(screen, " *** GAME OVER *** ", bannerX+50, bannerY)
		ebitenutil.DebugPrintAt(screen, " Press M for menu, Q to quit ", bannerX, bannerY+18)
	}
}
