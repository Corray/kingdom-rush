//go:build term

// Terminal (tcell) 渲染。V2 默认 build 用 ebiten,terminal 模式 opt-in:
//
//	go build -tags term
//
// V1.7 中 Renderer interface 实际未被 ebiten 复用(范式不同),V2 删除 interface,
// TermRenderer 仍为 terminal main 直接调用。
package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// termColor 把 UI-agnostic 的 RGB(entities.go) 转为 tcell.Color。
func termColor(c RGB) tcell.Color {
	return tcell.NewRGBColor(int32(c.R), int32(c.G), int32(c.B))
}

const (
	cellW      = 2
	gameTopRow = 1
)

// TermRenderer 是 V1.7 的 terminal 渲染器(tcell-based)。
// V2 删除了 Renderer interface(ebiten 范式与 tcell 不兼容,interface 没有
// 复用价值)。terminal mode 通过 build tag `term` 启用,term_main.go 直接
// 调用 TermRenderer 方法。
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

func (r *TermRenderer) Init() error            { return r.scr.Init() }
func (r *TermRenderer) Fini()                  { r.scr.Fini() }
func (r *TermRenderer) PollEvent() tcell.Event { return r.scr.PollEvent() }
func (r *TermRenderer) Sync()                  { r.scr.Sync() }

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
	colLockedFg = tcell.NewRGBColor(100, 100, 100)
	colHeroFg   = tcell.NewRGBColor(255, 200, 60) // V8 P4: 英雄金身 (Knight)
	// V10 P3: 职业配色 (与 ebiten 端一致)
	colArcherFg = tcell.NewRGBColor(110, 220, 110) // Archer 绿
	colRogueFg  = tcell.NewRGBColor(190, 120, 255) // Rogue 紫
)

// heroClassTermColor: V10 P3 — 职业字形配色 (Knight 金 / Archer 绿 / Rogue 紫)。
func heroClassTermColor(c *HeroClass) tcell.Color {
	switch c.Name {
	case "Archer":
		return colArcherFg
	case "Rogue":
		return colRogueFg
	}
	return colHeroFg
}

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
	} else if g.Phase == PhaseSkillTree {
		r.drawSkillTree(g)
	} else {
		r.drawGame(g)
	}
	r.scr.Show()
}

// drawSkillTree: V11 P3 — 技能树屏 (term)。职业纵列, '>' 标选中;
// 节点状态: [x] 已购 / [N*] 下一个 / 灰 = 锁定。
func (r *TermRenderer) drawSkillTree(g *Game) {
	stTitle := tcell.StyleDefault.Foreground(colTitleFg).Bold(true)
	stHelp := tcell.StyleDefault.Foreground(colHelpFg)
	stStar := tcell.StyleDefault.Foreground(colPrepFg).Bold(true)
	stDone := tcell.StyleDefault.Foreground(colStartFg)
	stLock := tcell.StyleDefault.Foreground(colLockedFg)
	stNext := tcell.StyleDefault.Foreground(colMenuFg).Bold(true)
	stMsg := tcell.StyleDefault.Foreground(colPrepFg)

	drawString(r.scr, 0, 0, " Gopher Defense  —  Skill Trees", stTitle)
	drawString(r.scr, 0, 1, fmt.Sprintf(" Stars: %d available / %d earned",
		g.Save.AvailableStars(), g.Save.TotalStars()), stStar)

	row := 3
	for ci := range heroClasses {
		c := &heroClasses[ci]
		tree := skillTrees[c.Name]
		lvl := g.Save.TreeLevel(c.Name)
		marker := "  "
		if ci == g.TreeClassIdx {
			marker = "> "
		}
		stClass := tcell.StyleDefault.Foreground(heroClassTermColor(c)).Bold(true)
		drawString(r.scr, 0, row, fmt.Sprintf(" %s%s  %d/%d", marker, c.Name, lvl, treeNodesPerClass), stClass)
		row++
		for ni := 0; ni < treeNodesPerClass; ni++ {
			node := tree[ni]
			var line string
			style := stLock
			switch {
			case ni < lvl:
				line = fmt.Sprintf("     [x] %-16s %s", node.Name, node.Desc)
				style = stDone
			case ni == lvl:
				line = fmt.Sprintf("     [%d*] %-16s %s", node.Price, node.Name, node.Desc)
				style = stNext
			default:
				line = fmt.Sprintf("     [%d*] %-16s %s", node.Price, node.Name, node.Desc)
			}
			drawString(r.scr, 0, row, line, style)
			row++
		}
		row++
	}
	drawString(r.scr, 0, row, " Left/Right: class | Space: learn | T/M: back", stHelp)
	if g.Msg != "" {
		drawString(r.scr, 0, row+1, " "+g.Msg, stMsg)
	}
}

