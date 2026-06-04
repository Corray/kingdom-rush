// 单元测试覆盖纯逻辑层:
//   - ExpandPath / ParseWave (level.go)
//   - Tower upgrade spec invariants (entities.go)
//   - TryAction (build / not enough gold / on-path / upgrade) (game.go)
//   - Update lives-dec on escape (game.go)
//
// 不覆盖: render layer (tcell screen 需 TTY) / main 事件循环 (集成测试)。
package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

// ============================================================
// ExpandPath
// ============================================================

func TestExpandPath_Horizontal(t *testing.T) {
	p, err := ExpandPath([][]int{{0, 5}, {3, 5}})
	if err != nil {
		t.Fatal(err)
	}
	want := []Point{{0, 5}, {1, 5}, {2, 5}, {3, 5}}
	if !pointsEqual(p, want) {
		t.Errorf("got %v want %v", p, want)
	}
}

func TestExpandPath_Vertical(t *testing.T) {
	p, err := ExpandPath([][]int{{2, 0}, {2, 3}})
	if err != nil {
		t.Fatal(err)
	}
	want := []Point{{2, 0}, {2, 1}, {2, 2}, {2, 3}}
	if !pointsEqual(p, want) {
		t.Errorf("got %v want %v", p, want)
	}
}

func TestExpandPath_BackwardsVertical(t *testing.T) {
	p, err := ExpandPath([][]int{{0, 3}, {0, 0}})
	if err != nil {
		t.Fatal(err)
	}
	want := []Point{{0, 3}, {0, 2}, {0, 1}, {0, 0}}
	if !pointsEqual(p, want) {
		t.Errorf("got %v want %v", p, want)
	}
}

func TestExpandPath_MultiSegment(t *testing.T) {
	p, err := ExpandPath([][]int{{0, 0}, {2, 0}, {2, 2}})
	if err != nil {
		t.Fatal(err)
	}
	want := []Point{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}}
	if !pointsEqual(p, want) {
		t.Errorf("got %v want %v", p, want)
	}
}

func TestExpandPath_NotAligned(t *testing.T) {
	if _, err := ExpandPath([][]int{{0, 0}, {1, 1}}); err == nil {
		t.Errorf("expected error for diagonal cp")
	}
}

func TestExpandPath_TooFew(t *testing.T) {
	if _, err := ExpandPath([][]int{{0, 0}}); err == nil {
		t.Errorf("expected error for single cp")
	}
}

func TestExpandPath_MalformedCP(t *testing.T) {
	if _, err := ExpandPath([][]int{{0, 0}, {1}}); err == nil {
		t.Errorf("expected error for malformed cp [1]")
	}
}

// ============================================================
// ParseWave
// ============================================================

func TestParseWave_Single(t *testing.T) {
	seq, err := ParseWave("n3")
	if err != nil {
		t.Fatal(err)
	}
	want := []EnemyKind{ENormal, ENormal, ENormal}
	if !kindsEqual(seq, want) {
		t.Errorf("got %v want %v", seq, want)
	}
}

func TestParseWave_Mixed(t *testing.T) {
	seq, err := ParseWave("n2 f1 n1")
	if err != nil {
		t.Fatal(err)
	}
	want := []EnemyKind{ENormal, ENormal, EFast, ENormal}
	if !kindsEqual(seq, want) {
		t.Errorf("got %v want %v", seq, want)
	}
}

func TestParseWave_Empty(t *testing.T) {
	if _, err := ParseWave(""); err == nil {
		t.Errorf("expected error for empty wave")
	}
}

func TestParseWave_UnknownKind(t *testing.T) {
	if _, err := ParseWave("z3"); err == nil {
		t.Errorf("expected error for unknown kind")
	}
}

func TestParseWave_InvalidCount(t *testing.T) {
	if _, err := ParseWave("nABC"); err == nil {
		t.Errorf("expected error for invalid count")
	}
}

func TestParseWave_NonPositiveCount(t *testing.T) {
	if _, err := ParseWave("n0"); err == nil {
		t.Errorf("expected error for n0")
	}
}

// ============================================================
// Tower spec invariants
// ============================================================

func TestTowerSpec_LevelsMonotonicDamage(t *testing.T) {
	for kind, spec := range towerSpecs {
		for i := 1; i < 3; i++ {
			if spec.Levels[i].Damage <= spec.Levels[i-1].Damage {
				t.Errorf("kind=%d lvl %d→%d damage not increasing (%d→%d)",
					kind, i, i+1, spec.Levels[i-1].Damage, spec.Levels[i].Damage)
			}
		}
	}
}

func TestTowerSpec_LevelsMonotonicRange(t *testing.T) {
	for kind, spec := range towerSpecs {
		for i := 1; i < 3; i++ {
			if spec.Levels[i].Range <= spec.Levels[i-1].Range {
				t.Errorf("kind=%d lvl %d→%d range not increasing (%.1f→%.1f)",
					kind, i, i+1, spec.Levels[i-1].Range, spec.Levels[i].Range)
			}
		}
	}
}

// ============================================================
// TryAction (build / upgrade)
// ============================================================

func TestTryAction_Build(t *testing.T) {
	g := newTestGame()
	g.Cursor = Point{5, 5} // off path
	g.TryAction()
	if len(g.Towers) != 1 {
		t.Fatalf("expected 1 tower, got %d", len(g.Towers))
	}
	wantGold := 100 - towerSpecs[TArcher].Levels[0].Cost
	if g.Gold != wantGold {
		t.Errorf("gold: got %d want %d", g.Gold, wantGold)
	}
	if g.Towers[0].Level != 1 {
		t.Errorf("level: got %d want 1", g.Towers[0].Level)
	}
}

func TestTryAction_NotEnoughGold(t *testing.T) {
	g := newTestGame()
	g.Gold = 30
	g.Cursor = Point{5, 5}
	g.TryAction()
	if len(g.Towers) != 0 {
		t.Errorf("expected 0 towers, got %d", len(g.Towers))
	}
}

func TestTryAction_OnPath(t *testing.T) {
	g := newTestGame()
	g.Cursor = g.Path[0]
	g.TryAction()
	if len(g.Towers) != 0 {
		t.Errorf("should not build on path")
	}
}

func TestTryAction_Upgrade(t *testing.T) {
	g := newTestGame()
	g.Gold = 1000
	g.Cursor = Point{5, 5}
	g.TryAction() // 建
	g.TryAction() // 升 2
	g.TryAction() // 升 3
	if g.Towers[0].Level != 3 {
		t.Errorf("expected level 3, got %d", g.Towers[0].Level)
	}
	// 第 4 次应该 max-level no-op
	prevGold := g.Gold
	g.TryAction()
	if g.Towers[0].Level != 3 {
		t.Errorf("should stay at max level 3, got %d", g.Towers[0].Level)
	}
	if g.Gold != prevGold {
		t.Errorf("gold should not change at max level: %d -> %d", prevGold, g.Gold)
	}
}

func TestTryAction_UpgradeInsufficientGold(t *testing.T) {
	g := newTestGame()
	g.Cursor = Point{5, 5}
	g.TryAction() // 建 50g, 剩 50
	g.Gold = 10   // 不够升级 40
	g.TryAction()
	if g.Towers[0].Level != 1 {
		t.Errorf("expected to stay level 1, got %d", g.Towers[0].Level)
	}
}

// ============================================================
// Update / Lives
// ============================================================

func TestUpdate_LivesDecOnEscape(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.Enemies = append(g.Enemies, &Enemy{
		Kind: ENormal, HP: 20, MaxHP: 20,
		PathIdx: float64(len(g.Path) - 1),
	})
	g.Update(0.1)
	if g.Lives != g.StartLives-1 {
		t.Errorf("lives: got %d want %d", g.Lives, g.StartLives-1)
	}
}

func TestUpdate_LostWhenLivesZero(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.Lives = 1
	g.Enemies = append(g.Enemies, &Enemy{
		Kind: ENormal, HP: 20, MaxHP: 20,
		PathIdx: float64(len(g.Path) - 1),
	})
	g.Update(0.1)
	if g.Phase != PhaseLost {
		t.Errorf("phase: got %v want PhaseLost", g.Phase)
	}
}

