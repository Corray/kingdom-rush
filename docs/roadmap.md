# kingdom-rush — Roadmap

> **定位：** 版本史（已收尾 era 的事实记录）+ 当前 active 版本的 phase 规划。
> 每个 era 收尾时打 annotated tag；phase 规划在动手前写入本文件，实现中按实际情况修订（修订留变更记录）。
>
> **维护方：** 用户拍方向，EL 维护内容。

---

## 版本史（V0 → V3，已收尾）

| Era | 主题 | 范围 | Tag | 收尾 commit |
|-----|------|------|-----|------------|
| V0 ~ V1.7 | Terminal TD | tcell 渲染、wave/塔/敌核心循环、多关卡、unlock 系统 + JSON 存档（~/.kingdom-rush）、33 tests | `v1.0` | `6e44720` |
| V2 ~ V2.7 | Ebiten 桌面 + WASM | 双 build（默认 ebiten / `-tags term` 保留 V1.7）、WASM + localStorage 存档、Makefile、攻击视觉特效、鼠标输入 + 射程圈 | `v2.0` | `20f5a87` |
| V3 Phase 1-6b + V3.6 | Sprite/UI 美化 + 内容 | Kenney CC0 sprite 全替换、bullet 飞行动画、HUD/menu 美化、Go Mono truetype 字体、程序化 path 绘制 + 描边、Spawner 敌型（第 5 种）、39 tests | `v3.0` | `b108a26` |
| V4 Phase 1-5 | 音频 + Game Feel | SFX 管线（9 音效 Kenney CC0）、BGM 双轨 + 音量档持久化（Juhani Junkala CC0）、走动/死亡动画（程序化插值+旋转+摆动）、伤害飘字 + 金币反馈、shake + boss 顿帧 + J 开关、72 tests | `v4.0` | `5936799` |

### V3 未竟项（不阻塞收尾，归入 backlog）

| 项 | 来源 | 去向 |
|----|------|------|
| 多字号字体（title 大字 / body 小字）| V3 Phase 5c commit 未做段 | V4 不做；未来 UI 需求触发时再做 |
| Tank enemy | V3.6 commit 跳过段——Kenney pack 仅 3 个 walker sprite 已占满 | 受素材限制搁置，换素材包时重评 |
| 新 Tower 型 | V3.6 commit 跳过段——同类 tower sprite 有限 | 同上 |

---

## V4 — 音频 + Game Feel（已收尾 2026-06-04，tag `v4.0`）

> Phase 1-5 全部完成（含可选 Phase 5，用户拍板做）。以下为规划原文存档 + 末尾收尾记录。

### V4 收尾记录（2026-06-04）

**完成度：** 5/5 phase，测试 39 → 72，全程三 build（desktop / term / WASM）保持绿。

**未竟项（不阻塞收尾）：**

| 项 | 来源 | 去向 |
|----|------|------|
| 独立 hit 音（弹着时序）| Phase 1 偏离记录——伤害即时结算，shoot+hit 同帧叠播浑浊 | 若未来做"bullet 飞行期间伤害延迟结算"再议 |
| 音色搭配 / bgmBaseVol 0.4 平衡微调 | Phase 1/2 待听感项，收尾时用户未提出调整 | 调参入口保留（映射表换文件 / 改常量一行）|
| 飘字像素级截图验证 | Phase 4 playwright tab 错乱，按迭代上限停手 | 逻辑单测 + 同模式渲染分支佐证；用户实玩可见 |

**验证遗留风险（盲点声明）：** shake 强度 / 顿帧时长 / 飘字密度均按"弱强度默认"设计但未经长时间实玩压测，若实际体验过强，J 键可整体关闭（shake+顿帧），飘字无开关——若被反馈干扰需补开关。

**下一版候选：** V5 — Gameplay 深度（卖塔 / targeting 策略 / Cannon AoE / 减速塔 / 主动技能），V4 方向决策时已标记。启动前需用户拍板。

---

### 以下为 V4 规划原文（2026-06-04 拍板时写入，存档不改）

### 方向决策记录

**拍板：** 用户 2026-06-04 从 4 个候选方向中选定「音频 + Game Feel」。

**选择理由：** V3 补完视觉后，音频是唯一完全缺失的感官通道（代码 grep 零命中）；gameplay 已有 3 塔 / 5 敌 / 10 关体量，但无声运行像 demo。Kenney 有 CC0 音频包可延续现有素材链，ebiten/audio 支持 WASM。

**备选方向 + 为什么没选（本期）：**

| 方向 | 内容 | 没选理由 |
|------|------|---------|
| Gameplay 深度 | 卖塔 / targeting 策略 / Cannon AoE / 减速塔 / 主动技能 | 每项都动核心循环测试面大；新塔型受 Kenney sprite 限制（V3.6 已实证）。**候选 V5 主题** |
| 内容扩展 | 关卡 11-20 / 难度模式 / endless / 星级 | 系统不变时新关卡是同配方重复，价值递减；适合在 gameplay 深度做完后做 |
| 混合精选 | 各方向抽核心项 | 无统一主题，phase 间无复用积累，工程上更碎 |

