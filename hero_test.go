// V8 Phase 1 测试: 英雄实体 + rally 移动 (纯逻辑)。
package main

import (
	"math"
	"testing"
)

// TestStepToward: 移动一步的三种情形 — 部分移动 / 到达(超出夹断) / 零距离。
func TestStepToward(t *testing.T) {
	// 部分移动: 从原点朝 (10,0) 走 3 → (3,0), 未到达
	nx, ny, arrived := stepToward(0, 0, 10, 0, 3)
	if nx != 3 || ny != 0 || arrived {
		t.Errorf("partial move: got (%v,%v) arrived=%v, want (3,0) arrived=false", nx, ny, arrived)
	}
	// 到达 (超出夹断): 剩余距离 4, 步长 10 → 落目标点, 到达
	nx, ny, arrived = stepToward(0, 0, 0, 4, 10)
	if nx != 0 || ny != 4 || !arrived {
		t.Errorf("overshoot clamp: got (%v,%v) arrived=%v, want (0,4) arrived=true", nx, ny, arrived)
	}
	// 零距离: 已在目标点 → 原地, 到达
	nx, ny, arrived = stepToward(5, 5, 5, 5, 3)
	if nx != 5 || ny != 5 || !arrived {
		t.Errorf("zero dist: got (%v,%v) arrived=%v, want (5,5) arrived=true", nx, ny, arrived)
	}
	// 对角线: 朝 (3,4) [距离 5] 走 2.5 → 走一半 (1.5,2.0)
	nx, ny, arrived = stepToward(0, 0, 3, 4, 2.5)
	if math.Abs(nx-1.5) > 1e-9 || math.Abs(ny-2.0) > 1e-9 || arrived {
		t.Errorf("diagonal: got (%v,%v) arrived=%v, want (1.5,2.0) arrived=false", nx, ny, arrived)
	}
}

// TestNewHero: 生成满血英雄, 集结点 = spawn (原地待命)。
func TestNewHero(t *testing.T) {
	h := newHero(8, 6)
	if h.X != 8 || h.Y != 6 {
		t.Errorf("spawn pos: got (%v,%v), want (8,6)", h.X, h.Y)
	}
	if h.RallyX != 8 || h.RallyY != 6 {
		t.Errorf("rally should default to spawn: got (%v,%v)", h.RallyX, h.RallyY)
	}
	if h.HP != knight.MaxHP || h.MaxHP != knight.MaxHP {
		t.Errorf("HP: got %d/%d, want %d full", h.HP, h.MaxHP, knight.MaxHP)
	}
	if !h.Alive() {
		t.Error("new hero should be alive")
	}
}

// TestHeroMoveTowardRally: 设集结点后逐帧移动, 最终到达 (确定性收敛)。
func TestHeroMoveTowardRally(t *testing.T) {
	h := newHero(0, 0)
	h.SetRally(10, 0)
	// 单帧 (dt=0.1): 走 Speed*0.1 = 0.4 cell
	h.moveStep(0.1)
	want := knight.Speed * 0.1
	if math.Abs(h.X-want) > 1e-9 {
		t.Errorf("after 1 frame: X=%v, want %v", h.X, want)
	}
	// 多帧推进直到到达 (上限防死循环)
	for i := 0; i < 1000 && (math.Abs(h.X-10) > 1e-9 || math.Abs(h.Y-0) > 1e-9); i++ {
		h.moveStep(0.1)
	}
	if math.Abs(h.X-10) > 1e-9 || math.Abs(h.Y) > 1e-9 {
		t.Errorf("did not converge to rally: got (%v,%v), want (10,0)", h.X, h.Y)
	}
}

