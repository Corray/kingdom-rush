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
| V7 Phase 1-3 | 发布版 | 改名 Gopher Defense（商标合规 + 存档零迁移）、README/LICENSE(MIT)/截图、GitHub Pages 上线（Actions 测试门禁 + 自动部署）、Release v7.0 + repo metadata、124 tests | `v7.0` | `be86c3c` |
| V8 Phase 1-5 | 英雄单位 | 可控英雄（光标+H 设 rally / 自动打地面敌 / 贴身阻挡 / 阵亡复活 / 飞行飞越）、首个非路径绑定实体、敌近战反击、两端渲染+HUD、仿真接入（英雄纯增量 Hard 17→19 零回归）、148 tests | `v8.0` | `601f7c1` |
| V9 Phase 1-4 | 英雄成长 | 关内 per-run 等级/XP（被动 XP 威胁加权——决策 B 两轮仿真证伪改判、升级提 HP/伤害/射程+回血）、AoE 横扫主动技能（L3 解锁 / `G` 键 / 8s 冷却 / 经 damageEnemy）、两端 HUD 等级/XP 条/技能冷却、仿真接入成长（平衡零回归 Hard 19/20）、test-infra headless 解耦、159 tests | `v9.0` | `a3c8be3` |
| V10 Phase 1-4 | 多英雄选择 | 三职业 Knight/Archer/Rogue（HeroClass 参数化重构、阻挡 per-class gating、菜单 `H` 选择 + `Save.HeroChoice` 零值兼容、程序化配色金/绿/紫）、per-class 仿真矩阵 + 一轮校准（三职业 Hard 19/20 净 +2 持平、难点关保留、差异落风格不落数值）、173 tests | `v10.0` | `f9dcd9c` |
| V11 Phase 1-4 | Meta 成长 | 跨局技能树（per-class 线性 4 节点、星货币 60 赚 vs 90 总价制造 build 取舍、`HeroBonus` 七字段 beginRun 快照、菜单 `T` 新 `PhaseSkillTree` 两端）、「perk 预算」替代「按关缩放」（持久层存树点不存等级 → 平衡冲击缩到一张 perk 表）、仿真上下界 + Knight 树一轮校准（满树三职业 Hard 20/20 earned power、无树基线零回归）、188 tests | `v11.0` | `cda346d` |
| V12 Phase 1-4 | 多入口 path | 单路径承重假设首次重构：`Paths [][]Point` + `Enemy.PathID` + `Level.CPS2`，57 触点机械替换零回归；全 20 关双路（cps2 汇合同终点）、spawn 均摊交替、汇流渲染去重两端；仿真校准（Normal 20/20、Hard 19/20 难点移 Lv9）、195 tests | `v12.0` | `4494364` |
| V13 | 新敌人类型 | 敌型 5→8：EShield（护甲减免 min 1）、ERegen（每秒回血）、EHealer（治疗附近盟友）；L8 起渐进混入 20 关 wave、endless 按波次解锁；202 tests | `v13.0` | `21cc588` |

### V3 未竟项（不阻塞收尾，归入 backlog）

| 项 | 来源 | 去向 |
|----|------|------|
| ~~多字号字体（title 大字 / body 小字）~~ | V3 Phase 5c commit 未做段 | **已消（2026-06-05, V7.2 M1）**：20pt 标题 face + drawTextBigCol |
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

**验证遗留风险（盲点声明）：** shake 强度 / 顿帧时长 / 飘字密度均按"弱强度默认"设计但未经长时间实玩压测，若实际体验过强，J 键可整体关闭（shake+顿帧），~~飘字无开关——若被反馈干扰需补开关~~（**已补：2026-06-05 V7.3 D2，飘字并入 J**）。

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
| term build 不支持 V5 新操作 | 卖塔/切策略/陨石未接 term 输入；~~且 term HUD 显示 4 塔但 '4' 键未接~~（**'4' 键已接：2026-06-05 V7.3 E1**）| V2 起"term 冻结 V1.7 体验"惯例（其余操作维持不接）|
| ~~Frost 塔 sprite 灰色原版~~ | 未做冰蓝 tint | **已消（2026-06-05 V7.3 D1）**：drawTileTint 冰蓝乘数 |
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

## V7 — 发布版（已收尾 2026-06-05，tag `v7.0`）

> Phase 1-3 全部完成。以下为规划原文存档 + 头部收尾记录。

### V7 收尾记录（2026-06-05）

**完成度：** 3/3 phase。**游戏正式公开：**
- 在线玩：https://corray.github.io/kingdom-rush/
- Release：https://github.com/Corray/kingdom-rush/releases/tag/v7.0
- push master 即自动重新部署（CI 测试门禁，不过不上线）

**关键执行记录：**
- 改名 6 处用户可见字符串 + 不改名单严格执行（存档 key/module 名——用户进度零迁移）；grep 白名单制验收
- 两项外向操作按拍板执行：转 public + Pages 开启
- CI 首跑失败 1 次（ebiten 包 init 需 DISPLAY）→ xvfb 修复（官方 CI 同款）
- License 终核逐文件映射（12 sfx + 2 bgm + sprites）全过

**未竟项（不阻塞收尾）：**

| 项 | 说明 | 去向 |
|----|------|------|
| itch.io 上传 | 规划时定 Pages only | 需要手动账号操作，用户自行或后续指令 |
| Release 无 desktop 二进制 | macOS unsigned 警告体验差 | 源码构建为主；要发二进制再议（goreleaser）|
| ~~Actions Node 20 弃用警告~~ | checkout@v4/setup-go@v5，2026-06-16 起强制 Node 24 | **已修（2026-06-05, e5b66bf）**：五 action 升当前大版本，run 验证警告消除 |
| docs/ 历史文档保留旧名 | 历史存档 immutable 惯例 | 不改（roadmap/audit 等是历史记录）|

**下一版候选（未拍板）：** a) 平衡打磨（公开后玩家/自玩反馈驱动）b) 宣传与分发（itch.io / r/golang 分享）c) 新机制探索。项目七个 era 全闭环，也可自然停在这里。

---

### 以下为 V7 规划原文（2026-06-04 启动时写入，存档不改）

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

## V7.2 — 界面美化反馈批次（已收尾 2026-06-05，tag `v7.2`）

> V7 公开发布后的首个反馈批次（对应 V7 收尾记录"下一版候选 a) 平衡打磨"），由用户实玩反馈"界面太丑"驱动。轻量批次，不走完整 era 模板。

**两批修复：**

| 批 | 内容 | commit |
|----|------|--------|
| 根因修复 | 草地竖条纹——`spriteGrass=24` 是泥土+草边过渡 tile（V3 Phase 2 选型错误潜伏 4 版本），换无缝纯草 tile 25 | `c5ff055` |
| 三层美化 | M1 菜单（20pt 标题 + 行内分段上色 + 总星数）/ M2 HUD（cost 实时红绿 + 信息分层）/ M3 地图（装饰确定性散布 7%）| `82129fa` |