func TestUpdate_NoOpWhenNotPlaying(t *testing.T) {
	g := newTestGame()
	g.Phase = PhaseLevelSelect
	startGold := g.Gold
	g.Update(1.0)
	if g.Gold != startGold {
		t.Errorf("gold changed during non-playing phase")
	}
}

// ============================================================
// LoadLevels (集成: 真实 yaml)
// ============================================================

func TestLoadLevels_10Levels(t *testing.T) {
	// V6 Phase 1: 10 → 20 关 (测试名保留防 diff 噪音, 语义为"全量加载")
	levels, err := LoadLevels()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(levels) != 20 {
		t.Errorf("expected 20 levels, got %d", len(levels))
	}
	for _, lv := range levels {
		if len(lv.Path) < 2 {
			t.Errorf("level %d: path too short (%d)", lv.ID, len(lv.Path))
		}
		if len(lv.Waves) < 1 {
			t.Errorf("level %d: no waves", lv.ID)
		}
	}
}

// ============================================================
// V1.6: 3 种塔 × 4 种敌人 (含飞行 / Boss)
// ============================================================

func TestParseWave_AllKinds(t *testing.T) {
	seq, err := ParseWave("n1 f1 g1 b1")
	if err != nil {
		t.Fatal(err)
	}
	want := []EnemyKind{ENormal, EFast, EGlider, EBoss}
	if !kindsEqual(seq, want) {
		t.Errorf("got %v want %v", seq, want)
	}
}

func TestTowerSpec_HitsFlying(t *testing.T) {
	cases := map[TowerKind]bool{
		TArcher: true,
		TCannon: false, // Cannon 打不到飞行
		TMagic:  true,
	}
	for kind, want := range cases {
		if got := towerSpecs[kind].HitsFlying; got != want {
			t.Errorf("tower %d HitsFlying: got %v want %v", kind, got, want)
		}
	}
}

func TestEnemySpec_Flying(t *testing.T) {
	if !enemySpecs[EGlider].Flying {
		t.Errorf("Glider should be Flying")
	}
	if enemySpecs[EBoss].Flying {
		t.Errorf("Boss should not be Flying")
	}
	if enemySpecs[ENormal].Flying || enemySpecs[EFast].Flying {
		t.Errorf("Normal/Fast should not be Flying")
	}
}

func TestUpdate_CannonSkipsFlying(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	// 建 Cannon 在 (1, 1) 邻近 path
	g.Towers = append(g.Towers, &Tower{Pos: Point{1, 1}, Kind: TCannon, Level: 1})
	// 放 Glider 在 path 中 (range 2.5 内)
	g.Enemies = append(g.Enemies, &Enemy{
		Kind: EGlider, HP: 18, MaxHP: 18, PathIdx: 1,
	})
	g.Update(0.1)
	// Cannon cd 1.5 已耗一帧,但目标 filter 飞行不击 → enemy HP 不变
	if g.Enemies[0].HP != 18 {
		t.Errorf("Glider HP should remain 18 (Cannon should skip), got %d", g.Enemies[0].HP)
	}
}

func TestUpdate_ArcherHitsFlying(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	// Archer 在 (1, 1) 接近 path
	g.Towers = append(g.Towers, &Tower{Pos: Point{1, 1}, Kind: TArcher, Level: 1})
	g.Enemies = append(g.Enemies, &Enemy{
		Kind: EGlider, HP: 18, MaxHP: 18, PathIdx: 1,
	})
	g.Update(0.1) // Archer cd 0.6, 触发一次射击
	if g.Enemies[0].HP == 18 {
		t.Errorf("Glider HP should be reduced (Archer hits flying)")
	}
}

// V3.6: Spawner 死时召唤 2 个 ENormal
func TestSpawner_DeathSpawnsNormals(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1 // 跳过 wave spawn (避免 Update 自动 spawn 干扰)
	g.Towers = []*Tower{{Pos: Point{1, 1}, Kind: TArcher, Level: 1}}
	// 单个 Spawner enemy in range, HP=1 → 一击致死 → 触发 spawn
	g.Enemies = []*Enemy{{Kind: ESpawner, HP: 1, MaxHP: 35, PathIdx: 1}}
	g.Update(0.7)
	if len(g.Enemies) != 3 {
		t.Errorf("expected 1 dead spawner + 2 normals = 3 total, got %d", len(g.Enemies))
	}
	normalsAlive := 0
	for _, e := range g.Enemies {
		if e.Kind == ENormal && !e.Dead && !e.Escaped {
			normalsAlive++
		}
	}
	if normalsAlive != 2 {
		t.Errorf("expected 2 alive normals, got %d", normalsAlive)
	}
}

func TestParseWave_Spawner(t *testing.T) {
	seq, err := ParseWave("s2")
	if err != nil {
		t.Fatal(err)
	}
	if len(seq) != 2 || seq[0] != ESpawner || seq[1] != ESpawner {
		t.Errorf("expected 2 ESpawner, got %v", seq)
	}
}

func TestUpdate_MagicHitsFlying(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.Towers = append(g.Towers, &Tower{Pos: Point{1, 1}, Kind: TMagic, Level: 1})
	g.Enemies = append(g.Enemies, &Enemy{
		Kind: EGlider, HP: 18, MaxHP: 18, PathIdx: 1,
	})
	g.Update(0.1)
	if g.Enemies[0].HP == 18 {
		t.Errorf("Glider HP should be reduced (Magic hits flying)")
	}
}

// ============================================================
// V2.6: 攻击视觉特效
// ============================================================

func TestEffects_AddOnShoot(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.Towers = []*Tower{{Pos: Point{1, 1}, Kind: TArcher, Level: 1}}
	g.Enemies = []*Enemy{{Kind: ENormal, HP: 20, MaxHP: 20, PathIdx: 1}}
	if len(g.Effects) != 0 {
		t.Fatalf("pre: expected 0 effects, got %d", len(g.Effects))
	}
	g.Update(0.1) // Archer cd 0.6, 但 cooldown 初始 0 → 第一帧即射击
	// V4 Phase 4 起: shoot + hit + 伤害飘字 = 3 effects
	if len(g.Effects) != 3 {
		t.Fatalf("expected 3 effects (shoot+hit+dmg text), got %d", len(g.Effects))
	}
	kinds := map[EffectKind]int{}
	for _, e := range g.Effects {
		kinds[e.Kind]++
	}
	if kinds[EShoot] != 1 || kinds[EHit] != 1 || kinds[EText] != 1 {
		t.Errorf("expected 1 EShoot + 1 EHit + 1 EText, got %v", kinds)
	}
}

func TestEffects_DecayTTL(t *testing.T) {
	effects := []*Effect{
		{Kind: EShoot, TTL: 0.10, MaxTTL: 0.15},
		{Kind: EHit, TTL: 0.25, MaxTTL: 0.30},
	}
	out := decayEffects(effects, 0.05)
	if len(out) != 2 {
		t.Errorf("both should survive: got %d", len(out))
	}
	if out[0].TTL != 0.05 {
		t.Errorf("shoot TTL: got %f want 0.05", out[0].TTL)
	}
}

func TestEffects_RemoveExpired(t *testing.T) {
	effects := []*Effect{
		{Kind: EShoot, TTL: 0.02, MaxTTL: 0.15}, // 减后过期
		{Kind: EHit, TTL: 0.25, MaxTTL: 0.30},   // 减后存活
	}
	out := decayEffects(effects, 0.05)
	if len(out) != 1 {
		t.Errorf("expected 1 surviving effect, got %d", len(out))
	}
	if out[0].Kind != EHit {
		t.Errorf("EHit should survive, got %v", out[0].Kind)
	}
}

func TestEffects_Alpha(t *testing.T) {
	e := &Effect{TTL: 0.15, MaxTTL: 0.3}
	if a := e.Alpha(); a < 0.49 || a > 0.51 {
		t.Errorf("Alpha at half TTL: got %f want 0.5", a)
	}
	e.TTL = 0
	if a := e.Alpha(); a != 0 {
		t.Errorf("Alpha at 0 TTL: got %f want 0", a)
	}
	e.TTL = e.MaxTTL
	if a := e.Alpha(); a != 1 {
		t.Errorf("Alpha at full TTL: got %f want 1", a)
	}
}

