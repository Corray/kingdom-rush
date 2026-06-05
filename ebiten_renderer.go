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
	"strings"
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
	// V4 Phase 4: HUD 金币闪烁 (金币增加时高亮 0.5s)
	lastGold  int
	lastPhase GamePhase
	goldFlash float64
	// V4 Phase 5: screen shake (lives 丢失) + boss 击杀顿帧
	lastLives     int
	lastBossKills int
	shake         float64       // shake 剩余时间
	hitStop       float64       // 顿帧剩余时间 (世界冻结, 渲染继续)
	offscreen     *ebiten.Image // shake 时整帧偏移用缓冲
	// V5 Phase 5: 陨石雨瞄准态 (UI 模式, 不入 game 逻辑)
	aiming bool
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
	// V4 Phase 5: 顿帧 — 世界冻结 (跳过 game.Update), 输入/渲染继续
	if eg.hitStop > 0 {
		eg.hitStop -= dt
		return nil
	}
	eg.game.Update(dt)
	// V4 Phase 1: drain SFX 队列并播放 (handleInput 的 build/upgrade
	// 与 Update 的 shoot/death/wave/win/lose 事件都在本帧消费)
	// V4 Phase 2: 主音量 = 存档音量档/10, BGM 状态机每帧驱动
	masterVol := float64(eg.game.Save.VolumeLevel()) / float64(maxVolume)
	playSounds(eg.game.DrainSounds(), masterVol)
	updateBGM(eg.game.Phase, masterVol, dt)
	// V4 Phase 4: 金币增加 → HUD 闪烁 (phase 切换时只同步不闪,
	// 防止 StartLevel 重置 Gold 触发假闪)
	// V4 Phase 5: lives 减少 → shake / boss 击杀 → 顿帧 (同模式观察,
	// JuiceOff 时不触发)
	if eg.game.Phase != eg.lastPhase {
		eg.lastPhase = eg.game.Phase
		eg.lastGold = eg.game.Gold
		eg.lastLives = eg.game.Lives
		eg.lastBossKills = eg.game.BossKills
		eg.goldFlash = 0
		eg.shake = 0
		eg.hitStop = 0
		eg.aiming = false // V5 Phase 5: phase 切换取消瞄准
	} else {
		if eg.game.Gold > eg.lastGold {
			eg.goldFlash = 0.5
		}
		if !eg.game.Save.JuiceOff {
			if eg.game.Lives < eg.lastLives {
				eg.shake = shakeDuration
			}
			if eg.game.BossKills > eg.lastBossKills {
				eg.hitStop = hitStopS
			}
		}
	}
	eg.lastGold = eg.game.Gold
	eg.lastLives = eg.game.Lives
	eg.lastBossKills = eg.game.BossKills
	if eg.goldFlash > 0 {
		eg.goldFlash -= dt
	}
	if eg.shake > 0 {
		eg.shake -= dt
	}
	return nil
}