**顺带：** Actions Node 24 升级（`e5b66bf`，限期项）；消 V3 未竟项"多字号字体"；README 截图 ×3 重拍；测试 124 → 126。

**方法论沉淀：** `browser_run_code_unsafe` 单段 playwright 代码全流程截图（真 CDP 鼠标 + 零调用间隔），根治三轮踩坑的"MCP 调用间隔跑不过游戏时钟"问题。

---

## V7.3 — 速效 QoL 包（已收尾 2026-06-05，tag `v7.3`）

> 优化矩阵（2026-06-05 全景审视）拍板的速效组合。

| 项 | 内容 |
|----|------|
| C1 暂停 | P 键，同顿帧语义（冻结/渲染输入继续/BGM 不断），实测 3s 冻结 |
| C2 倍速 | F 键 1x/2x，dt×2 全要素同步加速 |
| D1 Frost tint | drawTileTint 冰蓝乘数，场上+按钮一致（销 V5 未竟项）|
| D2 飘字开关 | 并入 J，渲染层过滤零逻辑变化（销 V4 盲点）|
| E1 term '4' | phase 分支接 Frost（销 V5 未竟项）|
| F2 badges | CI/Play/License + 键位表补 P/F |

---

## V7.4 — 工程包（已收尾 2026-06-05，tag `v7.4`）

> 三项全部交付。A1: save_core.go 唯一定义（126 测试零改动回归）；B1: wasm 15.4→12.9MB（踩坑 ×2 入档——net/http 在 wasm +8MB 改 syscall/js；"call to released function" 基线对比定案为 oto 既有无害行为）；C5: 129 tests，首跑产出 3 平衡发现 → 触发 V7.5。

**C5 平衡发现（2026-06-05 仿真数据）：**
1. Normal 太松——中等策略 20 关几乎全 3★（仅 Lv11 掉 1 命）
2. endless 曲线太缓——8 塔成型后撑到 wave 75
3. Hard 曲线倒挂——Lv5/8/11 wave1 失守 + 前期 1★，后期反而全 3★（后期 gold 补偿过度）

---

### 以下为 V7.4 规划原文（启动时写入，存档不改）

> 优化矩阵拍板的工程组合，三项独立交付：

| 项 | 范围 | AC |
|----|------|-----|
| **A1 save 双实现合并** | Save struct + 7 纯逻辑方法抽 shared `save_core.go`，IO 各留各（file / localStorage）——消除 5 轮 family edit 的结构性根因 | 既有测试零改动全过（回归证据）；双文件不再含重复定义 |
| **B1 BGM fetch 异步** | WASM build 的 BGM 不 embed 改运行时 fetch（桌面保持 embed）；workflow/Makefile 同步拷 bgm 到 web/ | wasm 体积 -2.6MB（15.4→~12.8）；BGM 未加载完游戏照跑、到达后起播 |
| **C5 平衡仿真器** | auto-player 贪心策略测试跑 20 关（Normal），输出每关 won/lives/星级——平衡改动从此有回归 | 20 关全部可通关断言；若仿真发现不可通关 → 触发关卡数值修正（仿真器的价值所在）|

**已知翻车场景：** C5 策略写太弱误报"不可通关"/太强掩盖问题——结论只作下限参考；B1 异步引入加载时序态，淡入状态机需容忍 player 迟到。

---

## V7.5 — 平衡调整（已收尾 2026-06-05，tag `v7.5`）

> C5 仿真器三发现驱动的首个平衡批次，两轮"改 → 仿真 → 修订"迭代。

**固化调整：**

| 杠杆 | 变更 | 效果 |
|------|------|------|
| endless HPScale | 5% → 10%/wave（主杠杆是 HP 不是预算）| 中等策略 75 → 42 波 |
| endless 预算 | +n²/6 超线性项 | 后期数量压力 |
| Hard LivesBonus | -2 → -1 | 救前期（Lv5 多撑一波），17/20 保持 |

**第一轮误判与改判（仿真器防错价值实证）：** L11-20 gold 降档对 Normal 无效（钱本来花不完）反而砸死 Lv11 → 撤销；**改判：Normal 中等策略全 3★ 是合理休闲基准，不调**——修正 C5 首跑"Normal 太松"的结论。

**留作设计特性：** Hard Lv8/11 难点关（wave1 重组成 × HP 1.4）；继续内推边际收益低，下一轮由真实手感数据驱动。

**终态画像：** Normal 20/20 全通 / Hard 17/20 曲线顺滑 / endless 42 波（≈8 分钟）。

---

## V8 — 英雄单位（已收尾 2026-06-09，tag `v8.0` @ `601f7c1`）

> V7 公开发布后首个**新机制 era**（区别于 V7.2-7.5 反馈/工程小批次）。用户拍板「gameplay 深度」三候选（多入口 path / 英雄单位 / 塔技能树）中选 **英雄单位**——补上游戏相比真 KR 最缺的「主动操作维度」。LMP L2，加法式扩展，不重写「单路径」承重假设。

### 启动决策（ADR 替代捕获，沿 V5/V6 先例不单设 ADR）

- **决策 A — 控制方案：复用光标 + `H` 键设集结点。** 真 KR 用鼠标拖拽设 rally，本作键盘网格——复用既有光标是唯一不新增移动输入、两端（ebiten+term）一致的方案。反方 WASD 直觉但与建造光标冲突需加模式切换，更糟。
- **决策 B — 阻挡机制：纳入本版，隔离成独立 Phase 3。** 阻挡（敌人贴住停步互殴）是 KR 招牌手感，但改敌人移动状态机是本版最大风险。前两 phase 先把英雄子系统跑通+测试锁定，阻挡风险单独验证；不值可砍 Phase 3 停在非阻挡 MVP。
- **存档：** 英雄为 per-run（同塔），**不动 save schema**；英雄解锁/等级留作未来版本。

### Phase 计划（重构先行 / 测试先行 / 增量验证）

| Phase | 范围 | AC |
|-------|------|-----|
| **P1 英雄实体+移动** | 新 `hero.go`：Hero struct（自由 float Pos / HP / 状态）+ HeroSpec（HP/速度/攻击/射程/冷却/复活）；rally 移动纯逻辑 | 英雄朝 rally 确定性移动；单测锁移动；build+test 绿 |
| **P2 战斗+复活** | 自动打射程内最近敌（**经 `damageEnemy` 统一入口**，守 killEnemy 家族约束）；邻敌对英雄扣血；0HP→死亡→`HeroRespawnCD` 复活（同 MeteorCD 模式）| 击杀走 damageEnemy（金币/特效触发）；死亡+复活；单测锁 |
| **P3 阻挡（风险隔离/可砍）** | Enemy 加 `blocked` 状态：英雄贴路径相邻时敌停步互殴；英雄死释放阻挡 | blocked 敌停步；双向扣血；单测锁 |
| **P4 渲染+输入两端** | ebiten：英雄 sprite+HP条+rally标记，`H` 键设 rally；term：英雄字形+键位；HUD：英雄 HP/复活倒计时 | 实玩可控；两端 build 绿；playwright 线上截图验手感 |
| **P5 仿真+平衡+收尾** | auto-player 仿真器纳入英雄；平衡数值；README/键位表；tag `v8.0` | 仿真 20 关回归绿；英雄不破坏平衡；门面齐 |

