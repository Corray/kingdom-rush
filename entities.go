// Tower / Enemy 实体与 spec table。
// Tower 支持 3 级升级,每级独立 spec(damage/range/cd/字符)。
package main

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

// ============================================================
// Point / Wave (基础类型)
// ============================================================

type Point struct{ X, Y int }

type WaveSpec struct {
	Enemies []EnemyKind
}

// ============================================================
// Tower
// ============================================================

type TowerKind int

const (
	TArcher TowerKind = iota
	TCannon
)

type TowerLevel struct {
	Cost     int // 升到该级的额外 cost (lvl 1 = 建造 cost)
	Damage   int
	Range    float64
	Cooldown float64
	Char1    rune
	Char2    rune
}

type TowerSpec struct {
	Name   string
	Color  tcell.Color
	Levels [3]TowerLevel // index 0 = level 1
}

var towerSpecs = map[TowerKind]TowerSpec{
	TArcher: {
		Name:  "Archer",
		Color: tcell.NewRGBColor(100, 180, 255),
		Levels: [3]TowerLevel{
			{Cost: 50, Damage: 8, Range: 3.5, Cooldown: 0.6, Char1: 'A', Char2: '1'},
			{Cost: 40, Damage: 14, Range: 4.0, Cooldown: 0.5, Char1: 'A', Char2: '2'},
			{Cost: 80, Damage: 22, Range: 4.5, Cooldown: 0.4, Char1: 'A', Char2: '3'},
		},
	},
	TCannon: {
		Name:  "Cannon",
		Color: tcell.NewRGBColor(180, 120, 255),
		Levels: [3]TowerLevel{
			{Cost: 80, Damage: 25, Range: 2.5, Cooldown: 1.5, Char1: 'C', Char2: '1'},
			{Cost: 60, Damage: 45, Range: 3.0, Cooldown: 1.3, Char1: 'C', Char2: '2'},
			{Cost: 100, Damage: 70, Range: 3.5, Cooldown: 1.0, Char1: 'C', Char2: '3'},
		},
	},
}

func TowerKinds() []TowerKind { return []TowerKind{TArcher, TCannon} }

type Tower struct {
	Pos      Point
	Kind     TowerKind
	Level    int // 1-3
	cooldown float64
}

func (t *Tower) Spec() TowerLevel {
	return towerSpecs[t.Kind].Levels[t.Level-1]
}

func (t *Tower) NextUpgradeCost() (int, bool) {
	if t.Level >= 3 {
		return 0, false
	}
	return towerSpecs[t.Kind].Levels[t.Level].Cost, true
}

// ============================================================
// Enemy
// ============================================================

type EnemyKind int

const (
	ENormal EnemyKind = iota
	EFast
)

type EnemySpec struct {
	HP     int
	Speed  float64
	Reward int
	Char1  rune
	Char2  rune
	Color  tcell.Color
}

var enemySpecs = map[EnemyKind]EnemySpec{
	ENormal: {
		HP: 20, Speed: 3.0, Reward: 10,
		Char1: 'o', Char2: 'o', Color: tcell.NewRGBColor(255, 100, 100),
	},
	EFast: {
		HP: 12, Speed: 5.5, Reward: 12,
		Char1: '>', Char2: '>', Color: tcell.NewRGBColor(255, 200, 80),
	},
}

type Enemy struct {
	PathIdx float64
	Kind    EnemyKind
	HP      int
	MaxHP   int
	Dead    bool
	Escaped bool
}

func (e *Enemy) Pos(path []Point) Point {
	i := int(math.Floor(e.PathIdx))
	if i >= len(path) {
		return path[len(path)-1]
	}
	if i < 0 {
		return path[0]
	}
	return path[i]
}
