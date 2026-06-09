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
	// V9 P1: 关内成长 (per-run, beginRun 重置为 1 级 0 XP)
	Level int
	XP    int
}

// V9 P1: 英雄关内成长参数 (per-run; P4 平衡校准)。
const (
	heroLevelCap    = 5
	heroHPPerLvl    = 25  // 每级 maxHP 增量 (120 → 220 @ lvl5)
	heroDmgPerLvl   = 5   // 每级 damage 增量 (15 → 35 @ lvl5)
	heroRangePerLvl = 0.1 // 每级 攻击射程增量 (1.8 → 2.2 @ lvl5)
)

// xpForNextLevel: 从 level 升到 level+1 所需 XP (击杀威胁加权累计, 见 enemyCost)。
// 线性递增: lvl1→2 需 6, 2→3 需 12 … 满级累计 60。
func xpForNextLevel(level int) int { return level * 6 }

// heroMaxHPFor: 指定等级的 maxHP。
func heroMaxHPFor(level int) int { return heroSpec.MaxHP + (level-1)*heroHPPerLvl }

// Damage: 当前等级攻击伤害。
func (h *Hero) Damage() int { return heroSpec.Damage + (h.Level-1)*heroDmgPerLvl }

// AttackRange: 当前等级攻击射程。
func (h *Hero) AttackRange() float64 {
	return heroSpec.Range + float64(h.Level-1)*heroRangePerLvl
}

// GainXP: 累计 XP 并按阈值升级 (升级回满血, 满级后 XP 锁 0)。返回升级后的等级
// (调用方据此判断是否升级以给反馈)。
func (h *Hero) GainXP(xp int) {
	if h.Level >= heroLevelCap {
		return
	}
	h.XP += xp
	for h.Level < heroLevelCap && h.XP >= xpForNextLevel(h.Level) {
		h.XP -= xpForNextLevel(h.Level)
		h.Level++
		h.MaxHP = heroMaxHPFor(h.Level)
		h.HP = h.MaxHP // 升级回满血 (KR 范式)
	}
	if h.Level >= heroLevelCap {
		h.XP = 0
	}
}

// DistTo: 英雄到格 (px,py) 的欧氏距离 (cell)。
func (h *Hero) DistTo(px, py float64) float64 {
	return math.Hypot(px-h.X, py-h.Y)
}

// newHero: 在 spawn 点生成满血 1 级英雄, 集结点 = spawn (原地待命)。
// V9: Level/XP 每关从 newHero 起 (per-run, beginRun 调用 = 关内重置)。
func newHero(spawnX, spawnY float64) *Hero {
	return &Hero{
		X: spawnX, Y: spawnY, RallyX: spawnX, RallyY: spawnY,
		HP: heroMaxHPFor(1), MaxHP: heroMaxHPFor(1),
		Level: 1, XP: 0,
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
