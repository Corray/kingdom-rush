//go:build js && wasm

// 存档 IO (WASM build): localStorage 持久化。
//
// Save struct 与纯逻辑方法在 save_core.go (V7.4 A1 合并, shared) —
// 此前本文件复刻全部定义, 5 轮 family edit 后消除。
//
// key: "kingdom-rush-save" (保留历史名, 用户进度零迁移)
// 容量 ~5 MB (远超所需 < 1 KB)
package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

const localStorageKey = "kingdom-rush-save"

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