// TestHeroSetRallyDeadGuard: 阵亡期间 (respawnCD>0) SetRally 被忽略。
func TestHeroSetRallyDeadGuard(t *testing.T) {
	h := newHero(5, 5)
	h.respawnCD = 3.0 // 模拟阵亡 (P2 才会真正设置)
	if h.Alive() {
		t.Error("hero with respawnCD>0 should not be alive")
	}
	h.SetRally(20, 20)
	if h.RallyX != 5 || h.RallyY != 5 {
		t.Errorf("dead hero rally must not change: got (%v,%v), want (5,5)", h.RallyX, h.RallyY)
	}
	// 阵亡期间 moveStep 不动
	h.moveStep(0.5)
	if h.X != 5 || h.Y != 5 {
		t.Errorf("dead hero must not move: got (%v,%v), want (5,5)", h.X, h.Y)
	}
}

// TestBeginRunSpawnsHero: 进入关卡后英雄生成在 path 中点, 满血在场。
func TestBeginRunSpawnsHero(t *testing.T) {
	g := newTestGame()
	if g.Hero == nil {
		t.Fatal("StartLevel should spawn a hero")
	}
	mid := g.Path[len(g.Path)/2]
	if g.Hero.X != float64(mid.X) || g.Hero.Y != float64(mid.Y) {
		t.Errorf("hero spawn: got (%v,%v), want path mid (%d,%d)",
			g.Hero.X, g.Hero.Y, mid.X, mid.Y)
	}
	if g.Hero.HP != knight.MaxHP || !g.Hero.Alive() {
		t.Errorf("hero should spawn full HP alive: HP=%d alive=%v", g.Hero.HP, g.Hero.Alive())
	}
}

// TestSetHeroRallyMovesHero: SetHeroRally 设集结点后, Update 驱动英雄移动。
func TestSetHeroRallyMovesHero(t *testing.T) {
	g := newTestGame()
	startX := g.Hero.X
	g.SetHeroRally(Point{X: int(startX) + 5, Y: int(g.Hero.Y)})
	g.Update(0.5) // 半秒 → 应朝集结点移动
	if g.Hero.X <= startX {
		t.Errorf("hero should move toward rally: startX=%v nowX=%v", startX, g.Hero.X)
	}
}

// injectEnemyAtHero: 在英雄所在 path index 注入一个敌人 (测试辅助)。
func injectEnemyAtHero(g *Game, kind EnemyKind, hp int) *Enemy {
	e := &Enemy{Kind: kind, HP: hp, MaxHP: hp, PathIdx: float64(len(g.Path) / 2)}
	g.Enemies = append(g.Enemies, e)
	return e
}

// TestHeroAttacksEnemy: 射程内敌人被英雄攻击 (经 damageEnemy, HP 下降)。
func TestHeroAttacksEnemy(t *testing.T) {
	g := newTestGame()
	e := injectEnemyAtHero(g, ENormal, 20)
	g.Update(0.1) // cooldown 初始 0 → 首帧即出手
	if e.HP >= 20 {
		t.Errorf("hero should damage enemy in range: HP=%d, want < 20", e.HP)
	}
	if e.HP != 20-knight.Damage {
		t.Errorf("expected one hit of %d: HP=%d", knight.Damage, e.HP)
	}
}

// TestHeroAttacksNearest: 两敌在场, 英雄打更近的那个。
func TestHeroAttacksNearest(t *testing.T) {
	g := newTestGame()
	near := injectEnemyAtHero(g, ENormal, 20) // 与英雄同格 (dist 0)
	far := &Enemy{Kind: ENormal, HP: 20, MaxHP: 20, PathIdx: float64(len(g.Path)/2) + 1}
	g.Enemies = append(g.Enemies, far)
	g.Update(0.1)
	if near.HP >= 20 {
		t.Errorf("nearest enemy should be hit: near.HP=%d", near.HP)
	}
	if far.HP != 20 {
		t.Errorf("farther enemy should be untouched: far.HP=%d", far.HP)
	}
}

