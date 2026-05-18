# ADR-002 — UI 库切换 tcell → ebiten + build tag 双 binary

**状态：** accepted
**决策日期：** 2026-05-18（V2 决策时点）
**回填日期：** 2026-05-18（V2.7 后,按 standard `large-module.md` ADR 联动回填）
**作者：** AI agent + 用户拍板
**关联 commit：** `c5fe11d` (V2)

---

## 上下文

V0~V1.7 全部基于 tcell terminal 渲染。V1.7 add 了 `Renderer` interface 作为
"future-ready" 抽象,但当时只有 `TermRenderer` 一种实现。

用户在 V1.7 commit 后选择推进 "Web/WASM Renderer (ebiten)" 作为下一步。
触发以下决策需求:
1. 是否完全替换 tcell? V1.7 用户已经习惯 terminal 体验
2. ebiten 是 callback-based(`Game.Update/Draw/Layout`),tcell 是 event-loop based,
   V1.7 `Renderer` interface 是否能复用?
3. `entities.go` 中 `TowerSpec.Color` 类型 `tcell.Color` 与 ebiten 不兼容,
   如何解耦?

## 决策

**三层决策合并为一个 ADR（虽是三个子决策,但本质都是"V1.7 → V2 UI 层架构调整"）：**

### 决策 1：默认 build = ebiten,terminal opt-in via `-tags term`

- `main.go` (no tag) 是 ebiten entry
- `term_main.go` (`//go:build term`) 保留 V1.7 main 逻辑
- `render.go` 也加 `//go:build term` (TermRenderer 仅 term build 参与)
- `ebiten_renderer.go` (`//go:build !term`) ebiten Game 实现

build 命令:
- `go build .` → ebiten desktop (默认)
- `go build -tags term -o ...term .` → terminal 模式

### 决策 2：删除 V1.7 的 `Renderer` interface

V1.7 加的 interface:
```go
type Renderer interface {
    Init() error
    Fini()
    PollEvent() tcell.Event   // ← tcell-specific!
    Sync()
    Draw(g *Game)
}
```

问题:
- `PollEvent` 返回 `tcell.Event` 是 tcell-specific 类型,无法被 ebiten 满足
- ebiten 是 callback-based,没有 `PollEvent` 等价物（事件由 `ebiten.RunGame`
  的 frame callback 内 `inpututil` 处理）
- interface 设计是"基于 tcell 范式" + "假设其他 UI 库也能 fit",但 ebiten
  范式根本不同 → interface 是过度设计(speculative generality)

V2 删除 interface,terminal main 直接调 `*TermRenderer` 方法,ebiten main
直接驱动 `*EbitenGame`。

### 决策 3：自定义 `RGB` 类型抽象颜色

`entities.go` 中 `TowerSpec.Color`/`EnemySpec.Color` 用 `tcell.Color`:
- entities 层不应依赖任何 UI 库
- terminal 用 `tcell.NewRGBColor(int32, int32, int32)`
- ebiten 用 `color.RGBA{R, G, B, A uint8}`

新加 `color.go`：
```go
type RGB struct { R, G, B uint8 }
```

各 renderer 转换:
- terminal: `termColor(RGB) → tcell.Color`
- ebiten: `ebitenColor(RGB) → color.RGBA`

## 反模式（已评估不采纳）

| 方案 | 不采纳理由 |
|------|----------|
| **A. 完全替换 tcell（删 terminal 模式）** | V1.7 用户已习惯 terminal 体验,完全替换 = 大幅 regression。失去 dev-friendly 快速验证渠道（terminal mode 启动 ms 级,无 GUI 依赖） |
| **B. 保留 Renderer interface,改造 ebiten 适配** | ebiten 是 callback-based,强行套 `PollEvent` 范式 = 假转换（poll 仅在 Update callback 内调度）。不解决真实问题,反增复杂度 |
| **C. 用 `image/color.Color` interface 而非自定义 RGB** | `image/color.Color` 是 RGBA(uint32) 接口,适配 tcell 需要拆 R/G/B（uint8 → int32）。自定义 RGB(R, G, B uint8) 三方都简单转换,避免 interface boxing 开销 |
| **D. 两个 main 文件用 cmd/ 子目录** | `cmd/kingdom-rush-gui/main.go` + `cmd/kingdom-rush-term/main.go`。改变 `go build .` 默认行为（需指定 cmd 路径）,与现有 Makefile / CI / 用户习惯不兼容 |
| **E. 自动启发式（检测 TTY → 选 terminal,否则 ebiten）** | os.Stdin TTY 检测在 GUI launcher / IDE 中 unreliable。用户控制权不应交给启发式 |

## 影响

### 直接影响

- 新增依赖：`github.com/hajimehoshi/ebiten/v2` + `inpututil` + `ebitenutil` + `vector`（间接 +5 包）
- binary 大小：terminal 4.8 MB → ebiten 14 MB（ebiten 引入大量 graphics deps）
- 新文件：`color.go`(16 行) / `ebiten_renderer.go`(~330 行) / `term_main.go`(ex-main.go 重命名)
- 改动文件：`entities.go`(tcell.Color → RGB) / `render.go`(加 build tag) / `main.go`(替换为 ebiten entry)

### 长期影响

- **正向**：解耦 UI 层后,加 WASM (ADR-003) / 加新 renderer impl 阻力低
- **正向**：terminal mode 作为 V1.7 兼容 + dev 快速验证渠道仍可用
- **风险**：双 binary 维护成本（两套 input/draw 逻辑）。但 game logic / save / level / entities 全共享,实际重复代码 < 20%
- **教训**：V1.7 的 `Renderer` interface 是过度设计的反面案例 — 在没有第二个 impl 时加抽象 layer,常常是"speculative generality"。规则建议: 等第二个 impl 真出现时再抽象,而不是预想

## 关联

- **commit**：`c5fe11d` (V2)
- **测试**：33/33 PASS（game/save/level/entities 层不依赖 UI,renderer 层手测）
- **前置 ADR**：ADR-001（YAML config 已经解耦 data 层）
- **后续 ADR**：ADR-003 (WASM 部署形态,基于本 ADR 的 build tag 模型扩展)
