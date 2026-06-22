// 攻击视觉特效系统 (V2.6)。
//
// Game.Update 中塔击中敌人时 push 两个 effect:
//   - EShoot: 从塔位到敌人位置的射线,fade out 0.15s
//   - EHit:   敌人位置的命中闪光圆,fade out 0.3s
//
// Game.Update 每帧 decay 所有 effect 的 TTL,过期移除。上限 200 防爆。
// 仅 ebiten 渲染消费 (terminal mode 忽略,保持 V1.7 体验)。
package main

import "strconv"

const maxEffects = 200

type EffectKind int

const (
	EShoot EffectKind = iota
	EHit
	EDeath // V4 Phase 3: 敌人死亡动画 (sprite 放大 + fade out)
	EText  // V4 Phase 4: 飘字 (伤害数字 / 金币获取, 上飘 + fade out)
)

type Effect struct {
	Kind   EffectKind
	From   Point     // EShoot: 塔位; EHit: 敌人位(同 To)
	To     Point     // 敌人位 (shoot 终点 / hit 中心)
	Color  RGB       // shoot 颜色随塔型, hit 通常红; EText 文字色
	Tower  TowerKind // V3 Phase 3b: 决定 bullet sprite (仅 EShoot 有意义)
	Enemy  EnemyKind // V4 Phase 3: 死亡 enemy 的 sprite (仅 EDeath 有意义)
	FX, FY float64   // V4 Phase 3/4: 插值路径坐标 cell float (EDeath / EText — 对齐平滑移动位置)
	Text   string    // V4 Phase 4: 飘字内容 (仅 EText 有意义)
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
func makeShootEffect(from, to Point, color RGB, tower TowerKind) *Effect {
	return &Effect{Kind: EShoot, From: from, To: to, Color: color, Tower: tower, TTL: 0.4, MaxTTL: 0.4}
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

// makeDeathEffect: V4 Phase 3 — 敌人死亡动画 (渲染为 enemy sprite
// 放大 + fade out)。fx/fy 用插值路径坐标, 与走动渲染位置无缝衔接。
func makeDeathEffect(fx, fy float64, kind EnemyKind) *Effect {
	return &Effect{
		Kind:   EDeath,
		Enemy:  kind,
		FX:     fx,
		FY:     fy,
		TTL:    0.35,
		MaxTTL: 0.35,
	}
}

// makeDamageText: V4 Phase 4 — 命中伤害数字 (浅黄, 上飘 0.5s)。
func makeDamageText(fx, fy float64, dmg int) *Effect {
	return &Effect{
		Kind:   EText,
		Text:   strconv.Itoa(dmg),
		Color:  RGB{255, 235, 140},
		FX:     fx,
		FY:     fy,
		TTL:    0.5,
		MaxTTL: 0.5,
	}
}

// makeHealText: V13 — 治疗飘字 "+Nhp" (绿色)。
func makeHealText(fx, fy float64, amount int) *Effect {
	return &Effect{
		Kind:   EText,
		Text:   "+" + strconv.Itoa(amount) + "hp",
		Color:  RGB{80, 255, 80},
		FX:     fx,
		FY:     fy - 0.3,
		TTL:    0.6,
		MaxTTL: 0.6,
	}
}

// makeGoldText: V4 Phase 4 — 击杀赏金 "+Ng" (金色, 上飘 0.8s 比伤害字长)。
func makeGoldText(fx, fy float64, amount int) *Effect {
	return &Effect{
		Kind:   EText,
		Text:   "+" + strconv.Itoa(amount) + "g",
		Color:  RGB{255, 200, 40},
		FX:     fx,
		FY:     fy,
		TTL:    0.8,
		MaxTTL: 0.8,
	}
}