### 已知翻车场景（盲点声明）

1. 决策 A 控制手感差 = 整特性废，仿真器**测不出手感**，只能 P4 实玩验证
2. 阻挡（P3）改敌人移动状态机——与 endless/spawner/boss 交互若出 bug，回归面比预期大（隔离独立 phase 即为此）
3. 英雄太强架空塔防 / 太弱无存在感——P5 仿真+实玩双校准

### V8 收尾记录（2026-06-08，实现完成）

**完成度：** 5/5 phase，测试 129 → 148（+19），全程三 build（desktop / term / wasm）+ go vet 保持绿。

**Phase 兑现：**

| Phase | commit | 交付 |
|-------|--------|------|
| P1 实体+移动 | `cdcff18` | hero.go（Hero/HeroSpec/stepToward 纯函数）+ beginRun 生成 + Update 驱动 + SetHeroRally |
| P2 战斗+复活 | `442ec2a` | 经 damageEnemy 攻击 + EnemySpec.Attack/Enemy.meleeCD 反击 + 阵亡复活 |
| P3 阻挡 | `3a7ad60` | 贴身阻挡非飞行敌停步互殴（最高风险 phase，与 endless/spawner/boss 全量绿）|
| P4 渲染+输入 | `6baa258` | 两端英雄绘制 + HUD + H 键；playwright 实测渲染/战斗/阻挡/阵亡复活 |
| P5 仿真+平衡+收尾 | （本提交）| autoPlay heroEnabled 参数化 + HeroNet 守护测试 + README + 收尾 |

**平衡结论（仿真实测，英雄不破坏既有平衡）：**

| 模式 | V7.5（无英雄）| V8（含英雄）| 判定 |
|------|------|------|------|
| Normal | 20/20 全 3★ | 20/20 全 3★ | 不变（已是满命天花板，英雄=冗余保险）|
| Hard | 17/20 | **19/20**（Lv11 仍失守 + Lv5 1★ / Lv7·13 2★）| 变易但仍有挑战 |
| Endless | 42 波 | 42 波 | 不变（失败由 HP 缩放主导，英雄贡献被吞）|

`TestBalance_HeroNetNonNegative` 实证：Hard 有英雄 19 / 无英雄 17 → **英雄增益 +2，且无英雄基线 = V7.5 的 17/20（底层平衡零回归，英雄是纯增量）**。

**手感验证（决策 A 风险点）：** playwright 本地 wasm 实测——英雄渲染 + HUD + H 键生效 + 实战扣血/击杀入账/敌聚集阻挡/阵亡复活倒计时全部观察到；英雄响应 rally 移动（离开隘口致敌涌过 GAME OVER，定位有策略意义）。**仿真器测不出"好不好玩 / 控制顺不顺手"——此项留用户实玩确认，确认后 tag `v8.0` + 部署。**

**英雄数值（占位经仿真验证保留，未重调）：** HP 120 / Speed 4.0 / Damage 15 / Range 1.8 / AttackCD 0.7s / RespawnS 12s；敌近战 ENormal 6 / EFast 4 / EGlider 0（飞行不近战）/ EBoss 25 / ESpawner 8。

**未竟项 / 留作未来版本：**
- 英雄解锁 / 等级 / XP（本版 per-run 满血固定数值，不入存档）
- 英雄专用 sprite（当前程序化金身标记，无 Gopher 英雄美术）
- 多英雄选择 / 英雄主动技能
- rally 引导线干净截图（playwright headless 方向键投递不稳，README 用战斗实拍图代替）

---

## V9 — 英雄成长（已收尾 2026-06-10，tag `v9.0` @ `a3c8be3`）

> V8 英雄上线后的深化方向首批。用户从「英雄深化」四子项（解锁/等级XP、专用 sprite、多英雄、技能树）中拍板 **英雄成长**——关内 XP/等级 + 主动技能，最能放大英雄的"用着有回报"钩子。LMP L2，全加法，不依赖新素材。

### 启动决策（ADR 替代捕获，沿 V5/V6/V8 先例不单设 ADR）

- **决策 A — 等级 per-run（不入存档）。** 每关从 1 级打起，关内击杀升级，关束清零。理由：跨局持久等级会让满级英雄碾压固定难度的早关（KR 靠等级上限+按关缩放压制，本作关卡固定），per-run 既有成长感又不破 V8 既验平衡。持久化/解锁留 V10。
- **决策 B — XP 来源 = 英雄自身击杀。** 仅英雄攻击的致命击给 XP（清晰归属，updateHero 内 target.Dead 判定），威胁加权复用 `enemyCost`（endless.go：Normal1/Fast2/Glider3/Spawner4/Boss10）。
- **决策 C — 主动技能 = AoE 横扫。** 英雄周围 AoE，某级解锁，`G` 键 + 冷却，照搬 Meteor 范式（统一 damageEnemy 结算）。近战群战定位 + 与阻挡协同（挡住一堆再横扫）。
- **砍/留 V10：** 专用 sprite（素材阻塞，Kenney 无 hero）/ 多英雄 / 技能树。

### Phase 计划（重构先行 / 测试先行 / 增量验证）

| Phase | 范围 | AC |
|-------|------|-----|
| **P1 XP+等级** | Hero 加 Level/XP；英雄击杀给 XP（威胁加权）；满阈值升级提 maxHP/damage/range + 回满血；per-run（beginRun 重置）；等级上限 | 击杀累 XP→升级→属性提升+回血；per-run 清零；单测锁曲线 |
| **P2 主动技能** | 某级解锁 AoE 横扫；冷却 + `G` 键触发（经 damageEnemy）；未解锁/冷却中拒绝 | 解锁后可释放打周围敌；冷却 gating；单测锁 |
| **P3 渲染+输入两端** | ebiten+term：HUD 英雄等级/XP 条/技能冷却；`G` 键接线；升级飘字反馈 | 实玩可见等级成长+技能；两端 build 绿；playwright 验 |
| **P4 仿真+平衡+收尾** | 仿真纳入英雄成长（升级后更强）；平衡校准（成长不破 Normal/Hard）；README/键位；tag `v9.0` | 仿真 20 关回归绿；成长后英雄不架空塔防；门面齐 |

### 已知翻车场景（盲点声明）

1. 升级曲线太快 → 英雄中期碾压架空塔防；太慢 → 升级无感。P4 仿真+实玩校准
2. AoE 横扫太强 → 英雄变 Meteor 平替；与 Cannon 溅射 + 阻挡叠加可能过载。P4 校准
3. per-run XP 来源仅"英雄击杀" → 若英雄被群殴速死/抢不到尾刀，XP 积累慢、升不上去（手感问题，仿真难测）

