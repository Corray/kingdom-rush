//go:build !js

// 存档 IO (native build): JSON 持久化到 ~/.kingdom-rush/save.json。
//
// Save struct 与纯逻辑方法在 save_core.go (V7.4 A1 合并, shared)。
// WASM IO 见 save_wasm.go (localStorage)。
//
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
