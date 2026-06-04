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

### V3 未竟项（不阻塞收尾，归入 backlog）

| 项 | 来源 | 去向 |
|----|------|------|
| 多字号字体（title 大字 / body 小字）| V3 Phase 5c commit 未做段 | V4 不做；未来 UI 需求触发时再做 |
| Tank enemy | V3.6 commit 跳过段——Kenney pack 仅 3 个 walker sprite 已占满 | 受素材限制搁置，换素材包时重评 |
| 新 Tower 型 | V3.6 commit 跳过段——同类 tower sprite 有限 | 同上 |

---

## V4 — 音频 + Game Feel（active，2026-06-04 拍板）

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

## 变更记录

| 日期 | 变更 |
|------|------|
| 2026-06-04 | 初建：V0→V3 版本史回填 + V3 收尾（v1.0/v2.0/v3.0 tag）+ V4 方向拍板（音频 + Game Feel）+ Phase 1-5 规划 |
