//go:build !term

// V4 Phase 1: SFX 播放层 (ebiten/audio)。
//
// assets/sfx/*.ogg (Kenney CC0, 见目录内 LICENSE-Kenney-audio.txt)
// embed 进 binary, initAudio 时一次性解码为 PCM bytes;
// 播放用 NewPlayerFromBytes — 同音效可叠播, fire-and-forget。
//
// 降级策略 (对齐 tilesheet / font 先例): 任何解码失败静默跳过该音效,
// 音频问题不阻断游戏。
//
// WASM 注意: 浏览器 autoplay 政策下, AudioContext 在首次用户手势
// (键盘/鼠标) 后才解锁; ebiten/audio 内部处理 resume, 此前的 Play
// 无声但不报错。
package main

import (
	"bytes"
	_ "embed"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const audioSampleRate = 44100

//go:embed assets/sfx/shoot_archer.ogg
var oggShootArcher []byte

//go:embed assets/sfx/shoot_cannon.ogg
var oggShootCannon []byte

//go:embed assets/sfx/shoot_magic.ogg
var oggShootMagic []byte

//go:embed assets/sfx/enemy_death.ogg
var oggEnemyDeath []byte

//go:embed assets/sfx/build.ogg
var oggBuild []byte

//go:embed assets/sfx/upgrade.ogg
var oggUpgrade []byte

//go:embed assets/sfx/sell.ogg
var oggSell []byte

//go:embed assets/sfx/wave_start.ogg
var oggWaveStart []byte

//go:embed assets/sfx/win.ogg
var oggWin []byte

//go:embed assets/sfx/lose.ogg
var oggLose []byte

var (
	audioCtx *audio.Context
	sfxPCM   map[SoundEvent][]byte
)

// initAudio: context 单例 + 全量 ogg → PCM 预解码。返回解码失败数 (0 = 全成功)。
func initAudio() int {
	audioCtx = audio.NewContext(audioSampleRate)
	raw := map[SoundEvent][]byte{
		SndShootArcher: oggShootArcher,
		SndShootCannon: oggShootCannon,
		SndShootMagic:  oggShootMagic,
		SndEnemyDeath:  oggEnemyDeath,
		SndBuild:       oggBuild,
		SndUpgrade:     oggUpgrade,
		SndSell:        oggSell,
		SndWaveStart:   oggWaveStart,
		SndWin:         oggWin,
		SndLose:        oggLose,
	}
	sfxPCM = make(map[SoundEvent][]byte, len(raw))
	failed := 0
	for ev, b := range raw {
		s, err := vorbis.DecodeWithSampleRate(audioSampleRate, bytes.NewReader(b))
		if err != nil {
			failed++
			continue
		}
		pcm, err := io.ReadAll(s)
		if err != nil {
			failed++
			continue
		}
		sfxPCM[ev] = pcm
	}
	return failed
}

// playSounds: drain 出的事件按主音量逐个播放 (未初始化 / 缺音效静默跳过)。
func playSounds(evs []SoundEvent, masterVol float64) {
	if audioCtx == nil || masterVol <= 0 {
		return
	}
	for _, ev := range evs {
		if pcm, ok := sfxPCM[ev]; ok {
			p := audioCtx.NewPlayerFromBytes(pcm)
			p.SetVolume(masterVol)
			p.Play()
		}
	}
}
