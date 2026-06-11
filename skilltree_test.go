// V11 P1 测试: 技能树持久层 — 节点表守护 / 预算算术 / 购买 gating / 零值兼容。
package main

import "testing"

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