func (eg *EbitenGame) Draw(screen *ebiten.Image) {
	// V4 Phase 5: shake 激活时画到离屏缓冲, 再整帧偏移 blit
	target := screen
	if eg.shake > 0 {
		if eg.offscreen == nil {
			eg.offscreen = ebiten.NewImage(windowW, windowH)
		}
		target = eg.offscreen
	}
	target.Fill(eColBg)
	if eg.game.Phase == PhaseLevelSelect {
		eg.drawLevelSelect(target)
	} else {
		eg.drawGame(target)
	}
	if target != screen {
		dx, dy := shakeOffset(eg.shake)
		screen.Fill(eColBg) // 偏移后露出的边缘用背景色补
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(dx, dy)
		screen.DrawImage(eg.offscreen, op)
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
	// V4 Phase 5: J 键开关屏幕反馈特效 (shake/顿帧)
	if inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		g.ToggleJuice()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.Phase == PhasePlaying {
		g.TryAction()
	}
	// V5 Phase 1: X 键 / 右键卖塔 (右键位置已由 hover 同步到 cursor)
	if g.Phase == PhasePlaying &&
		(inpututil.IsKeyJustPressed(ebiten.KeyX) ||
			inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)) {
		g.SellTower()
	}
	// V5 Phase 2: T 键循环切换光标处塔的 targeting 策略
	if g.Phase == PhasePlaying && inpututil.IsKeyJustPressed(ebiten.KeyT) {
		g.CycleTargeting()
	}
	// V5 Phase 5: R 键进入/取消陨石瞄准 (不用 Esc — 已绑退出)
	if g.Phase == PhasePlaying && inpututil.IsKeyJustPressed(ebiten.KeyR) {
		switch {
		case eg.aiming:
			eg.aiming = false
			g.Msg = "Meteor aim cancelled"
		case g.MeteorCD <= 0:
			eg.aiming = true
			g.Msg = "Meteor: click to strike (R to cancel)"
		default:
			g.Msg = fmt.Sprintf("Meteor on cooldown %.0fs", g.MeteorCD)
		}
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
		// V5 Phase 4: 4 键选 Frost (1/2/3/4 与按钮顺序一致)
		if inpututil.IsKeyJustPressed(ebiten.KeyDigit4) {
			g.Selected = TFrost
			g.Msg = "Selected Frost"
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
			// V5 Phase 5: 瞄准态点击 = 释放陨石 (优先于建塔/按钮)
			if eg.aiming {
				if cell, ok := pixelToCell(mx, my); ok && g.CastMeteor(cell) {
					eg.aiming = false
					if !g.Save.JuiceOff {
						eg.shake = shakeDuration // 落地震屏 (复用 V4 juice)
					}
				}
			} else if k, ok := towerButtonAt(mx, my); ok {
				g.Selected = k
				g.Msg = "Selected " + towerSpecs[k].Name
			} else if _, ok := pixelToCell(mx, my); ok {
				g.TryAction()
			}
		}
		return
	}
	if g.Phase == PhaseLevelSelect {
		// V6 Phase 2: D 键循环切换难度 (Normal → Hard → Easy)
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.CycleDifficulty()
		}
		// V6 Phase 3: E 键进入 endless (seed = 时间戳; 测试走注入)
		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			g.StartEndless(time.Now().UnixNano())
		}
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
		// V2.7: menu 阶段鼠标点击 row 启动关卡 (V6: 两列 hitbox)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			if i, ok := menuRowAtPixel(mx, my, len(g.Levels)); ok {
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

