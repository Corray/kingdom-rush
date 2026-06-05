//go:build js && wasm

// V7.4 B1: BGM 资产提供层 — WASM 运行时 fetch (不 embed)。
//
// 用 syscall/js 直调浏览器 fetch API — 不引 net/http (其在 wasm
// 拉进 TLS/HTTP2/crypto 整套栈, 实测 +8MB, 远超省下的 2.6MB)。
// 回调在 JS 事件循环执行, queueBGMData 走 channel 并发安全。
//
// js.Func 不 Release: 一次性回调共 4 个 (2 轨 × 2 级 then), 此前
// defer Release 自身触发 "call to released function" 时序错误 —
// 接受可忽略的固定泄漏换确定性。失败静默, BGM 缺失不阻断游戏。
package main

import "syscall/js"

func loadBGMAssets() {
	urls := map[bgmTrack]string{
		bgmMenu:   "bgm/menu.ogg",
		bgmBattle: "bgm/battle.ogg",
	}
	for tr, url := range urls {
		fetchBytes(url, func(data []byte) { queueBGMData(tr, data) })
	}
}

// fetchBytes: fetch(url) → arrayBuffer → []byte → cb。
func fetchBytes(url string, cb func([]byte)) {
	thenBuf := js.FuncOf(func(_ js.Value, args []js.Value) any {
		u8 := js.Global().Get("Uint8Array").New(args[0])
		data := make([]byte, u8.Length())
		js.CopyBytesToGo(data, u8)
		cb(data)
		return nil
	})
	thenResp := js.FuncOf(func(_ js.Value, args []js.Value) any {
		resp := args[0]
		if !resp.Get("ok").Bool() {
			return nil
		}
		resp.Call("arrayBuffer").Call("then", thenBuf)
		return nil
	})
	js.Global().Call("fetch", url).Call("then", thenResp)
}