// ============================================================
// V1.7: 存档 / unlock 系统
// ============================================================

func TestSave_IsUnlocked_DefaultOnlyL1(t *testing.T) {
	s := NewSave()
	if !s.IsUnlocked(1) {
		t.Errorf("Level 1 should always be unlocked")
	}
	if s.IsUnlocked(2) {
		t.Errorf("Level 2 should be locked initially")
	}
}

func TestSave_IsUnlocked_ChainAfterCompletion(t *testing.T) {
	s := NewSave()
	s.MarkCompleted(1)
	if !s.IsUnlocked(2) {
		t.Errorf("Level 2 should unlock after L1 completed")
	}
	if s.IsUnlocked(3) {
		t.Errorf("Level 3 should still be locked (L2 not done)")
	}
	s.MarkCompleted(2)
	if !s.IsUnlocked(3) {
		t.Errorf("Level 3 should unlock after L2 completed")
	}
}

func TestSave_Roundtrip(t *testing.T) {
	withTempSavePath(t, func() {
		s := NewSave()
		s.MarkCompleted(1)
		s.MarkCompleted(3)
		if err := StoreSave(s); err != nil {
			t.Fatal(err)
		}
		loaded, err := LoadSave()
		if err != nil {
			t.Fatal(err)
		}
		if !loaded.IsCompleted(1) || !loaded.IsCompleted(3) {
			t.Errorf("expected L1 and L3 completed, got %v", loaded.Completed)
		}
		if loaded.IsCompleted(2) {
			t.Errorf("L2 should not be completed")
		}
	})
}

func TestSave_LoadMissingFileNoError(t *testing.T) {
	withTempSavePath(t, func() {
		s, err := LoadSave()
		if err != nil {
			t.Fatalf("missing file should not error: %v", err)
		}
		if s.IsCompleted(1) {
			t.Errorf("empty save should have no completed")
		}
	})
}

func TestGame_StartLevelRejectsLocked(t *testing.T) {
	g := newTestGameMultiLevel()
	g.StartLevel(1) // Lv 2 (idx 1), L1 not completed → locked
	if g.Phase != PhaseLevelSelect {
		t.Errorf("locked level should not start, phase: %v", g.Phase)
	}
	if g.Msg == "" {
		t.Errorf("expected locked msg")
	}
}

func TestGame_VictoryMarksCompletedAndSaves(t *testing.T) {
	withTempSavePath(t, func() {
		g := newTestGame()
		g.prepTimer = 0
		g.spawned = 1 // 已 spawn 一个 (= wave size)
		// 模拟唯一敌人已 dead
		g.Enemies = []*Enemy{{Kind: ENormal, HP: 0, MaxHP: 20, Dead: true}}
		g.Update(0.1)
		if g.Phase != PhaseWon {
			t.Fatalf("expected PhaseWon, got %v", g.Phase)
		}
		if !g.Save.IsCompleted(1) {
			t.Errorf("L1 should be marked completed")
		}
		// 验证落盘: 重新 load 应该也有 L1 completed
		loaded, err := LoadSave()
		if err != nil {
			t.Fatal(err)
		}
		if !loaded.IsCompleted(1) {
			t.Errorf("L1 completion should persist to disk")
		}
	})
}

// ============================================================
// helpers
// ============================================================

func withTempSavePath(t *testing.T, fn func()) {
	t.Helper()
	old := savePathFn
	defer func() { savePathFn = old }()
	dir := t.TempDir()
	savePathFn = func() (string, error) {
		return dir + "/save.json", nil
	}
	fn()
}

func newTestGame() *Game {
	levels := []Level{{
		ID: 1, Name: "Test",
		StartGold: 100, StartLives: 5,
		Path:  []Point{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}},
		Waves: []WaveSpec{{Enemies: []EnemyKind{ENormal}}},
	}}
	g := NewGame(levels, NewSave())
	g.StartLevel(0)
	return g
}

func newTestGameMultiLevel() *Game {
	levels := []Level{
		{
			ID: 1, Name: "L1",
			StartGold: 100, StartLives: 5,
			Path:  []Point{{0, 0}, {1, 0}},
			Waves: []WaveSpec{{Enemies: []EnemyKind{ENormal}}},
		},
		{
			ID: 2, Name: "L2",
			StartGold: 100, StartLives: 5,
			Path:  []Point{{0, 0}, {1, 0}},
			Waves: []WaveSpec{{Enemies: []EnemyKind{ENormal}}},
		},
	}
	return NewGame(levels, NewSave())
}

func pointsEqual(a, b []Point) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func kindsEqual(a, b []EnemyKind) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ============================================================
// V4 Phase 1: SoundEvent 队列 (sound.go + game.go 触发点)
// 只测触发 (事件入队), 不测播放 (audio_player.go 是 !term 渲染层)。
// ============================================================

func drainHas(evs []SoundEvent, want SoundEvent) bool {
	for _, e := range evs {
		if e == want {
			return true
		}
	}
	return false
}

func TestSound_BuildPushesEvent(t *testing.T) {
	g := newTestGame()
	g.Cursor = Point{10, 10} // off-path
	g.TryAction()
	if !drainHas(g.DrainSounds(), SndBuild) {
		t.Errorf("build should push SndBuild")
	}
}

func TestSound_FailedBuildPushesNothing(t *testing.T) {
	g := newTestGame()
	g.Gold = 0
	g.Cursor = Point{10, 10}
	g.TryAction()
	if len(g.DrainSounds()) != 0 {
		t.Errorf("failed build (no gold) should push no sound")
	}
}

func TestSound_UpgradePushesEvent(t *testing.T) {
	g := newTestGame()
	g.Cursor = Point{10, 10}
	g.TryAction()   // build
	g.DrainSounds() // 清掉 build 事件
	g.Gold = 1000   // 保证升级买得起
	g.TryAction()   // upgrade (光标在塔上)
	if !drainHas(g.DrainSounds(), SndUpgrade) {
		t.Errorf("upgrade should push SndUpgrade")
	}
}

func TestSound_ShootAndDeathPush(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1
	// 敌人在塔射程内, HP 低保证一击致死
	g.Enemies = []*Enemy{{Kind: ENormal, HP: 1, MaxHP: 20, PathIdx: 1}}
	g.Towers = []*Tower{{Pos: Point{1, 1}, Kind: TArcher, Level: 1}}
	g.Update(0.05)
	evs := g.DrainSounds()
	if !drainHas(evs, SndShootArcher) {
		t.Errorf("archer shot should push SndShootArcher, got %v", evs)
	}
	if !drainHas(evs, SndEnemyDeath) {
		t.Errorf("kill should push SndEnemyDeath, got %v", evs)
	}
}

func TestSound_ShootKindPerTower(t *testing.T) {
	if shootSound(TArcher) != SndShootArcher ||
		shootSound(TCannon) != SndShootCannon ||
		shootSound(TMagic) != SndShootMagic {
		t.Errorf("shootSound mapping wrong")
	}
}

func TestSound_WaveStartOnPrepEnd(t *testing.T) {
	g := newTestGame() // StartLevel 后 prepTimer = wavePrepS
	g.Update(wavePrepS + 1)
	if !drainHas(g.DrainSounds(), SndWaveStart) {
		t.Errorf("prep timer 归零应 push SndWaveStart")
	}
	// 再 update 不应重复触发
	g.Update(0.1)
	if drainHas(g.DrainSounds(), SndWaveStart) {
		t.Errorf("wave start 每 wave 只触发一次")
	}
}

func TestSound_WinPushes(t *testing.T) {
	withTempSavePath(t, func() {
		g := newTestGame()
		g.prepTimer = 0
		g.spawned = 1
		g.Enemies = []*Enemy{{Kind: ENormal, HP: 0, MaxHP: 20, Dead: true}}
		g.Update(0.1)
		if g.Phase != PhaseWon {
			t.Fatalf("expected PhaseWon")
		}
		if !drainHas(g.DrainSounds(), SndWin) {
			t.Errorf("victory should push SndWin")
		}
	})
}