// cellPosF: V4 Phase 3 — 浮点 cell 坐标 → 像素 (平滑插值移动用)。
func cellPosF(fx, fy float64) (float32, float32) {
	return float32(fx * float64(cellPx)), float32(topBarH) + float32(fy*float64(cellPx))
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

// menuRowsPerCol: V6 Phase 1 — 两列布局, 每列关卡数。
const menuRowsPerCol = 10

// menuRowAtPixel 返回鼠标所在的 level select index (0..numLevels-1)。
// V3 Phase 5a: layout 重排 (startY 50→60, rowH 22→36 配合 box-style 渲染)
// V6 Phase 1: 两列布局 — 左列 0-9, 右列 10-19 (按 windowW/2 分界)
func menuRowAtPixel(mx, my, numLevels int) (int, bool) {
	const startY = 60
	const rowH = 36
	// V6 修复既有 bug: Go 负数整除向零截断, my 略小于 startY 时
	// (my-startY)/rowH == 0 误命中首行 (点 title 区误启动关卡)
	if my < startY {
		return 0, false
	}
	row := (my - startY) / rowH
	if row >= menuRowsPerCol {
		return 0, false
	}
	col := 0
	if mx >= windowW/2 {
		col = 1
	}
	i := col*menuRowsPerCol + row
	if i >= numLevels {
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
	// V7.2 M1: 大标题 (20pt 青) + 副标题/统计灰阶层次
	drawTextBigCol(screen, "Gopher Defense", 20, 6,
		color.RGBA{R: 120, G: 220, B: 255, A: 255})
	drawTextCol(screen, "Select Level", 220, 14,
		color.RGBA{R: 140, G: 140, B: 160, A: 255})
	// V4 Phase 2: 音量档显示 (右上角, -/= 调节)
	drawTextCol(screen,
		fmt.Sprintf("Vol %d/%d (-/=)", g.Save.VolumeLevel(), maxVolume),
		windowW-130, 10, color.RGBA{R: 170, G: 170, B: 190, A: 255})
	// V6 Phase 2: 难度显示 (右上角第二行, D 切换)
	drawTextCol(screen,
		fmt.Sprintf("Diff: %s (D)", g.Save.Difficulty.Spec().Name),
		windowW-130, 28, color.RGBA{R: 170, G: 170, B: 190, A: 255})

	completed, totalStars := 0, 0
	for _, lv := range g.Levels {
		if g.Save.IsCompleted(lv.ID) {
			completed++
		}
		totalStars += g.Save.StarsFor(lv.ID)
	}
	drawTextCol(screen,
		fmt.Sprintf("Progress %d/%d", completed, len(g.Levels)),
		20, 32, color.RGBA{R: 140, G: 140, B: 160, A: 255})
	drawTextCol(screen,
		fmt.Sprintf("Stars %d/%d", totalStars, len(g.Levels)*3),
		160, 32, color.RGBA{R: 255, G: 200, B: 40, A: 255})

	// V3 Phase 5a: 每关 row 用 box (panel + border)
	// 状态颜色: 通关绿底, 锁定灰底, 默认深底; hover unlocked 加金边
	// V6 Phase 1: 两列布局 (10 关/列), stats 压缩适配列宽
	mouseX, mouseY := ebiten.CursorPosition()
	const (
		rowStartY = 60
		rowH      = 36
		rowMargin = 2
		rowX      = 8
		colGap    = 8
	)
	colW := (windowW - rowX*2 - colGap) / 2

	for i, lv := range g.Levels {
		colIdx := i / menuRowsPerCol
		x := rowX + colIdx*(colW+colGap)
		y := rowStartY + (i%menuRowsPerCol)*rowH
		keyStr := "-" // 11-20 无键位, 鼠标点击
		if i < 9 {
			keyStr = fmt.Sprintf("%d", i+1)
		} else if i == 9 {
			keyStr = "0"
		}

		isCompleted := g.Save.IsCompleted(lv.ID)
		isUnlocked := g.Save.IsUnlocked(lv.ID)
		isHover := mouseY >= y && mouseY < y+rowH-rowMargin &&
			mouseX >= x && mouseX < x+colW

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
		fillRect(screen, float32(x), float32(y),
			float32(colW), float32(rowH-rowMargin), bgCol)

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
		strokeRect(screen, float32(x), float32(y),
			float32(colW), float32(rowH-rowMargin), borderCol, 2)

		// V6 Phase 4 星级 + V7.2 M1 行内分段上色 (等宽 7px/char 偏移)
		stars := strings.Repeat("*", g.Save.StarsFor(lv.ID))
		tx := x + 10
		seg := func(s string, col color.RGBA) {
			drawTextCol(screen, s, tx, y+12, col)
			tx += len(s) * 7
		}
		cKey := color.RGBA{R: 120, G: 220, B: 255, A: 255}  // 青
		cName := color.RGBA{R: 235, G: 235, B: 240, A: 255} // 白
		cDim := color.RGBA{R: 120, G: 120, B: 135, A: 255}  // 灰
		cGold := color.RGBA{R: 255, G: 200, B: 40, A: 255}  // 金
		cDone := color.RGBA{R: 90, G: 220, B: 110, A: 255}  // 绿
		if !isUnlocked {
			cKey, cName, cGold = cDim, cDim, cDim
		}
		seg(fmt.Sprintf("[%s] ", keyStr), cKey)
		switch {
		case isCompleted:
			seg("[DONE] ", cDone)
		case !isUnlocked:
			seg("[LOCK] ", cDim)
		default:
			seg("[    ] ", cDim)
		}
		seg(fmt.Sprintf("Lv%2d %-14s ", lv.ID, lv.Name), cName)
		seg(fmt.Sprintf("%-3s ", stars), cGold)
		seg(fmt.Sprintf("(w:%d g:%d l:%d)", len(lv.Waves), lv.StartGold, lv.StartLives), cDim)
	}

	// Help row + msg (V6: 按每列最大行数定位)
	helpY := rowStartY + menuRowsPerCol*rowH + 8
	// V6 Phase 3: endless 入口 + 纪录
	endlessLine := " [E] Endless mode — no best record yet"
	if g.Save.BestWave > 0 {
		endlessLine = fmt.Sprintf(" [E] Endless mode — best: wave %d", g.Save.BestWave)
	}
	drawTextCol(screen, endlessLine, 20, helpY,
		color.RGBA{R: 190, G: 160, B: 255, A: 255}) // 紫: 模式入口醒目
	drawTextCol(screen,
		" Click row to start (1-9/0 for Lv 1-10) | Q/Esc to quit",
		20, helpY+18, color.RGBA{R: 120, G: 120, B: 135, A: 255}) // 帮助行降噪
	if g.Msg != "" {
		drawTextCol(screen, " "+g.Msg, 20, helpY+36,
			color.RGBA{R: 255, G: 220, B: 80, A: 255}) // Msg 黄: 反馈醒目
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

	// title (V6 Phase 2: 非 Normal 难度标注; Phase 3: endless 专用格式)
	var title string
	if g.Endless {
		title = fmt.Sprintf(" Gopher Defense — Endless — Wave %d (best: %d) ",
			g.WaveIdx+1, g.Save.BestWave)
	} else {
		title = fmt.Sprintf(" Gopher Defense — Lv %d: %s — Wave %d/%d ",
			lv.ID, lv.Name, g.WaveIdx+1, len(lv.Waves))
	}
	if g.Save.Difficulty != DiffNormal {
		title += fmt.Sprintf("[%s] ", g.Save.Difficulty.Spec().Name)
	}
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

	// V7.2 M3: 装饰散布 (grass 之上, path/塔之下, 纯视觉)
	if tilesheet != nil {
		for _, ds := range g.Decor {
			dx, dy := cellPos(ds.Pos)
			drawTile(screen, decorSpriteID(ds.Kind), dx, dy)
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
	// V4 Phase 3: 程序化走动动画 — 平滑插值移动 (替换逐格跳动) +
	// 朝向旋转 (sprite 朝右为基准) + 垂直于行进方向的摆动
	for _, e := range g.Enemies {
		if e.Dead || e.Escaped {
			continue
		}
		lx, ly := pathLerp(g.Path, e.PathIdx)
		ddx, ddy := pathDir(g.Path, e.PathIdx)
		bob := float32(bobOffset(e.PathIdx))
		x, y := cellPosF(lx, ly)
		x += float32(-ddy) * bob // 摆动轴 = 行进方向的垂直向量
		y += float32(ddx) * bob
		if tilesheet != nil {
			scale := 1.0
			if e.Kind == EBoss {
				scale = 1.4
			}
			drawTileRot(screen, enemySpriteID(e.Kind), x, y, scale, 1.0,
				dirAngle(ddx, ddy))
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
		// V5 Phase 4: 减速状态 → 半透明冰蓝覆盖圈
		if e.SlowTimer > 0 {
			fillCircle(screen, x+float32(cellPx)/2, y+float32(cellPx)/2,
				float32(cellPx)/2-4, color.RGBA{R: 150, G: 200, B: 255, A: 90})
		}
		// HP bar 仍画 (sprite 之上, 跟随插值位置, 不旋转)
		hpRatio := float32(e.HP) / float32(e.MaxHP)
		if hpRatio < 0 {
			hpRatio = 0
		}
		barW := float32(cellPx) - 4
		fillRect(screen, x+2, y, barW, 3, eColHpBg)
		fillRect(screen, x+2, y, barW*hpRatio, 3, eColHpFg)
	}

	// V5 Phase 5: 陨石瞄准圈 (橙红, 优先于射程圈)
	if eg.aiming {
		cxr, cyr := cellPos(g.Cursor)
		vector.StrokeCircle(screen, cxr+float32(cellPx)/2, cyr+float32(cellPx)/2,
			float32(meteorRadius)*float32(cellPx), 2.5,
			color.RGBA{R: 255, G: 120, B: 40, A: 200}, true)
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
		} else if fx.Kind == EDeath {
			// V4 Phase 3: 死亡动画 — enemy sprite 放大 (1→1.6) + fade out
			dx, dy := cellPosF(fx.FX, fx.FY)
			if tilesheet != nil {
				grow := 1.0 + (1.0-fx.Alpha())*0.6
				drawTileAt(screen, enemySpriteID(fx.Enemy), dx, dy, grow, fx.Alpha())
			} else {
				r := float32(cellPx) / 3 * (1 + float32(1-fx.Alpha()))
				col := ebitenColor(enemySpecs[fx.Enemy].Color)
				col.A = uint8(255 * fx.Alpha())
				fillCircle(screen, dx+float32(cellPx)/2, dy+float32(cellPx)/2, r, col)
			}
		} else if fx.Kind == EText {
			// V4 Phase 4: 飘字 — 上飘 14px + fade out, 水平居中于 cell
			tx, ty := cellPosF(fx.FX, fx.FY)
			rise := (1 - fx.Alpha()) * 14
			tcol := color.RGBA{R: fx.Color.R, G: fx.Color.G, B: fx.Color.B,
				A: uint8(255 * fx.Alpha())}
			// gomono 12pt 字宽 ~7px, 按字数水平居中
			textX := int(tx) + cellPx/2 - len(fx.Text)*7/2
			drawTextCol(screen, fx.Text, textX, int(ty)-int(rise)-4, tcol)
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
	// $ icon + gold value (V4 Phase 4: 金币增加后 0.5s 内金色高亮)
	if tilesheet != nil {
		drawTileAt(screen, spriteGold, float32(curX), float32(statusY-2), 0.55, 1.0)
	}
	if eg.goldFlash > 0 {
		drawTextCol(screen, fmt.Sprintf("%d", g.Gold), curX+24, statusY,
			color.RGBA{R: 255, G: 200, B: 40, A: 255})
	} else {
		drawText(screen, fmt.Sprintf("%d", g.Gold), curX+24, statusY)
	}
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
	// msg (V7.2 M2: 黄色醒目)
	drawTextCol(screen, g.Msg, curX, statusY,
		color.RGBA{R: 255, G: 220, B: 80, A: 255})
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
		// V7.2 M2: 选中名金色; cost 买得起绿 / 买不起红
		nameCol := color.RGBA{R: 235, G: 235, B: 240, A: 255}
		if k == g.Selected {
			nameCol = color.RGBA{R: 255, G: 220, B: 80, A: 255}
		}
		costCol := color.RGBA{R: 90, G: 220, B: 110, A: 255}
		if g.Gold < spec.Levels[0].Cost {
			costCol = color.RGBA{R: 235, G: 100, B: 90, A: 255}
		}
		drawTextCol(screen,
			fmt.Sprintf("[%d]%s", int(k)+1, spec.Name),
			btnX+36, selY+4, nameCol)
		drawTextCol(screen,
			fmt.Sprintf("%dg", spec.Levels[0].Cost),
			btnX+36, selY+18, costCol)
		btnX += btnW + btnGap
	}

	// 升级提示 (右侧 of buttons) — V5 Phase 1: 加卖塔提示
	if atTower != nil {
		cost, can := atTower.NextUpgradeCost()
		var hint string
		if can {
			hint = fmt.Sprintf("Space=Upgrade %dg", cost)
		} else {
			hint = "MAX LEVEL"
		}
		hint += fmt.Sprintf(" | X=Sell +%dg | T=%s",
			sellRefund(atTower.Kind, atTower.Level), atTower.Target.Name())
		drawTextCol(screen, hint, btnX+8, selY+12,
			color.RGBA{R: 150, G: 210, B: 235, A: 255}) // V7.2 M2: 操作提示青}
	}

	// V5 Phase 5: 陨石冷却条 (按钮行最右; AC "冷却条显示")
	{
		labelX := windowW - 175
		label := "R:Meteor"
		if eg.aiming {
			label = "AIMING.."
		}
		drawText(screen, label, labelX, selY+8)
		barX := float32(labelX + 68)
		const barW, barH = 90, 8
		fillRect(screen, barX, float32(selY+12), barW, barH,
			color.RGBA{R: 40, G: 40, B: 60, A: 255})
		ready := 1 - g.MeteorCD/meteorCooldownS
		fillCol := color.RGBA{R: 255, G: 140, B: 40, A: 255} // 冷却中橙
		if g.MeteorCD <= 0 {
			fillCol = color.RGBA{R: 80, G: 220, B: 80, A: 255} // 就绪绿
		}
		fillRect(screen, barX, float32(selY+12), float32(barW*ready), barH, fillCol)
	}

	// help row
	helpY := selY + btnH + 8
	drawTextCol(screen,
		" Arrows/Mouse: move | 1-4: select | Space/Click: build/upgrade | M: menu | -/=: vol | J: fx | Q/Esc: quit",
		4, helpY, color.RGBA{R: 120, G: 120, B: 135, A: 255}) // V7.2 M2: 帮助行降噪

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
