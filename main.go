// kingdom-rush V0 — terminal ASCII 塔防 PoC
// 一图 / 一波 / 一种塔,验证核心循环。方向键移动光标,Space 放塔,Q/Esc 退出。
//
// V0.1 视觉升级:
//   - 每 game cell 渲染为 2 columns wide(补偿 terminal 字符宽高 1:2)
//   - path 用 ▒▒ + 棕色"土路"
//   - tower ▲▲ / enemy ●● 重伤 ✖✖ / cursor [] on-path ✗✗
package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	mapW  = 30
	mapH  = 15
	cellW = 2 // 每 game cell 渲染宽度 (字符 1:2 补偿)
	fps   = 30

	gameTopRow = 1 // row 0 = title

	startGold   = 100
	startLives  = 5
	towerCost   = 50
	waveSize    = 8
	spawnGapS   = 1.0
	enemyHP     = 20
	enemySpeed  = 3.0
	towerRange  = 3.5
	towerDmg    = 8
	towerCD     = 0.6
	goldPerKill = 10
)

// ============================================================
// Map / Path
// ============================================================

type Point struct{ X, Y int }

func buildPath() []Point {
	var p []Point
	for x := 0; x <= 7; x++ {
		p = append(p, Point{x, 3})
	}
	for y := 4; y <= 7; y++ {
		p = append(p, Point{7, y})
	}
	for x := 8; x <= 17; x++ {
		p = append(p, Point{x, 7})
	}
	for y := 6; y >= 3; y-- {
		p = append(p, Point{17, y})
	}
	for x := 18; x <= 29; x++ {
		p = append(p, Point{x, 3})
	}
	return p
}

var path = buildPath()

func pathContains(p Point) bool {
	for _, q := range path {
		if p == q {
			return true
		}
	}
	return false
}

// ============================================================
// Entities
// ============================================================

type Tower struct {
	Pos      Point
	cooldown float64
}

type Enemy struct {
	PathIdx float64
	HP      int
	MaxHP   int
	Dead    bool
	Escaped bool
}

func (e *Enemy) Pos() Point {
	i := int(math.Floor(e.PathIdx))
	if i >= len(path) {
		return path[len(path)-1]
	}
	if i < 0 {
		return path[0]
	}
	return path[i]
}

// ============================================================
// Game
// ============================================================

type Game struct {
	Towers  []*Tower
	Enemies []*Enemy
	Gold    int
	Lives   int
	Wave    int
	Cursor  Point
	Won     bool
	Lost    bool

	spawned    int
	spawnTimer float64
}

func NewGame() *Game {
	return &Game{
		Gold:   startGold,
		Lives:  startLives,
		Wave:   1,
		Cursor: Point{15, 10},
	}
}

func (g *Game) Update(dt float64) {
	if g.Won || g.Lost {
		return
	}

	if g.spawned < waveSize {
		g.spawnTimer += dt
		if g.spawnTimer >= spawnGapS {
			g.spawnTimer = 0
			g.Enemies = append(g.Enemies, &Enemy{HP: enemyHP, MaxHP: enemyHP})
			g.spawned++
		}
	}

	for _, e := range g.Enemies {
		if e.Dead || e.Escaped {
			continue
		}
		e.PathIdx += enemySpeed * dt
		if e.PathIdx >= float64(len(path)-1) {
			e.Escaped = true
			g.Lives--
			if g.Lives <= 0 {
				g.Lost = true
				return
			}
		}
	}

	for _, t := range g.Towers {
		t.cooldown -= dt
		if t.cooldown > 0 {
			continue
		}
		var target *Enemy
		maxIdx := -1.0
		for _, e := range g.Enemies {
			if e.Dead || e.Escaped {
				continue
			}
			ep := e.Pos()
			dx := float64(ep.X - t.Pos.X)
			dy := float64(ep.Y - t.Pos.Y)
			if math.Sqrt(dx*dx+dy*dy) <= towerRange && e.PathIdx > maxIdx {
				maxIdx = e.PathIdx
				target = e
			}
		}
		if target != nil {
			target.HP -= towerDmg
			t.cooldown = towerCD
			if target.HP <= 0 {
				target.Dead = true
				g.Gold += goldPerKill
			}
		}
	}

	if g.spawned >= waveSize {
		alive := false
		for _, e := range g.Enemies {
			if !e.Dead && !e.Escaped {
				alive = true
				break
			}
		}
		if !alive {
			g.Won = true
		}
	}
}

func (g *Game) TryPlaceTower() string {
	if g.Gold < towerCost {
		return fmt.Sprintf("Not enough gold (need %d)", towerCost)
	}
	if pathContains(g.Cursor) {
		return "Cannot build on path"
	}
	for _, t := range g.Towers {
		if t.Pos == g.Cursor {
			return "Tower already here"
		}
	}
	g.Towers = append(g.Towers, &Tower{Pos: g.Cursor})
	g.Gold -= towerCost
	return "Tower built"
}

func countAlive(es []*Enemy) int {
	n := 0
	for _, e := range es {
		if !e.Dead && !e.Escaped {
			n++
		}
	}
	return n
}

// ============================================================
// Render
// ============================================================

// putCell 在 game 坐标 (x,y) 上写 2 字符宽的一个 cell
func putCell(s tcell.Screen, x, y int, ch1, ch2 rune, style tcell.Style) {
	col := x * cellW
	row := y + gameTopRow
	s.SetContent(col, row, ch1, nil, style)
	s.SetContent(col+1, row, ch2, nil, style)
}

// drawString 在真实 (col, row) 写字符串
func drawString(s tcell.Screen, col, row int, str string, style tcell.Style) {
	for i, ch := range str {
		s.SetContent(col+i, row, ch, nil, style)
	}
}

