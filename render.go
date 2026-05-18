// Renderer interface + terminal (tcell) 实现。
// V1.5 仅 terminal impl,接口预留 web/WASM 适配。
package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

const (
	cellW      = 2
	gameTopRow = 1
)

// Renderer 是渲染层抽象。V2 web/WASM 实现新类型时实现此接口即可。
type Renderer interface {
	Init() error
	Fini()
	PollEvent() tcell.Event
	Sync()
	Draw(g *Game)
}

// ============================================================
// TermRenderer (tcell-based)
// ============================================================

type TermRenderer struct {
	scr tcell.Screen
}

func NewTermRenderer() (*TermRenderer, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	return &TermRenderer{scr: s}, nil
}

func (r *TermRenderer) Init() error             { return r.scr.Init() }
func (r *TermRenderer) Fini()                   { r.scr.Fini() }
func (r *TermRenderer) PollEvent() tcell.Event  { return r.scr.PollEvent() }
func (r *TermRenderer) Sync()                   { r.scr.Sync() }

var (
	colPathBg   = tcell.NewRGBColor(139, 69, 19)
	colPathFg   = tcell.NewRGBColor(210, 180, 140)
	colStartFg  = tcell.NewRGBColor(50, 255, 80)
	colEndFg    = tcell.NewRGBColor(255, 80, 80)
	colTitleFg  = tcell.NewRGBColor(120, 220, 255)
	colCursorFg = tcell.NewRGBColor(255, 220, 80)
	colStatusFg = tcell.NewRGBColor(220, 220, 220)
	colHelpFg   = tcell.NewRGBColor(140, 140, 140)
	colSelBg    = tcell.NewRGBColor(60, 60, 100)
	colPrepFg   = tcell.NewRGBColor(255, 200, 80)
	colMenuFg   = tcell.NewRGBColor(220, 220, 240)
)

func putCell(s tcell.Screen, x, y int, ch1, ch2 rune, style tcell.Style) {
	col := x * cellW
	row := y + gameTopRow
	s.SetContent(col, row, ch1, nil, style)
	s.SetContent(col+1, row, ch2, nil, style)
}

func drawString(s tcell.Screen, col, row int, str string, style tcell.Style) {
	for i, ch := range str {
		s.SetContent(col+i, row, ch, nil, style)
	}
}

func (r *TermRenderer) Draw(g *Game) {
	r.scr.Clear()
	if g.Phase == PhaseLevelSelect {
		r.drawLevelSelect(g)
	} else {
		r.drawGame(g)
	}
	r.scr.Show()
}

func (r *TermRenderer) drawLevelSelect(g *Game) {
	stTitle := tcell.StyleDefault.Foreground(colTitleFg).Bold(true)
	stMenu := tcell.StyleDefault.Foreground(colMenuFg)
	stHelp := tcell.StyleDefault.Foreground(colHelpFg)

	drawString(r.scr, 0, 0, " Kingdom Rush V1.5  —  Select Level", stTitle)
	drawString(r.scr, 0, 1, " ===================================", stHelp)

	for i, lv := range g.Levels {
		row := 3 + i
		key := i + 1
		keyStr := fmt.Sprintf("%d", key)
		if key == 10 {
			keyStr = "0"
		}
		line := fmt.Sprintf("   [%s]  Lv %2d  %-18s  (waves: %d  gold: %d  lives: %d)",
			keyStr, lv.ID, lv.Name, len(lv.Waves), lv.StartGold, lv.StartLives)
		drawString(r.scr, 0, row, line, stMenu)
	}
	drawString(r.scr, 0, 3+len(g.Levels)+1,
		" Press 1-9 or 0 to start | Q/Esc to quit", stHelp)
}