**已知翻车场景（推荐方向的反例）：** 纯体验向——V4 做完后 gameplay 战术深度短板仍在（targeting 写死最前优先、全塔单体伤害、不能卖塔）。如果玩家反馈集中在"不耐玩"而非"没声音"，本方向就是错排序。

### 现状基线（2026-06-04 实测）

- 音频：完全缺失（grep audio/sound 零命中）[已验证]
- 特效系统：`effects.go` + `decayEffects` 管线自 V2.6 存在，shared 逻辑（无 build tag）[已验证]
- 隔离模式：ebiten-only 文件用 `//go:build !term`，素材用 `go:embed`（sprite_loader.go / loader.go 先例）[已验证]
- 素材 attribution 先例：`assets/sprites/LICENSE-Kenney.txt` [已验证]
- 测试基线：39/39 PASS，不回退

### Phase 计划

#### Phase 1 — ebiten/audio 接入 + SFX 管线

- **范围：** audio context 单例 + `audio_loader.go`（`//go:build !term` + embed wav）；SFX 事件：shoot（按塔型 3 种）/ hit / death / build / upgrade / wave start / win / lose
- **关键架构决策：** game.go 是 shared 逻辑（term build 也编译），**不能 import ebiten/audio**。方案：game.go 维护纯数据事件队列 `g.SoundEvents []SoundEvent`（仿 Effects 管线先例），ebiten 侧每帧 drain + 播放，term 侧丢弃。测试只测触发（event 入队），不测播放
- **素材：** Kenney CC0 音频包（具体 pack 动手前按 research-first 核实官网 + license，落 `assets/sfx/LICENSE-*.txt`）
- **验收：** 桌面 + WASM 实测出声；`-tags term` build 仍过；测试 ≥ 39 且新增 SFX 触发测试
- **风险：** WASM 浏览器 autoplay 政策——音频需用户手势后解锁，ebiten 的处理机制是 Phase 1 第一个实测项 [推断，待验证]

#### Phase 2 — BGM + 音量控制

- **范围：** menu loop + 战斗 loop 两首 BGM（ogg/vorbis，文件体积远小于 wav）；切换淡入淡出；静音键 + 音量档；音量设置入存档
- **注意：** 存档结构扩展要同步改 `save.go` + `save_wasm.go` 双实现（fix-pattern-scan 家族：同一 struct 两份 IO 实现）
- **素材：** BGM 来源 TBD——Kenney 音乐类素材偏 jingle，循环 BGM 可能需 OpenGameArt CC0 等其他来源，**license 核实是本 phase 前置项**
- **验收：** menu/战斗切换无爆音；音量设置持久化（桌面 + WASM 双端）；WASM 首次交互后 BGM 正常起播

#### Phase 3 — Enemy 走动 / 死亡动画

- **范围：** 走动动画 + 死亡动画
- **关键约束：** Kenney pack walker sprite 仅 3 个且已占满（V3.6 实证）——**没有现成动画帧**。方案优先级：程序化动画（上下 bob 浮动 / 移动朝向翻转 / 死亡淡出 + 缩放）> sprite 帧动画（仅素材允许时）
- **死亡动画**走 effects 管线（decayEffects 先例）
- **验收：** 视觉实测（WASM + playwright 截图，V3 Phase 6b 先例）；性能无感知劣化（满屏敌人场景）

#### Phase 4 — 伤害飘字 + 金币反馈

- **范围：** Effect 类型扩 text 变体（truetype drawText 已就绪）；伤害数字上飘淡出；击杀 `+25g` 金色飘字；HUD 金币变化闪烁
- **验收：** 视觉实测；同屏大量飘字不卡顿（L10 末 wave 12+ 敌场景）

#### Phase 5 —（可选）screen shake / 击中顿帧

- **范围：** lives 丢失时轻微 shake；boss 击杀顿帧
- **默认弱强度**——过度使用伤体验；做完 Phase 1-4 后按手感决定做不做
- **验收：** 视觉实测 + 可一键关闭（设置项）

### 工程约定（延续 V3）

- commit 风格：`feat: V4 Phase N — <主题>`，body 含 verify 段（三 build + tests + 实测）
- 每 phase 独立可交付，按 incremental-verification 节奏推进
- 渲染/音频类改动必须实测验证（WASM + playwright 截图 / 实际听声），不止编译过
- V4 收尾时打 `v4.0` tag

---

## V5 — Gameplay 深度（active，2026-06-04 启动）

### 方向决策记录

**拍板：** V4 方向决策（2026-06-04）时已将本方向列为候选 V5 并记录内容；V4 收尾同日用户指示"启动 V5"，视为拍板生效。备选方向的对比矩阵见 V4 决策记录段，不重复。