var (
	colPathBg   = tcell.NewRGBColor(139, 69, 19)   // SaddleBrown
	colPathFg   = tcell.NewRGBColor(210, 180, 140) // Tan
	colStartFg  = tcell.NewRGBColor(50, 255, 80)   // bright green
	colEndFg    = tcell.NewRGBColor(255, 80, 80)   // bright red
	colTitleFg  = tcell.NewRGBColor(120, 220, 255) // cyan-ish
	colTowerFg  = tcell.NewRGBColor(100, 180, 255) // sky blue
	colEnemyFg  = tcell.NewRGBColor(255, 100, 100) // red
	colCursorFg = tcell.NewRGBColor(255, 220, 80)  // gold
	colStatusFg = tcell.NewRGBColor(220, 220, 220) // light gray
	colHelpFg   = tcell.NewRGBColor(140, 140, 140) // dim gray
)

func draw(s tcell.Screen, g *Game, msg string) {
	s.Clear()

	stPath := tcell.StyleDefault.Background(colPathBg).Foreground(colPathFg)
	stTower := tcell.StyleDefault.Foreground(colTowerFg).Bold(true)
	stCursor := tcell.StyleDefault.Foreground(colCursorFg).Bold(true)
	stTitle := tcell.StyleDefault.Foreground(colTitleFg).Bold(true)
	stStatus := tcell.StyleDefault.Foreground(colStatusFg)
	stHelp := tcell.StyleDefault.Foreground(colHelpFg)

	// title (row 0)
	title := fmt.Sprintf(" Kingdom Rush V0  —  Wave %d ", g.Wave)
	drawString(s, 0, 0, title, stTitle)

	// path
	for _, p := range path {
		putCell(s, p.X, p.Y, '▒', '▒', stPath)
	}
	// 起点 / 终点 (保留 path bg)
	stStart := stPath.Foreground(colStartFg).Bold(true)
	putCell(s, path[0].X, path[0].Y, ' ', 'S', stStart)
	end := path[len(path)-1]
	stEnd := stPath.Foreground(colEndFg).Bold(true)
	putCell(s, end.X, end.Y, ' ', 'E', stEnd)

	// towers
	for _, t := range g.Towers {
		putCell(s, t.Pos.X, t.Pos.Y, '▲', '▲', stTower)
	}

	// enemies
	for _, e := range g.Enemies {
		if e.Dead || e.Escaped {
			continue
		}
		p := e.Pos()
		ch1, ch2 := '●', '●'
		if e.HP <= e.MaxHP/3 {
			ch1, ch2 = '✖', '✖'
		}
		// 敌人在 path 上,保留棕底
		stE := tcell.StyleDefault.Background(colPathBg).Foreground(colEnemyFg).Bold(true)
		putCell(s, p.X, p.Y, ch1, ch2, stE)
	}

	// cursor
	cur1, cur2 := '[', ']'
	cStyle := stCursor
	if pathContains(g.Cursor) {
		cur1, cur2 = '✗', '✗'
		cStyle = tcell.StyleDefault.Background(colPathBg).Foreground(colEndFg).Bold(true)
	}
	putCell(s, g.Cursor.X, g.Cursor.Y, cur1, cur2, cStyle)

	// status (row = gameTopRow + mapH)
	statusRow := gameTopRow + mapH
	statusLine := fmt.Sprintf(" Gold: %-4d  Lives: %d/%d  Enemies: %d   %s",
		g.Gold, g.Lives, startLives, countAlive(g.Enemies), msg)
	drawString(s, 0, statusRow, statusLine, stStatus)
	drawString(s, 0, statusRow+1, " Arrows: move    Space: build (50g)    Q/Esc: quit", stHelp)

	// banner (居中 in 渲染坐标系: mapW*cellW = 60 cols)
	centerCol := mapW * cellW / 2
	centerRow := gameTopRow + mapH/2
	if g.Won {
		banner := " *** YOU WON *** "
		drawString(s, centerCol-len(banner)/2, centerRow, banner,
			tcell.StyleDefault.Foreground(colStartFg).Bold(true))
	}
	if g.Lost {
		banner := " *** GAME OVER *** "
		drawString(s, centerCol-len(banner)/2, centerRow, banner,
			tcell.StyleDefault.Foreground(colEndFg).Bold(true))
	}

	s.Show()
}

// ============================================================
// Main loop
// ============================================================

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer s.Fini()

	g := NewGame()
	msg := ""

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	evCh := make(chan tcell.Event, 16)
	go func() {
		for {
			ev := s.PollEvent()
			if ev == nil {
				return
			}
			select {
			case evCh <- ev:
			case <-ctx.Done():
				return
			}
		}
	}()

	ticker := time.NewTicker(time.Second / fps)
	defer ticker.Stop()

	last := time.Now()
	for {
		select {
		case ev := <-evCh:
			switch e := ev.(type) {
			case *tcell.EventKey:
				switch e.Key() {
				case tcell.KeyEscape:
					return
				case tcell.KeyUp:
					if g.Cursor.Y > 0 {
						g.Cursor.Y--
					}
				case tcell.KeyDown:
					if g.Cursor.Y < mapH-1 {
						g.Cursor.Y++
					}
				case tcell.KeyLeft:
					if g.Cursor.X > 0 {
						g.Cursor.X--
					}
				case tcell.KeyRight:
					if g.Cursor.X < mapW-1 {
						g.Cursor.X++
					}
				}
				switch e.Rune() {
				case ' ':
					msg = g.TryPlaceTower()
				case 'q', 'Q':
					return
				}
			case *tcell.EventResize:
				s.Sync()
			}
		case <-ticker.C:
			now := time.Now()
			dt := now.Sub(last).Seconds()
			last = now
			g.Update(dt)
			draw(s, g, msg)
		}
	}
}
