//go:build !js

// 存档系统(native build):JSON 持久化到 ~/.kingdom-rush/save.json
//
// WASM 版见 save_wasm.go (localStorage)。
//
// 规则:
//   - Level 1 always unlocked
//   - Level N (N > 1) unlocked iff Level (N-1) completed
//   - Completed map 记 levelID -> true
//   - 通关时 MarkCompleted + StoreSave(原子写: tmp + rename)
//   - 文件不存在 → 返回 empty save, 不报错(首次启动)
//
// 测试钩子: savePathFn 可被 override 指向 t.TempDir()。
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Save struct {
	Completed map[int]bool `json:"completed"`
	// V4 Phase 2: 音量档 0-10。指针区分"未设置"(nil, 旧存档无此字段
	// → 默认档) 与"显式 0"(静音)。家族注意: save_wasm.go 有同构定义。
	Volume *int `json:"volume,omitempty"`
	// V4 Phase 5: 屏幕反馈特效 (shake/顿帧) 关闭开关。bool 零值 =
	// 默认开启, 旧存档无字段自然兼容, 不需指针。家族: save_wasm.go。
	JuiceOff bool `json:"juice_off,omitempty"`
}

const (
	defaultVolume = 7
	maxVolume     = 10
)

func NewSave() Save {
	return Save{Completed: map[int]bool{}}
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

// === IO ===

var savePathFn = defaultSavePath

func defaultSavePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".kingdom-rush", "save.json"), nil
}

func LoadSave() (Save, error) {
	path, err := savePathFn()
	if err != nil {
		return NewSave(), err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return NewSave(), nil
		}
		return NewSave(), err
	}
	var s Save
	if err := json.Unmarshal(data, &s); err != nil {
		return NewSave(), fmt.Errorf("parse save: %w", err)
	}
	if s.Completed == nil {
		s.Completed = map[int]bool{}
	}
	return s, nil
}

func StoreSave(s Save) error {
	path, err := savePathFn()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	// 原子写: tmp + rename, 避免崩溃留半个文件
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}