// TestHeroKillGrantsGold: 英雄击杀经 damageEnemy → killEnemy, 入账金币。
func TestHeroKillGrantsGold(t *testing.T) {
	g := newTestGame()
	injectEnemyAtHero(g, ENormal, knight.Damage-1) // 一击毙
	goldBefore := g.Gold
	g.Update(0.1)
	if g.Gold <= goldBefore {
		t.Errorf("hero kill must grant gold via killEnemy: before=%d after=%d", goldBefore, g.Gold)
	}
}

// TestHeroIgnoresFlying: 飞行单位英雄打不到 (地面近战, 同 Cannon)。
func TestHeroIgnoresFlying(t *testing.T) {
	g := newTestGame()
	e := injectEnemyAtHero(g, EGlider, 18)
	g.Update(0.1)
	if e.HP != 18 {
		t.Errorf("hero must not damage flying enemy: HP=%d, want 18", e.HP)
	}
}

// TestEnemyDamagesHero: 接触敌按 meleeCD 节奏反击英雄 (扣血)。
func TestEnemyDamagesHero(t *testing.T) {
	g := newTestGame()
	injectEnemyAtHero(g, ENormal, 999) // 高血, 不被英雄秒杀, 持续接触
	hpBefore := g.Hero.HP
	g.Update(0.1) // 首帧: 敌 meleeCD=0 → 立即出手
	if g.Hero.HP != hpBefore-enemySpecs[ENormal].Attack {
		t.Errorf("enemy should hit hero once: HP %d → %d, want -%d",
			hpBefore, g.Hero.HP, enemySpecs[ENormal].Attack)
	}
}

// TestEnemyMeleeCadence: 敌人攻击英雄按 meleeCD 节奏, 不是每帧。
func TestEnemyMeleeCadence(t *testing.T) {
	g := newTestGame()
	injectEnemyAtHero(g, ENormal, 999)
	hpBefore := g.Hero.HP
	g.Update(0.1) // 首帧出手 (-Attack)
	g.Update(0.1) // 次帧 meleeCD 仍 >0 → 不出手
	if g.Hero.HP != hpBefore-enemySpecs[ENormal].Attack {
		t.Errorf("enemy must respect meleeCD (one hit in 2 frames): HP=%d, want %d",
			g.Hero.HP, hpBefore-enemySpecs[ENormal].Attack)
	}
}

// TestHeroDeathAndRespawn: 英雄归零 → 阵亡 → RespawnS 后满血复活在 path 中点。
func TestHeroDeathAndRespawn(t *testing.T) {
	g := newTestGame()
	g.Hero.X, g.Hero.Y = 0, 0 // 先把英雄挪到角落, 验证复活回 path 中点
	g.hurtHero(knight.MaxHP)
	if g.Hero.Alive() || g.Hero.HP != 0 {
		t.Fatalf("hero should be down: alive=%v HP=%d", g.Hero.Alive(), g.Hero.HP)
	}
	if g.Hero.respawnCD != knight.RespawnS {
		t.Errorf("respawnCD: got %v, want %v", g.Hero.respawnCD, knight.RespawnS)
	}
	// 复活前 SetRally 被忽略
	g.SetHeroRally(Point{X: 9, Y: 9})
	if g.Hero.RallyX == 9 {
		t.Error("dead hero rally must not change")
	}
	// 推进足够时间 → 复活
	g.updateHero(knight.RespawnS + 0.1)
	if !g.Hero.Alive() || g.Hero.HP != knight.MaxHP {
		t.Errorf("hero should revive full HP: alive=%v HP=%d", g.Hero.Alive(), g.Hero.HP)
	}
	mid := g.Path[len(g.Path)/2]
	if g.Hero.X != float64(mid.X) || g.Hero.Y != float64(mid.Y) {
		t.Errorf("hero should revive at path mid (%d,%d): got (%v,%v)",
			mid.X, mid.Y, g.Hero.X, g.Hero.Y)
	}
}

