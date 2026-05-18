//go:build js && wasm

// 存档系统(WASM build):localStorage 持久化。
//
// 浏览器没有 file system,用 localStorage 存 JSON。
// key: "kingdom-rush-save"
// 容量 ~5 MB (远超本游戏所需 < 1 KB)
//
// Save struct / IsCompleted / IsUnlocked / NewSave / MarkCompleted
// 复用 save.go 中的定义 (本文件不重复)。
//
// 注意: save.go 加了 build tag `!js`, 本文件 build tag `js && wasm`,
// 两者互斥。但 Save struct / NewSave / IsCompleted / IsUnlocked /
// MarkCompleted 这些纯逻辑定义在 save.go 中, 而 save.go 在 wasm
// build 中被排除 — 所以本文件需要复刻这些纯逻辑定义。
package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

const localStorageKey = "kingdom-rush-save"

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

func (s *Save) IsUnlocked(levelID int) bool {
	if levelID <= 1 {
		return true
	}
	return s.IsCompleted(levelID - 1)
}

// === localStorage IO ===

func LoadSave() (Save, error) {
	storage := js.Global().Get("localStorage")
	if !storage.Truthy() {
		// 极端情况:浏览器禁用 localStorage → 静默 empty save
		return NewSave(), nil
	}
	val := storage.Call("getItem", localStorageKey)
	if val.IsNull() || val.IsUndefined() {
		return NewSave(), nil
	}
	var s Save
	if err := json.Unmarshal([]byte(val.String()), &s); err != nil {
		return NewSave(), fmt.Errorf("parse save from localStorage: %w", err)
	}
	if s.Completed == nil {
		s.Completed = map[int]bool{}
	}
	return s, nil
}

func StoreSave(s Save) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	storage := js.Global().Get("localStorage")
	if !storage.Truthy() {
		return fmt.Errorf("localStorage unavailable")
	}
	storage.Call("setItem", localStorageKey, string(data))
	return nil
}
