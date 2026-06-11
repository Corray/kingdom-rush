// V11 P1 测试: 技能树持久层 — 节点表守护 / 预算算术 / 购买 gating / 零值兼容。
package main

import (
	"math"
	"testing"
)

// TestSkillTreeTable: 节点表守护 — 三树与职业表一一对应, 每树 4 节点,
// 定价 3/6/9/12 (单树 30), 三树总价 90 = 1.5× 可赚上限 60 (决策 C)。
func TestSkillTreeTable(t *testing.T) {
	if len(skillTrees) != len(heroClasses) {
		t.Fatalf("skill trees (%d) must match hero classes (%d)", len(skillTrees), len(heroClasses))
	}
	wantPrices := [treeNodesPerClass]int{3, 6, 9, 12}
	grand := 0
	for i := range heroClasses {
		name := heroClasses[i].Name
		tree, ok := skillTrees[name]
		if !ok {
			t.Fatalf("no skill tree for class %s", name)
		}
		for j, node := range tree {
			if node.Name == "" || node.Desc == "" {
				t.Errorf("%s node %d: empty name/desc", name, j)
			}
			if node.Price != wantPrices[j] {
				t.Errorf("%s node %d price = %d, want %d", name, j, node.Price, wantPrices[j])
			}
			grand += node.Price
		}
	}
	if grand != 90 {
		t.Errorf("grand total = %d, want 90 (1.5x earnable 60)", grand)
	}
}

// TestTreeLevelZeroValue: 旧存档/零值 → 全 0; 手改越界/负值 clamp。
func TestTreeLevelZeroValue(t *testing.T) {
	s := NewSave()
	for i := range heroClasses {
		if lvl := s.TreeLevel(heroClasses[i].Name); lvl != 0 {
			t.Errorf("zero-value save: %s tree level = %d, want 0", heroClasses[i].Name, lvl)
		}
	}
	if s.SpentStars() != 0 || s.AvailableStars() != s.TotalStars() {
		t.Errorf("zero-value save must have 0 spent: spent=%d", s.SpentStars())
	}
	s.TreeNodes = map[string]int{"Knight": 99, "Archer": -5}
	if s.TreeLevel("Knight") != treeNodesPerClass {
		t.Errorf("overflow must clamp to %d, got %d", treeNodesPerClass, s.TreeLevel("Knight"))
	}
	if s.TreeLevel("Archer") != 0 {
		t.Errorf("negative must clamp to 0, got %d", s.TreeLevel("Archer"))
	}
}

// TestStarsBudget: 预算算术 — 已赚 = Stars 求和, 已花 = 已购价格和, 余额 = 差。
func TestStarsBudget(t *testing.T) {
	s := NewSave()
	s.RecordStars(1, 3)
	s.RecordStars(2, 2)
	if s.TotalStars() != 5 {
		t.Fatalf("TotalStars = %d, want 5", s.TotalStars())
	}
	ok, name := s.BuyTreeNode("Knight") // 价 3, 余 5 → 成功
	if !ok || name != "Bulwark" {
		t.Fatalf("first buy should succeed with Bulwark, got ok=%v name=%q", ok, name)
	}
	if s.SpentStars() != 3 || s.AvailableStars() != 2 {
		t.Errorf("after buy: spent=%d avail=%d, want 3/2", s.SpentStars(), s.AvailableStars())
	}
	ok, _ = s.BuyTreeNode("Knight") // 下一节点价 6 > 余 2 → 拒绝
	if ok {
		t.Error("buy must be rejected when balance insufficient")
	}
	if s.TreeLevel("Knight") != 1 || s.SpentStars() != 3 {
		t.Errorf("rejected buy must not mutate: level=%d spent=%d", s.TreeLevel("Knight"), s.SpentStars())
	}
}

