// Level / Path / Wave 模型与解析逻辑。
//
// Level yaml schema:
//
//	id, name, start_gold, start_lives, cps ([[x,y],...]), waves (["n5 f1",...])
//
// ExpandPath: 把 control points 扩展为完整 path cells。相邻 cp 必须水平或
// 垂直对齐 (V1.5 不支持斜线)。
// ParseWave: 把 "n5 f2 n3" DSL 解析为 []EnemyKind。
package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Level struct {
	ID         int      `yaml:"id"`
	Name       string   `yaml:"name"`
	StartGold  int      `yaml:"start_gold"`
	StartLives int      `yaml:"start_lives"`
	CPS        [][]int  `yaml:"cps"`
	WavesRaw   []string `yaml:"waves"`

	// 派生字段(load 时填充, 不来自 yaml)
	Path  []Point    `yaml:"-"`
	Waves []WaveSpec `yaml:"-"`
}

func ExpandPath(cps [][]int) ([]Point, error) {
	if len(cps) < 2 {
		return nil, fmt.Errorf("need at least 2 control points, got %d", len(cps))
	}
	for i, cp := range cps {
		if len(cp) != 2 {
			return nil, fmt.Errorf("cp %d malformed: expected [x,y] got %v", i, cp)
		}
	}
	points := []Point{{cps[0][0], cps[0][1]}}
	for i := 1; i < len(cps); i++ {
		ax, ay := cps[i-1][0], cps[i-1][1]
		bx, by := cps[i][0], cps[i][1]
		switch {
		case ax == bx:
			step := 1
			if by < ay {
				step = -1
			}
			for y := ay + step; ; y += step {
				points = append(points, Point{ax, y})
				if y == by {
					break
				}
			}
		case ay == by:
			step := 1
			if bx < ax {
				step = -1
			}
			for x := ax + step; ; x += step {
				points = append(points, Point{x, ay})
				if x == bx {
					break
				}
			}
		default:
			return nil, fmt.Errorf("control points %d→%d must be axis-aligned: (%d,%d)→(%d,%d)",
				i-1, i, ax, ay, bx, by)
		}
	}
	return points, nil
}

func ParseWave(s string) ([]EnemyKind, error) {
	var out []EnemyKind
	tokens := strings.Fields(s)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty wave")
	}
	for _, tok := range tokens {
		if len(tok) < 2 {
			return nil, fmt.Errorf("invalid token %q", tok)
		}
		kindCh := tok[0]
		count, err := strconv.Atoi(tok[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid count in %q: %v", tok, err)
		}
		if count <= 0 {
			return nil, fmt.Errorf("non-positive count in %q", tok)
		}
		var kind EnemyKind
		switch kindCh {
		case 'n':
			kind = ENormal
		case 'f':
			kind = EFast
		case 'g':
			kind = EGlider
		case 'b':
			kind = EBoss
		case 's':
			kind = ESpawner
		default:
			return nil, fmt.Errorf("unknown enemy kind %q in %q", kindCh, tok)
		}
		for i := 0; i < count; i++ {
			out = append(out, kind)
		}
	}
	return out, nil
}

// FinalizeLevels 把 raw yaml 解析后的 Level slice 的 Path / Waves 字段填充
func FinalizeLevels(levels []Level) error {
	for i := range levels {
		path, err := ExpandPath(levels[i].CPS)
		if err != nil {
			return fmt.Errorf("level %d (%s): %w", levels[i].ID, levels[i].Name, err)
		}
		levels[i].Path = path
		levels[i].Waves = nil
		for j, raw := range levels[i].WavesRaw {
			seq, err := ParseWave(raw)
			if err != nil {
				return fmt.Errorf("level %d (%s) wave %d: %w", levels[i].ID, levels[i].Name, j, err)
			}
			levels[i].Waves = append(levels[i].Waves, WaveSpec{Enemies: seq})
		}
	}
	return nil
}
