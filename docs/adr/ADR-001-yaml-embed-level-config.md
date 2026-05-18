# ADR-001 — 关卡数据 YAML + `go:embed` 配置层

**状态：** accepted
**决策日期：** 2026-05-18（V1.5 决策时点）
**回填日期：** 2026-05-18（V2.7 后,按 standard `large-module.md` ADR 联动回填）
**作者：** AI agent + 用户拍板
**关联 commit：** `b2bc0cc` (V1.5)

---

## 上下文

V0~V1（597 行单文件 main.go）path 数据 hard-code 在 `buildPath()` 函数中,
waves 数据 hard-code 在 Go 数组中。

V1.5 要扩展到 10 关,如果继续 hard-code:
- 每关 path 数据写 Go slice literal → 10 × ~30 行 = 300+ 行噪声代码
- 每关 waves 数据 `[]EnemyKind{ENormal, ENormal, ...}` 极冗长
- 改关卡数据 = 改代码 + 重 build,迭代慢
- PM/设计师 无法独立编辑关卡数据（需开发参与）

## 决策

**关卡数据用 YAML 配置文件 + `go:embed` 编入 binary。**

- 数据位置：`assets/levels.yaml`
- 加载机制：`go:embed` 编译期内联,`gopkg.in/yaml.v3` 解析
- DSL 简化：
  - Path 用 control points `[[x,y],...]`,代码自动 expand 中间 cells（仅水平/垂直段）
  - Waves 用紧凑 DSL `"n5 f2 g1 b1"`,代码 parse 为 EnemyKind 序列
- 解析后 finalize 填充派生字段（`Level.Path` / `Level.Waves`）

## 反模式（已评估不采纳）

| 方案 | 不采纳理由 |
|------|----------|
| **A. 继续 hard-code Go data literals** | 10 关数据 ~600 行噪声,迭代慢,非开发者无法编辑 |
| **B. 外部 JSON 资源文件 (runtime fetch)** | 单 binary 部署优势丢失（需附 levels.yaml）,wasm 部署多一次 fetch |
| **C. SQLite / SQL database** | 过度工程,关卡数据天然适合声明式格式 |
| **D. Lua / 嵌入式脚本** | 引入运行时,体积膨胀,本游戏关卡数据是静态声明,无需脚本能力 |
| **E. JSON 替代 YAML** | JSON 不支持 comments,关卡 DSL 需注释解释 "n5 f2" 含义,YAML 更适合 |

## 影响

### 直接影响

- 新增依赖：`gopkg.in/yaml.v3`
- 新增文件：`assets/levels.yaml`(数据) / `loader.go`(embed + parse) / `level.go`(Level 模型 + DSL parser)
- 关卡编辑：只需改 YAML + 重 build,无需改 Go 代码
- Path 校验：control points 必须水平/垂直对齐,斜线会在 load 时报错（早期发现 bad data）

### 长期影响

- **正向**：后续加敌人类型只需 parser 加一个 case（V1.6 加 g/b 实际仅 2 行）
- **正向**：第三方贡献关卡门槛低（YAML 改即可,不需 Go 知识）
- **风险**：YAML 格式漂移（缩进错 / 类型错）会 load fail。Mitigation: tests/`TestLoadLevels_10Levels` 集成测试覆盖

## 关联

- **commit**：`b2bc0cc` (V1.5 实现)
- **测试**：`game_test.go` `TestParseWave_*` / `TestExpandPath_*` / `TestLoadLevels_10Levels`
- **后续 ADR**：ADR-002 (V2 UI 库切换) / ADR-003 (V2.5 WASM 部署)