// ============================================================
// V8 Phase 3: 阻挡机制 (敌人停步互殴)
// ============================================================

// TestHeroBlocksGroundEnemy: 贴身非飞行敌被阻挡 — 停步不前进。
func TestHeroBlocksGroundEnemy(t *testing.T) {
	g := newTestGame()
	e := injectEnemyAtHero(g, ENormal, 999) // 高血, 撑过测试时长不被秒
	g.Update(0.5)                            // 若不阻挡, ENormal speed 3.0 应前进 1.5
	if !e.Blocked {
		t.Error("ground enemy adjacent to hero should be Blocked")
	}
	if e.PathIdx != float64(len(g.Path)/2) {
		t.Errorf("blocked enemy must not advance: PathIdx=%v, want %v (不动)",
			e.PathIdx, len(g.Path)/2)
	}
}

// TestHeroDoesNotBlockFlying: 飞行敌飞越, 不被阻挡, 正常前进。
func TestHeroDoesNotBlockFlying(t *testing.T) {
	g := newTestGame()
	e := injectEnemyAtHero(g, EGlider, 999)
	start := e.PathIdx
	g.Update(0.1)
	if e.Blocked {
		t.Error("flying enemy must not be blocked")
	}
	if e.PathIdx <= start {
		t.Errorf("flying enemy should advance over hero: PathIdx %v → %v", start, e.PathIdx)
	}
}

// TestBlockReleasedOnHeroDeath: 英雄阵亡 → 阻挡解除, 敌恢复前进。
func TestBlockReleasedOnHeroDeath(t *testing.T) {
	g := newTestGame()
	e := injectEnemyAtHero(g, ENormal, 999)
	g.Update(0.1)
	if !e.Blocked {
		t.Fatal("enemy should be blocked while hero alive")
	}
	blockedIdx := e.PathIdx
	g.hurtHero(knight.MaxHP) // 英雄阵亡
	g.Update(0.1)
	if e.Blocked {
		t.Error("block must release when hero is down")
	}
	if e.PathIdx <= blockedIdx {
		t.Errorf("enemy should advance after hero death: %v → %v", blockedIdx, e.PathIdx)
	}
}

// TestBlockedEnemyExchangesDamage: 阻挡期间双向扣血 (P2+P3 联动)。
func TestBlockedEnemyExchangesDamage(t *testing.T) {
	g := newTestGame()
	e := injectEnemyAtHero(g, ENormal, 999)
	heroHP := g.Hero.HP
	g.Update(0.1)
	if !e.Blocked {
		t.Fatal("enemy should be blocked")
	}
	if e.HP >= 999 {
		t.Errorf("hero should damage blocked enemy: HP=%d", e.HP)
	}
	if g.Hero.HP >= heroHP {
		t.Errorf("blocked enemy should damage hero: HP %d → %d", heroHP, g.Hero.HP)
	}
}

// TestDistantEnemyNotBlocked: 射程外敌不受阻挡 (阻挡不是全场冻结)。
func TestDistantEnemyNotBlocked(t *testing.T) {
	g := newTestGame()
	// PathIdx 0 = (0,0), 距英雄 (3,0) = 3 cell > heroContactRange
	e := &Enemy{Kind: ENormal, HP: 999, MaxHP: 999, PathIdx: 0}
	g.Enemies = append(g.Enemies, e)
	g.Update(0.1)
	if e.Blocked {
		t.Error("distant enemy must not be blocked")
	}
	if e.PathIdx <= 0 {
		t.Errorf("distant enemy should advance freely: PathIdx=%v", e.PathIdx)
	}
}

// ============================================================
// V9 Phase 1: 英雄关内 XP + 等级成长
// ============================================================

