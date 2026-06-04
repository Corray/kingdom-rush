// V4 Phase 1: SFX 事件队列 (纯数据, 无 ebiten 依赖)。
//
// 架构: game.go 在状态变化处 pushSound, 渲染层每帧 DrainSounds 消费 —
//   - ebiten 侧: audio_player.go 解码 + 播放
//   - term 侧:   drain 后丢弃 (terminal build 无音频)
//
// game 逻辑层不 import 任何音频库, `-tags term` build 不受影响;
// 测试只测触发 (事件入队), 不测播放。
package main

type SoundEvent int

const (
	SndShootArcher SoundEvent = iota
	SndShootCannon
	SndShootMagic
	SndShootFrost // V5 Phase 4
	SndEnemyDeath
	SndBuild
	SndUpgrade
	SndSell // V5 Phase 1
	SndWaveStart
	SndWin
	SndLose
)

// maxSoundQueue: 消费者缺位 (如 term 侧漏 drain) 时队列不无限增长。
const maxSoundQueue = 64

func (g *Game) pushSound(s SoundEvent) {
	if len(g.SoundEvents) >= maxSoundQueue {
		return // 满则丢弃, 音频丢失不影响游戏正确性
	}
	g.SoundEvents = append(g.SoundEvents, s)
}

// DrainSounds 返回积压事件并清空队列 (每帧调用)。
func (g *Game) DrainSounds() []SoundEvent {
	evs := g.SoundEvents
	g.SoundEvents = nil
	return evs
}

// shootSound: 塔型 → 对应射击音效。
func shootSound(k TowerKind) SoundEvent {
	switch k {
	case TCannon:
		return SndShootCannon
	case TMagic:
		return SndShootMagic
	case TFrost:
		return SndShootFrost
	default:
		return SndShootArcher
	}
}
