// Tower / Enemy 实体与 spec table。
// Tower 支持 3 级升级,每级独立 spec(damage/range/cd/字符)。
//
// Spec.Color 用自定义 RGB 类型(color.go),renderer 各自转换,
// 此文件不依赖任何 UI 库。
package main

import (
	"math"
)

// ============================================================
// Point / Wave (基础类型)
// ============================================================

type Point struct{ X, Y int }

type WaveSpec struct {
	Enemies []EnemyKind
}

// ============================================================
// Tower
// ============================================================

type TowerKind int

const (
	TArcher TowerKind = iota
	TCannon
	TMagic
	TFrost  // V5 Phase 4: 低伤 + 命中减速
	TTesla  // V16: 链式闪电 — 弹射 2 次 50% 伤害
	TSniper // V16: 狙击塔 — 超远超高伤超慢
)

type TowerLevel struct {
	Cost     int // 升到该级的额外 cost (lvl 1 = 建造 cost)
	Damage   int
	Range    float64
	Cooldown float64
	Splash   float64 // V5 Phase 3: 溅射半径 (cell), 0 = 单体 (仅 Cannon 非零)
	Slow     float64 // V5 Phase 4: 命中减速系数 (速度乘数, 0 = 无; 仅 Frost 非零)
	Chain    int     // V16: 链式弹射次数 (0 = 无; TTesla 用)
	Char1    rune
	Char2    rune
}

type TowerSpec struct {
	Name       string
	Color      RGB
	HitsFlying bool // 能否打飞行单位
	Levels     [3]TowerLevel
	BranchB    *[2]TowerLevel // V16: 分支 B 的 L2+L3 (nil = 无分支)
	BranchName [2]string      // V16: [0]=A名, [1]=B名 (空 = 无分支)
}

var towerSpecs = map[TowerKind]TowerSpec{
	TArcher: {
		Name: "Archer", Color: rgb(100, 180, 255), HitsFlying: true,
		Levels: [3]TowerLevel{
			{Cost: 50, Damage: 8, Range: 3.5, Cooldown: 0.6, Char1: 'A', Char2: '1'},
			{Cost: 40, Damage: 14, Range: 4.0, Cooldown: 0.5, Char1: 'A', Char2: '2'},
			{Cost: 80, Damage: 22, Range: 4.5, Cooldown: 0.4, Char1: 'A', Char2: '3'},
		},
		BranchName: [2]string{"Marksman", "Rapidfire"},
		BranchB: &[2]TowerLevel{
			{Cost: 40, Damage: 10, Range: 3.5, Cooldown: 0.3, Char1: 'A', Char2: 'r'},
			{Cost: 80, Damage: 15, Range: 4.0, Cooldown: 0.2, Char1: 'A', Char2: 'R'},
		},
	},
	TCannon: {
		Name: "Cannon", Color: rgb(180, 120, 255), HitsFlying: false,
		Levels: [3]TowerLevel{
			{Cost: 80, Damage: 25, Range: 2.5, Cooldown: 1.5, Splash: 1.0, Char1: 'C', Char2: '1'},
			{Cost: 60, Damage: 45, Range: 3.0, Cooldown: 1.3, Splash: 1.2, Char1: 'C', Char2: '2'},
			{Cost: 100, Damage: 70, Range: 3.5, Cooldown: 1.0, Splash: 1.5, Char1: 'C', Char2: '3'},
		},
		BranchName: [2]string{"Mortar", "Gatling"},
		BranchB: &[2]TowerLevel{
			{Cost: 60, Damage: 18, Range: 3.0, Cooldown: 0.5, Char1: 'C', Char2: 'g'},
			{Cost: 100, Damage: 28, Range: 3.5, Cooldown: 0.35, Char1: 'C', Char2: 'G'},
		},
	},
	TMagic: {
		Name: "Magic", Color: rgb(220, 120, 255), HitsFlying: true,
		Levels: [3]TowerLevel{
			{Cost: 100, Damage: 18, Range: 3.0, Cooldown: 0.8, Char1: 'M', Char2: '1'},
			{Cost: 80, Damage: 30, Range: 3.5, Cooldown: 0.7, Char1: 'M', Char2: '2'},
			{Cost: 120, Damage: 50, Range: 4.0, Cooldown: 0.6, Char1: 'M', Char2: '3'},
		},
		BranchName: [2]string{"Archmage", "Enchanter"},
		BranchB: &[2]TowerLevel{
			{Cost: 80, Damage: 20, Range: 3.5, Cooldown: 0.8, Slow: 0.7, Char1: 'M', Char2: 'e'},
			{Cost: 120, Damage: 30, Range: 4.0, Cooldown: 0.7, Slow: 0.5, Char1: 'M', Char2: 'E'},
		},
	},
	// V5 Phase 4: 支援塔 — 伤害低, 价值在命中减速 (Slow = 速度乘数)
	TFrost: {
		Name: "Frost", Color: rgb(140, 200, 255), HitsFlying: true,
		Levels: [3]TowerLevel{
			{Cost: 70, Damage: 4, Range: 2.5, Cooldown: 0.8, Slow: 0.6, Char1: 'F', Char2: '1'},
			{Cost: 50, Damage: 8, Range: 3.0, Cooldown: 0.7, Slow: 0.5, Char1: 'F', Char2: '2'},
			{Cost: 90, Damage: 12, Range: 3.5, Cooldown: 0.6, Slow: 0.4, Char1: 'F', Char2: '3'},
		},
		BranchName: [2]string{"Deep Freeze", "Hailstorm"},
		BranchB: &[2]TowerLevel{
			{Cost: 50, Damage: 16, Range: 3.0, Cooldown: 0.6, Slow: 0.8, Char1: 'F', Char2: 'h'},
			{Cost: 90, Damage: 28, Range: 3.5, Cooldown: 0.5, Slow: 0.85, Char1: 'F', Char2: 'H'},
		},
	},
	TTesla: {
		Name: "Tesla", Color: rgb(180, 220, 255), HitsFlying: true,
		Levels: [3]TowerLevel{
			{Cost: 90, Damage: 12, Range: 3.0, Cooldown: 1.0, Chain: 2, Char1: 'T', Char2: '1'},
			{Cost: 70, Damage: 20, Range: 3.5, Cooldown: 0.9, Chain: 2, Char1: 'T', Char2: '2'},
			{Cost: 110, Damage: 32, Range: 4.0, Cooldown: 0.8, Chain: 3, Char1: 'T', Char2: '3'},
		},
	},
	// V16 收尾校准: cost 120→100 / cd 3.0→2.6 系 (仿真 sniper-mix 16/20
	// 坍塌 — 单发溢出+铺塔慢的机会成本过高; 伤害不动保持反 Boss 定位)
	TSniper: {
		Name: "Sniper", Color: rgb(200, 160, 100), HitsFlying: false,
		Levels: [3]TowerLevel{
			{Cost: 100, Damage: 40, Range: 6.0, Cooldown: 2.6, Char1: 'S', Char2: '1'},
			{Cost: 90, Damage: 70, Range: 7.0, Cooldown: 2.2, Char1: 'S', Char2: '2'},
			{Cost: 140, Damage: 120, Range: 8.0, Cooldown: 1.8, Char1: 'S', Char2: '3'},
		},
	},
}

