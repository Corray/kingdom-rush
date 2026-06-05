//go:build !term

// V4 Phase 2: BGM 播放层 (ebiten/audio infinite loop)。
//
// 两轨道: menu (Stage Select) / battle (Stage 1), Juhani Junkala CC0,
// 见 assets/bgm/LICENSE-JuhaniJunkala.txt。
//
// 与 SFX (audio_player.go) 的区别: BGM 流式解码 + InfiniteLoop 常驻
// Player, 不预解码全量 PCM (两轨 ~2.6MB ogg 解码后会膨胀 ~10×)。
//
// 切换状态机 (淡出 → 换轨 → 淡入, 防爆音):
//
//	每帧 updateBGM(phase, masterVol, dt):
//	  期望轨 != 当前轨 → fade 递减到 0 → Pause 当前 + Rewind/Play 新轨
//	  期望轨 == 当前轨 → fade 递增到 1
//	实际音量 = masterVol × fade × bgmBaseVol
//
// phase 映射: LevelSelect → menu / Playing → battle / Won·Lost → 无
// (停 BGM 让 win/lose jingle 独奏)。
package main

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type bgmTrack int

const (
	bgmNone bgmTrack = iota
	bgmMenu
	bgmBattle
)

const (
	bgmFadeSpeed = 2.0 // fade 系数变化速率/秒 (全程 ~0.5s)
	bgmBaseVol   = 0.4 // BGM 相对 SFX 压低, 不盖音效
)

var (
	bgmPlayers map[bgmTrack]*audio.Player
	bgmCurrent = bgmNone
	bgmFade    = 0.0
)

// bgmFor: 游戏 phase → 期望 BGM 轨道 (纯函数, 可测)。
func bgmFor(p GamePhase) bgmTrack {
	switch p {
	case PhaseLevelSelect:
		return bgmMenu
	case PhasePlaying:
		return bgmBattle
	default:
		return bgmNone
	}
}

// V7.4 B1: BGM 数据异步到达通道 — 资产提供层 (embed / fetch) 与
// 解码注册解耦。fetch goroutine 只投递 bytes, 解码与 map 写入
// 一律在主线程 updateBGM poll (无锁无 data race)。
type bgmArrival struct {
	tr   bgmTrack
	data []byte
}

var bgmDataCh = make(chan bgmArrival, 2)

// queueBGMData: 资产提供层投递入口 (native 同步 / wasm fetch goroutine)。
func queueBGMData(tr bgmTrack, data []byte) {
	select {
	case bgmDataCh <- bgmArrival{tr: tr, data: data}:
	default: // 容量 2 (轨道数), 重复投递丢弃
	}
}

// initBGM: 建空 player 表 + 触发资产加载 (native 立即投递 /
// wasm 后台 fetch)。返回非 0 表示 audio 不可用。
func initBGM() int {
	if audioCtx == nil {
		return 2
	}
	bgmPlayers = make(map[bgmTrack]*audio.Player, 2)
	loadBGMAssets()
	return 0
}

// registerArrivedBGM: 主线程消费到达的轨道数据 — 解码 + 建 player;
// 若该轨正是当前期望轨 (迟到场景), 注册即起播。
func registerArrivedBGM(m bgmArrival) {
	s, err := vorbis.DecodeWithSampleRate(audioSampleRate, bytes.NewReader(m.data))
	if err != nil {
		return // 静默: BGM 缺失不阻断游戏
	}
	loop := audio.NewInfiniteLoop(s, s.Length())
	p, err := audioCtx.NewPlayer(loop)
	if err != nil {
		return
	}
	bgmPlayers[m.tr] = p
	if m.tr == bgmCurrent {
		_ = p.Rewind()
		p.Play()
	}
}

// updateBGM: 每帧驱动切换状态机 + 应用音量。
func updateBGM(phase GamePhase, masterVol, dt float64) {
	if bgmPlayers == nil {
		return
	}
	// V7.4 B1: poll 异步到达的轨道 (每帧最多一条, 足够)
	select {
	case m := <-bgmDataCh:
		registerArrivedBGM(m)
	default:
	}
	target := bgmFor(phase)
	if bgmCurrent != target {
		// 淡出当前轨, 归零后切轨
		bgmFade -= bgmFadeSpeed * dt
		if bgmFade <= 0 {
			bgmFade = 0
			if p := bgmPlayers[bgmCurrent]; p != nil {
				p.Pause()
			}
			bgmCurrent = target
			if p := bgmPlayers[bgmCurrent]; p != nil {
				_ = p.Rewind() // 失败仅影响起播位置, 不阻断
				p.Play()
			}
		}
	} else if bgmFade < 1 {
		bgmFade += bgmFadeSpeed * dt
		if bgmFade > 1 {
			bgmFade = 1
		}
	}
	if p := bgmPlayers[bgmCurrent]; p != nil {
		p.SetVolume(masterVol * bgmFade * bgmBaseVol)
	}
}
