// V8 Phase 1: 英雄单位 — 首个非路径绑定实体 (自由 float 坐标)。
//
// 设计: 玩家用光标 + H 键设集结点 (RallyX/Y), 英雄朝集结点移动 (P1)
// 并自动攻击射程内敌人 (P2), 贴敌时阻挡 (P3)。
// 纯逻辑无 ebiten 依赖, 可测。per-run 实体 (同塔), 不入存档。
//
// V10 P1: heroSpec 单例 → HeroClass 参数化 (数值 + 成长 + 技能 + 阻挡
// 全部 per-class)。Knight = V8/V9 原数值, 行为零变化; Archer/Rogue P2 加。
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

// HeroClass: V10 — 英雄职业参数包。基础数值 (HeroSpec) + 关内成长曲线 (V9 P1)
// + 主动技能参数 (V9 P2, 统一「自身 AoE」范式只调参 — V10 决策 D)
// + 阻挡 gating (V10 决策 B)。
type HeroClass struct {
	Name string
	HeroSpec
	// 关内成长 (per-run)
	LevelCap    int
	HPPerLvl    int     // 每级 maxHP 增量
	DmgPerLvl   int     // 每级 damage 增量
	RangePerLvl float64 // 每级 攻击射程增量
	XPBase      int     // 升级阈值系数: lvl→lvl+1 需 lvl×XPBase XP
	// 主动技能: 自身周围 AoE (经 damageEnemy, 不打飞行)
	AbilityLevel     int     // 解锁等级
	AbilityCooldownS float64 // 冷却 (s)
	AbilityRadius    float64 // AoE 半径 (cell)
	AbilityDmgMul    int     // 伤害 = 当前 Damage() × 此值
	// V10 决策 B: 是否贴身阻挡非飞行敌 (近战职业 true, 远程 false)
	Blocks bool
}

// xpForNext: 从 level 升到 level+1 所需 XP (击杀威胁加权累计, 见 enemyCost)。
// Knight: 线性递增 lvl1→2 需 6, 2→3 需 12 … 满级累计 60。
func (c *HeroClass) xpForNext(level int) int { return level * c.XPBase }

// maxHPFor: 指定等级的 maxHP。
func (c *HeroClass) maxHPFor(level int) int { return c.MaxHP + (level-1)*c.HPPerLvl }

// heroClasses: 职业表。index 0 = Knight (默认, Save.HeroChoice 零值兼容)。
//
// 三职业定位 (V10 决策 A, 数值 P4 仿真校准):
//   - Knight: 近战坦 + 阻挡 (V8/V9 原数值, 零回归基线)。Speed 4.0 介于
//     ENormal(3.0) 与 EFast(5.5) — 追得上普通敌、追不上快敌, 守位定位
//   - Archer: 远程 (3.5 ≈ Archer 塔) 不阻挡 — 风筝输出, 换掉隘口控制
//   - Rogue: 高速 (5.5 = EFast, 追得上快敌) 低耐高攻速 + 阻挡 — 游走截击
var heroClasses = []HeroClass{
	{
		Name:     "Knight",
		HeroSpec: HeroSpec{MaxHP: 120, Speed: 4.0, Damage: 15, Range: 1.8, AttackCD: 0.7, RespawnS: 12.0},
		LevelCap: 5, HPPerLvl: 25, DmgPerLvl: 5, RangePerLvl: 0.1, XPBase: 6,
		AbilityLevel: 3, AbilityCooldownS: 8.0, AbilityRadius: 2.0, AbilityDmgMul: 3,
		Blocks: true,
	},
	{
		Name:     "Archer",
		HeroSpec: HeroSpec{MaxHP: 80, Speed: 4.5, Damage: 12, Range: 3.5, AttackCD: 0.5, RespawnS: 12.0},
		LevelCap: 5, HPPerLvl: 15, DmgPerLvl: 4, RangePerLvl: 0.15, XPBase: 6,
		AbilityLevel: 3, AbilityCooldownS: 8.0, AbilityRadius: 3.0, AbilityDmgMul: 2,
		Blocks: false, // 决策 B: 远程不肉搏, 不守隘口 — 换高射程输出
	},
	{
		Name:     "Rogue",
		HeroSpec: HeroSpec{MaxHP: 90, Speed: 5.5, Damage: 9, Range: 1.5, AttackCD: 0.35, RespawnS: 10.0},
		LevelCap: 5, HPPerLvl: 18, DmgPerLvl: 3, RangePerLvl: 0.1, XPBase: 6,
		AbilityLevel: 3, AbilityCooldownS: 6.0, AbilityRadius: 1.5, AbilityDmgMul: 4,
		Blocks: true,
	},
}