// TestBuyTreeNodeGating: 未知职业拒绝; 预算充足时连买到满后拒绝。
func TestBuyTreeNodeGating(t *testing.T) {
	s := NewSave()
	for lv := 1; lv <= 20; lv++ {
		s.RecordStars(lv, 3) // 60 星满预算
	}
	if ok, _ := s.BuyTreeNode("Paladin"); ok {
		t.Error("unknown class must be rejected")
	}
	for i := 0; i < treeNodesPerClass; i++ {
		if ok, _ := s.BuyTreeNode("Rogue"); !ok {
			t.Fatalf("buy %d should succeed with full budget", i+1)
		}
	}
	if s.TreeLevel("Rogue") != treeNodesPerClass {
		t.Fatalf("Rogue tree should be maxed, got %d", s.TreeLevel("Rogue"))
	}
	if ok, _ := s.BuyTreeNode("Rogue"); ok {
		t.Error("maxed tree must reject further buys")
	}
	if s.SpentStars() != 30 || s.AvailableStars() != 30 {
		t.Errorf("full Rogue tree: spent=%d avail=%d, want 30/30", s.SpentStars(), s.AvailableStars())
	}
	// 余 30 不够再满第二棵 (30) 之外的第三棵 → 取舍实锤: 买满 Knight 后归零
	for i := 0; i < treeNodesPerClass; i++ {
		if ok, _ := s.BuyTreeNode("Knight"); !ok {
			t.Fatalf("Knight buy %d should succeed (exactly 30 left)", i+1)
		}
	}
	if s.AvailableStars() != 0 {
		t.Errorf("two full trees must drain 60 stars exactly, avail=%d", s.AvailableStars())
	}
	if ok, _ := s.BuyTreeNode("Archer"); ok {
		t.Error("third tree must be unaffordable (决策 C 取舍)")
	}
}

// TestSaveRoundtripTreeNodes: TreeNodes 经 JSON 存取往返 (守 json tag)。
func TestSaveRoundtripTreeNodes(t *testing.T) {
	withTempSavePath(t, func() {
		s := NewSave()
		s.TreeNodes = map[string]int{"Knight": 2, "Rogue": 1}
		if err := StoreSave(s); err != nil {
			t.Fatal(err)
		}
		loaded, err := LoadSave()
		if err != nil {
			t.Fatal(err)
		}
		if loaded.TreeLevel("Knight") != 2 || loaded.TreeLevel("Rogue") != 1 || loaded.TreeLevel("Archer") != 0 {
			t.Errorf("roundtrip tree levels = K%d/R%d/A%d, want 2/1/0",
				loaded.TreeLevel("Knight"), loaded.TreeLevel("Rogue"), loaded.TreeLevel("Archer"))
		}
	})
}

// ============================================================
// V11 Phase 2: HeroBonus 效果接线
// ============================================================

// TestTreeNodesHaveEffects: 节点表守护 — 每个节点 Bonus 至少一个非零字段
// (防"有价无效"的空 perk)。
func TestTreeNodesHaveEffects(t *testing.T) {
	for name, tree := range skillTrees {
		for j, node := range tree {
			if node.Bonus == (HeroBonus{}) {
				t.Errorf("%s node %d (%s): zero-effect perk", name, j, node.Name)
			}
		}
	}
}

// TestHeroBonusAggregation: 聚合 = 前 N 节点效果求和; 零购买/未知职业 = 零值。
func TestHeroBonusAggregation(t *testing.T) {
	s := NewSave()
	if s.HeroBonusFor("Knight") != (HeroBonus{}) {
		t.Error("zero purchases must aggregate to zero bonus")
	}
	if s.HeroBonusFor("Paladin") != (HeroBonus{}) {
		t.Error("unknown class must aggregate to zero bonus")
	}
	s.TreeNodes = map[string]int{"Knight": 2}
	got := s.HeroBonusFor("Knight")
	want := HeroBonus{MaxHP: 30, Damage: 4} // Bulwark + Sharpened Blade
	if got != want {
		t.Errorf("Knight lvl2 bonus = %+v, want %+v", got, want)
	}
	s.TreeNodes["Rogue"] = 4
	gr := s.HeroBonusFor("Rogue")
	wantR := HeroBonus{Speed: 0.7, Damage: 2, AbilityCDReduceS: 2, RespawnReduceS: 4}
	if gr != wantR {
		t.Errorf("Rogue full bonus = %+v, want %+v", gr, wantR)
	}
}