func (r *TermRenderer) drawLevelSelect(g *Game) {
	stTitle := tcell.StyleDefault.Foreground(colTitleFg).Bold(true)
	stMenu := tcell.StyleDefault.Foreground(colMenuFg)
	stDone := tcell.StyleDefault.Foreground(colStartFg).Bold(true)
	stLock := tcell.StyleDefault.Foreground(colLockedFg)
	stHelp := tcell.StyleDefault.Foreground(colHelpFg)
	stMsg := tcell.StyleDefault.Foreground(colPrepFg)

	drawString(r.scr, 0, 0, " Gopher Defense  —  Select Level", stTitle)
	drawString(r.scr, 0, 1, " ===================================", stHelp)

	completed := 0
	for _, lv := range g.Levels {
		if g.Save.IsCompleted(lv.ID) {
			completed++
		}
	}
	progress := fmt.Sprintf(" Progress: %d / %d cleared", completed, len(g.Levels))
	drawString(r.scr, 0, 2, progress, stHelp)
	// V10 P3: 当前英雄职业 (H 切换); V11 P3: T 技能树入口
	heroClass := &heroClasses[g.Save.HeroClassIdx()]
	drawString(r.scr, 0, 3, fmt.Sprintf(" Hero: %s (H to change, T: skill trees)", heroClass.Name),
		tcell.StyleDefault.Foreground(heroClassTermColor(heroClass)).Bold(true))

	for i, lv := range g.Levels {
		row := 5 + i
		key := i + 1
		keyStr := fmt.Sprintf("%d", key)
		if key == 10 {
			keyStr = "0"
		}
		status := "[    ]"
		style := stMenu
		if g.Save.IsCompleted(lv.ID) {
			status = "[DONE]"
			style = stDone
		} else if !g.Save.IsUnlocked(lv.ID) {
			status = "[LOCK]"
			style = stLock
		}
		line := fmt.Sprintf("   [%s] %s  Lv %2d  %-18s  (waves:%d  gold:%d  lives:%d)",
			keyStr, status, lv.ID, lv.Name, len(lv.Waves), lv.StartGold, lv.StartLives)
		drawString(r.scr, 0, row, line, style)
	}
	helpRow := 5 + len(g.Levels) + 1
	drawString(r.scr, 0, helpRow,
		" Press 1-9 or 0 to start unlocked level | Q/Esc to quit", stHelp)
	if g.Msg != "" {
		drawString(r.scr, 0, helpRow+1, " "+g.Msg, stMsg)
	}
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
	title := fmt.Sprintf(" KR V1.7 — Lv %d: %s — Wave %d/%d ",
		lv.ID, lv.Name, g.WaveIdx+1, len(lv.Waves))
	drawString(s, 0, 0, title, stTitle)

	// path — V12: 遍历所有 path (单路 = 1 条; 汇流段重合 cell 同字符叠加)
	for _, path := range g.Paths {
		for _, p := range path {
			putCell(s, p.X, p.Y, '▒', '▒', stPath)
		}
	}
	// 起点 / 终点 markers — V12 P1: 暂用 Paths[0] (单路行为不变); P2 多起点
	p0 := g.Paths[0]
	putCell(s, p0[0].X, p0[0].Y, ' ', 'S',
		stPath.Foreground(colStartFg).Bold(true))
	end := p0[len(p0)-1]
	putCell(s, end.X, end.Y, ' ', 'E', stPath.Foreground(colEndFg).Bold(true))

	// towers
	for _, t := range g.Towers {
		spec := towerSpecs[t.Kind]
		lvl := t.Spec()
		st := tcell.StyleDefault.Foreground(termColor(spec.Color)).Bold(true)
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
		st := tcell.StyleDefault.Background(colPathBg).Foreground(termColor(spec.Color)).Bold(true)
		p := e.Pos(g.Paths[e.PathID])
		putCell(s, p.X, p.Y, ch1, ch2, st)
	}

	// V8 P4: 英雄 (自由坐标取整; 职业首字母+职业色 / 阵亡 '+' 灰墓碑)
	if g.Hero != nil {
		hx, hy := int(g.Hero.X+0.5), int(g.Hero.Y+0.5)
		if g.Hero.Alive() {
			glyph := rune(g.Hero.Class.Name[0]) // V10 P3: K/A/R
			putCell(s, hx, hy, glyph, glyph,
				tcell.StyleDefault.Foreground(heroClassTermColor(g.Hero.Class)).Bold(true))
		} else {
			putCell(s, hx, hy, '+', '+',
				tcell.StyleDefault.Foreground(colLockedFg))
		}
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

	// status row (V8 P4: 加英雄 HP / 复活倒计时)
	statusRow := gameTopRow + mapH
	heroStr := ""
	if g.Hero != nil {
		if g.Hero.Alive() {
			heroStr = fmt.Sprintf("%s:L%d %d/%d ", g.Hero.Class.Name, g.Hero.Level, g.Hero.HP, g.Hero.MaxHP)
			if g.Hero.AbilityReady() { // V9: 横扫就绪提示
				heroStr += "[G!] "
			}
		} else {
			heroStr = fmt.Sprintf("Hero:DOWN %.0fs ", g.Hero.respawnCD)
		}
	}
	leftStatus := fmt.Sprintf(" Gold: %-4d Lives: %d/%d Enemies: %-3d %s %s",
		g.Gold, g.Lives, g.StartLives, g.CountAliveEnemies(), heroStr, g.Msg)
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
		st := tcell.StyleDefault.Foreground(termColor(spec.Color))
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
		" Arrows: move  1/2: select  Space: build/upgrade  H: rally hero  G: cleave  M: menu  Q/Esc: quit",
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
