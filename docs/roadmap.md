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
| V5 Phase 1-5 | Gameplay 深度 | 卖塔（退款 70%）、targeting 策略（First/Last/Strong）、Cannon AoE 溅射、状态效果系统 + Frost 塔（第 4 塔型）、陨石雨主动技能（R 瞄准 + 25s 冷却）、killEnemy/pickTarget 重构、102 tests | `v5.0` | `8ee77f5` |
| V6 Phase 1-4 | 内容扩展 | 关卡 11-20（菜单两列）、难度三档（Normal/Hard/Easy）、Endless mode（预算制生成器 + seed 确定性）、星级评分、beginRun/spawnEnemy 重构、124 tests | `v6.0` | `d2d4bf1` |

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

## V5 — Gameplay 深度（已收尾 2026-06-04，tag `v5.0`）

> Phase 1-5 全部完成。以下为规划原文存档 + 头部收尾记录。

### V5 收尾记录（2026-06-04）

**完成度：** 5/5 phase，测试 72 → 102，三 build 全程绿。战术短板五项全清（卖塔 / targeting / AoE / 状态效果 / 主动操作）。

**工程纪律兑现：**
- 重构先行 ×2（pickTarget / killEnemy）均以"既有测试零改动全过"作回归证据
- killEnemy 统一击杀路径家族约束三次兑现（溅射杀 / 陨石杀 Spawner 均触发召唤，单测锁定）
- 测试先行实际拦截 1 个实现 bug（卖塔退款 IEEE754 截断）
- playwright tab 复用工具债（V4 P4 起 ≥2 例）以"换端口 = 新 origin"workaround 关闭

**未竟项（不阻塞收尾）：**

| 项 | 说明 | 去向 |
|----|------|------|
| 数值平衡未实玩压测 | 退款 0.7 / 溅射 0.5 / 减速 0.6-0.4 / 陨石 60 伤 25s | 全常量化，实玩反馈一行可调 |
| term build 不支持 V5 新操作 | 卖塔/切策略/陨石未接 term 输入；且 term HUD 显示 4 塔但 '4' 键未接（仅 1/2/3）| V2 起"term 冻结 V1.7 体验"惯例；若要补 Frost 选择是 term_main 一行 case |
| Frost 塔 sprite 灰色原版 | 未做冰蓝 tint（视觉区分靠形状 + 按钮文字 + 减速圈）| 若反馈混淆再加 tint 管线 |
| 陨石释放瞬间视觉未截图 | 火海 + shake 组合效果（game-over 抢跑）| 逻辑全单测；用户实玩可见 |

**下一版候选：** V6 — 内容扩展（关卡 11-20 / 难度模式 / endless / 星级评分），V4 方向矩阵中的第三项——gameplay 系统已深化，"同配方重复"的反对理由已减弱。启动前需用户拍板。

---

### 以下为 V5 规划原文（2026-06-04 启动时写入，存档不改）

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

## V6 — 内容扩展（已收尾 2026-06-04，tag `v6.0`）

> Phase 1-4 全部完成。以下为规划原文存档 + 头部收尾记录。

### V6 收尾记录（2026-06-04）

**完成度：** 4/4 phase，测试 102 → 124。游戏终态：20 关 + 3 难度 + endless + 星级，4 塔 5 敌全系统三端（桌面/终端/浏览器）。

**纪律兑现：**
- 内容数据测试兜底落地（20 关完整性 / unlock 链 / 难度单调性——Spawner 加权）
- seed 注入确定性铁律（首次引入随机性即立规，20-seed 门禁测试）
- save family 第 3/4/5 轮零失同步（Difficulty / BestWave / Stars）
- 统一施加点 ×2 新增（newEnemy 难度系数 / spawnEnemy endless 缩放——召唤物不漏）
- 测试拦截既有 bug 第 3 例（V3 Phase 5a 起的负数整除 hitbox 误命中）

**未竟项（不阻塞收尾）：**

| 项 | 说明 | 去向 |
|----|------|------|
| 新关卡曲线 + endless 平衡未实玩压测 | AI 设计内容的固有限制 | 数值全常量化，实玩反馈一行可调 |
| term 菜单 20 关在矮终端可能溢出 | 规划时记录，未实际验证 term 渲染 | term 冻结惯例；用户报障再处理 |
| V6 新操作未接 term（D/E 键）| 惯例延伸（V5 起累计：卖塔/策略/陨石/难度/endless）| term 定位已事实降级为"V1.7 兼容保留"|
| 星级/纪录纯 ASCII 呈现 | gomono 字形覆盖保守决策 | 若引入图标字体或 sprite 星标再美化 |