### V9 收尾记录（2026-06-10）

**完成度：** 4/4 phase，测试 148 → 159（+11），三 build（desktop / term / wasm）+ go vet 全程绿。另产 test-infra 解耦（`579f7a6`：2 个 ebiten-only 测试剥到 `//go:build !term`，逻辑测试支持 headless `go test -tags term`——无窗口服务器环境 GLFW init panic 的根治）。

**Phase 兑现：**

| Phase | commit | 交付 |
|-------|--------|------|
| P1 XP+等级 | `02b6fce` | Hero Level/XP + GainXP 升级回血/封顶 + Damage()/AttackRange() 按级缩放 + beginRun 重置 |
| P2 主动技能 | `32804b1` | cleave（L3 解锁 / 8s 冷却 / 半径 2 / Damage×3，经 damageEnemy 统一结算；抽 heroAwardKillXP 单一 XP-grant 点）|
| P3 渲染+输入两端 | `f0d42f6` | 两端 HUD 等级/XP 条/技能冷却 + `G` 键接线；playwright 验出升级偏慢（5 杀仍 L1）→ P4 调曲线 |
| P4 仿真+平衡+收尾 | `a3c8be3` | 仿真纳入成长（HeroLevelReport 诊断 + simStrategy 接 cleave）+ 决策 B 改判 + README/键位 |

**决策 B 两轮证伪改判（同 V7.5 gold 降档反转性质，数据驱动留档）：**

- 启动预设「XP = 英雄自身击杀」→ 仿真证被塔抢尾刀饿死：15/20 关英雄停 L1，cleave@L3 够不着，盲点 3 实锤
- 改「附近助攻」→ 仍 12/20 关停 L1（塔在英雄射程外击杀）
- 终定「**被动 XP**」：英雄在场则每个击杀给 XP（killEnemy 唯一击杀点统一给，威胁加权 enemyCost 保留）。理由：英雄的定位动机在战斗/阻挡而非抢刀，被动化不削弱定位
- 结果成长弧：首关即可升至 L3 解锁 cleave；按关卡难度 L1→L3-L5，关内成长可感知

**平衡结论（仿真实测，成长 + cleave 上界零破坏）：**

| 模式 | V8 | V9（成长 + cleave）| 判定 |
|------|------|------|------|
| Normal | 20/20 全 3★ | 20/20 全 3★ | 不变 |
| Hard | 19/20 | 19/20（Lv11 仍失守，同 squeaker）| 不变（同难点）|
| Endless | 42 波 | 44 波 | 微增 |

HeroNet 增益 +2 保持，无英雄基线仍 17/20（零回归）。**per-run 自限生效**：每关从 L1 打起，难点波发生在英雄弱期，成长救不了场——决策 A 的平衡假设被仿真证实。

**未竟项 / 留作 V10：**
- 英雄专用 sprite（素材阻塞，Kenney 无 hero）
- 多英雄选择 / 技能树
- 跨局持久化等级/解锁（决策 A 留档：需配套按关缩放压制才能开）
- cleave AoE 视觉干净截图（无塔 smoke 升不到 L3，逻辑已全单测覆盖）

---

## V10 — 多英雄选择（已收尾 2026-06-11，tag `v10.0` @ `f9dcd9c`）

> V9 英雄成长收尾后，用户从英雄深化剩余子项（专用 sprite / 多英雄 / 技能树 / 跨局持久化）拍板 **多英雄**——代码层最顺的下一刀（`heroSpec` 单例引用面仅 2 处 game.go），纯加法复用 V8 战斗/阻挡 + V9 成长。技能树+持久化捆绑留 V11 候选「meta 成长」（需配套按关缩放，平衡重做级风险）；真 sprite 继续素材阻塞留档。LMP L2。

### 启动决策（ADR 替代捕获，沿 V5~V9 先例不单设 ADR）

- **决策 A — 阵容三人（用户拍板）。** 现有英雄定名 **Knight**（近战坦/阻挡，数值不动 = 零回归基线）；新增 **Archer**（远程 ~3.5 ≈ Archer 塔射程 / 低 HP / 不阻挡，风筝输出）+ **Rogue**（高速 5.5 追得上 EFast / 低 HP / 高攻速 / 阻挡，游走截击）。风险隔离：Rogue 最后加，手感重叠可砍。
- **决策 B — 阻挡按职业 gating。** `HeroClass.Blocks` 字段。Archer 不阻挡是刻意取舍（会肉搏的弓手既怪又抹平职业差异）；反方风险 = 不阻挡英雄「没存在感」，P4 仿真+实玩双验。
- **决策 C — 选择入口沿 difficulty 先例。** 菜单键位循环切换 + `Save.HeroChoice` 零值兼容持久化（0 = Knight，旧存档行为不变）。
- **决策 D — 技能统一「自身 AoE」范式。** per-hero 只调参数（半径/倍率/冷却），不做异构技能（瞄准型留技能树 era）——控制本版复杂度。
- **视觉：** 程序化配色区分（Knight 金 / Archer 绿 / Rogue 紫），真 sprite 留档。

### Phase 计划（重构先行 / 测试先行 / 增量验证）

| Phase | 范围 | AC |
|-------|------|-----|
| **P1 参数化重构** | `heroSpec` 单例 → `HeroClass`（spec + 成长增量 + 技能参数 + Blocks）；Hero 持 class 引用；Knight = 现数值 | 零行为变化，159 tests 全绿不改断言 |
| **P2 新英雄** | Archer/Rogue 数值定义；阻挡 per-class gating；技能参数化生效 | 三职业差异可测（射程/速度/阻挡/技能）；单测锁每职业 |
| **P3 选择UI+存档+视觉两端** | 菜单循环切换键 + `Save.HeroChoice`；程序化配色 + HUD 职业名；ebiten+term | 实玩可选可见；旧存档零值 = Knight；两端 build 绿 |
| **P4 仿真+平衡+收尾** | 仿真矩阵 per-hero × Normal/Hard/Endless；HeroNet 扩 per-class；README/键位；tag `v10.0` | 三职业各自不破平衡且互有差异；Knight 基线零回归；门面齐 |

### 已知翻车场景（盲点声明）

1. Archer 不阻挡 → 可能「没存在感」（输出被塔淹没、又不守隘口）；仿真测得出净增益、测不出存在感，留实玩
2. 三职业手感差异做不出来 = 换皮；Rogue 与 Knight 重叠风险最高（都近战阻挡），数值上必须拉开（速度/攻速 vs 坦度）
3. 平衡矩阵 ×3：某职业在 Hard 净增益为负（弱于无英雄）= 设计失败信号，需重调而非接受

### V10 收尾记录（2026-06-11）

**完成度：** 4/4 phase，测试 159 → 173（+14，headless 171），三 build + vet 全程绿，每 phase 一 commit。

**Phase 兑现：**