func TestSound_LosePushes(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.Lives = 1
	// 敌人直接走到终点 (path len 6, PathIdx 超出即 escape)
	g.Enemies = []*Enemy{{Kind: ENormal, HP: 20, MaxHP: 20, PathIdx: 4.9}}
	g.spawned = 1
	g.Update(1.0) // speed 1.x × 1s 足够越过末端
	if g.Phase != PhaseLost {
		t.Fatalf("expected PhaseLost, lives=%d", g.Lives)
	}
	if !drainHas(g.DrainSounds(), SndLose) {
		t.Errorf("defeat should push SndLose")
	}
}

func TestSound_DrainClears(t *testing.T) {
	g := newTestGame()
	g.pushSound(SndBuild)
	if len(g.DrainSounds()) != 1 {
		t.Fatalf("first drain should return 1 event")
	}
	if len(g.DrainSounds()) != 0 {
		t.Errorf("second drain should be empty")
	}
}

func TestSound_QueueCap(t *testing.T) {
	g := newTestGame()
	for i := 0; i < maxSoundQueue*2; i++ {
		g.pushSound(SndBuild)
	}
	if n := len(g.DrainSounds()); n != maxSoundQueue {
		t.Errorf("queue should cap at %d, got %d", maxSoundQueue, n)
	}
}

// ============================================================
// V4 Phase 2: 音量档 (save.go) + AdjustVolume (game.go) + bgmFor (bgm.go)
// ============================================================

func TestVolume_DefaultWhenUnset(t *testing.T) {
	s := NewSave()
	if s.VolumeLevel() != defaultVolume {
		t.Errorf("unset volume should default %d, got %d", defaultVolume, s.VolumeLevel())
	}
}

func TestVolume_OldSaveJSONCompat(t *testing.T) {
	// 旧存档无 volume 字段 → 默认档, 不能误判为 0 (静音)
	var s Save
	if err := json.Unmarshal([]byte(`{"completed":{"1":true}}`), &s); err != nil {
		t.Fatal(err)
	}
	if s.VolumeLevel() != defaultVolume {
		t.Errorf("old save should default volume %d, got %d", defaultVolume, s.VolumeLevel())
	}
	if !s.IsCompleted(1) {
		t.Errorf("old save completed should survive")
	}
}

func TestVolume_SetClamp(t *testing.T) {
	s := NewSave()
	s.SetVolumeLevel(-3)
	if s.VolumeLevel() != 0 {
		t.Errorf("clamp low: want 0, got %d", s.VolumeLevel())
	}
	s.SetVolumeLevel(99)
	if s.VolumeLevel() != maxVolume {
		t.Errorf("clamp high: want %d, got %d", maxVolume, s.VolumeLevel())
	}
}

func TestVolume_ExplicitZeroPersists(t *testing.T) {
	// 显式静音 (0) 与未设置 (默认 7) 必须区分 — 指针字段语义
	withTempSavePath(t, func() {
		s := NewSave()
		s.SetVolumeLevel(0)
		if err := StoreSave(s); err != nil {
			t.Fatal(err)
		}
		loaded, err := LoadSave()
		if err != nil {
			t.Fatal(err)
		}
		if loaded.VolumeLevel() != 0 {
			t.Errorf("explicit mute should persist as 0, got %d", loaded.VolumeLevel())
		}
	})
}

func TestVolume_AdjustPersistsAndClamps(t *testing.T) {
	withTempSavePath(t, func() {
		g := newTestGame()
		g.AdjustVolume(-1) // 7 → 6
		if g.Save.VolumeLevel() != 6 {
			t.Fatalf("want 6, got %d", g.Save.VolumeLevel())
		}
		if g.Msg == "" {
			t.Errorf("volume change should set Msg")
		}
		loaded, err := LoadSave()
		if err != nil {
			t.Fatal(err)
		}
		if loaded.VolumeLevel() != 6 {
			t.Errorf("volume should persist, got %d", loaded.VolumeLevel())
		}
		for i := 0; i < 20; i++ {
			g.AdjustVolume(+1)
		}
		if g.Save.VolumeLevel() != maxVolume {
			t.Errorf("repeated +1 should clamp at %d", maxVolume)
		}
	})
}

func TestBGM_TrackForPhase(t *testing.T) {
	cases := []struct {
		phase GamePhase
		want  bgmTrack
	}{
		{PhaseLevelSelect, bgmMenu},
		{PhasePlaying, bgmBattle},
		{PhaseWon, bgmNone},
		{PhaseLost, bgmNone},
	}
	for _, c := range cases {
		if got := bgmFor(c.phase); got != c.want {
			t.Errorf("bgmFor(%v) = %v, want %v", c.phase, got, c.want)
		}
	}
}

// ============================================================
// V4 Phase 3: 走动动画纯数学 (anim.go) + 死亡动画 effect
// ============================================================

func TestAnim_PathLerpMidpoint(t *testing.T) {
	path := []Point{{0, 0}, {1, 0}, {1, 1}}
	x, y := pathLerp(path, 0.5)
	if x != 0.5 || y != 0 {
		t.Errorf("lerp(0.5) = (%v,%v), want (0.5,0)", x, y)
	}
	x, y = pathLerp(path, 1.5)
	if x != 1 || y != 0.5 {
		t.Errorf("lerp(1.5) = (%v,%v), want (1,0.5)", x, y)
	}
}

func TestAnim_PathLerpClamps(t *testing.T) {
	path := []Point{{2, 3}, {3, 3}}
	x, y := pathLerp(path, -1)
	if x != 2 || y != 3 {
		t.Errorf("lerp(-1) should clamp to start, got (%v,%v)", x, y)
	}
	x, y = pathLerp(path, 99)
	if x != 3 || y != 3 {
		t.Errorf("lerp(99) should clamp to end, got (%v,%v)", x, y)
	}
}

func TestAnim_PathDir(t *testing.T) {
	path := []Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
	cases := []struct {
		idx    float64
		dx, dy int
	}{
		{0.5, 1, 0},  // 向右
		{1.5, 0, 1},  // 向下
		{2.5, -1, 0}, // 向左
		{99, -1, 0},  // 越界 clamp 到末段
	}
	for _, c := range cases {
		dx, dy := pathDir(path, c.idx)
		if dx != c.dx || dy != c.dy {
			t.Errorf("pathDir(%v) = (%d,%d), want (%d,%d)", c.idx, dx, dy, c.dx, c.dy)
		}
	}
}

func TestAnim_BobBounded(t *testing.T) {
	for i := 0.0; i < 10; i += 0.1 {
		if b := bobOffset(i); b > bobAmp || b < -bobAmp {
			t.Fatalf("bob(%v) = %v exceeds ±%v", i, b, bobAmp)
		}
	}
	if bobOffset(0) != 0 {
		t.Errorf("bob(0) should be 0")
	}
}

func TestAnim_DirAngle(t *testing.T) {
	if dirAngle(1, 0) != 0 {
		t.Errorf("right should be 0 rad")
	}
	if a := dirAngle(0, 1); a < 1.5 || a > 1.6 {
		t.Errorf("down should be ~π/2, got %v", a)
	}
	if dirAngle(0, 0) != 0 {
		t.Errorf("zero dir should default 0")
	}
}

func TestDeathEffect_PushedOnKill(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1
	g.Enemies = []*Enemy{{Kind: EFast, HP: 1, MaxHP: 10, PathIdx: 1.5}}
	g.Towers = []*Tower{{Pos: Point{1, 1}, Kind: TArcher, Level: 1}}
	g.Update(0.05)
	var death *Effect
	for _, fx := range g.Effects {
		if fx.Kind == EDeath {
			death = fx
		}
	}
	if death == nil {
		t.Fatalf("kill should push EDeath effect")
	}
	if death.Enemy != EFast {
		t.Errorf("death effect should carry enemy kind EFast, got %v", death.Enemy)
	}
	// Update 顺序: 先 move 后 shoot — 死亡位置 = 移动后的插值位置
	wantFX := 1.5 + enemySpecs[EFast].Speed*0.05
	if death.FX != wantFX || death.FY != 0 {
		t.Errorf("death pos should be (%v,0), got (%v,%v)", wantFX, death.FX, death.FY)
	}
}