// TestHeroLevelCurve: 喂 XP 逐级升级 — 属性提升 + 升级回满血 + 满级封顶。
func TestHeroLevelCurve(t *testing.T) {
	h := newHero(0, 0)
	if h.Level != 1 || h.XP != 0 {
		t.Fatalf("new hero should be lvl1 0xp, got lvl%d xp%d", h.Level, h.XP)
	}
	base := knight.MaxHP
	h.HP = 1 // 模拟受伤, 验证升级回满血
	h.GainXP(knight.xpForNext(1))
	if h.Level != 2 {
		t.Errorf("should reach lvl2 after %d xp, got %d", knight.xpForNext(1), h.Level)
	}
	if h.MaxHP != base+knight.HPPerLvl {
		t.Errorf("lvl2 maxHP = %d, want %d", h.MaxHP, base+knight.HPPerLvl)
	}
	if h.HP != h.MaxHP {
		t.Errorf("level up should heal to full: HP=%d MaxHP=%d", h.HP, h.MaxHP)
	}
	// 巨量 XP → 封顶
	h.GainXP(100000)
	if h.Level != knight.LevelCap {
		t.Errorf("should cap at lvl%d, got %d", knight.LevelCap, h.Level)
	}
	if h.XP != 0 {
		t.Errorf("capped hero XP should lock 0, got %d", h.XP)
	}
}

// TestHeroLevelStatsScale: 高等级 Damage / AttackRange / maxHP 按级递增。
func TestHeroLevelStatsScale(t *testing.T) {
	h := newHero(0, 0)
	if h.Damage() != knight.Damage || math.Abs(h.AttackRange()-knight.Range) > 1e-9 {
		t.Errorf("lvl1 stats should equal base: dmg=%d range=%v", h.Damage(), h.AttackRange())
	}
	h.Level = 3
	if h.Damage() != knight.Damage+2*knight.DmgPerLvl {
		t.Errorf("lvl3 damage = %d, want %d", h.Damage(), knight.Damage+2*knight.DmgPerLvl)
	}
	if math.Abs(h.AttackRange()-(knight.Range+2*knight.RangePerLvl)) > 1e-9 {
		t.Errorf("lvl3 range = %v, want %v", h.AttackRange(), knight.Range+2*knight.RangePerLvl)
	}
	if knight.maxHPFor(3) != knight.MaxHP+2*knight.HPPerLvl {
		t.Errorf("lvl3 maxHP = %d, want %d", knight.maxHPFor(3), knight.MaxHP+2*knight.HPPerLvl)
	}
}

// TestHeroXPFromKill: 英雄击杀给威胁加权 XP (enemyCost)。
func TestHeroXPFromKill(t *testing.T) {
	g := newTestGame()
	injectEnemyAtHero(g, ENormal, knight.Damage-1) // 一击毙
	if g.Hero.XP != 0 {
		t.Fatal("hero should start with 0 XP")
	}
	g.Update(0.1) // 英雄秒杀 → 得 enemyCost[ENormal] XP
	if g.Hero.XP != enemyCost[ENormal] {
		t.Errorf("hero should gain %d XP from kill, got %d", enemyCost[ENormal], g.Hero.XP)
	}
}

// TestHeroPerRunReset: 新关重置等级/XP (per-run, beginRun → newHero)。
func TestHeroPerRunReset(t *testing.T) {
	g := newTestGame()
	g.Hero.Level = 4
	g.Hero.XP = 3
	g.StartLevel(0) // 重开 → beginRun → newHero
	if g.Hero.Level != 1 || g.Hero.XP != 0 {
		t.Errorf("per-run: new battle should reset to lvl1 0xp, got lvl%d xp%d",
			g.Hero.Level, g.Hero.XP)
	}
}