**主题：** 战术深度——V4 收尾时的已知短板（targeting 写死最前优先、全塔单体伤害、不能卖塔、无状态效果、无主动操作）逐项补齐。

**已知翻车场景（反例）：** 每个 phase 都动核心战斗循环（`game.go` shoot/move），回归风险高于 V4 的渲染层改动——对策：pickTarget / killEnemy 等先重构为纯函数再扩展，测试先行。平衡性（退款比例 / 溅射衰减 / 减速幅度）数值全部常量化，实玩后可一行调。

### 现状基线（2026-06-04 实测）

- targeting：shoot 循环内联"最远 PathIdx 优先"，无策略选择 [已验证: game.go]
- 伤害：全塔单体即时结算；击杀处理（音效/动画/赏金/Spawner 召唤）内联在 shoot 循环 [已验证]
- 塔操作：只有建/升（TryAction），无卖
- 敌人：无状态效果字段（speed 直接查 spec）
- `TowerLevel.Cost` 为逐级增量 → 卖塔退款可纯函数累计 [已验证: entities.go:35]
- 测试基线 72/72；term build 输入层保持 V1.7 体验（V2 起惯例，V5 新操作仅接 ebiten 输入）

### Phase 计划

#### Phase 1 — 卖塔

- **范围：** 光标在塔上按 X（或右键）卖塔；退款 = 已投入（逐级 Cost 累计）× `sellRefundRate`（常量 0.7）；塔移除 + 金币入账（HUD flash 自动生效）+ Msg + 卖塔 SFX（Kenney interface pack 补 1 个音效，沿用既有下载链路）
- **验收：** 退款数值纯函数单测；卖后塔消失不再射击；term build 不受影响

#### Phase 2 — Targeting 策略

- **范围：** Tower 加 `Targeting` enum（First 默认 = 现状 / Last / Strong 最高 HP）；先把内联目标选择重构为纯函数 `pickTarget`（行为不变回归），再扩策略；光标在塔上按 T 循环切换；塔上 HUD 提示当前策略
- **验收：** 三策略各自选对目标的单测；重构后旧行为回归（First == 现状）

#### Phase 3 — Cannon AoE 溅射

- **范围：** `TowerLevel` 加 `Splash float64`（仅 Cannon 非零，半径 cell 单位，随级别增长）；主目标全额伤害，半径内其他敌人 50% 溅射（`splashFactor` 常量）；**前置重构：** 击杀处理（音效/动画/赏金/Spawner）抽 `killEnemy` helper——溅射致死与主目标致死走同一路径，防 fix-pattern-scan 家族漏（Spawner 召唤逻辑只能有一份）
- **验收：** 溅射多杀单测（含溅射杀 Spawner 触发召唤）；Archer/Magic 无溅射回归

#### Phase 4 — 状态效果系统 + 减速载体

- **范围：** Enemy 加状态效果（先做 slow：`slowFactor` + `slowTimer`，move 循环按系数减速，过期恢复，不叠加取最强）；**载体待 phase 开工时定**——选项 a) 新塔型（前置：PIL 裁 tilesheet 调研可用 tower sprite，V3.6 实证风险）b) Magic 塔升级附带减速 c) 现有 sprite 换色变体。框架先行，载体后定
- **验收：** 减速生效/过期/不叠加单测；满屏减速场景性能无感知劣化

#### Phase 5 — 主动技能（陨石雨）

- **范围：** R 键进入瞄准 → 点击释放：半径内全敌 AoE 伤害 + 火焰特效 + shake 复用（Phase 5 juice）+ 冷却（`meteorCooldownS` 常量，HUD 显示冷却条）；再按 R 取消瞄准（**不用 Esc——已绑退出**）
- **验收：** 范围伤害/冷却 gating 单测；瞄准态 UI 可见；冷却条显示

### 工程约定（延续 V4）

- commit 风格 `feat: V5 Phase N — <主题>`，body 含 verify 段（三 build + tests + 实测）
- 核心循环改动测试先行：重构步（pickTarget / killEnemy）单独验证行为不变，再加新功能
- 平衡性数值全常量化；WASM 验证优先静态可见对象（V4 Phase 4 教训）
- V5 收尾打 `v5.0` tag

---

## 变更记录

| 日期 | 变更 |
|------|------|
| 2026-06-04 | 初建：V0→V3 版本史回填 + V3 收尾（v1.0/v2.0/v3.0 tag）+ V4 方向拍板（音频 + Game Feel）+ Phase 1-5 规划 |
| 2026-06-04 | V4 收尾：tag `v4.0`（5936799），版本史补 V4 行，收尾记录（未竟项 ×3 + 盲点声明），V5 候选标记（Gameplay 深度，待拍板）|
| 2026-06-04 | V5 启动：用户拍板 Gameplay 深度，Phase 1-5 规划（卖塔 → targeting 策略 → Cannon AoE → 状态效果/减速 → 主动技能）；核心循环改动定"重构先行 + 测试先行"纪律 |
