package main

import "testing"

func TestAchievement_UnlockIdempotent(t *testing.T) {
	s := NewSave()
	if s.HasAchievement("first_blood") {
		t.Fatal("new save should have no achievements")
	}
	if !s.UnlockAchievement("first_blood") {
		t.Error("first unlock should return true")
	}
	if !s.HasAchievement("first_blood") {
		t.Error("should have achievement after unlock")
	}
	if s.UnlockAchievement("first_blood") {
		t.Error("duplicate unlock should return false")
	}
}

func TestAchievement_Count(t *testing.T) {
	s := NewSave()
	if s.AchievementCount() != 0 {
		t.Errorf("new save count = %d, want 0", s.AchievementCount())
	}
	s.UnlockAchievement("first_blood")
	s.UnlockAchievement("hunter_100")
	if s.AchievementCount() != 2 {
		t.Errorf("count = %d, want 2", s.AchievementCount())
	}
}

func TestAchievement_CheckProgressionUnlocks(t *testing.T) {
	s := NewSave()
	s.MarkCompleted(1)
	got := s.CheckAchievements()
	found := false
	for _, id := range got {
		if id == "first_blood" {
			found = true
		}
	}
	if !found {
		t.Errorf("completing level 1 should unlock first_blood, got %v", got)
	}
}

func TestAchievement_KillMilestones(t *testing.T) {
	s := NewSave()
	s.TotalKills = 99
	got := s.CheckAchievements()
	for _, id := range got {
		if id == "hunter_100" {
			t.Error("99 kills should not unlock hunter_100")
		}
	}
	s.TotalKills = 100
	got = s.CheckAchievements()
	found := false
	for _, id := range got {
		if id == "hunter_100" {
			found = true
		}
	}
	if !found {
		t.Error("100 kills should unlock hunter_100")
	}
}

func TestAchievement_FlawlessAndHard(t *testing.T) {
	s := NewSave()
	s.HasFlawless = true
	s.HasHardClear = true
	got := s.CheckAchievements()
	hasFlawless, hasHard := false, false
	for _, id := range got {
		if id == "flawless" {
			hasFlawless = true
		}
		if id == "hard_clear" {
			hasHard = true
		}
	}
	if !hasFlawless {
		t.Error("should unlock flawless")
	}
	if !hasHard {
		t.Error("should unlock hard_clear")
	}
}

func TestAchievement_EndlessMilestones(t *testing.T) {
	s := NewSave()
	s.BestWave = 25
	got := s.CheckAchievements()
	has10, has25, has50 := false, false, false
	for _, id := range got {
		switch id {
		case "endless_10":
			has10 = true
		case "endless_25":
			has25 = true
		case "endless_50":
			has50 = true
		}
	}
	if !has10 || !has25 {
		t.Error("wave 25 should unlock endless_10 and endless_25")
	}
	if has50 {
		t.Error("wave 25 should not unlock endless_50")
	}
}

func TestAchievement_AllDefinitionsHaveUniqueIDs(t *testing.T) {
	seen := map[string]bool{}
	for _, a := range achievements {
		if seen[a.ID] {
			t.Errorf("duplicate achievement ID: %s", a.ID)
		}
		seen[a.ID] = true
		if a.Name == "" || a.Desc == "" {
			t.Errorf("achievement %s missing name or desc", a.ID)
		}
	}
}

func TestAchievement_CompletedCount(t *testing.T) {
	s := NewSave()
	if s.CompletedCount() != 0 {
		t.Error("new save should have 0 completed")
	}
	s.MarkCompleted(1)
	s.MarkCompleted(5)
	if s.CompletedCount() != 2 {
		t.Errorf("completed = %d, want 2", s.CompletedCount())
	}
}
