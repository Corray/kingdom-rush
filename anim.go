// V4 Phase 3: 程序化走动动画的纯数学层 (无 ebiten 依赖, 可测)。
//
// 背景: Kenney pack 没有 walker 动画帧 (V3.6 实证, 3 个 walker sprite
// 已被 Normal/Fast/Boss 占满), roadmap 定程序化方案:
//   - pathLerp:   渲染层平滑插值移动 (游戏逻辑仍用 e.Pos 取整, 不变)
//   - pathDir:    当前段行进方向 → 朝向旋转 (sprite 默认朝右, 实测
//     tilesheet tile246/248 右侧有箭头)
//   - bobOffset:  垂直于行进方向的摆动 (随 PathIdx 推进, 速度快摆得快)
package main

import "math"

const (
	bobFreq = 14.0 // rad / path-cell, 一个摆动周期 ~0.45 cell
	bobAmp  = 2.0  // 摆动振幅 px
)

// pathLerp: PathIdx → 平滑插值路径坐标 (cell 单位 float)。
// 越界 clamp 到端点 (与 Enemy.Pos 的 clamp 语义一致)。
func pathLerp(path []Point, idx float64) (float64, float64) {
	if len(path) == 0 {
		return 0, 0
	}
	if idx <= 0 {
		return float64(path[0].X), float64(path[0].Y)
	}
	i := int(math.Floor(idx))
	if i >= len(path)-1 {
		last := path[len(path)-1]
		return float64(last.X), float64(last.Y)
	}
	frac := idx - float64(i)
	a, b := path[i], path[i+1]
	return float64(a.X) + float64(b.X-a.X)*frac,
		float64(a.Y) + float64(b.Y-a.Y)*frac
}

// pathDir: 当前段方向 (dx, dy ∈ {-1,0,1})。path 不足两点 → 朝右。
func pathDir(path []Point, idx float64) (int, int) {
	if len(path) < 2 {
		return 1, 0
	}
	i := int(math.Floor(idx))
	if i < 0 {
		i = 0
	}
	if i >= len(path)-1 {
		i = len(path) - 2
	}
	return path[i+1].X - path[i].X, path[i+1].Y - path[i].Y
}

// bobOffset: 行进摆动偏移 px (渲染时叠加在垂直于行进方向的轴上)。
func bobOffset(pathIdx float64) float64 {
	return math.Sin(pathIdx*bobFreq) * bobAmp
}

// dirAngle: 方向 → sprite 旋转角 (弧度)。sprite 默认朝右: (1,0) → 0。
func dirAngle(dx, dy int) float64 {
	if dx == 0 && dy == 0 {
		return 0
	}
	return math.Atan2(float64(dy), float64(dx))
}
