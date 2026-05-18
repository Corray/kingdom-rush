# ADR-003 — WASM build + localStorage 存档拆分

**状态：** accepted
**决策日期：** 2026-05-18（V2.5 决策时点）
**回填日期：** 2026-05-18（V2.7 后,按 standard `large-module.md` ADR 联动回填）
**作者：** AI agent + 用户拍板
**关联 commit：** `c623f72` (V2.5)

---

## 上下文

V2 ebiten 桌面 binary 工作良好。用户要求加 web/WASM 部署路径(浏览器内运行)。

ebiten 官方支持 WASM (`GOOS=js GOARCH=wasm go build`),Game logic / save / level
等代码层无需改动。但**存档系统**在 WASM 环境会失败:

- `save.go` 用 `os.ReadFile` / `os.WriteFile` / `os.MkdirAll`
- WASM 浏览器环境无 file system(syscall/js 不暴露 fs API)
- 直接 WASM build,save 调用会 runtime panic 或 silent fail

## 决策

**用 Go build tag 拆分 save 实现,WASM 用 localStorage,native 保持 file IO。**

### 决策 1：build tag 互斥拆分

- `save.go`：加 `//go:build !js`（除 wasm 外都参与编译）
- `save_wasm.go`（新）：`//go:build js && wasm`,用 `syscall/js` 调 localStorage

互斥关系：
| Build target | save.go | save_wasm.go |
|---|---|---|
| `go build .`（ebiten desktop）| ✓ | ✗ |
| `go build -tags term`（terminal）| ✓ | ✗ |
| `GOOS=js GOARCH=wasm go build .`（WASM）| ✗ | ✓ |

### 决策 2：Save struct 在 save_wasm.go 复刻定义

由于 `save.go` 在 WASM build 中被 build tag 排除,**`Save` struct / `NewSave` /
`IsCompleted` / `IsUnlocked` / `MarkCompleted` 等纯逻辑也消失**。WASM build
需要这些定义,所以在 `save_wasm.go` 中复刻一份相同定义。

牺牲 DRY 换 build tag 简单性（替代方案见反模式 C）。

### 决策 3：localStorage 直接用 syscall/js

```go
storage := js.Global().Get("localStorage")
storage.Call("getItem", "kingdom-rush-save")
storage.Call("setItem", "kingdom-rush-save", jsonString)
```

key 固定 `"kingdom-rush-save"`,JSON 序列化（同 native 格式,容量 < 1 KB << 5 MB 限制）。

### 决策 4：HTML wrapper + wasm_exec.js 拷贝

- `web/index.html`（commit 入版本控制）：minimum HTML + 加载 wasm_exec.js + .wasm
- `web/wasm_exec.js`（gitignore）：从 `$GOROOT/lib/wasm/wasm_exec.js` 拷贝,
  Go 工具链自带。版本绑定 Go runtime,不应 commit
- `web/kingdom-rush.wasm`（gitignore）：build artifact

### 决策 5：Makefile 一键 build

- `make build` / `make build-term` / `make wasm` / `make serve`
- `wasm` target 自动拷贝 wasm_exec.js + 编译 .wasm
- `serve` target 启 `python3 -m http.server 8080`

## 反模式（已评估不采纳）

| 方案 | 不采纳理由 |
|------|----------|
| **A. 运行时检测 platform（reflect / build-time-injected var）** | runtime 检测 = save logic 内含 if-else 分支,WASM/native 两条路径都参与编译。但 native build 无 `syscall/js`（编译期 fail）,反向亦然。**build tag 是 Go 解决跨平台分歧的标准方式**,不绕弯 |
| **B. 统一 file IO abstraction（io/fs.FS interface）** | 加 `FS` interface + native impl + wasm-localStorage impl + 依赖注入。增 ~200 行 boilerplate 换 < 100 行 code 解耦。代价 > 收益（YAGNI） |
| **C. 拆纯逻辑到 save_common.go（无 tag）** | `Save struct` / `NewSave` / `IsCompleted` 等放 `save_common.go`（无 tag,两边都参与）,`save.go` / `save_wasm.go` 只各自 IO。**问题**：`game_test.go` 用 `savePathFn` 钩子,savePathFn 是 `save.go` (native) 特有,test 仍然只在 native build 跑。拆 common 文件不减少最终复杂度,反而增文件数 |
| **D. WASM 中 save noop（不持久化）** | 浏览器关闭后进度丢失,玩家体验差。localStorage 是浏览器侧 standard,5 MB 容量远超本游戏需要（< 1 KB）|
| **E. IndexedDB 替代 localStorage** | localStorage 同步 API 简单（同 file IO 心智模型）。IndexedDB 异步 + transactional,本游戏存档量 < 1 KB 不需要这种复杂度 |

## 影响

### 直接影响

- 新文件：`save_wasm.go`(87 行) / `web/index.html`(53 行) / `Makefile`(50 行)
- 改动：`save.go` 加 `//go:build !js` / `.gitignore` 加 web 产物
- 浏览器使用：`make serve` → `http://localhost:8080`,点 canvas 获焦键盘
- 生产部署：`web/` 目录扔任何静态 host（GitHub Pages / Netlify / Cloudflare / nginx）

### 长期影响

- **正向**：3 build target（ebiten desktop / terminal / WASM）共享 game logic,
  发布渠道扩展无限制
- **正向**：localStorage save 跨 native/WASM 体验一致（进度持久化）
- **风险**：localStorage 单 origin 隔离,跨设备不同步。Mitigation: 未来加云端
  存档（V3? 需要后端,本 ADR 不涉及）
- **风险**：浏览器隐私模式 / 禁用 localStorage → save silently no-op。本游戏
  接受此 degradation,代码层 `storage.Truthy()` 守卫不 panic

## 关联

- **commit**：`c623f72` (V2.5)
- **build matrix**：
  - `go build .` → 14 MB(ebiten desktop)
  - `go build -tags term .` → 4.8 MB(terminal V1.7)
  - `GOOS=js GOARCH=wasm go build .` → 14 MB(`.wasm`)
- **前置 ADR**：
  - ADR-001（YAML config:基础数据层与平台无关,WASM 直接 reuse）
  - ADR-002（build tag 模型已建立,本 ADR 扩展第三个 target）
- **未来扩展**：
  - 云端存档（需后端）
  - IndexedDB 升级（如果存档容量超 localStorage 限制）
