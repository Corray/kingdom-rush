// V6 Phase 2: 难度模式 (纯数据, 无 ebiten 依赖)。
//
// 系数表作用点 (全部经统一入口, 防散落):
//   - 敌 HP:      newEnemy (game.go — spawn + Spawner 召唤共用)
//   - 击杀赏金:   killEnemy (game.go)
//   - 起始 lives: StartLevel (game.go)
//
// DiffNormal = 0 (零值): 旧存档无 difficulty 字段 → Normal,
// 与 TargetFirst / JuiceOff 同模式, 既有测试自然回归。
package main

type Difficulty int

const (
	DiffNormal Difficulty = iota // 零值 = 现状行为 (系数全 1)
	DiffHard
	DiffEasy
)

const difficultyCount = 3

type DifficultySpec struct {
	Name       string
	HPMul      float64 // 敌 HP 系数
	RewardMul  float64 // 击杀赏金系数
	LivesBonus int     // 起始 lives 增减
}

var difficultySpecs = map[Difficulty]DifficultySpec{
	DiffNormal: {Name: "Normal", HPMul: 1.0, RewardMul: 1.0, LivesBonus: 0},
	DiffHard:   {Name: "Hard", HPMul: 1.4, RewardMul: 0.8, LivesBonus: -1}, // V7.5: -2→-1 救前期 (仿真: Lv5/8/11 wave1 失守)
	DiffEasy:   {Name: "Easy", HPMul: 0.7, RewardMul: 1.2, LivesBonus: 3},
}

// Spec: 系数查表, 非法值 (手改存档等) 回退 Normal — 防 HPMul 0 炸局。
func (d Difficulty) Spec() DifficultySpec {
	if s, ok := difficultySpecs[d]; ok {
		return s
	}
	return difficultySpecs[DiffNormal]
}

// Next: D 键循环 Normal → Hard → Easy →。
func (d Difficulty) Next() Difficulty {
	return (d + 1) % difficultyCount
}