func (r *TermRenderer) drawGame(g *Game) {
	s := r.scr

	stPath := tcell.StyleDefault.Background(colPathBg).Foreground(colPathFg)
	stCursor := tcell.StyleDefault.Foreground(colCursorFg).Bold(true)
	stTitle := tcell.StyleDefault.Foreground(colTitleFg).Bold(true)
	stStatus := tcell.StyleDefault.Foreground(colStatusFg)
	stHelp := tcell.StyleDefault.Foreground(colHelpFg)
	stPrep := tcell.StyleDefault.Foreground(colPrepFg).Bold(true)

	lv := g.currentLevel()
	if lv == nil {
		return
	}

	// title
	title := fmt.Sprintf(" KR V1.5 — Lv %d: %s — Wave %d/%d ",
		lv.ID, lv.Name, g.WaveIdx+1, len(lv.Waves))
	drawString(s, 0, 0, title, stTitle)

	// path
	for _, p := range g.Path {
		putCell(s, p.X, p.Y, '▒', '▒', stPath)
	}
	putCell(s, g.Path[0].X, g.Path[0].Y, ' ', 'S',
		stPath.Foreground(colStartFg).Bold(true))
	end := g.Path[len(g.Path)-1]
	putCell(s, end.X, end.Y, ' ', 'E', stPath.Foreground(colEndFg).Bold(true))

	// towers
	for _, t := range g.Towers {
		spec := towerSpecs[t.Kind]
		lvl := t.Spec()
		st := tcell.StyleDefault.Foreground(spec.Color).Bold(true)
		putCell(s, t.Pos.X, t.Pos.Y, lvl.Char1, lvl.Char2, st)
	}

	// enemies
	for _, e := range g.Enemies {
		if e.Dead || e.Escaped {
			continue
		}
		spec := enemySpecs[e.Kind]
		ch1, ch2 := spec.Char1, spec.Char2
		if e.HP <= e.MaxHP/3 {
			ch1, ch2 = 'x', 'x'
		}
		st := tcell.StyleDefault.Background(colPathBg).Foreground(spec.Color).Bold(true)
		p := e.Pos(g.Path)
		putCell(s, p.X, p.Y, ch1, ch2, st)
	}

	// cursor (path / tower / empty 三态)
	cur1, cur2 := '[', ']'
	cStyle := stCursor
	onPath := g.pathContains(g.Cursor)
	var atTower *Tower
	for _, t := range g.Towers {
		if t.Pos == g.Cursor {
			atTower = t
			break
		}
	}
	if onPath {
		cur1, cur2 = 'x', 'x'
		cStyle = tcell.StyleDefault.Background(colPathBg).Foreground(colEndFg).Bold(true)
	} else if atTower != nil {
		cur1, cur2 = '<', '>'
		cStyle = stCursor
	}
	putCell(s, g.Cursor.X, g.Cursor.Y, cur1, cur2, cStyle)

	// status row
	statusRow := gameTopRow + mapH
	leftStatus := fmt.Sprintf(" Gold: %-4d Lives: %d/%d Enemies: %-3d  %s",
		g.Gold, g.Lives, g.StartLives, g.CountAliveEnemies(), g.Msg)
	drawString(s, 0, statusRow, leftStatus, stStatus)
	if g.prepTimer > 0 && g.Phase == PhasePlaying {
		prepMsg := fmt.Sprintf(" PREP: %.1fs ", g.prepTimer)
		drawString(s, mapW*cellW-len(prepMsg), statusRow, prepMsg, stPrep)
	}

	// 塔选择行 (+ 升级提示如果光标在塔上)
	selRow := statusRow + 1
	col := 1
	for _, k := range TowerKinds() {
		spec := towerSpecs[k]
		label := fmt.Sprintf(" [%d] %s %dg ", int(k)+1, spec.Name, spec.Levels[0].Cost)
		st := tcell.StyleDefault.Foreground(spec.Color)
		if k == g.Selected {
			st = st.Background(colSelBg).Bold(true)
		}
		drawString(s, col, selRow, label, st)
		col += len(label) + 1
	}
	if atTower != nil {
		cost, can := atTower.NextUpgradeCost()
		var hint string
		if can {
			hint = fmt.Sprintf(" Space=Upgrade (%dg) ", cost)
		} else {
			hint = " MAX LEVEL "
		}
		drawString(s, col, selRow, hint, stCursor)
	}

	// help row
	drawString(s, 0, statusRow+2,
		" Arrows: move  1/2: select  Space: build/upgrade  M: menu  Q/Esc: quit",
		stHelp)

	// banner
	centerCol := mapW * cellW / 2
	centerRow := gameTopRow + mapH/2
	if g.Phase == PhaseWon {
		banner := fmt.Sprintf(" *** VICTORY *** Level %d cleared! ", lv.ID)
		drawString(s, centerCol-len(banner)/2, centerRow, banner,
			tcell.StyleDefault.Foreground(colStartFg).Bold(true))
		drawString(s, centerCol-15, centerRow+1, " Press M for menu, Q to quit ", stHelp)
	}
	if g.Phase == PhaseLost {
		banner := " *** GAME OVER *** "
		drawString(s, centerCol-len(banner)/2, centerRow, banner,
			tcell.StyleDefault.Foreground(colEndFg).Bold(true))
		drawString(s, centerCol-15, centerRow+1, " Press M for menu, Q to quit ", stHelp)
	}
}