func TowerKinds() []TowerKind { return []TowerKind{TArcher, TCannon, TMagic, TFrost, TTesla, TSniper} }

type Tower struct {
	Pos      Point
	Kind     TowerKind
	Level    int        // 1-3
	Branch   int        // V16: 0=分支A(默认), 1=分支B
	Target   TargetMode // V5 Phase 2: targeting 策略 (零值 First = 原行为)
	cooldown float64
}

func (t *Tower) Spec() TowerLevel {
	spec := towerSpecs[t.Kind]
	if t.Branch == 1 && spec.BranchB != nil && t.Level >= 2 {
		return spec.BranchB[t.Level-2]
	}
	return spec.Levels[t.Level-1]
}

func (t *Tower) NextUpgradeCost() (int, bool) {
	if t.Level >= 3 {
		return 0, false
	}
	spec := towerSpecs[t.Kind]
	if t.Branch == 1 && spec.BranchB != nil && t.Level >= 1 {
		return spec.BranchB[t.Level-1].Cost, true
	}
	return spec.Levels[t.Level].Cost, true
}

// towerInvested: V5 Phase 1 — 建造 + 已升级的累计投入 (卖塔退款基数)。
// Cost 字段是逐级增量 (lvl1 = 建造 cost), 累加 [0, level)。
func towerInvested(kind TowerKind, level int) int {
	levels := towerSpecs[kind].Levels
	total := 0
	for i := 0; i < level && i < len(levels); i++ {
		total += levels[i].Cost
	}
	return total
}

// ============================================================
// Enemy
// ============================================================

type EnemyKind int

const (
	ENormal  EnemyKind = iota
	EFast
	EGlider  // 飞行单位:Cannon 打不到
	EBoss    // 大型敌人:HP 高 / 速度慢 / 奖励高
	ESpawner // 召唤者:死时 spawn 2 个 ENormal at same PathIdx (V3.6)
	EShield  // V13: 护甲兵 — 每次受击减免 Armor 伤害 (min 1)
	ERegen   // V13: 回血兵 — 持续回复 HP
	EHealer  // V13: 治疗者 — 定期治疗附近盟友
)

