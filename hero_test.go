// V8 Phase 1 测试: 英雄实体 + rally 移动 (纯逻辑)。
package main

import (
	"math"
	"testing"
)

// TestStepToward: 移动一步的三种情形 — 部分移动 / 到达(超出夹断) / 零距离。
func TestStepToward(t *testing.T) {
	// 部分移动: 从原点朝 (10,0) 走 3 → (3,0), 未到达
	nx, ny, arrived := stepToward(0, 0, 10, 0, 3)
	if nx != 3 || ny != 0 || arrived {
		t.Errorf("partial move: got (%v,%v) arrived=%v, want (3,0) arrived=false", nx, ny, arrived)
	}
	// 到达 (超出夹断): 剩余距离 4, 步长 10 → 落目标点, 到达
	nx, ny, arrived = stepToward(0, 0, 0, 4, 10)
	if nx != 0 || ny != 4 || !arrived {
		t.Errorf("overshoot clamp: got (%v,%v) arrived=%v, want (0,4) arrived=true", nx, ny, arrived)
	}
	// 零距离: 已在目标点 → 原地, 到达
	nx, ny, arrived = stepToward(5, 5, 5, 5, 3)
	if nx != 5 || ny != 5 || !arrived {
		t.Errorf("zero dist: got (%v,%v) arrived=%v, want (5,5) arrived=true", nx, ny, arrived)
	}
	// 对角线: 朝 (3,4) [距离 5] 走 2.5 → 走一半 (1.5,2.0)
	nx, ny, arrived = stepToward(0, 0, 3, 4, 2.5)
	if math.Abs(nx-1.5) > 1e-9 || math.Abs(ny-2.0) > 1e-9 || arrived {
		t.Errorf("diagonal: got (%v,%v) arrived=%v, want (1.5,2.0) arrived=false", nx, ny, arrived)
	}
}

// TestNewHero: 生成满血英雄, 集结点 = spawn (原地待命)。
func TestNewHero(t *testing.T) {
	h := newHero(8, 6)
	if h.X != 8 || h.Y != 6 {
		t.Errorf("spawn pos: got (%v,%v), want (8,6)", h.X, h.Y)
	}
	if h.RallyX != 8 || h.RallyY != 6 {
		t.Errorf("rally should default to spawn: got (%v,%v)", h.RallyX, h.RallyY)
	}
	if h.HP != heroSpec.MaxHP || h.MaxHP != heroSpec.MaxHP {
		t.Errorf("HP: got %d/%d, want %d full", h.HP, h.MaxHP, heroSpec.MaxHP)
	}
	if !h.Alive() {
		t.Error("new hero should be alive")
	}
}

// TestHeroMoveTowardRally: 设集结点后逐帧移动, 最终到达 (确定性收敛)。
func TestHeroMoveTowardRally(t *testing.T) {
	h := newHero(0, 0)
	h.SetRally(10, 0)
	// 单帧 (dt=0.1): 走 Speed*0.1 = 0.4 cell
	h.moveStep(0.1)
	want := heroSpec.Speed * 0.1
	if math.Abs(h.X-want) > 1e-9 {
		t.Errorf("after 1 frame: X=%v, want %v", h.X, want)
	}
	// 多帧推进直到到达 (上限防死循环)
	for i := 0; i < 1000 && (math.Abs(h.X-10) > 1e-9 || math.Abs(h.Y-0) > 1e-9); i++ {
		h.moveStep(0.1)
	}
	if math.Abs(h.X-10) > 1e-9 || math.Abs(h.Y) > 1e-9 {
		t.Errorf("did not converge to rally: got (%v,%v), want (10,0)", h.X, h.Y)
	}
}

// TestHeroSetRallyDeadGuard: 阵亡期间 (respawnCD>0) SetRally 被忽略。
func TestHeroSetRallyDeadGuard(t *testing.T) {
	h := newHero(5, 5)
	h.respawnCD = 3.0 // 模拟阵亡 (P2 才会真正设置)
	if h.Alive() {
		t.Error("hero with respawnCD>0 should not be alive")
	}
	h.SetRally(20, 20)
	if h.RallyX != 5 || h.RallyY != 5 {
		t.Errorf("dead hero rally must not change: got (%v,%v), want (5,5)", h.RallyX, h.RallyY)
	}
	// 阵亡期间 moveStep 不动
	h.moveStep(0.5)
	if h.X != 5 || h.Y != 5 {
		t.Errorf("dead hero must not move: got (%v,%v), want (5,5)", h.X, h.Y)
	}
}

// TestBeginRunSpawnsHero: 进入关卡后英雄生成在 path 中点, 满血在场。
func TestBeginRunSpawnsHero(t *testing.T) {
	g := newTestGame()
	if g.Hero == nil {
		t.Fatal("StartLevel should spawn a hero")
	}
	mid := g.Path[len(g.Path)/2]
	if g.Hero.X != float64(mid.X) || g.Hero.Y != float64(mid.Y) {
		t.Errorf("hero spawn: got (%v,%v), want path mid (%d,%d)",
			g.Hero.X, g.Hero.Y, mid.X, mid.Y)
	}
	if g.Hero.HP != heroSpec.MaxHP || !g.Hero.Alive() {
		t.Errorf("hero should spawn full HP alive: HP=%d alive=%v", g.Hero.HP, g.Hero.Alive())
	}
}

// TestSetHeroRallyMovesHero: SetHeroRally 设集结点后, Update 驱动英雄移动。
func TestSetHeroRallyMovesHero(t *testing.T) {
	g := newTestGame()
	startX := g.Hero.X
	g.SetHeroRally(Point{X: int(startX) + 5, Y: int(g.Hero.Y)})
	g.Update(0.5) // 半秒 → 应朝集结点移动
	if g.Hero.X <= startX {
		t.Errorf("hero should move toward rally: startX=%v nowX=%v", startX, g.Hero.X)
	}
}
