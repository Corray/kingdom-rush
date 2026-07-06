// V11 P1: 技能树 — 跨局 meta 成长持久层 (星 → per-class perk 节点)。
//
// 设计 (roadmap V11 决策 A-E; V17 决策 A 扩展):
//   - 货币 = 通关星 (Save.Stars; V17: 30 关 × 3★ = 90 可赚上限), 老存档存量星直接可用
//   - 每职业线性 5 节点 (V11 4 节点 + V17 capstone): 前一节点已购才能买下一个
//   - 定价 3/6/9/12/10 (单树满 40 星), 三树总价 120 ≈ 1.33× 预算 → 必须选 build
//     (V17 决策 A: 30 关后可赚 90 追平 V11 三树总价, 加 capstone 恢复取舍)
//   - 购买只在菜单发生 (P3 入口); 效果经 HeroBonus 在 beginRun 快照进 Hero (P2)
//   - 无树基线零回归: 不买节点 = V10 行为完全不变
package main

// HeroBonus: perk 增量聚合 (V11 P2)。零值 = 无 perk = V10 行为。
// beginRun 时按职业已购节点求和后快照进 Hero (关内不变)。
type HeroBonus struct {
	MaxHP            int
	Damage           int
	Range            float64
	Speed            float64
	RespawnReduceS   float64 // 复活缩短 (正值 = 减时, 应用处 clamp ≥1s)
	AbilityRadius    float64
	AbilityCDReduceS float64 // 技能冷却缩短 (正值 = 减时, 应用处 clamp ≥1s)
}

// add: 聚合两份增量。
func (b HeroBonus) add(o HeroBonus) HeroBonus {
	return HeroBonus{
		MaxHP:            b.MaxHP + o.MaxHP,
		Damage:           b.Damage + o.Damage,
		Range:            b.Range + o.Range,
		Speed:            b.Speed + o.Speed,
		RespawnReduceS:   b.RespawnReduceS + o.RespawnReduceS,
		AbilityRadius:    b.AbilityRadius + o.AbilityRadius,
		AbilityCDReduceS: b.AbilityCDReduceS + o.AbilityCDReduceS,
	}
}

// TreeNode: 技能树节点。
type TreeNode struct {
	Name  string
	Desc  string
	Price int
	Bonus HeroBonus
}

const treeNodesPerClass = 5

// skillTrees: per-class 线性节点表 (key = HeroClass.Name, 与 Save.TreeNodes
// 的 key 一致)。前 4 节点 = V11 原值不动 (老存档零回归); 第 5 节点 =
// V17 capstone (组合 perk, 价 10)。
var skillTrees = map[string][treeNodesPerClass]TreeNode{
	"Knight": {
		{Name: "Bulwark", Desc: "+30 max HP", Price: 3, Bonus: HeroBonus{MaxHP: 30}},
		{Name: "Sharpened Blade", Desc: "+6 damage", Price: 6, Bonus: HeroBonus{Damage: 6}},
		{Name: "Wide Cleave", Desc: "+0.6 cleave radius", Price: 9, Bonus: HeroBonus{AbilityRadius: 0.6}},
		{Name: "Undying", Desc: "respawn 6s faster", Price: 12, Bonus: HeroBonus{RespawnReduceS: 6}},
		{Name: "Warlord", Desc: "+20 HP & +4 damage", Price: 10, Bonus: HeroBonus{MaxHP: 20, Damage: 4}},
	},
	"Archer": {
		{Name: "Eagle Eye", Desc: "+0.6 attack range", Price: 3, Bonus: HeroBonus{Range: 0.6}},
		{Name: "Fleet Foot", Desc: "+0.8 move speed", Price: 6, Bonus: HeroBonus{Speed: 0.8}},
		{Name: "Piercing Arrows", Desc: "+3 damage", Price: 9, Bonus: HeroBonus{Damage: 3}},
		{Name: "Storm of Arrows", Desc: "+0.8 volley radius", Price: 12, Bonus: HeroBonus{AbilityRadius: 0.8}},
		{Name: "Windrunner", Desc: "+0.5 range & +2 damage", Price: 10, Bonus: HeroBonus{Range: 0.5, Damage: 2}},
	},
	"Rogue": {
		{Name: "Shadow Step", Desc: "+0.7 move speed", Price: 3, Bonus: HeroBonus{Speed: 0.7}},
		{Name: "Twin Fangs", Desc: "+2 damage", Price: 6, Bonus: HeroBonus{Damage: 2}},
		{Name: "Blade Flurry", Desc: "ability 2s faster", Price: 9, Bonus: HeroBonus{AbilityCDReduceS: 2}},
		{Name: "Cheat Death", Desc: "respawn 4s faster", Price: 12, Bonus: HeroBonus{RespawnReduceS: 4}},
		{Name: "Phantom", Desc: "+0.7 speed & ability 2s faster", Price: 10, Bonus: HeroBonus{Speed: 0.7, AbilityCDReduceS: 2}},
	},
}

// HeroBonusFor: 该职业已购节点的效果聚合 (beginRun 快照用)。
func (s *Save) HeroBonusFor(className string) HeroBonus {
	var b HeroBonus
	tree, ok := skillTrees[className]
	if !ok {
		return b
	}
	for i := 0; i < s.TreeLevel(className); i++ {
		b = b.add(tree[i].Bonus)
	}
	return b
}

// TreeLevel: 该职业已购节点数, clamp [0, treeNodesPerClass]
// (nil map / 旧存档 → 0; 手改存档越界 → 封顶)。
func (s *Save) TreeLevel(className string) int {
	n := s.TreeNodes[className]
	if n < 0 {
		return 0
	}
	if n > treeNodesPerClass {
		return treeNodesPerClass
	}
	return n
}

// TotalStars: 累计已赚星 (per-level 最高星求和)。
func (s *Save) TotalStars() int {
	total := 0
	for _, v := range s.Stars {
		total += v
	}
	return total
}

// SpentStars: 已购节点总花费 (跨三职业)。
func (s *Save) SpentStars() int {
	spent := 0
	for name, tree := range skillTrees {
		for i := 0; i < s.TreeLevel(name); i++ {
			spent += tree[i].Price
		}
	}
	return spent
}

// AvailableStars: 可花余额 = 已赚 - 已花。
func (s *Save) AvailableStars() int {
	return s.TotalStars() - s.SpentStars()
}

// NextTreeNode: 该职业下一个可购节点。已满/未知职业返回 (nil, -1)。
func (s *Save) NextTreeNode(className string) (*TreeNode, int) {
	tree, ok := skillTrees[className]
	if !ok {
		return nil, -1
	}
	lvl := s.TreeLevel(className)
	if lvl >= treeNodesPerClass {
		return nil, -1
	}
	return &tree[lvl], lvl
}

// BuyTreeNode: 购买该职业下一节点。gating: 职业存在 + 未满 + 余额够。
// 成功返回 (true, 节点名); 调用方负责 StoreSave 持久化 (P3 入口统一做)。
func (s *Save) BuyTreeNode(className string) (bool, string) {
	node, _ := s.NextTreeNode(className)
	if node == nil {
		return false, ""
	}
	if s.AvailableStars() < node.Price {
		return false, ""
	}
	if s.TreeNodes == nil {
		s.TreeNodes = map[string]int{}
	}
	s.TreeNodes[className] = s.TreeLevel(className) + 1
	return true, node.Name
}