| Phase | commit | 交付 |
|-------|--------|------|
| P1 参数化重构 | `94968f2` | HeroClass 参数包（数值/成长/技能/Blocks），测试纯机械符号替换零断言变化 |
| P2 新英雄 | `3ae2a18` | Archer（远程 3.5 / 不阻挡）+ Rogue（速 5.5 / 攻速 0.35 / 阻挡 / 10s 复活）+ 8 测 |
| P3 选择UI+存档+视觉 | `45c276b` | 菜单 H 循环 + `Save.HeroChoice` 零值兼容 + 两端配色金/绿/紫 + K/A/R 字形 |
| P4 仿真矩阵+校准 | `f9dcd9c` | autoPlayClass 参数化 + per-class Normal/Hard 矩阵 + 净非负守护 |

**平衡校准（盲点 3 实证反向，一轮收敛）：** 预警的是「某职业净负」，实际首跑 Archer/Rogue 净增益 +3 强于 Knight +2、Hard 双双 20/20 抹平 V7.5 难点关 Lv11 → 调参（Archer dmg 12→9 / Rogue 9→7）后三职业 Hard 全 19/20、净 +2 持平、Normal 全通。**职业差异落在风格（射程/速度/阻挡/复活/技能形状），不落数值强度。**

**决策 B 行为差异实证（playwright 冒烟）：** Archer 无塔裸守 → 敌全漏 GAME OVER；V9 同场景 Knight 可独守 wave 1。职业选择有真实策略后果，非换皮。

**未竟项 / 留档：**
- Rogue 与 Knight 实玩手感重叠度（用户确认满意，但长期观察项）
- 真 sprite（素材阻塞不变）；技能树 + 持久化 → V11「Meta 成长」
- cleave / rally 干净截图债（V8/V9 遗留，依旧未还）

---

## V11 — Meta 成长（已收尾 2026-06-16，tag `v11.0` @ `cda346d`）

> 英雄深化收官之作，用户拍板（候选对比：多入口 path 换轨道 / itch.io 分发 / 音画打磨）。**核心设计转向：用「perk 预算」替代「按关缩放」**——此前评估认为 meta 成长须按关缩放配套（动 20 关平衡面），前提是持久化*原始等级*；改判为持久层只存**技能树点**：V9 per-run 关内等级原样保留，跨局积累的是星星 → 技能树节点，perk 有顶（仿真上下界验证），平衡冲击从「重做 20 关」缩到「校准一张 perk 表」。LMP L2-L3。

### 启动决策（ADR 替代捕获，沿 V5~V10 先例不单设 ADR）

- **决策 A — 货币 = 星。** `Save.Stars` 现成（20 关 × 3★ = 60 可赚上限），老存档存量星直接可用（奖励既有进度；反方风险 = 一次性爆买冲击，接受）。
- **决策 B — 树形态 = per-class 线性 4 节点（用户拍板）。** 每职业一条小树 ×3 = 12 节点，职业纵深延续 V10 风格差异；列表式 UI 不画图形树（term 可操作）。
- **决策 C — 定价制造取舍。** 单职业满树 30 星（3/6/9/12），三树总价 90 ≈ 1.5× 可赚上限 60 → 全解锁不可达，必须选 build。
- **决策 D — 满树允许碾压 Hard（用户拍板，哲学转向留档）。** V7.5/V10 两代守护的「Lv11 难点关」让位于 meta 成长的奖励本质——花 60 星征服难关是 earned power。**硬约束不变：无树基线 = V10 行为零回归**（不买节点/旧存档体验完全不变，Hard 仍 19/20）。
- **决策 E — 入口 = 菜单 `T` 键，新 `PhaseSkillTree`。** endless 之后首个新增 Phase；方向键选节点 + Space 购买 + M 返回。
- **效果快照时机：** 购买只在菜单发生 → beginRun 时把已购 perk 聚合成 `HeroBonus` 快照进 Hero，关内不变。

### Phase 计划（重构先行 / 测试先行 / 增量验证）

| Phase | 范围 | AC |
|-------|------|-----|
| **P1 持久层+购买逻辑** | `Save.TreeNodes`（map[职业名]已购数，零值兼容）；节点表（名称/描述/价格）；TotalStars/Spent/Available 预算；BuyTreeNode gating | 预算算术单测锁；不足/越界/已满拒绝；roundtrip；旧存档零行为变化 |
| **P2 效果接线** | `HeroBonus` 聚合（HP/伤害/射程/速度/复活/技能参数增量）；Hero accessors 加成；beginRun 快照 | 每节点效果单测锁；零购买 = 零 bonus 守护 |
| **P3 树 UI 两端** | 菜单 `T` 入口 + `PhaseSkillTree` 列表屏（职业切换/节点导航/购买/星余额）；ebiten+term | 实玩可买可见生效；两端 build 绿；playwright 验 |
| **P4 仿真+平衡+收尾** | 仿真矩阵：无树基线（零回归守护）+ 满树上界 × 三职业 × Normal/Hard；perk 幅度校准；README/键位 | 无树 = V10 数据零回归；满树 Normal 全通、Hard 允许 20/20；门面齐 |

### 已知翻车场景（盲点声明）

1. 校准矩阵比 V10 大（树状态 × 职业 × 难度），可能不止一轮收敛；V10 一轮命中有运气成分
2. 新 Phase 状态机边角：树屏中途退出/半购买状态/与 endless·难度键的交互——menu 输入分支首次变三态
3. perk 太弱 = 白做（买了无感），太强 = 无树/满树两个世界（新玩家被劝退）；幅度全靠 P4 仿真+实玩双校准
4. 老玩家 60 星一次性爆买 → 直接满树体验，错过渐进成长弧（接受，见决策 A）

### V11 收尾记录（2026-06-16）

**完成度：** 4/4 phase，测试 173 → 188（+15，headless 186），三 build（desktop / term / wasm）+ go vet 全程绿，每 phase 一 commit。

**Phase 兑现：**

| Phase | commit | 交付 |
|-------|--------|------|
| P1 持久层+购买 | `7db4abe` | `skilltree.go`（星预算 + `BuyTreeNode` gating + `Save.TreeNodes` 零值兼容），「60 星恰好两棵树、第三棵永远差钱」取舍写进测试 |
| P2 效果接线 | `9548569` | `HeroBonus` 七字段聚合 + Hero 派生数值方法化 + beginRun 快照；零 bonus = V10 基线零回归守护测试 |
| P3 树 UI 两端 | `548d530` | 新 `PhaseSkillTree` + 菜单 `T` 入口 + 两端列表屏（职业色/选中金边/已购绿/可买亮/锁定灰）；playwright 冒烟 |
| P4 仿真+校准 | `cda346d` | `autoPlayTree` 上下界矩阵 + Knight 树一轮校准 |

