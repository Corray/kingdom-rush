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
	levels, err := LoadLevels()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(levels) != 10 {
		t.Errorf("expected 10 levels, got %d", len(levels))
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
