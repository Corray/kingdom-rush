// V11 P1: 技能树 — 跨局 meta 成长持久层 (星 → per-class perk 节点)。
//
// 设计 (roadmap V11 决策 A-E):
//   - 货币 = 通关星 (Save.Stars, 20 关 × 3★ = 60 可赚上限), 老存档存量星直接可用
//   - 每职业线性 4 节点: 前一节点已购才能买下一个 (线性树 → "已购数"即完整状态)
//   - 定价 3/6/9/12 (单树满 30 星), 三树总价 90 ≈ 1.5× 预算 → 必须选 build
//   - 购买只在菜单发生 (P3 入口); 效果经 HeroBonus 在 beginRun 快照进 Hero (P2)
//   - 无树基线零回归: 不买节点 = V10 行为完全不变
package main

// TreeNode: 技能树节点。Bonus 效果字段 P2 接线。
type TreeNode struct {
	Name  string
	Desc  string
	Price int
}

const treeNodesPerClass = 4

// skillTrees: per-class 线性节点表 (key = HeroClass.Name, 与 Save.TreeNodes
// 的 key 一致)。效果数值 P2 定义 + P4 仿真校准, Desc 先行锁定方向。
var skillTrees = map[string][treeNodesPerClass]TreeNode{
	"Knight": {
		{Name: "Bulwark", Desc: "+max HP", Price: 3},
		{Name: "Sharpened Blade", Desc: "+damage", Price: 6},
		{Name: "Wide Cleave", Desc: "+cleave radius", Price: 9},
		{Name: "Undying", Desc: "faster respawn", Price: 12},
	},
	"Archer": {
		{Name: "Eagle Eye", Desc: "+attack range", Price: 3},
		{Name: "Fleet Foot", Desc: "+move speed", Price: 6},
		{Name: "Piercing Arrows", Desc: "+damage", Price: 9},
		{Name: "Storm of Arrows", Desc: "+volley radius", Price: 12},
	},
	"Rogue": {
		{Name: "Shadow Step", Desc: "+move speed", Price: 3},
		{Name: "Twin Fangs", Desc: "+damage", Price: 6},
		{Name: "Blade Flurry", Desc: "faster fan of blades", Price: 9},
		{Name: "Cheat Death", Desc: "faster respawn", Price: 12},
	},
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
