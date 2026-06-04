// V6 Phase 3: Endless mode — 预算制 wave 生成器 (纯逻辑, 可测)。
//
// 确定性约束: 项目首次引入随机性, rng 一律注入 (rand.New(seed)),
// 不用全局 rand — 同 seed 同序列, 测试可锁。
//
// 预算模型: wave n 预算 = base + n×inc, 按敌型 cost 选购直到花完
// (cost 超预算时回退 Normal, 预算恰好归零)。敌型池随 wave 解锁,
// 超过 10 波后额外 HP 缩放 (endlessHPScale)。
package main

import "math/rand"

const (
	endlessBaseBudget = 20
	endlessBudgetInc  = 6
	endlessStartGold  = 180
	endlessStartLives = 8
)

// enemyCost: 预算权重 (大致按威胁度)。
var enemyCost = map[EnemyKind]int{
	ENormal:  1,
	EFast:    2,
	EGlider:  3,
	ESpawner: 4,
	EBoss:    10,
}

// endlessCPs: 专用 path — 全图蛇形 (长路线适合持久战, 同 L14 风格)。
var endlessCPs = [][]int{
	{0, 1}, {28, 1}, {28, 4}, {1, 4}, {1, 7}, {28, 7},
	{28, 10}, {1, 10}, {1, 13}, {29, 13},
}

// endlessPool: wave n (1-based) 可用敌型 (逐步解锁)。
func endlessPool(n int) []EnemyKind {
	pool := []EnemyKind{ENormal, EFast}
	if n >= 3 {
		pool = append(pool, EGlider)
	}
	if n >= 5 {
		pool = append(pool, ESpawner)
	}
	if n >= 8 {
		pool = append(pool, EBoss)
	}
	return pool
}

// genEndlessWave: wave n 的敌人序列。预算花到恰好归零 → 总威胁度
// 严格等于预算公式, 强度单调性由公式保证。
func genEndlessWave(n int, rng *rand.Rand) []EnemyKind {
	budget := endlessBaseBudget + n*endlessBudgetInc
	pool := endlessPool(n)
	var out []EnemyKind
	for budget > 0 {
		k := pool[rng.Intn(len(pool))]
		c := enemyCost[k]
		if c > budget {
			k, c = ENormal, enemyCost[ENormal]
		}
		out = append(out, k)
		budget -= c
	}
	return out
}

// endlessHPScale: 超过 10 波后敌人 HP 额外缩放 (+5%/波)。
func endlessHPScale(n int) float64 {
	if n <= 10 {
		return 1.0
	}
	return 1.0 + 0.05*float64(n-10)
}