// TestHeroBonusStats: perk 反映到全部派生数值 + clamp 防护。
func TestHeroBonusStats(t *testing.T) {
	b := HeroBonus{MaxHP: 30, Damage: 4, Range: 0.6, Speed: 0.8,
		RespawnReduceS: 4, AbilityRadius: 0.6, AbilityCDReduceS: 2}
	h := newHeroWithBonus(knight, b, 0, 0)
	if h.MaxHP != knight.MaxHP+30 || h.HP != h.MaxHP {
		t.Errorf("spawn HP = %d/%d, want %d full", h.HP, h.MaxHP, knight.MaxHP+30)
	}
	if h.Damage() != knight.Damage+4 {
		t.Errorf("damage = %d, want %d", h.Damage(), knight.Damage+4)
	}
	if math.Abs(h.AttackRange()-(knight.Range+0.6)) > 1e-9 {
		t.Errorf("range = %v, want %v", h.AttackRange(), knight.Range+0.6)
	}
	if math.Abs(h.MoveSpeed()-(knight.Speed+0.8)) > 1e-9 {
		t.Errorf("speed = %v, want %v", h.MoveSpeed(), knight.Speed+0.8)
	}
	if math.Abs(h.RespawnTime()-(knight.RespawnS-4)) > 1e-9 {
		t.Errorf("respawn = %v, want %v", h.RespawnTime(), knight.RespawnS-4)
	}
	if math.Abs(h.AbilityCooldown()-(knight.AbilityCooldownS-2)) > 1e-9 {
		t.Errorf("ability CD = %v, want %v", h.AbilityCooldown(), knight.AbilityCooldownS-2)
	}
	if math.Abs(h.AbilityRange()-(knight.AbilityRadius+0.6)) > 1e-9 {
		t.Errorf("ability radius = %v, want %v", h.AbilityRange(), knight.AbilityRadius+0.6)
	}
	// clamp: 巨幅缩减不得把时长压到 1s 以下
	hc := newHeroWithBonus(knight, HeroBonus{RespawnReduceS: 99, AbilityCDReduceS: 99}, 0, 0)
	if hc.RespawnTime() != 1 || hc.AbilityCooldown() != 1 {
		t.Errorf("clamp: respawn=%v abilityCD=%v, want 1/1", hc.RespawnTime(), hc.AbilityCooldown())
	}
}

// TestZeroBonusZeroRegression: 零 perk 英雄派生数值 = V10 职业基线 (硬约束守护)。
func TestZeroBonusZeroRegression(t *testing.T) {
	for i := range heroClasses {
		c := &heroClasses[i]
		h := newHeroOf(c, 0, 0)
		if h.MaxHP != c.maxHPFor(1) || h.Damage() != c.Damage ||
			math.Abs(h.AttackRange()-c.Range) > 1e-9 ||
			math.Abs(h.MoveSpeed()-c.Speed) > 1e-9 ||
			math.Abs(h.RespawnTime()-c.RespawnS) > 1e-9 ||
			math.Abs(h.AbilityCooldown()-c.AbilityCooldownS) > 1e-9 ||
			math.Abs(h.AbilityRange()-c.AbilityRadius) > 1e-9 {
			t.Errorf("%s: zero-bonus hero deviates from class baseline", c.Name)
		}
	}
}

// TestBeginRunAppliesBonus: beginRun 按存档快照 perk; 升级后 perk 仍叠加。
func TestBeginRunAppliesBonus(t *testing.T) {
	g := newTestGame()
	g.Save.TreeNodes = map[string]int{"Knight": 1} // Bulwark +30 HP
	g.StartLevel(0)
	if g.Hero.MaxHP != knight.MaxHP+30 {
		t.Fatalf("beginRun should snapshot perk: maxHP=%d, want %d", g.Hero.MaxHP, knight.MaxHP+30)
	}
	g.Hero.GainXP(knight.xpForNext(1)) // → lvl2
	want := knight.maxHPFor(2) + 30
	if g.Hero.MaxHP != want {
		t.Errorf("perk must persist through level-up: maxHP=%d, want %d", g.Hero.MaxHP, want)
	}
}
