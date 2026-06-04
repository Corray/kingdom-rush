// V5 Phase 2: 塔 targeting 策略 (纯逻辑, 无 ebiten 依赖, 可测)。
//
// 重构来源: game.go Update shoot 循环的内联目标选择 (V1 起写死
// "path 最前优先") 抽为 pickTarget 纯函数 — TargetFirst 行为与
// 原内联逻辑等价 (回归测试锁定), Last/Strong 为 V5 新增策略。
//
// Tower.Target 零值 = TargetFirst → 既有代码/测试自然兼容。
package main

import "math"

type TargetMode int

const (
	TargetFirst  TargetMode = iota // path 最前 (默认, 原行为)
	TargetLast                     // path 最后 (刚进场的)
	TargetStrong                   // 当前 HP 最高 (平手取最前)
)

const targetModeCount = 3

// Name: HUD / Msg 显示名。
func (m TargetMode) Name() string {
	switch m {
	case TargetLast:
		return "Last"
	case TargetStrong:
		return "Strong"
	default:
		return "First"
	}
}

// Next: T 键循环切换。
func (m TargetMode) Next() TargetMode {
	return (m + 1) % targetModeCount
}

// pickTarget: 从射程内可命中敌人中按塔的策略选目标, 无候选返回 nil。
// 过滤规则与原内联逻辑一致: 跳过 dead/escaped + 飞行单位需 HitsFlying。
func pickTarget(t *Tower, enemies []*Enemy, path []Point) *Enemy {
	lvl := t.Spec()
	spec := towerSpecs[t.Kind]
	var best *Enemy
	for _, e := range enemies {
		if e.Dead || e.Escaped {
			continue
		}
		if enemySpecs[e.Kind].Flying && !spec.HitsFlying {
			continue
		}
		ep := e.Pos(path)
		dx := float64(ep.X - t.Pos.X)
		dy := float64(ep.Y - t.Pos.Y)
		if math.Sqrt(dx*dx+dy*dy) > lvl.Range {
			continue
		}
		if best == nil || betterTarget(t.Target, e, best) {
			best = e
		}
	}
	return best
}

// betterTarget: e 是否比当前 best 更符合策略 m。
func betterTarget(m TargetMode, e, best *Enemy) bool {
	switch m {
	case TargetLast:
		return e.PathIdx < best.PathIdx
	case TargetStrong:
		if e.HP != best.HP {
			return e.HP > best.HP
		}
		return e.PathIdx > best.PathIdx // 平手取最前
	default: // TargetFirst
		return e.PathIdx > best.PathIdx
	}
}
