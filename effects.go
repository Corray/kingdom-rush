// 攻击视觉特效系统 (V2.6)。
//
// Game.Update 中塔击中敌人时 push 两个 effect:
//   - EShoot: 从塔位到敌人位置的射线,fade out 0.15s
//   - EHit:   敌人位置的命中闪光圆,fade out 0.3s
//
// Game.Update 每帧 decay 所有 effect 的 TTL,过期移除。上限 200 防爆。
// 仅 ebiten 渲染消费 (terminal mode 忽略,保持 V1.7 体验)。
package main

const maxEffects = 200

type EffectKind int

const (
	EShoot EffectKind = iota
	EHit
)

type Effect struct {
	Kind   EffectKind
	From   Point // EShoot: 塔位; EHit: 敌人位(同 To)
	To     Point // 敌人位 (shoot 终点 / hit 中心)
	Color  RGB   // shoot 颜色随塔型, hit 通常红
	TTL    float64
	MaxTTL float64
}

// Alpha 返回当前 fade-out 透明度 0..1 (1=刚生成, 0=过期)
func (e *Effect) Alpha() float64 {
	if e.MaxTTL <= 0 {
		return 0
	}
	a := e.TTL / e.MaxTTL
	if a < 0 {
		return 0
	}
	if a > 1 {
		return 1
	}
	return a
}

// decayEffects 减 TTL, 移除过期 effect, 强制上限
func decayEffects(effects []*Effect, dt float64) []*Effect {
	out := effects[:0]
	for _, e := range effects {
		e.TTL -= dt
		if e.TTL > 0 {
			out = append(out, e)
		}
	}
	// 防爆: 超 maxEffects 丢早期
	if len(out) > maxEffects {
		out = out[len(out)-maxEffects:]
	}
	return out
}

// V3 Phase 3: Shoot TTL 0.15 → 0.4 让 bullet 飞行可见
// EHit 也调到 0.4 与 shoot 同步 fade
func makeShootEffect(from, to Point, color RGB) *Effect {
	return &Effect{Kind: EShoot, From: from, To: to, Color: color, TTL: 0.4, MaxTTL: 0.4}
}

func makeHitEffect(at Point) *Effect {
	return &Effect{
		Kind:   EHit,
		From:   at,
		To:     at,
		Color:  RGB{255, 220, 80},
		TTL:    0.4,
		MaxTTL: 0.4,
	}
}
