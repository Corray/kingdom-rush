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
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

//go:embed assets/bgm/menu.ogg
var oggBgmMenu []byte

//go:embed assets/bgm/battle.ogg
var oggBgmBattle []byte

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

// initBGM: 流式解码 + InfiniteLoop 包装为常驻 Player (暂停态)。
// 依赖 initAudio 先建好 audioCtx。返回失败轨道数。
func initBGM() int {
	if audioCtx == nil {
		return 2
	}
	raw := map[bgmTrack][]byte{
		bgmMenu:   oggBgmMenu,
		bgmBattle: oggBgmBattle,
	}
	bgmPlayers = make(map[bgmTrack]*audio.Player, len(raw))
	failed := 0
	for tr, b := range raw {
		s, err := vorbis.DecodeWithSampleRate(audioSampleRate, bytes.NewReader(b))
		if err != nil {
			failed++
			continue
		}
		loop := audio.NewInfiniteLoop(s, s.Length())
		p, err := audioCtx.NewPlayer(loop)
		if err != nil {
			failed++
			continue
		}
		bgmPlayers[tr] = p
	}
	return failed
}

// updateBGM: 每帧驱动切换状态机 + 应用音量。
func updateBGM(phase GamePhase, masterVol, dt float64) {
	if bgmPlayers == nil {
		return
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
