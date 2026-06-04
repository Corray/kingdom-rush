// 单元测试覆盖纯逻辑层:
//   - ExpandPath / ParseWave (level.go)
//   - Tower upgrade spec invariants (entities.go)
//   - TryAction (build / not enough gold / on-path / upgrade) (game.go)
//   - Update lives-dec on escape (game.go)
//
// 不覆盖: render layer (tcell screen 需 TTY) / main 事件循环 (集成测试)。
package main

import "testing"

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
	// shoot + hit = 2 effects
	if len(g.Effects) != 2 {
		t.Fatalf("expected 2 effects (shoot+hit), got %d", len(g.Effects))
	}
	kinds := map[EffectKind]int{}
	for _, e := range g.Effects {
		kinds[e.Kind]++
	}
	if kinds[EShoot] != 1 || kinds[EHit] != 1 {
		t.Errorf("expected 1 EShoot + 1 EHit, got %v", kinds)
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