**核心改判（设计转向，本 era 的胜负手）：** 「perk 预算」替代「按关缩放」——原评估认为 meta 成长须按关缩放配套（持久化*原始等级* → 满级碾压固定难度早关 → 动 20 关平衡面），属最重候选。改判为持久层只存**技能树点**，V9 per-run 关内等级原样保留，跨局积累的是星 → perk（有顶），平衡冲击从「重做 20 关」缩到「校准一张 12 节点 perk 表」——最重的刀变成单日交付。

**平衡校准（盲点 3「太弱白做」实证，一轮收敛）：** 首跑 Knight 满树打 Hard 增益 +0（花满 30 星无感，Archer/Rogue 已 +1）→ 调 Knight 树（Sharpened Blade +4→+6 dmg / Undying 复活 -4s→-6s）→ 三职业满树统一 Hard 20/20（+1，征服 Lv11 难点关）。**无树基线 19/20 零回归**（决策 D earned power 兑现 + 硬约束守住）。

**平衡结论：**

| 状态 | Normal | Hard | 守护 |
|------|--------|------|------|
| 无树（= V10 基线）| 20/20 | 19/20 | 零回归硬约束 |
| 满树（三职业）| 20/20 | 20/20 | earned power，不破 Normal |

**未竟项 / 留档：**
- 老玩家 60 星爆买体验（接受，决策 A 留档）
- 真 sprite（素材阻塞不变）；cleave / rally 干净截图债（V8/V9/V10 累积未还）
- 技能树图形化呈现（当前列表式，term 兼容优先）

**英雄深化四子项收官：** V9 成长 / V10 多英雄 / V11 技能树+持久化全部交付，仅「专用 sprite」因素材阻塞悬留。

---

## V12 — 多入口 path（已收尾 2026-06-22，tag `v12.0` @ `4494364`）

> 英雄深化收官后，用户从广度/对外候选拍板 **多入口 path**——gameplay 广度最大的一刀，也是从 V0 埋下的「单路径」承重假设的首次正面重构。研究发现承重面收敛：`Enemy.Pos`/`pickTarget`/`pathLerp`/`pathDir` 在 V4/V5 重构时**已参数化 `path []Point`**（不读全局），故多 path 的核心计算无需重写，触点（57 处/7 文件）多为「调用处实参从 `g.Path` 换 `g.Paths[e.PathID]`」的机械替换。LMP **L3**（架构级承重假设重构）。

### 启动决策（ADR 替代捕获，沿 V5~V11 先例不单设 ADR）

- **决策 A — 范围：全 20 关重做双路（用户拍板）。** 最大 gameplay 变化，也是最大回归面——20 关数据全重写 + 平衡全重校。反方（加法式可选第二条、风险隔离）被否，用户要彻底的多路体验。
- **决策 B — 分路：自动均摊。** spawn 按 `spawned` 计数轮流分配 `PathID`（偶 → path0 / 奇 → path1）。确定性可测，**wave DSL 零改动**（老 wave 数据全复用，只新增第二条 cps）。
- **决策 C — 终点：共享汇流（用户拍板）。** 实现避开「汇点索引映射」几何：**数据层重合尾段**——两条 cps 末尾汇合到同一终点 cell，`ExpandPath` 各自展开后尾段 cells 重合；每个敌人仍沿自己的完整 path 走（含重合尾段），`escaped` 判定（`PathIdx ≥ len(path)-1`）不变；渲染去重（`pathLookup` 本是 set）画成一条。汇流复杂度降到与独立终点几乎相同。
- **决策 D — endless 保持单路。** `Paths = [][]Point{path}`，`PathID` 恒 0，均摊降级为全 0。多路只用于手工 20 关，endless 生成器不动（降范围）。
- **决策 E — 数据模型：** `Game.Path []Point` → `Paths [][]Point`；`Enemy` 加 `PathID int`；`Level` 新增可选 `CPS2`（空 = 单路，endless/兼容用）；`pathLookup` = 所有 path cells 并集（塔/decor 避让，天然去重）。

### Phase 计划（重构先行 / 测试先行 / 增量验证）

| Phase | 范围 | AC |
|-------|------|-----|
| **P1 数据模型多路化（零行为变化）** | `Paths [][]Point` + `Enemy.PathID` + `Level.CPS2`；单路退化（`Paths=[]{path}` 全 PathID=0）；57 触点 `e.Pos(g.Path)`→`e.Pos(g.Paths[e.PathID])`；Spawner 召唤继承父 PathID | 188 tests 全绿**不改断言**（单路退化 = 零回归硬约束）|
| **P2 spawn 均摊 + 汇流渲染两端** | spawn 按计数分配 PathID（单路恒 0）；渲染遍历所有 path 去重 cells + 多起点 S / 共享终点 E；ebiten+term | 双路 fixture 敌均摊单测；两端 build 绿；playwright 验汇流 |
| **P3 全 20 关双路数据 + 汇流** | 每关 yaml 加 cps2 汇合到终点；level 数据测试（双路连通 + 尾段真重合 + 汇点一致）| 20 关皆双路且汇流；path 校验测试兜底 |
| **P4 仿真+平衡+收尾** | 双路削弱英雄阻挡（只守一路）+ 塔覆盖局部 → 平衡全重校；仿真 20 关双路回归；README/收尾；tag `v12.0` | 双路 Normal 全通、Hard 合理、英雄净非负；门面齐 |

### 已知翻车场景（盲点声明）

1. **全 20 关平衡重校**——敌分两路密度减半，但塔/英雄只覆盖局部；每关都可能要调，校准轮次远超 V10/V11（V10 一轮、V11 一轮的运气可能用完）
2. **汇流靠数据纪律**——两条 cps 末尾必须真对齐到同 cells，否则渲染出两个终点 / 连通断裂；P3 level 数据测试强制兜底（尾段重合 + 汇点一致断言）
3. **英雄阻挡只能守一路**——多路天然削弱英雄（V8-V11 的英雄价值部分架空）；设计上合理（迫使取舍）但可能过度，P4 仿真 + 实玩双校
4. **汇流段双倍密度**——两路敌在尾段叠加，尾段可能过载 / 成为唯一有效防守点，使「多入口」沦为「单出口」假多路；趣味性存疑，实玩验
5. **57 触点机械替换易漏**——`e.Pos(g.Path)` 漏改一处 = 该敌按错 path 取位（错位/穿模）；P1 grep 清零 + 单路零回归测试兜底

### V12 收尾记录（2026-06-22）

**完成度：** 4/4 phase，测试 188 → 195（+7），三 build（desktop / term / wasm）+ go vet 全程绿。

**Phase 兑现：**

| Phase | commit | 交付 |
|-------|--------|------|
| P1 数据模型多路化 | `c3b9fbe` | `Paths [][]Point` + `Enemy.PathID` + `Level.CPS2`；57 触点机械替换；188 tests 零断言变化（单路退化 = 零回归） |
| P2 spawn 均摊+汇流渲染 | `7634da3` | spawn 按计数交替分配 PathID；渲染遍历所有 path + 多起点 S / 共享终点 E；6 multipath 测试 |
| P3 全 20 关双路数据 | `4494364` | 20 关各添加 `cps2`（第二入口，末尾汇合同终点）；`TestLevels_DualPathIntegrity` 校验双路连通+汇流+地图范围 |
| P4 仿真+平衡校准 | `4494364`（同 P3 合并提交） | 修复 `rankedTowerSpots` 仅覆盖 `Paths[0]` → 覆盖所有路径；Normal 20/20 全通 |