**下一版候选（未拍板，三选或另提）：**
a) 平衡打磨版——用户实玩反馈驱动的数值/手感调整批次
b) 发布版——itch.io / GitHub Pages 部署 + README + 截图 + 玩法说明
c) 新机制探索——多入口 path / 英雄单位 / 塔技能树等

---

### 以下为 V6 规划原文（2026-06-04 启动时写入，存档不改）

### 方向决策记录

**拍板：** V4 方向矩阵第三项，V5 收尾标记候选（"gameplay 深化后同配方重复的反对理由已减弱"），用户 2026-06-04 指示启动。

**主题：** 把 V5 攒下的系统深度变成可玩内容量——新关卡吃满 4 塔 5 敌 + 卖塔/策略/AoE/减速/陨石的组合空间，再加重玩性结构（难度 / endless / 星级）。

**已知翻车场景（反例）：** 内容设计（关卡/wave/平衡数值）由 AI 产出但无法实玩调参——对策：数值保守起步 + 全常量化 + 关卡数据完整性测试兜底（path 合法 / DSL 可解析），手感调整交给用户实玩反馈循环。

### 现状基线（2026-06-04 实测）

- 10 关 levels.yaml（cps path + waves DSL "n/f/g/b/s"），DSL 已覆盖 5 敌型 [已验证]
- 菜单单列 rowH 36 / startY 60，windowH 620 → 20 关需两列（单列 780 溢出）；点击 hitbox `menuRowAt` 需同步 [已验证: ebiten_renderer.go]
- term 菜单线性列出全部关卡——20 关在矮终端可能溢出（term 冻结惯例下记录不阻塞）
- 项目零 math/rand 使用——endless 生成器引入时 seed 注入保证确定性可测 [已验证: grep]
- Save 已有三字段（completed/volume/juice_off），双实现 family 纪律已跑 2 轮

### Phase 计划

#### Phase 1 — 关卡 11-20 + 菜单两列

- **范围：** 10 新关——path 形状多样化（螺旋/双 U/长直道/密集折返），wave 设计吃满 V5 系统（Spawner 群配 AoE、飞行波配 Frost/Archer、Boss 串配陨石时机）；难度曲线衔接 L10；菜单两列布局（10/列）+ `menuRowAt` hitbox 同步；键盘 1-9/0 保留选前 10 关，11-20 鼠标点击选择
- **验收：** 关卡数据完整性测试（每关 path 解析合法 + waves 非空 + unlock 链连续）；两列菜单截图（换端口方案）；L11-20 至少 1 关 smoke 通关路径（测试模拟）

#### Phase 2 — 难度模式

- **范围：** easy/normal/hard 三档系数（敌 HP / 奖励 / 起始 lives，常量表）；菜单 D 键循环切换 + 标题栏显示；存档记偏好（family edit ×3 轮）
- **验收：** 系数应用单测（同关不同难度的敌 HP/lives 差异）；持久化往返；normal == 现状回归

#### Phase 3 — Endless mode

- **范围：** 菜单 E 键进入；预算制 wave 生成器（wave N 预算 = base + N×增量，按敌型 cost 加权选购，超池后 HP 缩放），**seed 注入**保证确定性；best wave 计数入存档；HUD 显示 Wave N（无上限）
- **验收：** 生成器确定性单测（同 seed 同序列）+ 强度单调性（预算递增）；best wave 持久化；lose 后回菜单显示纪录

#### Phase 4 — 星级评分

- **范围：** 通关按剩命比例定星（3★ 满命 / 2★ ≥70% / 1★ 通关）；存档 per-level 取 max 不降级（family edit）；菜单行显示 ★★☆
- **验收：** 定星纯函数单测（边界 70%）；取 max 不降级；菜单显示截图

### 工程约定（延续 V5 + 内容向新增）

- 内容数据也走测试：levels.yaml 完整性测试兜底（防手写 yaml 错格式 / path 断裂 / unlock 链断号）
- save 扩展继续 family 纪律（save.go + save_wasm.go 同步 + 注释互指）
- 浏览器验证用换端口方案（V5 P5 确立）
- 平衡数值常量化；V6 收尾打 `v6.0` tag

---

## V7 — 发布版（active，2026-06-04 启动）

### 方向决策记录

