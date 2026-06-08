// V8 Phase 1: 英雄单位 — 首个非路径绑定实体 (自由 float 坐标)。
//
// 设计: 玩家用光标 + H 键设集结点 (RallyX/Y), 英雄朝集结点移动 (P1)
// 并自动攻击射程内敌人 (P2), 贴敌时阻挡 (P3)。
// 纯逻辑无 ebiten 依赖, 可测。per-run 实体 (同塔), 不入存档。
package main

import "math"

type HeroSpec struct {
	MaxHP    int
	Speed    float64 // 移动速度 (cell/s)
	Damage   int     // P2: 攻击伤害
	Range    float64 // P2: 攻击射程 (cell)
	AttackCD float64 // P2: 攻击间隔 (s)
	RespawnS float64 // P2: 复活时间 (s)
}

// heroSpec: 单英雄数值 (P5 平衡校准)。Speed 4.0 介于 ENormal(3.0) 与
// EFast(5.5) — 追得上普通敌、追不上快敌, 鼓励守位而非追击。
var heroSpec = HeroSpec{
	MaxHP: 120, Speed: 4.0, Damage: 15, Range: 1.8, AttackCD: 0.7, RespawnS: 12.0,
}

type Hero struct {
	X, Y      float64 // 自由坐标 (cell, 非路径绑定)
	RallyX    float64 // 集结点 (玩家 H 键设定)
	RallyY    float64
	HP        int
	MaxHP     int
	cooldown  float64 // P2: 攻击冷却剩余 (>0 不能出手)
	respawnCD float64 // P2: >0 = 阵亡复活倒计时 (此期间不在场)
}

// DistTo: 英雄到格 (px,py) 的欧氏距离 (cell)。
func (h *Hero) DistTo(px, py float64) float64 {
	return math.Hypot(px-h.X, py-h.Y)
}

// newHero: 在 spawn 点生成满血英雄, 集结点 = spawn (原地待命)。
func newHero(spawnX, spawnY float64) *Hero {
	return &Hero{
		X: spawnX, Y: spawnY, RallyX: spawnX, RallyY: spawnY,
		HP: heroSpec.MaxHP, MaxHP: heroSpec.MaxHP,
	}
}

// Alive: 是否在场 (阵亡复活期间为 false)。
func (h *Hero) Alive() bool { return h.respawnCD <= 0 }

// SetRally: 设集结点 (玩家 H 键)。阵亡期间忽略 (不能指挥死人)。
func (h *Hero) SetRally(x, y float64) {
	if h.respawnCD > 0 {
		return
	}
	h.RallyX, h.RallyY = x, y
}

// stepToward: 从 (x,y) 朝 (tx,ty) 移动 dist 距离, 不超过目标。
// 返回新坐标 + 是否到达 (剩余距离 ≤ dist 即到达, 直接落目标点)。
func stepToward(x, y, tx, ty, dist float64) (float64, float64, bool) {
	dx, dy := tx-x, ty-y
	d := math.Hypot(dx, dy)
	if d <= dist || d == 0 {
		return tx, ty, true
	}
	return x + dx/d*dist, y + dy/d*dist, false
}

// moveStep: P1 — 朝集结点移动一帧。阵亡期间不动。
func (h *Hero) moveStep(dt float64) {
	if h.respawnCD > 0 {
		return
	}
	h.X, h.Y, _ = stepToward(h.X, h.Y, h.RallyX, h.RallyY, heroSpec.Speed*dt)
}
