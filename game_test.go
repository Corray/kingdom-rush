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
// helpers
// ============================================================

func newTestGame() *Game {
	levels := []Level{{
		ID: 1, Name: "Test",
		StartGold: 100, StartLives: 5,
		Path:  []Point{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}},
		Waves: []WaveSpec{{Enemies: []EnemyKind{ENormal}}},
	}}
	g := NewGame(levels)
	g.StartLevel(0)
	return g
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