**拍板：** V6 收尾三候选中用户选 b) 发布版（2026-06-04）。同时拍板两项外向决策：
1. **改名发布**——"Kingdom Rush" 是 Ironhide 注册商标，对外名改为 **Gopher Defense**（Go 吉祥物 + TD），标注 "inspired by Kingdom Rush"；仓库名/module 名保留
2. **转 public + GitHub Pages**——Actions 自动构建部署，push 即发布

**已知翻车场景（反例）：** 改名只改门面、游戏内残留旧名 → grep 验收兜底；Pages 部署后 WASM 加载路径/缓存问题 → 线上 URL playwright 实测（不只看 workflow 绿）。

### 现状基线（2026-06-04 实测）

- 仓库 PRIVATE / 无 README / index.html 停留 V2.5 门面（键位缺 V5/V6 全部新键）[已验证]
- 改名涉及面：用户可见字符串 6 处（term/ebiten menu + HUD title + window title + index.html ×3）[已验证: grep]
- **不改名单**：存档路径 `~/.kingdom-rush` / localStorage key / module 名 / wasm 文件名——改存档 key = 丢用户进度，技术标识符与对外名解耦
- 素材 license 全 CC0 已溯源（assets/*/LICENSE-*），发布合规基础齐

### Phase 计划

#### Phase 1 — 改名 + README + 门面

- **范围：** 6 处用户可见字符串改 Gopher Defense；README.md（特性/截图/键位表/三端构建指南/素材 attribution/版本史链接/"inspired by" 标注/存档位置说明）；截图 2-3 张入库 `docs/screenshots/`
- **验收：** `grep "Kingdom Rush"` 用户可见层零残留（注释/存档 key 白名单除外）；README 完整可读；截图入库

#### Phase 2 — Pages 部署工程

- **范围：** index.html 重写（新名/全键位表/loading 美化/内联 favicon 修 404）；GitHub Actions workflow（build wasm `-ldflags "-s -w"` 减体积 + actions/deploy-pages）；执行 `gh repo edit --visibility public`（已拍板授权）+ 开 Pages
- **验收：** workflow 绿；Pages URL 浏览器可玩（playwright 线上实测，注意这次没有换端口问题——线上 URL 天然新 origin）

#### Phase 3 — 发布收尾

- **范围：** GitHub Release v7.0（release notes 链 Pages URL + 版本史摘要）；仓库 metadata（description / homepage / topics）；license 清单终核；全量回归
- **验收：** Release 页可见；repo 门面齐；124+ tests 绿

### 工程约定

- 改名 = 纯字符串替换，不动逻辑；每 phase 三 build + 全量测试照常
- V7 收尾打 `v7.0` tag

---

## 变更记录

| 日期 | 变更 |
|------|------|
| 2026-06-04 | 初建：V0→V3 版本史回填 + V3 收尾（v1.0/v2.0/v3.0 tag）+ V4 方向拍板（音频 + Game Feel）+ Phase 1-5 规划 |
| 2026-06-04 | V4 收尾：tag `v4.0`（5936799），版本史补 V4 行，收尾记录（未竟项 ×3 + 盲点声明），V5 候选标记（Gameplay 深度，待拍板）|
| 2026-06-04 | V5 启动：用户拍板 Gameplay 深度，Phase 1-5 规划（卖塔 → targeting 策略 → Cannon AoE → 状态效果/减速 → 主动技能）；核心循环改动定"重构先行 + 测试先行"纪律 |
| 2026-06-04 | V5 收尾：tag `v5.0`（8ee77f5），版本史补 V5 行，收尾记录（纪律兑现 ×4 + 未竟项 ×4），V6 候选标记（内容扩展，待拍板）|
| 2026-06-04 | V6 启动：用户拍板内容扩展，Phase 1-4 规划（关卡 11-20 + 菜单两列 → 难度模式 → endless → 星级）；内容数据测试兜底纪律 + endless seed 注入约束 |
| 2026-06-04 | V6 收尾：tag `v6.0`（d2d4bf1），版本史补 V6 行，收尾记录（纪律兑现 ×5 + 未竟项 ×4），下一版三候选记录（平衡打磨 / 发布 / 新机制，未拍板）|
| 2026-06-04 | V7 启动：用户拍板发布版 + 两项外向决策（改名 Gopher Defense / 转 public+Pages）；Phase 1-3 规划（改名+README → Pages 部署 → Release 收尾）；存档 key 不改名单明确 |
