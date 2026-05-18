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
	cellPx     = 28
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

	// path
	for _, p := range g.Path {
		x, y := cellPos(p)
		fillRect(screen, x, y, float32(cellPx), float32(cellPx), eColPathBg)
	}
	// start / end
	if len(g.Path) > 0 {
		sx, sy := cellPos(g.Path[0])
		fillRect(screen, sx, sy, float32(cellPx), float32(cellPx), eColStart)
		ebitenutil.DebugPrintAt(screen, "S", int(sx)+cellPx/2-3, int(sy)+cellPx/2-6)

		end := g.Path[len(g.Path)-1]
		ex, ey := cellPos(end)
		fillRect(screen, ex, ey, float32(cellPx), float32(cellPx), eColEnd)
		ebitenutil.DebugPrintAt(screen, "E", int(ex)+cellPx/2-3, int(ey)+cellPx/2-6)
	}

	// towers
	for _, t := range g.Towers {
		spec := towerSpecs[t.Kind]
		lvl := t.Spec()
		x, y := cellPos(t.Pos)
		cx := x + float32(cellPx)/2
		cy := y + float32(cellPx)/2
		fillCircle(screen, cx, cy, float32(cellPx)/2-2, ebitenColor(spec.Color))
		label := fmt.Sprintf("%c%c", lvl.Char1, lvl.Char2)
		ebitenutil.DebugPrintAt(screen, label, int(cx)-6, int(cy)-6)
	}

	// enemies
	for _, e := range g.Enemies {
		if e.Dead || e.Escaped {
			continue
		}
		spec := enemySpecs[e.Kind]
		p := e.Pos(g.Path)
		x, y := cellPos(p)
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
		// HP bar
		hpRatio := float32(e.HP) / float32(e.MaxHP)
		if hpRatio < 0 {
			hpRatio = 0
		}
		barW := float32(cellPx) - 4
		fillRect(screen, x+2, y, barW, 3, eColHpBg)
		fillRect(screen, x+2, y, barW*hpRatio, 3, eColHpFg)
	}

	// V2.6: 攻击视觉特效 (在 cursor / status 之前画, 在 enemy 之上)
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
			vector.StrokeLine(screen, cx1, cy1, cx2, cy2, 2, c, true)
		} else { // EHit: fade 圆, 半径随 fade 减小
			r := float32(cellPx)/2*float32(fx.Alpha()) + 2
			fillCircle(screen, cx2, cy2, r, c)
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
