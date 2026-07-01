// Save 数据结构与纯逻辑 (V7.4 A1 — 从 save.go/save_wasm.go 双份
// 重复定义合并而来)。
//
// 背景: 双实现各自定义 Save struct + 7 个方法, V4-V6 间经历 5 轮
// "改一处必须同步另一处"的 family edit (volume/juice_off/difficulty
// /best_wave/stars), 结构性风险。本文件无 build tag — 两个 build
// 共享唯一定义, IO 层各留各 (save.go 文件 / save_wasm.go localStorage)。
//
// 规则:
//   - Level 1 always unlocked
//   - Level N (N > 1) unlocked iff Level (N-1) completed
//   - 偏好字段全部零值兼容 (旧存档无字段 → 默认行为)
package main

type Save struct {
	Completed map[int]bool `json:"completed"`
	// V4 Phase 2: 音量档 0-10。指针区分"未设置"(nil → 默认档)
	// 与"显式 0"(静音)。
	Volume *int `json:"volume,omitempty"`
	// V4 Phase 5: 屏幕反馈特效关闭开关 (V7.3 起含飘字)。零值 = 开。
	JuiceOff bool `json:"juice_off,omitempty"`
	// V6 Phase 2: 难度偏好。零值 = DiffNormal; 非法值 Spec() 回退。
	Difficulty Difficulty `json:"difficulty,omitempty"`
	// V6 Phase 3: endless 最佳纪录 (已清 wave 数, 取 max)。
	BestWave int `json:"best_wave,omitempty"`
	// V6 Phase 4: per-level 最高星 (levelID → 1-3)。nil map 读安全。
	Stars map[int]int `json:"stars,omitempty"`
	// V10 P3: 英雄职业选择 (heroClasses index)。零值 0 = Knight
	// (旧存档无字段 → 默认行为不变); 越界经 HeroClassIdx 回退 0。
	HeroChoice int `json:"hero_choice,omitempty"`
	// V11 P1: 技能树已购节点数 (职业名 → 0-4, 线性树)。nil map 读安全
	// (旧存档 → 全 0 = 无 perk, V10 行为零回归)。预算/购买逻辑见 skilltree.go。
	TreeNodes map[string]int `json:"tree_nodes,omitempty"`
	// V15: 成就系统
	Achievements   map[string]bool `json:"achievements,omitempty"`
	TotalKills     int             `json:"total_kills,omitempty"`
	BossKillsTotal int             `json:"boss_kills_total,omitempty"`
	WonClasses     map[int]bool    `json:"won_classes,omitempty"`
	HasFlawless    bool            `json:"has_flawless,omitempty"`
	HasHardClear   bool            `json:"has_hard_clear,omitempty"`
	HardCompleted  map[int]bool    `json:"hard_completed,omitempty"`
}

const (
	defaultVolume = 7
	maxVolume     = 10
)

func NewSave() Save {
	return Save{Completed: map[int]bool{}}
}

func (s *Save) MarkCompleted(levelID int) {
	if s.Completed == nil {
		s.Completed = map[int]bool{}
	}
	s.Completed[levelID] = true
}

func (s *Save) IsCompleted(levelID int) bool {
	return s.Completed[levelID]
}

// IsUnlocked: Level 1 总是 unlocked; Level N > 1 unlocked iff Level N-1 completed
func (s *Save) IsUnlocked(levelID int) bool {
	if levelID <= 1 {
		return true
	}
	return s.IsCompleted(levelID - 1)
}

// VolumeLevel: 当前音量档 0-10, 未设置返回默认档。
func (s *Save) VolumeLevel() int {
	if s.Volume == nil {
		return defaultVolume
	}
	return *s.Volume
}

// SetVolumeLevel: 设音量档, clamp 到 [0, maxVolume]。
func (s *Save) SetVolumeLevel(v int) {
	if v < 0 {
		v = 0
	}
	if v > maxVolume {
		v = maxVolume
	}
	s.Volume = &v
}

// RecordStars: V6 Phase 4 — 记录通关星级, 取 max 不降级。
func (s *Save) RecordStars(levelID, stars int) {
	if s.Stars == nil {
		s.Stars = map[int]int{}
	}
	if stars > s.Stars[levelID] {
		s.Stars[levelID] = stars
	}
}

// StarsFor: 该关最高星 (未通关 / 旧存档 → 0)。
func (s *Save) StarsFor(levelID int) int {
	return s.Stars[levelID]
}

// HeroClassIdx: V10 P3 — 合法化职业 index (负值/越界回退 0 = Knight,
// 防手改存档/未来减职业)。
func (s *Save) HeroClassIdx() int {
	if s.HeroChoice < 0 || s.HeroChoice >= len(heroClasses) {
		return 0
	}
	return s.HeroChoice
}