func TestDeathEffect_Expires(t *testing.T) {
	fx := makeDeathEffect(1, 1, ENormal)
	effects := decayEffects([]*Effect{fx}, 0.4) // TTL 0.35 < 0.4
	if len(effects) != 0 {
		t.Errorf("death effect should expire after TTL")
	}
}

// ============================================================
// V4 Phase 4: 飘字 effect (伤害数字 + 击杀赏金)
// ============================================================

func TestTextEffect_Makers(t *testing.T) {
	d := makeDamageText(1, 2, 12)
	if d.Kind != EText || d.Text != "12" || d.FX != 1 || d.FY != 2 {
		t.Errorf("damage text wrong: %+v", d)
	}
	gtx := makeGoldText(0, 0, 25)
	if gtx.Kind != EText || gtx.Text != "+25g" {
		t.Errorf("gold text wrong: %+v", gtx)
	}
	if gtx.MaxTTL <= d.MaxTTL {
		t.Errorf("gold text should outlive damage text")
	}
}

func TestTextEffect_DamageOnHit(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1
	// 高 HP 保证只命中不击杀
	g.Enemies = []*Enemy{{Kind: EBoss, HP: 500, MaxHP: 500, PathIdx: 1}}
	g.Towers = []*Tower{{Pos: Point{1, 1}, Kind: TArcher, Level: 1}}
	g.Update(0.05)
	var dmgText *Effect
	for _, fx := range g.Effects {
		if fx.Kind == EText {
			dmgText = fx
		}
	}
	if dmgText == nil {
		t.Fatalf("hit should push damage text")
	}
	wantDmg := towerSpecs[TArcher].Levels[0].Damage
	if dmgText.Text != fmt.Sprintf("%d", wantDmg) {
		t.Errorf("damage text = %q, want %d", dmgText.Text, wantDmg)
	}
}

func TestTextEffect_GoldOnKill(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1
	g.Enemies = []*Enemy{{Kind: ENormal, HP: 1, MaxHP: 20, PathIdx: 1}}
	g.Towers = []*Tower{{Pos: Point{1, 1}, Kind: TArcher, Level: 1}}
	g.Update(0.05)
	found := false
	want := fmt.Sprintf("+%dg", enemySpecs[ENormal].Reward)
	for _, fx := range g.Effects {
		if fx.Kind == EText && fx.Text == want {
			found = true
		}
	}
	if !found {
		t.Errorf("kill should push gold text %q, effects: %d", want, len(g.Effects))
	}
}

// ============================================================
// V4 Phase 5: screen shake / 顿帧 (anim.go 纯数学 + game.go 触发源)
// ============================================================

func TestShake_OffsetBoundedAndZero(t *testing.T) {
	if dx, dy := shakeOffset(0); dx != 0 || dy != 0 {
		t.Errorf("expired shake should be (0,0), got (%v,%v)", dx, dy)
	}
	if dx, dy := shakeOffset(-1); dx != 0 || dy != 0 {
		t.Errorf("negative remaining should be (0,0)")
	}
	for r := 0.01; r <= shakeDuration+0.1; r += 0.02 {
		dx, dy := shakeOffset(r)
		if dx > shakeMaxAmp || dx < -shakeMaxAmp || dy > shakeMaxAmp || dy < -shakeMaxAmp {
			t.Fatalf("shake(%v) = (%v,%v) exceeds ±%v", r, dx, dy, shakeMaxAmp)
		}
	}
}

func TestJuice_BossKillCounter(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1
	g.Enemies = []*Enemy{{Kind: EBoss, HP: 1, MaxHP: 100, PathIdx: 1}}
	g.Towers = []*Tower{{Pos: Point{1, 1}, Kind: TArcher, Level: 1}}
	g.Update(0.05)
	if g.BossKills != 1 {
		t.Errorf("boss kill should increment BossKills, got %d", g.BossKills)
	}
}

func TestJuice_NormalKillNoCount(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1
	g.Enemies = []*Enemy{{Kind: ENormal, HP: 1, MaxHP: 20, PathIdx: 1}}
	g.Towers = []*Tower{{Pos: Point{1, 1}, Kind: TArcher, Level: 1}}
	g.Update(0.05)
	if g.BossKills != 0 {
		t.Errorf("normal kill should not increment BossKills, got %d", g.BossKills)
	}
}

func TestJuice_TogglePersists(t *testing.T) {
	withTempSavePath(t, func() {
		g := newTestGame()
		if g.Save.JuiceOff {
			t.Fatalf("juice should default ON (JuiceOff=false)")
		}
		g.ToggleJuice()
		if !g.Save.JuiceOff || g.Msg == "" {
			t.Errorf("toggle should set JuiceOff + Msg")
		}
		loaded, err := LoadSave()
		if err != nil {
			t.Fatal(err)
		}
		if !loaded.JuiceOff {
			t.Errorf("JuiceOff should persist")
		}
		g.ToggleJuice() // 再切回
		if g.Save.JuiceOff {
			t.Errorf("second toggle should restore ON")
		}
	})
}

// ============================================================
// V5 Phase 1: 卖塔 (towerInvested + SellTower)
// ============================================================

func TestSell_TowerInvested(t *testing.T) {
	// Archer: lvl1=50, lvl2=+40, lvl3=+80
	cases := []struct {
		level int
		want  int
	}{
		{1, 50}, {2, 90}, {3, 170},
	}
	for _, c := range cases {
		if got := towerInvested(TArcher, c.level); got != c.want {
			t.Errorf("invested(Archer, lvl%d) = %d, want %d", c.level, got, c.want)
		}
	}
	if got := towerInvested(TArcher, 99); got != 170 {
		t.Errorf("over-level should clamp to full sum, got %d", got)
	}
}

func TestSell_RefundAndRemoval(t *testing.T) {
	withTempSavePath(t, func() {
		g := newTestGame()
		g.Cursor = Point{10, 10}
		g.TryAction() // build Archer lvl1 (gold 100 → 50)
		g.DrainSounds()
		g.SellTower() // refund = 50 × 0.7 = 35
		if len(g.Towers) != 0 {
			t.Fatalf("tower should be removed after sell")
		}
		if g.Gold != 85 { // 50 + 35
			t.Errorf("gold = %d, want 85 (50 + 35 refund)", g.Gold)
		}
		if !drainHas(g.DrainSounds(), SndSell) {
			t.Errorf("sell should push SndSell")
		}
		// 金币飘字
		found := false
		for _, fx := range g.Effects {
			if fx.Kind == EText && fx.Text == "+35g" {
				found = true
			}
		}
		if !found {
			t.Errorf("sell should push gold text +35g")
		}
		if g.Msg == "" {
			t.Errorf("sell should set Msg")
		}
	})
}

func TestSell_UpgradedTowerRefund(t *testing.T) {
	withTempSavePath(t, func() {
		g := newTestGame()
		g.Gold = 1000
		g.Cursor = Point{10, 10}
		g.TryAction() // build lvl1
		g.TryAction() // upgrade → lvl2 (invested 90)
		g.DrainSounds()
		goldBefore := g.Gold
		g.SellTower() // refund = 90 × 0.7 = 63
		if g.Gold != goldBefore+63 {
			t.Errorf("upgraded sell refund = %d, want %d", g.Gold-goldBefore, 63)
		}
	})
}

func TestSell_EmptyCellNoEffect(t *testing.T) {
	g := newTestGame()
	g.Cursor = Point{10, 10}
	goldBefore := g.Gold
	g.SellTower()
	if g.Gold != goldBefore {
		t.Errorf("empty sell should not change gold")
	}
	if len(g.DrainSounds()) != 0 {
		t.Errorf("empty sell should push no sound")
	}
	if g.Msg == "" {
		t.Errorf("empty sell should still give Msg feedback")
	}
}