// knight: 默认职业别名 (V8/V9 原数值, 测试基线)。
var knight = &heroClasses[0]

type Hero struct {
	Class     *HeroClass // V10: 职业参数 (数值/成长/技能/阻挡)
	X, Y      float64    // 自由坐标 (cell, 非路径绑定)
	RallyX    float64    // 集结点 (玩家 H 键设定)
	RallyY    float64
	HP        int
	MaxHP     int
	cooldown  float64 // P2: 攻击冷却剩余 (>0 不能出手)
	respawnCD float64 // P2: >0 = 阵亡复活倒计时 (此期间不在场)
	// V9 P1: 关内成长 (per-run, beginRun 重置为 1 级 0 XP)
	Level int
	XP    int
	// V9 P2: 主动技能 (AoE 横扫) 冷却剩余 (>0 不可释放)
	abilityCD float64
}

// AbilityUnlocked: 是否已达解锁等级。
func (h *Hero) AbilityUnlocked() bool { return h.Level >= h.Class.AbilityLevel }

// AbilityReady: 在场存活 + 已解锁 + 不在冷却。
func (h *Hero) AbilityReady() bool {
	return h.Alive() && h.AbilityUnlocked() && h.abilityCD <= 0
}

// Damage: 当前等级攻击伤害。
func (h *Hero) Damage() int { return h.Class.Damage + (h.Level-1)*h.Class.DmgPerLvl }

// AttackRange: 当前等级攻击射程。
func (h *Hero) AttackRange() float64 {
	return h.Class.Range + float64(h.Level-1)*h.Class.RangePerLvl
}

// GainXP: 累计 XP 并按阈值升级 (升级回满血, 满级后 XP 锁 0)。返回升级后的等级
// (调用方据此判断是否升级以给反馈)。
func (h *Hero) GainXP(xp int) {
	if h.Level >= h.Class.LevelCap {
		return
	}
	h.XP += xp
	for h.Level < h.Class.LevelCap && h.XP >= h.Class.xpForNext(h.Level) {
		h.XP -= h.Class.xpForNext(h.Level)
		h.Level++
		h.MaxHP = h.Class.maxHPFor(h.Level)
		h.HP = h.MaxHP // 升级回满血 (KR 范式)
	}
	if h.Level >= h.Class.LevelCap {
		h.XP = 0
	}
}

// DistTo: 英雄到格 (px,py) 的欧氏距离 (cell)。
func (h *Hero) DistTo(px, py float64) float64 {
	return math.Hypot(px-h.X, py-h.Y)
}

// newHero: 在 spawn 点生成满血 1 级英雄, 集结点 = spawn (原地待命)。
// V9: Level/XP 每关从 newHero 起 (per-run, beginRun 调用 = 关内重置)。
// V10 P1: 默认 Knight; P3 起由 Save.HeroChoice 经 newHeroOf 选职业。
func newHero(spawnX, spawnY float64) *Hero {
	return newHeroOf(knight, spawnX, spawnY)
}

// newHeroOf: 指定职业生成英雄。
func newHeroOf(c *HeroClass, spawnX, spawnY float64) *Hero {
	return &Hero{
		Class: c,
		X:     spawnX, Y: spawnY, RallyX: spawnX, RallyY: spawnY,
		HP: c.maxHPFor(1), MaxHP: c.maxHPFor(1),
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
	h.X, h.Y, _ = stepToward(h.X, h.Y, h.RallyX, h.RallyY, h.Class.Speed*dt)
}
