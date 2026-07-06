// V15: 成就系统 — 定义 + 检查逻辑 (纯数据, 不依赖 UI)。
package main

type Achievement struct {
	ID   string
	Name string
	Desc string
}

var achievements = []Achievement{
	// 进度 (V17: conqueror 条件保持 20 关不动 — 防打破已解锁语义;
	// 30 关全通走新增 second_dawn)
	{ID: "first_blood", Name: "First Blood", Desc: "Clear level 1"},
	{ID: "halfway", Name: "Halfway There", Desc: "Clear 10 levels"},
	{ID: "conqueror", Name: "Conqueror", Desc: "Clear 20 levels"},
	{ID: "second_dawn", Name: "Second Dawn", Desc: "Clear all 30 levels"},
	{ID: "perfectionist", Name: "Perfectionist", Desc: "Earn 60 stars"},
	{ID: "star_master", Name: "Star Master", Desc: "Earn 90 stars (all 3-star)"},
	// 战斗
	{ID: "hunter_100", Name: "Hunter", Desc: "Kill 100 enemies (lifetime)"},
	{ID: "hunter_500", Name: "Veteran", Desc: "Kill 500 enemies (lifetime)"},
	{ID: "hunter_2000", Name: "Legend", Desc: "Kill 2000 enemies (lifetime)"},
	{ID: "boss_slayer", Name: "Boss Slayer", Desc: "Kill 10 bosses (lifetime)"},
	// 英雄
	{ID: "all_classes", Name: "Jack of All Trades", Desc: "Win with all 3 hero classes"},
	{ID: "max_tree", Name: "Specialist", Desc: "Max out a skill tree (4 nodes)"},
	// 策略
	{ID: "flawless", Name: "Flawless", Desc: "Clear any level without losing a life"},
	{ID: "hard_clear", Name: "Iron Will", Desc: "Clear any level on Hard"},
	{ID: "hard_all", Name: "Unstoppable", Desc: "Clear the first 20 levels on Hard"},
	// Endless
	{ID: "endless_10", Name: "Survivor", Desc: "Reach wave 10 in Endless"},
	{ID: "endless_25", Name: "Endurance", Desc: "Reach wave 25 in Endless"},
	{ID: "endless_50", Name: "Immortal", Desc: "Reach wave 50 in Endless"},
}

var achievementByID = func() map[string]*Achievement {
	m := make(map[string]*Achievement, len(achievements))
	for i := range achievements {
		m[achievements[i].ID] = &achievements[i]
	}
	return m
}()

func (s *Save) HasAchievement(id string) bool {
	return s.Achievements[id]
}

func (s *Save) UnlockAchievement(id string) bool {
	if s.Achievements == nil {
		s.Achievements = map[string]bool{}
	}
	if s.Achievements[id] {
		return false
	}
	s.Achievements[id] = true
	return true
}

func (s *Save) AchievementCount() int {
	n := 0
	for _, v := range s.Achievements {
		if v {
			n++
		}
	}
	return n
}

func (s *Save) CompletedCount() int {
	n := 0
	for _, v := range s.Completed {
		if v {
			n++
		}
	}
	return n
}

func (s *Save) ClassWins() map[int]bool {
	return s.WonClasses
}

// CheckAchievements evaluates all unlock conditions, returns newly unlocked IDs.
func (s *Save) CheckAchievements() []string {
	var unlocked []string
	check := func(id string, cond bool) {
		if cond && s.UnlockAchievement(id) {
			unlocked = append(unlocked, id)
		}
	}

	completed := s.CompletedCount()
	check("first_blood", completed >= 1)
	check("halfway", completed >= 10)
	check("conqueror", completed >= 20)
	check("second_dawn", completed >= 30)
	check("perfectionist", s.TotalStars() >= 60)
	check("star_master", s.TotalStars() >= 90)

	check("hunter_100", s.TotalKills >= 100)
	check("hunter_500", s.TotalKills >= 500)
	check("hunter_2000", s.TotalKills >= 2000)
	check("boss_slayer", s.BossKillsTotal >= 10)

	if len(s.WonClasses) >= 3 {
		check("all_classes", true)
	}
	for _, lvl := range s.TreeNodes {
		if lvl >= treeNodesPerClass {
			check("max_tree", true)
			break
		}
	}

	check("flawless", s.HasFlawless)
	check("hard_clear", s.HasHardClear)

	hardAll := true
	if s.HardCompleted == nil {
		hardAll = false
	} else {
		for i := 1; i <= 20; i++ {
			if !s.HardCompleted[i] {
				hardAll = false
				break
			}
		}
	}
	check("hard_all", hardAll)

	check("endless_10", s.BestWave >= 10)
	check("endless_25", s.BestWave >= 25)
	check("endless_50", s.BestWave >= 50)

	return unlocked
}