type EnemySpec struct {
	HP     int
	Speed  float64
	Reward int
	Flying bool // true = 飞行(Cannon 不能 target)
	Attack int  // V8 P2: 近战攻击英雄的伤害 (飞行 = 0, 不近战; P3 阻挡时按 meleeCD 节奏出手)
	Armor  int  // V13: 每次受击减免伤害 (min 1); 0 = 无护甲
	Regen  int  // V13: 每秒回复 HP; 0 = 无回血
	HealCD float64 // V13: 治疗间隔秒 (0 = 非治疗者)
	Char1  rune
	Char2  rune
	Color  RGB
}

var enemySpecs = map[EnemyKind]EnemySpec{
	ENormal: {
		HP: 20, Speed: 3.0, Reward: 10, Flying: false, Attack: 6,
		Char1: 'o', Char2: 'o', Color: rgb(255, 100, 100),
	},
	EFast: {
		HP: 12, Speed: 5.5, Reward: 12, Flying: false, Attack: 4,
		Char1: '>', Char2: '>', Color: rgb(255, 200, 80),
	},
	EGlider: {
		HP: 18, Speed: 4.0, Reward: 18, Flying: true, Attack: 0,
		Char1: '~', Char2: '~', Color: rgb(100, 220, 220),
	},
	EBoss: {
		HP: 150, Speed: 1.5, Reward: 50, Flying: false, Attack: 25,
		Char1: 'B', Char2: 'B', Color: rgb(255, 60, 200),
	},
	ESpawner: {
		HP: 35, Speed: 2.5, Reward: 25, Flying: false, Attack: 8,
		Char1: 'S', Char2: 'p', Color: rgb(120, 220, 100),
	},
	EShield: {
		HP: 30, Speed: 2.0, Reward: 20, Flying: false, Attack: 10, Armor: 4,
		Char1: '[', Char2: ']', Color: rgb(160, 160, 180),
	},
	ERegen: {
		HP: 40, Speed: 2.5, Reward: 22, Flying: false, Attack: 6, Regen: 3,
		Char1: 'R', Char2: 'g', Color: rgb(80, 255, 80),
	},
	EHealer: {
		HP: 20, Speed: 2.5, Reward: 30, Flying: false, Attack: 4, HealCD: 2.0,
		Char1: '+', Char2: '+', Color: rgb(255, 255, 100),
	},
}

type Enemy struct {
	PathIdx float64
	PathID  int // V12: 走哪条 path (索引 g.Paths; 单路恒 0)
	Kind    EnemyKind
	HP      int
	MaxHP   int
	Dead    bool
	Escaped bool
	// V5 Phase 4: 减速状态效果 (SlowTimer > 0 时速度 × SlowFactor)
	SlowFactor float64
	SlowTimer  float64
	// V8 P2: 近战英雄的出手冷却 (>0 时不能再攻击英雄); P3 阻挡用同字段
	meleeCD float64
	// V8 P3: 被英雄阻挡 (贴身停步互殴; 渲染/测试观察, 每帧由 move loop 刷新)
	Blocked bool
	// V13: 治疗者出手冷却
	healCD float64
	// V13: 回血计时器 (每秒累积)
	regenAcc float64
	// V14: Boss 行为计时器
	bossCD     float64
	bossShield float64 // >0 时无敌
	bossCharge float64 // >0 时冲锋 (速度 ×3)
}

// slowDurationS: V5 Phase 4 — 单次减速持续时间 (命中刷新)。
const slowDurationS = 1.5

// ApplySlow: 施加减速。不叠加取最强 (系数更小 = 更慢 = 更强),
// 任何命中都刷新持续时间。
func (e *Enemy) ApplySlow(factor float64) {
	if e.SlowTimer <= 0 || factor < e.SlowFactor {
		e.SlowFactor = factor
	}
	e.SlowTimer = slowDurationS
}

// EffectiveSpeed: 当前速度 (减速生效时打折)。
func (e *Enemy) EffectiveSpeed() float64 {
	speed := enemySpecs[e.Kind].Speed
	if e.bossCharge > 0 {
		speed *= 3.0
	}
	if e.SlowTimer > 0 {
		return speed * e.SlowFactor
	}
	return speed
}

func (e *Enemy) Pos(path []Point) Point {
	i := int(math.Floor(e.PathIdx))
	if i >= len(path) {
		return path[len(path)-1]
	}
	if i < 0 {
		return path[0]
	}
	return path[i]
}
