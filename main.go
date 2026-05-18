// kingdom-rush V1.5 entry
//
// 流程: LevelSelect (1-9/0 选关) → Playing → Won/Lost → (M 返回菜单)
//
// 控制:
//   Menu:   1-9, 0 选关 | Q/Esc 退出
//   Game:   Arrows 移光标 | 1/2 选塔型 | Space 建塔/升塔
//           M 回菜单 | Q/Esc 退出
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	levels, err := LoadLevels()
	if err != nil {
		fmt.Fprintln(os.Stderr, "load levels:", err)
		os.Exit(1)
	}
	if len(levels) == 0 {
		fmt.Fprintln(os.Stderr, "no levels loaded")
		os.Exit(1)
	}

	save, err := LoadSave()
	if err != nil {
		// 加载失败不阻断启动(可能是 corrupt save),fallback empty save 让玩家继续
		fmt.Fprintln(os.Stderr, "warning: load save failed:", err, "(starting fresh)")
		save = NewSave()
	}

	r, err := NewTermRenderer()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := r.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer r.Fini()

	g := NewGame(levels, save)

	evCh := make(chan tcell.Event, 16)
	done := make(chan struct{})
	go func() {
		for {
			ev := r.PollEvent()
			if ev == nil {
				return
			}
			select {
			case evCh <- ev:
			case <-done:
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
			if quit := handleEvent(g, ev, r); quit {
				close(done)
				return
			}
		case <-ticker.C:
			now := time.Now()
			dt := now.Sub(last).Seconds()
			last = now
			g.Update(dt)
			r.Draw(g)
		}
	}
}

func handleEvent(g *Game, ev tcell.Event, r Renderer) (quit bool) {
	switch e := ev.(type) {
	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyEscape:
			return true
		case tcell.KeyUp:
			if g.Phase == PhasePlaying {
				g.MoveCursor(0, -1)
			}
		case tcell.KeyDown:
			if g.Phase == PhasePlaying {
				g.MoveCursor(0, 1)
			}
		case tcell.KeyLeft:
			if g.Phase == PhasePlaying {
				g.MoveCursor(-1, 0)
			}
		case tcell.KeyRight:
			if g.Phase == PhasePlaying {
				g.MoveCursor(1, 0)
			}
		}
		switch e.Rune() {
		case 'q', 'Q':
			return true
		case 'm', 'M':
			g.BackToMenu()
		case ' ':
			if g.Phase == PhasePlaying {
				g.TryAction()
			}
		case '1':
			if g.Phase == PhaseLevelSelect {
				g.StartLevel(0)
			} else if g.Phase == PhasePlaying {
				g.Selected = TArcher
				g.Msg = "Selected Archer"
			}
		case '2':
			if g.Phase == PhaseLevelSelect {
				g.StartLevel(1)
			} else if g.Phase == PhasePlaying {
				g.Selected = TCannon
				g.Msg = "Selected Cannon"
			}
		case '3':
			if g.Phase == PhaseLevelSelect {
				g.StartLevel(2)
			} else if g.Phase == PhasePlaying {
				g.Selected = TMagic
				g.Msg = "Selected Magic"
			}
		case '4', '5', '6', '7', '8', '9':
			if g.Phase == PhaseLevelSelect {
				idx := int(e.Rune() - '1')
				g.StartLevel(idx)
			}
		case '0':
			if g.Phase == PhaseLevelSelect {
				g.StartLevel(9) // Level 10
			}
		}
	case *tcell.EventResize:
		r.Sync()
	}
	return false
}
