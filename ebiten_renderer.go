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
	// V3 Phase 5c: 加载 truetype font (失败 fallback ebitenutil 7×13)
	if err := loadGameFont(); err != nil {
		fmt.Println("warning: font load failed:", err, "(falling back to bitmap font)")
	}
	// V4 Phase 1: SFX 解码 (失败不致命, 缺失音效静默跳过)
	if n := initAudio(); n > 0 {
		fmt.Println("warning:", n, "sfx failed to decode (those sounds muted)")
	}
	// V4 Phase 2: BGM 常驻 player (依赖 initAudio 的 audioCtx)
	if n := initBGM(); n > 0 {
		fmt.Println("warning:", n, "bgm track(s) failed to load (music muted)")
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
	// V4 Phase 1: drain SFX 队列并播放 (handleInput 的 build/upgrade
	// 与 Update 的 shoot/death/wave/win/lose 事件都在本帧消费)
	// V4 Phase 2: 主音量 = 存档音量档/10, BGM 状态机每帧驱动
	masterVol := float64(eg.game.Save.VolumeLevel()) / float64(maxVolume)
	playSounds(eg.game.DrainSounds(), masterVol)
	updateBGM(eg.game.Phase, masterVol, dt)
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
	// V4 Phase 2: 音量档 -/= (全 phase 生效, 0 档 = 静音)
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		g.AdjustVolume(-1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		g.AdjustVolume(+1)
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
		// V3 Phase 5b: 左键先 check tower select button, 否则 game area TryAction
		mx, my := ebiten.CursorPosition()
		if p, ok := pixelToCell(mx, my); ok {
			g.Cursor = p
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if k, ok := towerButtonAt(mx, my); ok {
				g.Selected = k
				g.Msg = "Selected " + towerSpecs[k].Name
			} else if _, ok := pixelToCell(mx, my); ok {
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
// V3 Phase 5a: layout 重排 (startY 50→60, rowH 22→36 配合 box-style 渲染)
func menuRowAtPixel(my, numLevels int) (int, bool) {
	const startY = 60
	const rowH = 36
	i := (my - startY) / rowH
	if i < 0 || i >= numLevels {
		return 0, false
	}
	return i, true
}

// V3 Phase 5b: tower select button hitbox (与 drawGame 的 button 渲染保持一致)
// 返回 (kind, true) 表示鼠标在某 button 内, (0, false) 表示不在任何 button 内
func towerButtonAt(mx, my int) (TowerKind, bool) {
	statusY := topBarH + gameAreaH + 8
	selY := statusY + 22
	const (
		btnW   = 100
		btnH   = 32
		btnGap = 8
	)
	btnX := 8
	for _, k := range TowerKinds() {
		if mx >= btnX && mx < btnX+btnW &&
			my >= selY && my < selY+btnH {
			return k, true
		}
		btnX += btnW + btnGap
	}
	return 0, false
}

// ============================================================
// Level Select
// ============================================================

func (eg *EbitenGame) drawLevelSelect(screen *ebiten.Image) {
	g := eg.game

	// V3 Phase 5a: title panel 顶部 (50 px)
	fillRect(screen, 0, 0, float32(windowW), 50,
		color.RGBA{R: 15, G: 15, B: 25, A: 230})
	strokeRect(screen, 0, 0, float32(windowW), 50,
		color.RGBA{R: 60, G: 60, B: 80, A: 255}, 1)
	drawText(screen,
		" Kingdom Rush V3  —  Select Level", 20, 10)
	// V4 Phase 2: 音量档显示 (右上角, -/= 调节)
	drawText(screen,
		fmt.Sprintf("Vol %d/%d (-/=)", g.Save.VolumeLevel(), maxVolume),
		windowW-130, 10)

	completed := 0
	for _, lv := range g.Levels {
		if g.Save.IsCompleted(lv.ID) {
			completed++
		}
	}
	drawText(screen,
		fmt.Sprintf(" Progress: %d / %d cleared", completed, len(g.Levels)),
		20, 28)

	// V3 Phase 5a: 每关 row 用 box (panel + border)
	// 状态颜色: 通关绿底, 锁定灰底, 默认深底; hover unlocked 加金边
	_, mouseY := ebiten.CursorPosition()
	mouseX, _ := ebiten.CursorPosition()
	const (
		rowStartY = 60
		rowH      = 36
		rowMargin = 2
		rowX      = 8
	)
	rowW := windowW - rowX*2

	for i, lv := range g.Levels {
		y := rowStartY + i*rowH
		key := i + 1
		keyStr := fmt.Sprintf("%d", key)
		if key == 10 {
			keyStr = "0"
		}

		isCompleted := g.Save.IsCompleted(lv.ID)
		isUnlocked := g.Save.IsUnlocked(lv.ID)
		isHover := mouseY >= y && mouseY < y+rowH-rowMargin &&
			mouseX >= rowX && mouseX < rowX+rowW

		// Box bg
		var bgCol color.RGBA
		switch {
		case isCompleted:
			bgCol = color.RGBA{R: 30, G: 60, B: 30, A: 220} // 深绿
		case !isUnlocked:
			bgCol = color.RGBA{R: 35, G: 35, B: 40, A: 220} // 深灰 (locked)
		default:
			bgCol = color.RGBA{R: 40, G: 40, B: 60, A: 220} // 默认 unlocked
		}
		if isHover && isUnlocked {
			bgCol = color.RGBA{R: 60, G: 80, B: 110, A: 240}
		}
		fillRect(screen, float32(rowX), float32(y),
			float32(rowW), float32(rowH-rowMargin), bgCol)

		// Box border
		var borderCol color.RGBA
		switch {
		case isCompleted:
			borderCol = color.RGBA{R: 80, G: 200, B: 80, A: 255} // 亮绿
		case !isUnlocked:
			borderCol = color.RGBA{R: 80, G: 80, B: 80, A: 255} // 灰
		default:
			borderCol = color.RGBA{R: 150, G: 150, B: 180, A: 255}
		}
		if isHover && isUnlocked {
			borderCol = color.RGBA{R: 255, G: 220, B: 80, A: 255} // 金
		}
		strokeRect(screen, float32(rowX), float32(y),
			float32(rowW), float32(rowH-rowMargin), borderCol, 2)

		// Status indicator (left side)
		status := "[    ]"
		if isCompleted {
			status = "[DONE]"
		} else if !isUnlocked {
			status = "[LOCK]"
		}

		line := fmt.Sprintf("[%s]  %s  Lv %2d  %-18s    (waves:%d  gold:%d  lives:%d)",
			keyStr, status, lv.ID, lv.Name, len(lv.Waves), lv.StartGold, lv.StartLives)
		drawText(screen, line, rowX+12, y+12)
	}

	// Help row + msg
	helpY := rowStartY + len(g.Levels)*rowH + 8
	drawText(screen,
		" Click row or press 1-9 / 0 to start | Q/Esc to quit",
		20, helpY)
	if g.Msg != "" {
		drawText(screen, " "+g.Msg, 20, helpY+18)
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

	// V3 Phase 4: top title bar panel background
	fillRect(screen, 0, 0, float32(windowW), float32(topBarH),
		color.RGBA{R: 15, G: 15, B: 25, A: 230})
	strokeRect(screen, 0, 0, float32(windowW), float32(topBarH),
		color.RGBA{R: 60, G: 60, B: 80, A: 255}, 1)

	// title
	title := fmt.Sprintf(" KR V3 — Lv %d: %s — Wave %d/%d ",
		lv.ID, lv.Name, g.WaveIdx+1, len(lv.Waves))
	drawText(screen, title, 8, 8)
	// V4 Phase 2: 音量档显示 (右上角, -/= 调节)
	drawText(screen,
		fmt.Sprintf("Vol %d/%d", g.Save.VolumeLevel(), maxVolume),
		windowW-80, 8)

	// V3 Phase 6: ALL cells grass tile (背景统一, path overlay 之上)
	if tilesheet != nil {
		for y := 0; y < mapH; y++ {
			for x := 0; x < mapW; x++ {
				px, py := cellPos(Point{X: x, Y: y})
				drawTile(screen, spriteGrass, px, py)
			}
		}
	}

	// V3 Phase 6: procedural path overlay — 中心圆 + 邻居方向矩形
	// 自动适配任何 corner/T/直/端点 形状, 无需 sprite catalog
	// Phase 6b: 两 pass outline — 先全 path 画外扩 outlineW 的 edge 色,
	// 再叠正常尺寸 fill 色, 留出的环带即描边。pass 间不能交错:
	// 单 loop 内邻 cell 的 edge 会盖掉本 cell 已画的 fill, seam 处留深色横纹。
	pathCol := color.RGBA{R: 139, G: 69, B: 19, A: 255}     // SaddleBrown
	pathEdgeCol := color.RGBA{R: 100, G: 50, B: 14, A: 255} // 深棕 outline
	const roadW = 22
	const roadHalf = roadW / 2.0
	const outlineW = 3
	// drawPathShape: 单 cell 的程序化路面形状, grow = 外扩像素 (edge pass 用)
	drawPathShape := func(p Point, grow float32, col color.RGBA) {
		x, y := cellPos(p)
		cx := x + float32(cellPx)/2
		cy := y + float32(cellPx)/2
		half := float32(roadHalf) + grow

		// 中心圆 (略大于 road width, 与方向矩形混合形成圆角)
		fillCircle(screen, cx, cy, half+1, col)

		// 方向连接矩形 (从 cell 边到中心; 邻 cell 画自己那半, seam 处衔接)
		if g.pathContains(Point{X: p.X - 1, Y: p.Y}) {
			fillRect(screen, x, cy-half, float32(cellPx)/2, half*2, col)
		}
		if g.pathContains(Point{X: p.X + 1, Y: p.Y}) {
			fillRect(screen, cx, cy-half, float32(cellPx)/2, half*2, col)
		}
		if g.pathContains(Point{X: p.X, Y: p.Y - 1}) {
			fillRect(screen, cx-half, y, half*2, float32(cellPx)/2, col)
		}
		if g.pathContains(Point{X: p.X, Y: p.Y + 1}) {
			fillRect(screen, cx-half, cy, half*2, float32(cellPx)/2, col)
		}
	}
	for _, p := range g.Path {
		drawPathShape(p, outlineW, pathEdgeCol) // pass 1: edge
	}
	for _, p := range g.Path {
		drawPathShape(p, 0, pathCol) // pass 2: fill
	}

	// 起点 / 终点 markers — 在 path overlay 之上, 圆色块 + 字符标签
	if len(g.Path) > 0 {
		sx, sy := cellPos(g.Path[0])
		fillCircle(screen, sx+float32(cellPx)/2, sy+float32(cellPx)/2,
			float32(cellPx)/3, eColStart)
		drawText(screen, "S", int(sx)+cellPx/2-3, int(sy)+cellPx/2-6)

		end := g.Path[len(g.Path)-1]
		ex, ey := cellPos(end)
		fillCircle(screen, ex+float32(cellPx)/2, ey+float32(cellPx)/2,
			float32(cellPx)/3, eColEnd)
		drawText(screen, "E", int(ex)+cellPx/2-3, int(ey)+cellPx/2-6)
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
			drawText(screen, label, int(cx)-6, int(cy)-6)
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

	// V3 Phase 4: bottom UI panel background + 边框
	panelY := float32(topBarH + gameAreaH)
	fillRect(screen, 0, panelY, float32(windowW), float32(bottomBarH),
		color.RGBA{R: 15, G: 15, B: 25, A: 230})
	strokeRect(screen, 0, panelY, float32(windowW), float32(bottomBarH),
		color.RGBA{R: 60, G: 60, B: 80, A: 255}, 1)

	// V3 Phase 4: status row 用 sprite icon for Gold / Lives
	statusY := topBarH + gameAreaH + 8
	curX := 8
	// $ icon + gold value
	if tilesheet != nil {
		drawTileAt(screen, spriteGold, float32(curX), float32(statusY-2), 0.55, 1.0)
	}
	drawText(screen, fmt.Sprintf("%d", g.Gold), curX+24, statusY)
	curX += 80
	// + icon + lives
	if tilesheet != nil {
		drawTileAt(screen, spriteLives, float32(curX), float32(statusY-2), 0.55, 1.0)
	}
	drawText(screen,
		fmt.Sprintf("%d/%d", g.Lives, g.StartLives), curX+24, statusY)
	curX += 90
	// enemies count
	drawText(screen,
		fmt.Sprintf("Enemies:%-3d", g.CountAliveEnemies()), curX, statusY)
	curX += 110
	// msg
	drawText(screen, g.Msg, curX, statusY)
	// prep timer (right edge)
	if g.prepTimer > 0 && g.Phase == PhasePlaying {
		prepMsg := fmt.Sprintf("PREP: %.1fs", g.prepTimer)
		drawText(screen, prepMsg, windowW-90, statusY)
	}

	// V3 Phase 4: 塔选择按钮 (button style + mini tower sprite + 选中 gold border)
	selY := statusY + 22
	const (
		btnW   = 100
		btnH   = 32
		btnGap = 8
	)
	btnX := 8

	var atTower *Tower
	for _, t := range g.Towers {
		if t.Pos == g.Cursor {
			atTower = t
			break
		}
	}

	for _, k := range TowerKinds() {
		spec := towerSpecs[k]
		bgCol := color.RGBA{R: 40, G: 40, B: 60, A: 220}
		borderCol := color.RGBA{R: 120, G: 120, B: 150, A: 255}
		if k == g.Selected {
			bgCol = color.RGBA{R: 80, G: 80, B: 140, A: 240}
			borderCol = color.RGBA{R: 255, G: 220, B: 80, A: 255}
		}
		fillRect(screen, float32(btnX), float32(selY),
			float32(btnW), float32(btnH), bgCol)
		strokeRect(screen, float32(btnX), float32(selY),
			float32(btnW), float32(btnH), borderCol, 2)
		// mini tower sprite (32×32 at left)
		if tilesheet != nil {
			drawTileAt(screen, towerSpriteID(k, 1),
				float32(btnX+1), float32(selY+1), 1.0, 1.0)
		}
		// label 2 行: [N]Name 在上, Cost 在下
		drawText(screen,
			fmt.Sprintf("[%d]%s", int(k)+1, spec.Name),
			btnX+36, selY+4)
		drawText(screen,
			fmt.Sprintf("%dg", spec.Levels[0].Cost),
			btnX+36, selY+18)
		btnX += btnW + btnGap
	}

	// 升级提示 (右侧 of buttons)
	if atTower != nil {
		cost, can := atTower.NextUpgradeCost()
		var hint string
		if can {
			hint = fmt.Sprintf("Space=Upgrade %dg", cost)
		} else {
			hint = "MAX LEVEL"
		}
		drawText(screen, hint, btnX+8, selY+12)
	}

	// help row
	helpY := selY + btnH + 8
	drawText(screen,
		" Arrows/Mouse: move | 1/2/3: select | Space/Click: build/upgrade | M: menu | Q/Esc: quit",
		4, helpY)

	// banner
	bannerX := windowW/2 - 130
	bannerY := topBarH + gameAreaH/2
	if g.Phase == PhaseWon {
		msg := fmt.Sprintf(" *** VICTORY *** Level %d cleared! ", lv.ID)
		drawText(screen, msg, bannerX, bannerY)
		drawText(screen, " Press M for menu, Q to quit ", bannerX, bannerY+18)
	}
	if g.Phase == PhaseLost {
		drawText(screen, " *** GAME OVER *** ", bannerX+50, bannerY)
		drawText(screen, " Press M for menu, Q to quit ", bannerX, bannerY+18)
	}
}