func TestSell_SoldTowerStopsShooting(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1
	g.Cursor = Point{1, 1}
	g.Towers = []*Tower{{Pos: Point{1, 1}, Kind: TArcher, Level: 1}}
	g.Enemies = []*Enemy{{Kind: ENormal, HP: 100, MaxHP: 100, PathIdx: 1}}
	g.SellTower()
	g.Effects = nil // 清掉卖塔飘字, 只看射击产物
	g.Update(0.05)
	for _, fx := range g.Effects {
		if fx.Kind == EShoot {
			t.Fatalf("sold tower must not shoot")
		}
	}
	if g.Enemies[0].HP != 100 {
		t.Errorf("enemy should be untouched after tower sold")
	}
}

// ============================================================
// V5 Phase 2: targeting 策略 (pickTarget 重构回归 + 三策略)
// ============================================================

// 三敌人 fixture: 前(idx 3, HP 20) / 中(idx 2, HP 99) / 后(idx 1, HP 5)
func targetingFixture() (*Tower, []*Enemy, []Point) {
	path := []Point{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}}
	tower := &Tower{Pos: Point{2, 1}, Kind: TArcher, Level: 1}
	enemies := []*Enemy{
		{Kind: ENormal, HP: 20, MaxHP: 20, PathIdx: 3},
		{Kind: ENormal, HP: 99, MaxHP: 99, PathIdx: 2},
		{Kind: ENormal, HP: 5, MaxHP: 20, PathIdx: 1},
	}
	return tower, enemies, path
}

func TestTargeting_FirstMatchesLegacyBehavior(t *testing.T) {
	// 重构回归: 零值 TargetFirst 必须选 path 最前 (原内联逻辑)
	tower, enemies, path := targetingFixture()
	if got := pickTarget(tower, enemies, path); got != enemies[0] {
		t.Errorf("First should pick front-most (idx 3), got %+v", got)
	}
}

func TestTargeting_LegacyFilters(t *testing.T) {
	// 重构回归: dead/escaped/飞行/射程过滤与原逻辑一致
	path := []Point{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}}
	cannon := &Tower{Pos: Point{2, 1}, Kind: TCannon, Level: 1} // HitsFlying=false
	enemies := []*Enemy{
		{Kind: EGlider, HP: 18, MaxHP: 18, PathIdx: 3},             // 飞行 → cannon 跳过
		{Kind: ENormal, HP: 20, MaxHP: 20, PathIdx: 2, Dead: true}, // dead 跳过
		{Kind: ENormal, HP: 20, MaxHP: 20, PathIdx: 1},
	}
	if got := pickTarget(cannon, enemies, path); got != enemies[2] {
		t.Errorf("cannon should skip flying+dead, pick idx1 normal, got %+v", got)
	}
	// 射程外 → nil
	far := &Tower{Pos: Point{2, 9}, Kind: TArcher, Level: 1} // dy=9 > range 3.5
	if got := pickTarget(far, enemies, path); got != nil {
		t.Errorf("out of range should return nil, got %+v", got)
	}
}

func TestTargeting_Last(t *testing.T) {
	tower, enemies, path := targetingFixture()
	tower.Target = TargetLast
	if got := pickTarget(tower, enemies, path); got != enemies[2] {
		t.Errorf("Last should pick rear-most (idx 1), got %+v", got)
	}
}

func TestTargeting_Strong(t *testing.T) {
	tower, enemies, path := targetingFixture()
	tower.Target = TargetStrong
	if got := pickTarget(tower, enemies, path); got != enemies[1] {
		t.Errorf("Strong should pick HP 99, got %+v", got)
	}
	// HP 平手 → 取最前
	enemies[0].HP = 99
	if got := pickTarget(tower, enemies, path); got != enemies[0] {
		t.Errorf("Strong tie should pick front-most, got %+v", got)
	}
}

func TestTargeting_NextCycles(t *testing.T) {
	m := TargetFirst
	seq := []TargetMode{TargetLast, TargetStrong, TargetFirst}
	for i, want := range seq {
		m = m.Next()
		if m != want {
			t.Fatalf("cycle step %d = %v, want %v", i, m, want)
		}
	}
}

func TestTargeting_CycleOnTower(t *testing.T) {
	g := newTestGame()
	g.Cursor = Point{10, 10}
	g.TryAction() // build
	g.CycleTargeting()
	if g.Towers[0].Target != TargetLast {
		t.Errorf("cycle should switch First → Last, got %v", g.Towers[0].Target)
	}
	if g.Msg == "" {
		t.Errorf("cycle should set Msg")
	}
	g.Cursor = Point{5, 5} // 空地
	g.CycleTargeting()
	if g.Towers[0].Target != TargetLast {
		t.Errorf("empty-cell cycle should not touch other towers")
	}
}

func TestTargeting_StrongShootsHighHP(t *testing.T) {
	// 集成: Strong 塔实际打中高 HP 敌人
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1
	g.Towers = []*Tower{{Pos: Point{1, 1}, Kind: TArcher, Level: 1, Target: TargetStrong}}
	weak := &Enemy{Kind: ENormal, HP: 5, MaxHP: 20, PathIdx: 3}
	strong := &Enemy{Kind: EBoss, HP: 150, MaxHP: 150, PathIdx: 1}
	g.Enemies = []*Enemy{weak, strong}
	g.Update(0.05)
	if strong.HP != 150-towerSpecs[TArcher].Levels[0].Damage {
		t.Errorf("strong enemy should be hit, HP = %d", strong.HP)
	}
	if weak.HP != 5 {
		t.Errorf("weak enemy should be untouched, HP = %d", weak.HP)
	}
}

// ============================================================
// V5 Phase 3: Cannon AoE 溅射 (damageEnemy/killEnemy 重构 + Splash)
// ============================================================

func TestSplash_SpecOnlyCannon(t *testing.T) {
	for _, k := range []TowerKind{TArcher, TMagic} {
		for i, lv := range towerSpecs[k].Levels {
			if lv.Splash != 0 {
				t.Errorf("%s lvl%d should have no splash, got %v",
					towerSpecs[k].Name, i+1, lv.Splash)
			}
		}
	}
	for i, lv := range towerSpecs[TCannon].Levels {
		if lv.Splash <= 0 {
			t.Errorf("Cannon lvl%d should have splash, got %v", i+1, lv.Splash)
		}
	}
}

// 溅射 fixture: cannon (3,1); 主目标 idx3=(3,0); 近邻 idx2.5→(2,0)
// 距主目标 1 ≤ Splash 1.0; 远敌 idx1=(1,0) 距 2 > 1.0
func splashFixture() *Game {
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1
	g.Towers = []*Tower{{Pos: Point{3, 1}, Kind: TCannon, Level: 1}}
	g.Enemies = []*Enemy{
		{Kind: ENormal, HP: 100, MaxHP: 100, PathIdx: 3},   // 主目标 (最前)
		{Kind: ENormal, HP: 100, MaxHP: 100, PathIdx: 2.5}, // 溅射圈内
		{Kind: ENormal, HP: 100, MaxHP: 100, PathIdx: 1},   // 圈外
	}
	return g
}

func TestSplash_CannonHitsCluster(t *testing.T) {
	g := splashFixture()
	g.Update(0.05)
	// Cannon lvl1: 主目标 25, 溅射 round(25×0.5)=13
	if g.Enemies[0].HP != 75 {
		t.Errorf("main target HP = %d, want 75", g.Enemies[0].HP)
	}
	if g.Enemies[1].HP != 87 {
		t.Errorf("splash victim HP = %d, want 87 (25×0.5 round=13)", g.Enemies[1].HP)
	}
	if g.Enemies[2].HP != 100 {
		t.Errorf("out-of-splash enemy HP = %d, want 100", g.Enemies[2].HP)
	}
}