**平衡结论（仿真实测，双路削弱英雄阻挡但整体平衡健康）：**

| 模式 | V11（单路） | V12（双路） | 判定 |
|------|-------------|-------------|------|
| Normal | 20/20 全 3★ | 20/20 全通（18 关 3★） | 稳 |
| Hard | 19/20 (Lv11 难点) | 19/20 (Lv9 难点) | 难点关移位，仍有挑战 |
| 无英雄基线 | 17/20 | 15/20 | 双路削弱英雄（预期，迫使取舍） |

**盲点兑现：**
- 盲点 1（平衡重校）：auto-player bug 修后**零轮手调**收敛（仿真器路径覆盖修正 = 根因，不是数值问题）
- 盲点 2（汇流数据纪律）：`TestLevels_DualPathIntegrity` 20 关全过，尾段真重合
- 盲点 3（英雄只守一路）：无英雄基线 15/20（-2），设计意图实现（迫使取舍）
- 盲点 5（57 触点漏改）：P1 grep 清零 + 单路零回归 188 tests 兜底

**未竟项：** 无。

---

## V13 — 新敌人类型（已收尾 2026-06-22，tag `v13.0` @ `21cc588`）

> V12 双路完成后的 gameplay 深度扩展。从「新敌人 / 真 sprite / Boss 机制」三候选中选 **新敌人类型**——纯加法，每种克制不同防守策略，与双路系统协同丰富战术组合。LMP L2。

### 启动决策（沿 V5~V12 先例不单设 ADR）

- **决策 A — 三种新敌人，各克制一种策略。** EShield 克制 Archer 速射；ERegen 克制分散火力；EHealer 克制被动防御。
- **决策 B — 渐进引入。** L8 起 Shield，L10 起 Regen，L13 起 Healer，L20 全家桶。替换部分现有敌人而非增加总量，保守不破平衡。
- **决策 C — endless 按波次解锁。** wave 4 Shield / 6 Regen / 7 Healer / 8 Boss，预算权重 Shield 3 / Regen 3 / Healer 4。

### 三种新敌人

| 类型 | DSL | HP | 速度 | 护甲 | 回血/s | 治疗CD | 奖励 | 克制 |
|------|-----|-----|------|------|--------|--------|------|------|
| EShield | `d` | 30 | 2.0 | 4 | — | — | 20 | Archer 低伤速射（8-4=4 实际伤害） |
| ERegen | `r` | 40 | 2.5 | — | 3 | — | 22 | 分散火力（需集火或爆发击杀） |
| EHealer | `h` | 20 | 2.5 | — | — | 2.0s | 30 | 被动防御（治疗半径 2，+5hp 最低血量盟友） |

### 机制实现

- **护甲**：`damageEnemy` 统一入口减免 `Armor` 伤害（min 1），高伤穿甲（Cannon 25-4=21）
- **回血**：Update 循环每秒 +`Regen` HP（不超 MaxHP）
- **治疗**：Update 循环按 `HealCD` 周期治疗半径内最低血量盟友 +5hp，绿色飘字

### V13 收尾记录（2026-06-22）

**完成度：** 一次性交付（entities + 机制 + wave 数据 + 测试），测试 195 → 202（+7），三 build + vet 绿。

**新增测试：** ParseWave V13 DSL / Shield 护甲减免 ×3（标准/最低 1/高伤穿甲）/ Regen 回血 + 不超 Max / Healer 治疗盟友。

**平衡结论（仿真 202 tests 全过）：**

| 模式 | V12 | V13 | 判定 |
|------|-----|-----|------|
| Normal | 20/20 | 20/20 | 不变 |
| Hard | 19/20 | 19/20 | 不变 |
| 三职业净非负 | K19/A20/R19 | K19/A20/R19 | 不变 |

**替换策略验证：** 新敌人替换部分旧敌人（不增总量），auto-player 策略自然应对（Cannon/Magic 高伤穿甲 Shield，塔覆盖面杀 Regen/Healer），平衡零手调。