// TestHeroRespawnKeepsLevel: 阵亡复活保留等级 (关内死不掉级), HP 回当级满血。
func TestHeroRespawnKeepsLevel(t *testing.T) {
	g := newTestGame()
	g.Hero.GainXP(knight.xpForNext(1)) // → lvl2
	if g.Hero.Level != 2 {
		t.Fatalf("setup: should be lvl2, got %d", g.Hero.Level)
	}
	lvMaxHP := g.Hero.MaxHP
	g.hurtHero(g.Hero.MaxHP)               // 阵亡
	g.updateHero(knight.RespawnS + 0.1)  // 复活
	if g.Hero.Level != 2 {
		t.Errorf("respawn should keep level: got %d", g.Hero.Level)
	}
	if g.Hero.HP != lvMaxHP {
		t.Errorf("respawn should heal to leveled max %d, got %d", lvMaxHP, g.Hero.HP)
	}
}

// ============================================================
// V9 Phase 2: 英雄主动技能 (AoE 横扫)
// ============================================================

// TestHeroAbilityLockedBeforeLevel: 未达解锁等级时技能锁定, 释放失败。
func TestHeroAbilityLockedBeforeLevel(t *testing.T) {
	g := newTestGame()
	g.Hero.Level = knight.AbilityLevel - 1
	if g.Hero.AbilityUnlocked() {
		t.Error("ability should be locked below unlock level")
	}
	if g.CastHeroAbility() {
		t.Error("cast should fail when locked")
	}
}

// TestHeroAbilityHitsGround: 解锁后横扫打半径内地面敌 (伤害 = Damage()×mul)。
func TestHeroAbilityHitsGround(t *testing.T) {
	g := newTestGame()
	g.Hero.Level = knight.AbilityLevel
	if !g.Hero.AbilityReady() {
		t.Fatal("ability should be ready at unlock level, no cooldown")
	}
	e := injectEnemyAtHero(g, ENormal, 999)
	if !g.CastHeroAbility() {
		t.Fatal("cast should succeed at unlock level")
	}
	wantDmg := g.Hero.Damage() * knight.AbilityDmgMul
	if e.HP != 999-wantDmg {
		t.Errorf("cleave dmg: HP=%d, want %d (Damage %d × %d)",
			e.HP, 999-wantDmg, g.Hero.Damage(), knight.AbilityDmgMul)
	}
}

// TestHeroAbilityCooldownGating: 释放后进冷却, 冷却中再释放失败, 随时间衰减。
func TestHeroAbilityCooldownGating(t *testing.T) {
	g := newTestGame()
	g.Hero.Level = knight.AbilityLevel
	if !g.CastHeroAbility() {
		t.Fatal("first cast should succeed")
	}
	if g.Hero.abilityCD <= 0 {
		t.Error("cooldown should be set after cast")
	}
	if g.CastHeroAbility() {
		t.Error("second cast within cooldown should fail")
	}
	cd := g.Hero.abilityCD
	g.Update(0.5)
	if g.Hero.abilityCD >= cd {
		t.Errorf("abilityCD should decay: %v → %v", cd, g.Hero.abilityCD)
	}
}

// TestHeroAbilityIgnoresFlying: 横扫不打飞行 (同英雄近战定位)。
func TestHeroAbilityIgnoresFlying(t *testing.T) {
	g := newTestGame()
	g.Hero.Level = knight.AbilityLevel
	e := injectEnemyAtHero(g, EGlider, 999)
	g.CastHeroAbility()
	if e.HP != 999 {
		t.Errorf("cleave must not hit flying: HP=%d", e.HP)
	}
}

// TestHeroAbilityKillGrantsGold: 横扫击杀经 damageEnemy → killEnemy 入账金币。
func TestHeroAbilityKillGrantsGold(t *testing.T) {
	g := newTestGame()
	g.Hero.Level = knight.AbilityLevel
	injectEnemyAtHero(g, ENormal, 1) // 横扫秒杀
	goldBefore := g.Gold
	g.CastHeroAbility()
	if g.Gold <= goldBefore {
		t.Errorf("cleave kill must grant gold via killEnemy: %d → %d", goldBefore, g.Gold)
	}
}