func TestSplash_KillsSpawnerTriggersSummon(t *testing.T) {
	// 家族验收点: 溅射杀与主杀同路径 — Spawner 被溅死必须召唤
	g := splashFixture()
	g.Enemies[1] = &Enemy{Kind: ESpawner, HP: 1, MaxHP: 35, PathIdx: 2.5}
	goldBefore := g.Gold
	g.Update(0.05)
	if !g.Enemies[1].Dead {
		t.Fatalf("spawner should die from splash (HP 1 < 13)")
	}
	// Update 先 move 后 shoot: spawner 死时 PathIdx 已推进, 召唤物继承死亡时位置
	spawnerIdx := g.Enemies[1].PathIdx
	normals := 0
	for _, e := range g.Enemies {
		if e.Kind == ENormal && e.PathIdx == spawnerIdx {
			normals++
		}
	}
	if normals != 2 {
		t.Errorf("splash-killed spawner should summon 2 normals at its death PathIdx %v, got %d", spawnerIdx, normals)
	}
	if g.Gold != goldBefore+enemySpecs[ESpawner].Reward {
		t.Errorf("splash kill should pay reward, gold %d → %d", goldBefore, g.Gold)
	}
}

func TestSplash_FlyingImmuneToCannonSplash(t *testing.T) {
	g := splashFixture()
	g.Enemies[1] = &Enemy{Kind: EGlider, HP: 100, MaxHP: 100, PathIdx: 2.5}
	g.Update(0.05)
	if g.Enemies[1].HP != 100 {
		t.Errorf("flying unit must be immune to cannon splash, HP = %d", g.Enemies[1].HP)
	}
}

func TestSplash_ArcherNoSplashRegression(t *testing.T) {
	g := splashFixture()
	g.Towers = []*Tower{{Pos: Point{3, 1}, Kind: TArcher, Level: 1}}
	g.Update(0.05)
	// Archer 单体: 只有主目标掉血
	if g.Enemies[0].HP != 100-towerSpecs[TArcher].Levels[0].Damage {
		t.Errorf("archer main target HP = %d", g.Enemies[0].HP)
	}
	if g.Enemies[1].HP != 100 || g.Enemies[2].HP != 100 {
		t.Errorf("archer must not splash: HP %d / %d",
			g.Enemies[1].HP, g.Enemies[2].HP)
	}
}

func TestSplash_BossSplashKillCountsForHitStop(t *testing.T) {
	// killEnemy 统一路径回归: 溅射杀 Boss 也要计入 BossKills
	g := splashFixture()
	g.Enemies[1] = &Enemy{Kind: EBoss, HP: 1, MaxHP: 150, PathIdx: 2.5}
	g.Update(0.05)
	if g.BossKills != 1 {
		t.Errorf("splash boss kill should count, BossKills = %d", g.BossKills)
	}
}

// ============================================================
// V5 Phase 4: 状态效果系统 (slow) + Frost 塔
// ============================================================

func TestSlow_SpecOnlyFrost(t *testing.T) {
	for _, k := range []TowerKind{TArcher, TCannon, TMagic} {
		for i, lv := range towerSpecs[k].Levels {
			if lv.Slow != 0 {
				t.Errorf("%s lvl%d should have no slow, got %v",
					towerSpecs[k].Name, i+1, lv.Slow)
			}
		}
	}
	for i, lv := range towerSpecs[TFrost].Levels {
		if lv.Slow <= 0 || lv.Slow >= 1 {
			t.Errorf("Frost lvl%d slow should be in (0,1), got %v", i+1, lv.Slow)
		}
	}
	if len(TowerKinds()) != 4 {
		t.Errorf("TowerKinds should include Frost, got %d", len(TowerKinds()))
	}
}

func TestSlow_EffectiveSpeed(t *testing.T) {
	e := &Enemy{Kind: ENormal, HP: 20, MaxHP: 20}
	base := enemySpecs[ENormal].Speed
	if e.EffectiveSpeed() != base {
		t.Errorf("no slow: speed = %v, want %v", e.EffectiveSpeed(), base)
	}
	e.ApplySlow(0.5)
	if e.EffectiveSpeed() != base*0.5 {
		t.Errorf("slowed: speed = %v, want %v", e.EffectiveSpeed(), base*0.5)
	}
}

func TestSlow_ExpiresAndRecovers(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1
	e := &Enemy{Kind: ENormal, HP: 20, MaxHP: 20, PathIdx: 0}
	e.ApplySlow(0.5)
	g.Enemies = []*Enemy{e}
	base := enemySpecs[ENormal].Speed
	g.Update(0.1) // 减速期: 移动 0.5×base×0.1
	wantIdx := base * 0.5 * 0.1
	if e.PathIdx != wantIdx {
		t.Errorf("slowed move: idx = %v, want %v", e.PathIdx, wantIdx)
	}
	// 推完剩余 timer (1.5s) → 恢复
	for i := 0; i < 20; i++ {
		g.Update(0.1)
	}
	if e.SlowTimer > 0 {
		t.Fatalf("slow should expire, timer = %v", e.SlowTimer)
	}
	if e.EffectiveSpeed() != base {
		t.Errorf("speed should recover to %v, got %v", base, e.EffectiveSpeed())
	}
}

func TestSlow_NoStackTakesStrongest(t *testing.T) {
	e := &Enemy{Kind: ENormal, HP: 20, MaxHP: 20}
	e.ApplySlow(0.4) // 强
	e.SlowTimer = 0.5
	e.ApplySlow(0.6) // 弱: 不应覆盖系数, 但刷新时间
	if e.SlowFactor != 0.4 {
		t.Errorf("weaker slow must not override, factor = %v", e.SlowFactor)
	}
	if e.SlowTimer != slowDurationS {
		t.Errorf("any hit should refresh timer, got %v", e.SlowTimer)
	}
	e.ApplySlow(0.3) // 更强: 覆盖
	if e.SlowFactor != 0.3 {
		t.Errorf("stronger slow should override, factor = %v", e.SlowFactor)
	}
}

func TestSlow_ExpiredThenWeakerApplies(t *testing.T) {
	// 过期后弱减速也能生效 (不被残留旧系数挡住)
	e := &Enemy{Kind: ENormal, HP: 20, MaxHP: 20}
	e.ApplySlow(0.3)
	e.SlowTimer = 0 // 模拟过期
	e.ApplySlow(0.6)
	if e.SlowFactor != 0.6 {
		t.Errorf("after expiry weaker slow should apply, factor = %v", e.SlowFactor)
	}
}

func TestSlow_FrostShotApplies(t *testing.T) {
	// 集成: Frost 塔命中施加减速 + 正常掉血
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1
	g.Towers = []*Tower{{Pos: Point{1, 1}, Kind: TFrost, Level: 1}}
	e := &Enemy{Kind: EBoss, HP: 150, MaxHP: 150, PathIdx: 1}
	g.Enemies = []*Enemy{e}
	g.Update(0.05)
	if e.SlowTimer <= 0 || e.SlowFactor != towerSpecs[TFrost].Levels[0].Slow {
		t.Errorf("frost hit should slow: timer=%v factor=%v", e.SlowTimer, e.SlowFactor)
	}
	if e.HP != 150-towerSpecs[TFrost].Levels[0].Damage {
		t.Errorf("frost hit should damage, HP = %d", e.HP)
	}
	if !drainHas(g.DrainSounds(), SndShootFrost) {
		t.Errorf("frost shot should push SndShootFrost")
	}
}

func TestSlow_ArcherDoesNotSlow(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1
	g.Towers = []*Tower{{Pos: Point{1, 1}, Kind: TArcher, Level: 1}}
	e := &Enemy{Kind: EBoss, HP: 150, MaxHP: 150, PathIdx: 1}
	g.Enemies = []*Enemy{e}
	g.Update(0.05)
	if e.SlowTimer > 0 {
		t.Errorf("archer must not slow")
	}
}

// ============================================================
// V5 Phase 5: 陨石雨主动技能 (CastMeteor + 冷却)
// ============================================================