**未竟项：** 无。

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
| 2026-06-05 | V7 收尾：tag `v7.0`（be86c3c），版本史补 V7 行，收尾记录（执行记录 ×4 + 未竟项 ×4），游戏公开上线；下一版三候选（平衡打磨/宣传分发/新机制，未拍板，可自然停）|
| 2026-06-05 | V7 未竟项销账：Actions Node 20 弃用警告已修（e5b66bf，五 action 升级 + run 验证），限期项清零 |
| 2026-06-05 | V7.2 收尾：tag `v7.2`（82129fa），界面美化反馈批次归档（条纹根因 + M1/M2/M3），V3 未竟项"多字号"销账，126 tests |
| 2026-06-05 | V7.3 收尾：tag `v7.3`（a5a8582），速效 QoL 包（暂停/倍速/Frost tint/飘字开关/term '4'/badges），销 V4 盲点 ×1 + V5 未竟项 ×2 |
| 2026-06-05 | V7.4 启动：工程包三项（save 合并 / BGM 异步 / 平衡仿真器），优化矩阵拍板 |
| 2026-06-05 | V7.4 收尾：tag `v7.4`（ee75ea0），三项交付 + C5 三平衡发现入档 |
| 2026-06-05 | V7.5 启动：平衡调整（用户拍板）——主杠杆 = 后期关 gold 降档（修 Normal 后期松 + Hard 倒挂同根因）+ endless 预算超线性化；数据驱动迭代（改 → 仿真 → 调），目标 Normal 全通且星级有分布 / endless 20-35 波 / Hard 曲线顺滑 |
| 2026-06-05 | V7.5 收尾：tag `v7.5`（1b966e3）。注：启动时预设的"gold 降档"主杠杆被第一轮仿真证伪并撤销（Normal 全 3★ 改判为合理基准），最终杠杆 = endless HPScale + Hard lives——规划假设被数据修正的首例，留档 |
| 2026-06-08 | V8 启动：用户拍板新机制 = 英雄单位（gameplay 深度三候选中选）。决策 A 控制 = 光标+H 设 rally；决策 B 阻挡纳入并隔离 P3。5 Phase 规划（实体+移动 → 战斗+复活 → 阻挡 → 渲染+输入两端 → 仿真+平衡+收尾）；英雄 per-run 不动 save schema；沿 V5/V6 先例不单设 ADR |
| 2026-06-08 | V8 实现完成：5/5 phase 交付（cdcff18/442ec2a/3a7ad60/6baa258 + 601f7c1），测试 129→148。平衡仿真证英雄是纯增量（Hard 17→19，无英雄基线零回归）|
| 2026-06-09 | V8 部署上线：用户选「直接部署线上玩着确认」→ push master（Actions run 27178141887 绿，test gate 148 过 + WASM 13.3MB + deploy-pages），线上冒烟确认英雄上线。**tag `v8.0` 待用户线上实玩手感确认**（满意则 tag 归档；需调走 V8.1 tuning 批次）|
| 2026-06-09 | V8 收尾：用户实玩确认手感满意 → tag `v8.0` @ `601f7c1`，版本史表补 V8 行 + 收尾记录定稿，README/spool 同步。八 era 闭环（v1.0~v8.0 / 148 tests）|
| 2026-06-09 | V9 启动：用户从英雄深化四子项拍板「英雄成长」。决策 A 等级 per-run（不入存档不破平衡）/ B XP=英雄击杀威胁加权 / C 主动技能=AoE 横扫 G 键照搬 Meteor。4 Phase（XP+等级 → 主动技能 → 渲染+输入 → 仿真+平衡+收尾）；sprite/多英雄/技能树留 V10 |
| 2026-06-09 | V9 P1-P4 实现完成 + 部署：决策 B 经两轮仿真证伪改判（自身击杀→助攻→**被动 XP**，前两者被塔抢尾刀饿死 12-15/20 关停 L1）；test-infra 解耦 ebiten-only 测试支持 headless。push master（Actions 27190988283 绿，test gate 159 过 + WASM 13.32MB），线上冒烟确认。**tag `v9.0` 待用户线上实玩成长/cleave 手感确认** |
| 2026-06-10 | V9 收尾：用户实玩确认成长/cleave 手感满意 → tag `v9.0` @ `a3c8be3`，版本史表补 V9 行 + 收尾记录定稿（决策 B 改判留档），README/spool 同步。九 era 闭环（v1.0~v9.0 / 159 tests）|
| 2026-06-11 | V10 启动：用户拍板英雄深化 → 多英雄（技能树+持久化捆绑留 V11「meta 成长」/ 真 sprite 素材阻塞留档）。阵容三人 Knight+Archer+Rogue（用户拍板）；决策 B 阻挡 per-class gating（Archer 不阻挡）/ C 选择入口沿 difficulty 先例 + `Save.HeroChoice` 零值兼容 / D 技能统一自身 AoE 只调参。4 Phase（参数化重构 → 新英雄 → 选择UI+存档+视觉 → 仿真平衡矩阵+收尾）|
| 2026-06-11 | V10 P1-P4 实现完成 + 部署：P1 重构（`94968f2`，157 headless 零断言变化）→ P2 新英雄（`3ae2a18`）→ P3 选择UI+存档+视觉（`45c276b`）→ P4 仿真矩阵+校准（`f9dcd9c`）。**首跑翻车点 3 实证反向**：Archer/Rogue 净增益 +3 强于 Knight +2、Hard 20/20 抹平 Lv11 难点关 → 一轮调参（Archer dmg 12→9 / Rogue 9→7）三职业归一 19/20 净 +2，差异落在风格不落数值。playwright 冒烟：菜单 H 循环/职业色/A 字本体/Archer 不阻挡（无塔 GAME OVER vs V9 Knight 独守，决策 B 行为差异实证）。173 tests。**tag `v10.0` 待用户线上实玩三职业手感确认** |
| 2026-06-11 | V10 收尾：用户实玩确认三职业手感满意 → tag `v10.0` @ `f9dcd9c`，版本史表补 V10 行 + 收尾记录定稿，README/spool 同步。十 era 闭环（v1.0~v10.0 / 173 tests）。V11 方向用户拍板「Meta 成长」（技能树+持久化捆绑）|
| 2026-06-11 | V11 启动：核心改判「perk 预算」替代「按关缩放」（持久层存技能树点不存等级，per-run 等级保留）。决策 A 货币=星（60 可赚）/ B per-class 线性 4 节点（用户拍板）/ C 总价 90 = 1.5× 预算制造取舍 / D 满树允许碾压 Hard（用户拍板哲学转向，无树基线零回归硬约束）/ E 菜单 T 键 + PhaseSkillTree。4 Phase（持久层+购买 → 效果接线 → 树UI两端 → 仿真上下界+校准）|
| 2026-06-11 | V11 P1-P4 实现完成 + 部署：P1 持久层+购买（`7db4abe`，60 星恰好两棵树的取舍实锤入测）→ P2 HeroBonus 七字段接线（含零 bonus = V10 基线零回归守护）→ P3 PhaseSkillTree 两端 + playwright 冒烟（三列职业色/选中金边/0 星拒绝反馈/localStorage 跨会话恢复）→ P4 仿真上下界。**盲点 3「太弱白做」实证**：首跑 Knight 满树 Hard +0（30 星无感）→ 校准（+6 dmg / 复活 -6s）三职业满树统一 20/20（+1，earned power 兑现），无树基线 19/20 不变。188 tests。**tag `v11.0` 待用户线上实玩技能树手感确认** |
| 2026-06-16 | V11 收尾：用户实玩确认技能树手感满意 → tag `v11.0` @ `cda346d`，版本史表补 V11 行 + 收尾记录定稿（perk 预算改判 + 盲点 3 校准留档），README/spool 同步。十一 era 闭环（v1.0~v11.0 / 188 tests）。英雄深化四子项收官（成长/多英雄/技能树+持久化交付，仅专用 sprite 素材阻塞悬留）|
| 2026-06-17 | V12 启动：用户拍板多入口 path（gameplay 广度，单路径承重假设首次重构，LMP L3）。研究证承重面收敛（Pos/pickTarget/pathLerp 已参数化 path）。决策 A 全 20 关双路（用户拍板，最大回归面）/ B 自动均摊（wave DSL 零改动）/ C 共享汇流（用户拍板，数据层重合尾段免几何）/ D endless 保持单路 / E Paths[][]Point + Enemy.PathID + Level.CPS2。4 Phase（数据模型多路化零回归 → spawn 均摊+汇流渲染 → 全 20 关双路数据 → 仿真+平衡全重校+收尾）|
| 2026-06-22 | V12 收尾：P3+P4 合并交付（`4494364`），20 关 cps2 全入 + `TestLevels_DualPathIntegrity` 校验 + `rankedTowerSpots` 多路覆盖修复。Normal 20/20、Hard 19/20（Lv9 难点）、无英雄 15/20（双路削弱预期）。tag `v12.0` @ `4494364`，195 tests |
| 2026-06-22 | V13 一次性交付（`21cc588`）：敌型 5→8（EShield 护甲 / ERegen 回血 / EHealer 治疗），damageEnemy 加 armor 减免 / Update 加 regen+healer AI / wave DSL 扩展 d/r/h / L8 起渐进混入 20 关 / endless 按波次解锁。仿真 202 tests 全过（Normal 20/20、Hard 19/20 不变）。tag `v13.0` @ `21cc588` |
