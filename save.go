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
}

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