func TestMeteor_DamagesAllInRadius(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1
	g.Enemies = []*Enemy{
		{Kind: EBoss, HP: 150, MaxHP: 150, PathIdx: 2},   // (2,0) 圈内
		{Kind: EGlider, HP: 100, MaxHP: 100, PathIdx: 3}, // (3,0) 圈内 — 飞行也炸
		{Kind: ENormal, HP: 100, MaxHP: 100, PathIdx: 5}, // (5,0) 距 (2,0)=3 圈外
	}
	if !g.CastMeteor(Point{2, 0}) {
		t.Fatalf("cast should succeed when ready")
	}
	if g.Enemies[0].HP != 150-meteorDamage {
		t.Errorf("boss in radius HP = %d, want %d", g.Enemies[0].HP, 150-meteorDamage)
	}
	if g.Enemies[1].HP != 100-meteorDamage {
		t.Errorf("flying in radius must be hit (meteor falls from sky), HP = %d", g.Enemies[1].HP)
	}
	if g.Enemies[2].HP != 100 {
		t.Errorf("out of radius should be untouched, HP = %d", g.Enemies[2].HP)
	}
	if !drainHas(g.DrainSounds(), SndMeteor) {
		t.Errorf("cast should push SndMeteor")
	}
	if g.Msg == "" {
		t.Errorf("cast should set Msg")
	}
}

func TestMeteor_CooldownGating(t *testing.T) {
	g := newTestGame()
	g.prepTimer = 0
	if !g.CastMeteor(Point{2, 0}) {
		t.Fatalf("first cast should succeed")
	}
	if g.MeteorCD != meteorCooldownS {
		t.Fatalf("cast should start cooldown, CD = %v", g.MeteorCD)
	}
	if g.CastMeteor(Point{2, 0}) {
		t.Fatalf("second cast during cooldown must fail")
	}
	// 冷却随 Update 衰减, 冷却完可再放。
	// fixture 注意: 每帧重置存活敌人, 防止 wave 清场 → PhaseWon →
	// Update 早退冷却停摆 (顺带避免 Won 路径写真实存档)
	for i := 0; i < int(meteorCooldownS/0.1)+2; i++ {
		g.Enemies = []*Enemy{{Kind: ENormal, HP: 999, MaxHP: 999, PathIdx: 0}}
		g.spawned = 1
		g.Update(0.1)
	}
	if g.MeteorCD != 0 {
		t.Fatalf("cooldown should reach 0, got %v", g.MeteorCD)
	}
	if !g.CastMeteor(Point{2, 0}) {
		t.Errorf("cast after cooldown should succeed")
	}
}

func TestMeteor_KillGoesUnifiedPath(t *testing.T) {
	// 家族再验证: 陨石杀 Spawner 必须召唤 (damageEnemy → killEnemy 路径)
	g := newTestGame()
	g.prepTimer = 0
	g.spawned = 1
	g.Enemies = []*Enemy{{Kind: ESpawner, HP: 1, MaxHP: 35, PathIdx: 2}}
	goldBefore := g.Gold
	g.CastMeteor(Point{2, 0})
	if !g.Enemies[0].Dead {
		t.Fatalf("spawner should die from meteor")
	}
	normals := 0
	for _, e := range g.Enemies {
		if e.Kind == ENormal {
			normals++
		}
	}
	if normals != 2 {
		t.Errorf("meteor-killed spawner should summon 2 normals, got %d", normals)
	}
	if g.Gold != goldBefore+enemySpecs[ESpawner].Reward {
		t.Errorf("meteor kill should pay reward")
	}
}

func TestMeteor_PhaseGating(t *testing.T) {
	g := newTestGame()
	g.Phase = PhaseLevelSelect
	if g.CastMeteor(Point{2, 0}) {
		t.Errorf("cast outside playing phase must fail")
	}
}

func TestMeteor_LevelStartResets(t *testing.T) {
	g := newTestGame()
	g.MeteorCD = 10
	g.StartLevel(0)
	if g.MeteorCD != 0 {
		t.Errorf("level start should reset meteor cooldown, got %v", g.MeteorCD)
	}
}

// ============================================================
// V6 Phase 1: 关卡 11-20 数据完整性 + 菜单两列 hitbox
// ============================================================

func TestLevels_IntegrityAll(t *testing.T) {
	levels, err := LoadLevels()
	if err != nil {
		t.Fatalf("LoadLevels failed: %v", err)
	}
	if len(levels) != 20 {
		t.Fatalf("want 20 levels, got %d", len(levels))
	}
	for i, lv := range levels {
		// unlock 链连续 (IsUnlocked 依赖 ID = 前一关 ID+1)
		if lv.ID != i+1 {
			t.Errorf("level[%d] ID = %d, unlock 链要求连续", i, lv.ID)
		}
		if len(lv.Path) < 2 {
			t.Errorf("Lv%d path too short: %d", lv.ID, len(lv.Path))
		}
		// 起终点约定: x=0 进场, x=29 离场 (全 20 关一致)
		if lv.Path[0].X != 0 {
			t.Errorf("Lv%d path should start at x=0, got %v", lv.ID, lv.Path[0])
		}
		if lv.Path[len(lv.Path)-1].X != mapW-1 {
			t.Errorf("Lv%d path should end at x=%d, got %v", lv.ID, mapW-1, lv.Path[len(lv.Path)-1])
		}
		// 全 path cell 在地图内
		for _, p := range lv.Path {
			if p.X < 0 || p.X >= mapW || p.Y < 0 || p.Y >= mapH {
				t.Errorf("Lv%d path cell %v out of map %dx%d", lv.ID, p, mapW, mapH)
			}
		}
		if len(lv.Waves) < 3 {
			t.Errorf("Lv%d should have ≥3 waves, got %d", lv.ID, len(lv.Waves))
		}
		for w, wave := range lv.Waves {
			if len(wave.Enemies) == 0 {
				t.Errorf("Lv%d wave %d empty", lv.ID, w+1)
			}
		}
		if lv.StartGold < 100 || lv.StartLives < 3 {
			t.Errorf("Lv%d 起始资源异常: gold %d lives %d", lv.ID, lv.StartGold, lv.StartLives)
		}
	}
}

func TestLevels_DifficultyRampsUp(t *testing.T) {
	// 粗粒度难度单调性: 后 10 关总敌量 ≥ 前 10 关均值 (防 yaml 抄错)
	levels, err := LoadLevels()
	if err != nil {
		t.Fatal(err)
	}
	// 实战敌量加权: Spawner 死后召唤 2 normals → 计 3
	total := func(lv Level) int {
		n := 0
		for _, w := range lv.Waves {
			for _, k := range w.Enemies {
				if k == ESpawner {
					n += 3
				} else {
					n++
				}
			}
		}
		return n
	}
	firstHalfSum := 0
	for _, lv := range levels[:10] {
		firstHalfSum += total(lv)
	}
	avg := firstHalfSum / 10
	for _, lv := range levels[10:] {
		if total(lv) < avg {
			t.Errorf("Lv%d 总敌量 %d < 前 10 关均值 %d (难度曲线倒挂?)", lv.ID, total(lv), avg)
		}
	}
}

func TestMenu_RowAtPixelTwoColumns(t *testing.T) {
	const startY, rowH = 60, 36
	cases := []struct {
		mx, my int
		want   int
		wantOK bool
	}{
		{10, startY + 5, 0, true},                     // 左列首行
		{10, startY + 9*rowH + 5, 9, true},            // 左列末行
		{windowW - 10, startY + 5, 10, true},          // 右列首行 = Lv11
		{windowW - 10, startY + 9*rowH + 5, 19, true}, // 右列末行 = Lv20
		{10, startY - 5, 0, false},                    // 上方界外
		{10, startY + 10*rowH + 5, 0, false},          // 下方界外
	}
	for _, c := range cases {
		got, ok := menuRowAtPixel(c.mx, c.my, 20)
		if ok != c.wantOK || (ok && got != c.want) {
			t.Errorf("menuRowAtPixel(%d,%d) = (%d,%v), want (%d,%v)",
				c.mx, c.my, got, ok, c.want, c.wantOK)
		}
	}
	// 10 关存档兼容: 右列点击不应命中 (numLevels=10)
	if _, ok := menuRowAtPixel(windowW-10, startY+5, 10); ok {
		t.Errorf("右列在只有 10 关时不应命中")
	}
}
