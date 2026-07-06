# kingdom-rush — Spool

会话级日志（项目层）。每个主题完成时通过 `/spool` 追加。会话末用 `/wrap-up` 兜底。

详见 standard `docs/concepts/spool-fallback.md`。

---

## 2026-05-18 — audit task 闭环 + adjudication 全套落地（B + C）

**入口:** 用户问 "还有什么没做",列覆盖矩阵识别 7 行 backlog（A 提交 / B task 闭环 / C adjudication / D spool / E PM 端 Issue close / F mock release-history / G bootstrap 缺 ADR）,前 4 行进 backlog,后 3 行丢弃。本次推进 B + C + D。

### B — audit task 文件闭环

- `docs/audit/task-2026-05-15-issue-process.md` 状态 `进行中` → `已完成`（完成时间 2026-05-15 18:33,闭环回填 2026-05-18）
- 14 个 checkbox 全勾,任务日志从 `audit/2026-05-15-issue-process-log.md` 抽要点回填
- 路径误写修正: `audit/problem-registry.md` → `docs/problems/problem-registry.md`
- Act 段记 3 个收尾反馈,本次违规自身（checkbox 未及时勾选）候选 P-004,待用户判是否升级

### C — review adjudication（IPR-001/002/003）

- 用户拍板:全按 AI 候选 verdict（confirmed / dismissed / deferred）
- 落地 7 项: registry 状态流转 ×2 + audit 报告勘误段 + KR-FB-005 occurrences 1→2 + FB-006 衍生 + handoff status approved + 归档 completed/2026-05/
- handoff 文件: `docs/handoff/completed/2026-05/2026-05-18-adjudication-issue-process-2026-05-15.md`
- FB-006 候选: `docs/feedback/2026-05-18-audit-skill-exemption-detection-fb.md`（audit skill 维度 3 缺豁免 comment 识别能力）

### 关键发现（事实级）

1. **audit 自荐基于过期判断** — audit 报告 L126 推荐"升 KR-FB-005 severity → high",实际 KR-FB-005 已是 high（feedback 文件 L240）。原因:audit 跑的时间和 FB 升级时间错位,audit skill 无机制感知 FB 当前 severity。[已验证: 两文件交叉]
2. **audit skill 维度 3 缺豁免 comment 识别能力** — IPR-002 是 false positive 主因。豁免 comment #4458847614 已贴但 audit 机械检查不识别。跨项目普适。[二手: IPR-002 / KR-FB-003]
3. **task-lifecycle "立即勾选" 在 AI 主导执行时易被跳过** — kingdom-rush 首次 audit task 拖 3 天才闭环,违反规则本身（候选 P-004,层级项目级 / 规则级未拍板）。[已验证: task 文件 commit history]

### 未做遗留（非本会话推进）

- **A（最高优先）** — 全部改动 uncommitted: CLAUDE.md / docs/feedback/{,新 FB} / docs/release-history.md / task 文件 / findings-registry / problem-registry / audit 报告 / KR-FB-005 occurrences / handoff 归档 移动。下次会话起点。
- (a) standard 仓库 SKILL.md issue Step 6 commit message 关键词禁用列表 — 非本项目 scope,待 standard 维护方
- (b) `docs/problems/fb-index.md` 仍是骨架,KR-FB-001~005 + FB-006 未录入 — 独立 backlog
- (c) IPR-003 / P-003 deferred 触发条件: active phase 引入第一个有逻辑分支函数时改 confirmed → fixing — 长期 backlog
- (d) IPR-001 在 `confirmed` vs `merged-into KR-FB-005` 之间存在边界歧义,standard 文档对 "FB 实例 vs 独立 finding" 未明示 — 已在 handoff 标注盲点

### 下次会话建议起点

按优先级: **A（commit + push）** > fb-index 补录 > standard 回传 > P-004 拍板（task-lifecycle 违规升级）


## 2026-06-03 — backlog 清零（V3 Phase 6 提交 + issue #1 闭环 + fb-index 补录 + registry 批处理）

**入口:** 用户问 "还有什么没做",列 10 行覆盖矩阵（4 保留 / 2 已完成移出 / 2 丢弃 / 1 改判 / 1 盲点）,逐项推进 1→4 全部完成。

### 1 — commit + push（d225fb8 / 1b52c85）

- `feat: V3 Phase 6` — path 程序化绘制（中心圆+方向矩形 overlay 替换 dirt tile, grass 全铺, S/E 圆形 marker）,三 build + 39/39 tests 验证后提交
- `docs(feedback)` — 2 个 FB 文件「上报记录」段（2026-05-22 已上报 standard#11）,闭环 2026-05-18 遗留项 (a)
- 未做: path outline polish（pathEdgeCol 占位）,留 backlog

### 2 — issue #1 显式 close（comment #4611226585）

- 4 字段 close comment（commits / CHANGELOG 不涉及 / release-log 指向 build #1 / 关联 KR-FB-005+P-001）
- timeline 核实: 2026-05-15 曾被 commit `0859892` `closes #1` auto-close（10:04Z）→ reopen（10:10Z）,与 KR-FB-005 记录一致 [已验证: gh api timeline]
- 显式 `gh issue close`,不靠关键词

### 3 — fb-index 补录（c730116）

- KR-FB-001~005 + FB-006 全部入索引,删示例骨架,统计 0 → 6
- 编号决策: 保留原编号不重号（已被 registry / standard#11 交叉引用）,新条目起 ADR-008 `FB-YYYYMMDD-<hash>` 双轨
- 闭环 2026-05-18 遗留项 (b)

### 4 — registry 状态流转批处理 ×4（2a4c946, 用户授权"按你想法先做"）

- **P-004 入册** confirmed / 规则级（task-lifecycle checkbox 纪律在 AI 主导执行时失效,跨项目普适）
- **P-001 / IPR-001 → resolved**（issue #1 闭环 + git-workflow §5.2 沉淀 + standard#11 上报）
- **P-003 / IPR-003 → resolved**（前提不成立: V1.6 起测试累积,39/39 PASS,跳过 fixing）
- **KR-FB-005 → applied**（git-workflow §5.2 = 修复建议的规则化落地,本机分发副本 L89-105 验证;SKILL.md Step 6 是否补段未验证——无 standard 仓库访问权限）
- 五处联动: problem-registry / findings-registry / audit 报告勘误段 / FB 源文件 / fb-index

### 关键发现（事实级）

1. **standard#11 当前 gh 账号（Corray）不可见** — chatlabs-ai/agent-dev-standard 私有,需切 bcorraychen;本次改用本机分发 rule 副本做一手验证绕开。[已验证: gh GraphQL 报错]
2. **deferred 条目缺"触发条件达成"巡检机制** — P-003 触发条件（引入第一个逻辑分支函数）在 V1.6（约 2 周前）就满足,无人发现直到本次全量盘点。deferred 状态自带触发条件但没有定期比对的钩子,候选改进项（暂不立 FB,1 例）。

### 未做遗留

- （可选）V3 Phase 6 path outline polish — pathEdgeCol 占位未用
- P-004 是规则级 → 后续可考虑回传 standard（与 KR-FB 批次同通道,攒下一批一起报）
- V3 是否到 Phase 6 收尾 / 有无 Phase 7+ — 无规划文档,待用户口头定

### 5 — V3 Phase 6b path outline polish（b108a26, 补做）

- drawPathShape 闭包 + 两 pass 描边（edge 全画完再叠 fill, 防 seam 横纹）, 消 pathEdgeCol 占位
- 视觉验证: WASM + playwright 截图, Lv 1 直线 + Lv 4 corner 均连续无伪影; localStorage 注入存档解锁 Lv 4 测 corner
- 本条目并入 2026-06-03 会话「未做遗留」第 1 项 → 已完成

## 2026-06-04 — V3 收尾 + V4 规划

- **V3 收尾:** 打 annotated tag v1.0/v2.0/v3.0（三 era 边界落 git 原生载体）+ docs/roadmap.md 版本史回填 + V3 未竟项归档（多字号 / Tank / 新塔型,后两项 Kenney sprite 限制）
- **V4 拍板:** 用户从 4 候选（音频+GameFeel / Gameplay 深度 / 内容扩展 / 混合精选）选定**音频 + Game Feel**;Gameplay 深度列候选 V5 主题
- **V4 Phase 1-5 规划落 roadmap.md**（161becf）: SFX 管线（SoundEvents 纯数据队列解耦 term build）→ BGM+音量 → 走动/死亡动画（程序化优先,无现成帧）→ 飘字 → (可选) shake
- 关键风险标注: WASM autoplay 政策（Phase 1 首个实测项）/ BGM 素材来源 TBD（license 核实是 Phase 2 前置）/ save 双实现家族

## 2026-06-04 — V4 Phase 1: SFX 管线（30ff577）

- **架构落地**: sound.go 纯数据队列（shared）+ audio_player.go 播放层（!term）+ game.go 七触发点 + term drain。roadmap 既定方案无偏移
- **素材**: Kenney CC0 三包（interface/impact/digital-audio）官网直链下载, 9 ogg 88KB embed, LICENSE 映射表落 assets/sfx/
- **偏离 roadmap 一处**: 不另设 hit 音（伤害即时结算, shoot+hit 同帧叠播浑浊）, 留 Phase 5 弹着时序再议 — 已记 commit body
- **WASM autoplay 实测**（roadmap 标的首个风险项）: 手势后 AudioContext probe = "running", 解锁机制 ebiten 内部处理, 无需游戏侧代码 — 风险关闭
- **验证**: 49/49 tests（+10 触发测试）/ 三 build / WASM 实玩全音效路径 0 错误; 真实听感待用户本机
- 测试技巧沉淀: Go fmt.Println 在 WASM 落浏览器 console → decode 失败可被 playwright 机器观测

### V4 Phase 2 — BGM + 音量控制（189b4bd）

- BGM 双轨: Juhani Junkala "4 Chiptunes (Adventure)" CC0（OGA 页面 + 包内 INFO.txt 双重核实）, Stage Select→menu / Stage 1→battle, chiptune 与现有像素风统一 — roadmap "BGM 来源 TBD" 风险关闭
- bgm.go 流式解码 + InfiniteLoop（区别于 SFX 全量 PCM 驻留）+ 淡入淡出状态机（fade→0 切轨 →1, ~0.5s）
- 音量档: save 双实现同构扩展 Volume *int（指针区分旧存档未设置 vs 显式静音 0）— roadmap 预警的 fix-pattern-scan 家族, 两份同步改 + 注释互指
- -/= 键全 phase 调节, HUD 右上常显, 立即 StoreSave
- 验证: 55/55 tests / 三 build（wasm 16.4MB, +2.6MB embed）/ WASM 实测旧存档兼容（无 volume 字段→默认 7）+ 调节持久化（localStorage volume:5, completed 保留）+ 切轨往返 0 错误
- 遗留: bgmBaseVol 0.4 音量平衡待用户听感; 临时盘 ENOSPC 教训 — 大文件下载注意 curl timeout 加大 + 及时清 /tmp

### V4 Phase 3 — enemy 走动/死亡动画（5fc195f）

- 程序化三件套（零新素材, roadmap 既定方案）: pathLerp 平滑插值（替换逐格跳动, 仅渲染层）+ 四向旋转（基准朝向"朝右"经 tilesheet PIL 裁剪实证, 不是猜的）+ 垂直行进轴 bob 摆动
- 死亡动画走 effects 管线: EDeath + 插值坐标字段, 死在被击中处无跳变
- 65/65 tests（+7）; 1 处测试预期修正 — Update 先 move 后 shoot, 死亡位置=移动后位置, 实现对测试错
- WASM 截图裁剪放大验证: 水平段朝右+亚格点+bob 高度差 / Lv4 竖直段横转 90° ✓; 死亡动画 0.35s 窗口 still 没逮住 — 逻辑单测覆盖+渲染同 EHit 模式, 动态待用户实玩
- 方法沉淀: sprite 朝向这类"素材事实"用 PIL 裁 tilesheet 直接看, 别按记忆猜; playwright 截图坐标 = 视口 1200 宽 ≠ 游戏 840 宽, 裁剪要乘缩放

### V4 Phase 4 — 伤害飘字 + 金币反馈（2f1e06a）

- EText effect: 命中飘伤害数字（浅黄 0.5s）+ 击杀飘 "+Ng"（金色 0.8s, 右偏 0.3 cell 错开）, 复用 decay 管线 + 插值坐标; drawTextCol 带色字体变体
- HUD 金币闪烁: 增加才闪（花钱不闪）, phase 切换帧防假闪
- 68/68 tests（+3, 另更新 1 处旧断言 2→3 effects — 设计变更）
- **验证缺口（如实）**: 飘字像素级截图未捕获 — playwright 会话 tab 状态错乱（fresh navigate 后仍渲染旧帧）, 4 次截图失败按迭代上限停手。逻辑单测 + 渲染同 EDeath/EHit 模式 + 首轮对局 0 console error; 待用户实玩验收
- 教训沉淀: playwright MCP 长会话（多次 navigate/close 循环）后页面状态不可信, 渲染验证应在会话早期一次性完成, 或换 fresh browser context; 截图驱动的"战斗瞬间"类验证成本高, 优先选静态可见的验证对象（如 Phase 2 的 Vol 显示 / Phase 3 的持续行走状态）

### V4 Phase 5 — screen shake + boss 顿帧 + J 开关（5936799）

- shake: 离屏缓冲整帧偏移, shakeOffset 纯函数（4px 线性衰减 0.3s, 确定性无随机）; 非 shake 帧零开销
- 顿帧: BossKills 计数器 + 渲染层观察增量（沿用 goldFlash 模式, 不加事件队列）, hitStop 0.12s 跳 game.Update
- J 开关: JuiceOff bool omitempty（零值=默认开, 旧存档自然兼容 — 与 Volume 指针方案对比是适配字段语义选的）, save 双实现 family 同步
- 72/72 tests（+4）; WASM fresh context 只验静态对象（Phase 4 教训生效）: J 键 E2E localStorage 持久化 + omitempty 往返 ✓
- V4 Phase 1-5 全部完成 — 待用户整体验收后 V4 收尾（tag v4.0 + roadmap 更新）

## 2026-06-04 — V4 收尾（d874b53 + tag v4.0）

- tag v4.0 @ 5936799 已推送（v1.0~v4.0 四 era 边界齐）
- roadmap: 版本史补 V4 行 / V4 section 转已收尾（规划原文存档 + 收尾记录前置）/ 变更记录
- 未竟项归档 ×3: 独立 hit 音（弹着时序触发再议）/ 音色与 bgmBaseVol 调参入口 / 飘字像素级截图缺口
- 盲点声明入档: shake/顿帧/飘字密度未经长时实玩压测; 飘字无独立开关
- V5 候选: Gameplay 深度（卖塔/targeting/AoE/减速塔/主动技能）— 待用户拍板
- 跨项目遗留提醒: P-004（规则级, task-lifecycle checkbox 纪律）仍待攒批回传 standard（需 bcorraychen 账号）

## 2026-06-04 — V5 启动（26810c2）

- 主题: Gameplay 深度（V4 候选标记 → 用户拍板启动）, Phase 1-5 规划落 roadmap
- Phase: 卖塔 → targeting 策略 → Cannon AoE → 状态效果/减速 → 陨石雨主动技能
- V5 特有风险定纪律: 每 phase 动核心战斗循环 → 重构先行（pickTarget/killEnemy 纯函数化, 行为不变回归）+ 测试先行 + 平衡数值全常量
- 关键预判: Phase 3 击杀路径家族（Spawner 召唤逻辑必须唯一份）/ Phase 4 载体三选项待 sprite 调研 / Phase 5 不用 Esc（已绑退出）

### V5 Phase 1 — 卖塔（0473c40）

- towerInvested 逐级累计 + sellRefund 单一出处（两处算式防分叉）+ SellTower（退款/移除/SndSell/金币飘字/Msg）
- **测试先暴露实现 bug**: 90×0.7=62.999… int() 截断 → math.Round 修（测试对实现错, 与 V4P3 的"实现对测试错"成对照）
- X 键/右键仅 ebiten 层 + 仅 PhasePlaying; sell.ogg ← Kenney back_002（CC0 重新核实）
- 77/77 tests; AC 三条全测试覆盖（截图不在 AC 内）
- **工具债确认**: playwright tab 无法真实 reload（带 query 串导航仍复用旧状态）— V4P4 同根因二次复发, 达 formalization 类型 A 阈值（≥2 例）→ 候选改进: 游戏驱动类验证改用每次 spawn 新 browser context 或弃浏览器改 ebiten 离屏渲染测试; 下次 phase 验证前先决策

### V5 Phase 2 — targeting 策略（3cdacff）

- 重构先行落地: pickTarget 抽纯函数, First 零值 == 原"最前优先", 2 个回归测试先锁行为再扩策略 — roadmap 定的纪律首次完整执行
- 三策略: First/Last/Strong（平手取最前, 确定性 tie-break）; T 键循环 + 塔上提示 + Msg
- 84/84 tests（+7, 含集成: Strong 塔实打高 HP 放过低 HP）
- 浏览器验证跳过（同构接线 + AC 无截图项 + 工具债未决策）

### V5 Phase 3 — Cannon AoE 溅射（2db6b96）

- 前置重构落地: damageEnemy/killEnemy 唯一击杀路径, 注释立家族约束（新伤害来源禁止内联复制击杀逻辑）; 既有 84 测试零改动全过 = 重构回归证据
- 溅射: Cannon Splash 1.0/1.2/1.5, 50% 折扣, 飞行过滤一致; range 快照防召唤物连锁被溅
- 90/90 tests; 家族验收点全中: 溅射杀 Spawner 召唤 ✓ / 溅射杀 Boss 计数 ✓
- 又一次"先 move 后 shoot"测试预期错（V4P3 同款, 第 2 次）— 模式: 写时序相关断言前先查 Update 内操作顺序

### V5 Phase 4 — 状态效果系统 + Frost 塔（3c8ddb7）

- 载体决策: PIL 裁 tilesheet 调研 → tile 227/251 未占用, "塔 sprite 耗尽"风险（V3.6 实证）解除 → 选 a) 新塔型 TFrost
- 状态框架: ApplySlow 语义三规则（不叠加取最强 / 命中刷新 / 过期后弱效可生效）+ EffectiveSpeed 统一入口
- spec 驱动架构红利兑现: TowerKinds 扩 4 后, UI 按钮 / 卖塔 / targeting / term 渲染零改动自动适配
- 97/97 tests（+7）; 三 build ✓

### V5 Phase 5 — 陨石雨主动技能（8ee77f5）

- CastMeteor: 半径 2 全敌 60 伤（含飞行）, 25s 冷却, 经 damageEnemy 统一路径（家族约束第三次兑现: 陨石杀 Spawner 召唤 ✓）
- R 瞄准态（ebiten UI 层）+ 橙圈 + 点击释放 + shake 复用 + HUD 冷却条（绿满/橙进度/AIMING）
- 102/102 tests（+5）; 修 1 处 fixture（1-wave 测试关通关致冷却停摆 — 每帧重置敌人防清场）
- **工具债 workaround 确认**: 换端口 = 新 origin = 强制新文档, 破解 playwright tab 复用 — 瞄准 UI/冷却条/Frost 按钮一帧全验; 后续游戏驱动验证用"每次换端口"方案, 工具债关闭

## 2026-06-04 — V5 收尾（3f46b48 + tag v5.0）

- tag v5.0 @ 8ee77f5 已推送（v1.0~v5.0 五 era 齐）; roadmap 版本史 + 收尾记录 + 变更记录
- 纪律兑现归档: 重构先行 ×2 / killEnemy 家族约束三次兑现 / 测试拦截 IEEE754 bug / playwright 工具债关闭（换端口方案）
- 未竟项 ×4 归档: 数值平衡待实玩 / term 不支持 V5 操作（含 4 塔显示 vs '4' 键未接的不一致）/ Frost 未 tint / 陨石瞬间未截图
- V6 候选: 内容扩展（V4 方向矩阵第三项, gameplay 深化后反对理由减弱）— 待拍板
- 单日跨度回顾: 今日完成 V4 收尾 → V5 启动 → V5 五个 phase → V5 收尾, 测试 72→102, tag ×2

## 2026-06-04 — V6 启动（规划落 roadmap）

- 主题: 内容扩展（V4 矩阵第三项 → V5 收尾候选 → 拍板）, Phase 1-4: 关卡 11-20+菜单两列 → 难度 → endless → 星级
- 规划前提实测: windowH 620 单列 20 关溢出 → 两列必做; 零 rand 使用 → endless seed 注入约束; save family 第 3/4 轮预告
- 内容向纪律: levels.yaml 数据完整性测试兜底 — AI 内容无法实玩调参, 数值保守+常量化+测试, 手感交用户反馈

### V6 Phase 1 — 关卡 11-20 + 菜单两列（ca43c0c）

- 10 新关每关一个 V5 系统主题（飞行波/Fast 海/Spawner 海/长途/制空/短直道/Boss 复合/极限 Fast/攻坚/终局）, lives 递减 gold 补偿
- 菜单两列 + menuRowAtPixel(mx,my) 双列 hitbox; **测试暴露并修复 V3P5a 起的既有 bug**（Go 负数整除向零截断 → 点 title 区误启动 Lv1）— 测试先行第二次拦截真 bug
- 数据完整性测试兜底落地（20 关/链连续/起终点/边界/敌量单调性 — Spawner 加权 ×3）
- 105/105 tests; WASM 两列菜单截图 ✓（换端口方案稳定复用第 2 次）

### V6 Phase 2 — 难度模式（提交见 git log）

- difficulty.go 三档系数表, DiffNormal=0 零值兼容（第 3 次复用该模式）; Spec() 非法值回退防炸局
- newEnemy 敌人生成唯一入口（wave spawn + 召唤共用）— killEnemy 后第二个统一施加点家族约束
- save family 第 3 轮（Difficulty 字段）; 111/111 tests; WASM D 键 E2E（Diff: Hard 显示 + localStorage difficulty:1）

### V6 Phase 3 — Endless mode（提交见 git log）

- endless.go 预算制生成器: 预算恰好花完（总威胁度==公式 → 单调性由构造保证）+ 敌型池逐波解锁 + 超 10 波 HP 缩放
- 确定性铁律立规: 首次引入随机性, rng 注入不用全局 — 同 seed 可测（20 seeds 门禁测试）
- beginRun 抽取（StartLevel/StartEndless 共用, 零改动回归）+ spawnEnemy 统一封装（第 3 个统一施加点: 难度+endless 缩放, 召唤物不漏）
- save family 第 4 轮（BestWave, 清场即记防 M 退出丢纪录）
- 119/119 tests; WASM E 键全要素截图 ✓

### V6 Phase 4 — 星级评分（提交见 git log）

- starsFor 纯函数（3★满命/2★≥70%/1★通关, 边界测试 ×8）+ RecordStars 取 max（save family 第 5 轮, nil map 安全 = 旧存档零迁移）
- ASCII '*' 决策: gomono U+2605 字形覆盖不确定, 不赌 tofu
- 124/124 tests; WASM 注入存档验证菜单 ***/*/空 三态显示 ✓
- V6 Phase 1-4 全部完成 — 待用户验收后收尾（tag v6.0）

## 2026-06-04 — V6 收尾（d9b56d7 + tag v6.0）

- tag v6.0 @ d2d4bf1 已推送（v1.0~v6.0 六 era 边界齐）; roadmap 版本史 + 收尾记录 + 变更记录
- 纪律兑现归档 ×5: 内容数据测试兜底 / seed 确定性铁律 / save family 三轮零失同步 / 统一施加点 ×2 / 测试拦截既有 bug 第 3 例
- 未竟项 ×4: 内容平衡待实玩 / term 20 关溢出未验 / term 未接 V6 操作（定位事实降级 V1.7 兼容保留）/ ASCII 星标
- 下一版三候选未拍板: a) 平衡打磨 b) 发布版（部署+README）c) 新机制探索
- 单日总账: V4 收尾起 → V5 全程 → V6 全程, 测试 72→124, tag ×3（v4/v5/v6）, 游戏终态 20 关+3 难度+endless+星级

## 2026-06-04 — V7 启动（发布版规划）

- 拍板 ×3: 发布版方向 + 改名 Gopher Defense（商标风险, "inspired by" 标注）+ 转 public/Pages
- Phase 1-3: 改名+README+截图 → index.html+Actions+Pages → Release+metadata
- 关键边界: 存档 key/module 名不改（丢进度风险）, 技术标识符与对外名解耦; 改名验收 = grep 零残留（白名单制）

## 2026-06-04 — V7 Phase 1: 改名 + README + LICENSE + 截图（提交见 git log）

- Gopher Defense 改名 6 处用户可见字符串, 不改名单执行（存档 key/module 名 — 进度零迁移）; grep 白名单制验收
- README 英文门面（特性/键位/构建/credits 全链 license）+ MIT LICENSE + 截图 ×3（letterbox 公式裁黑边）
- 截图体力活教训再证: MCP 调用间隔大, "等交火"类拍摄改用 endless 重生场景 + 单次长 evaluate 内完成建塔等待

## 2026-06-04 — V7 Phase 2: Pages 上线（4a50aad + be86c3c）

- index.html 重写（全键位/内联 favicon 修历史 404/页脚 credits+标注/实测体积文案 5.9MB gzip）
- Actions workflow: 测试门禁（xvfb — ebiten 包 init 需 DISPLAY, 首跑失败后修, ebiten 官方 CI 同款）+ WASM -s -w（16.4→15.4MB）+ deploy-pages
- 仓库转 public（用户已拍板）+ Pages workflow 模式开启
- **线上实测 ✓**: https://corray.github.io/kingdom-rush/ — playwright 打开线上 URL, canvas 加载/loading 消失/菜单完整渲染, 游戏正式公开可玩
- push master 即自动发布, 测试不过不部署

## 2026-06-05 — V7 Phase 3: Release 收尾（tag v7.0 + Release 页）

- License 终核: 12 sfx + 2 bgm + sprites 全映射 ✓ 根 MIT ✓
- repo metadata: description/homepage(Pages URL)/topics ×6
- tag v7.0 @ be86c3c + GitHub Release "Gopher Defense v7.0 — first public release"（notes 含玩链/特性/构建/版本史/license 链）
- 全量回归 124/124 + 三 build ✓
- 剩余: roadmap V7 收尾归档（版本史行 + 收尾记录）

## 2026-06-05 — V7 收尾（d4e0ad8）— 七 era 全闭环

- roadmap 版本史 + 收尾记录归档; tag v7.0 已在 Phase 3 打（规划即如此）
- 未竟项 ×4 归档, 其中 Actions Node 20 弃用警告有 2026-06-16 限期 — 候选小修挂账
- 下一版三候选未拍板（平衡打磨/宣传分发/新机制）, 项目可自然停
- 项目总账: v1.0→v7.0 七 tag / 33→124 tests / terminal→桌面→WASM→公开上线 / 全素材 CC0 溯源 / save 6 字段全程向后兼容
- 跨项目遗留仍挂: P-004 规则级回传 standard（需 bcorraychen）

## 2026-06-05 — Actions Node 20 弃用警告修复（e5b66bf）

- 五 action 升当前大版本（gh api 实查版本号, 非凭记忆）: checkout v6 / setup-go v6 / configure-pages v6 / upload-pages-artifact v5 / deploy-pages v5
- run 26990581814 全绿 + annotation 零输出（警告消除）+ 线上 200
- V7 未竟项中唯一限期项（2026-06-16）销账, roadmap 同步

## 2026-06-05 — P-004 回传 standard（#32）

- 跨项目遗留清账: P-004 (规则级) → chatlabs-ai/agent-dev-standard#32, 账号 bcorraychen 提完即切回 Corray
- 发现升级: 与 artifact-based-handoff §中断恢复 已载失真模式正反互证（"勾了没写" vs "写了没勾"）→ issue 里点出"checkbox 自觉同步双向都不可靠", 整改建议三层（纪律/审查/结构, 结构层治本）
- 前批 #11 已 CLOSED（维护方处理过）→ 新发现开新 issue 的判断正确
- 至此 kingdom-rush 无任何挂账: registry 全 resolved/dismissed/已回传, V1-V7 闭环, 线上运行

## 2026-06-05 — 草地条纹修复（c5ff055, 用户反馈"界面太丑"）

- 诊断: spriteGrass=24 是泥土+草边过渡 tile（PIL 平铺测试实证）→ 全画面橙底绿竖条纹, V3 Phase 2 选型错误潜伏 4 个版本; tile 25 = 无缝纯草, 一常量修复
- README battle/endless 截图重拍（新图含交火/飘字/金币全要素）
- **方法论突破**: browser_run_code_unsafe 单段 playwright 代码全流程（真 CDP 鼠标+零调用间隔）, 根治截图抢跑 — 此前 V4P4/V5P5/本次共 3 轮踩"MCP 间隔跑不过游戏时钟", 自此关闭
- 已部署线上; 待用户确认: 背景修复后是否满意, 或继续 M 级美化（菜单/HUD）

## 2026-06-05 — V7.2 UI 美化批次 M1+M2+M3（82129fa）

- M1 菜单: 20pt 标题 face（消 V3 P5c"多字号"未竟项）+ 关卡行分段上色（等宽 7px/char 偏移技法）+ 锁定行变灰 + 总星数统计
- M2 HUD: 塔按钮 cost 实时红绿（实用信息可视化）+ Msg 黄/提示青/帮助灰三层
- M3 地图: 装饰散布（tile 134/131/132/137, PIL 调研）, genDecor 确定性 seed + 7% 密度 + 避路; 126/126 tests（+2）
- README 截图 ×3 重拍（run_code_unsafe 单段流程二次复用, 稳定）
- 已部署线上; "界面太丑"两批修复完成（条纹根因 + 三层美化）

## 2026-06-05 — V7.2 收尾归档（3abd61d + tag v7.2）

- 用户确认界面更新生效 → tag v7.2 @ 82129fa; roadmap 加轻量批次小节（不走完整 era 模板 — 公开发布后的反馈批次形态首例）
- V3 未竟项"多字号字体"销账（潜伏 4 版本, V7.2 M1 消）
- 批次形态确认: V7 后进入"反馈驱动小批次"节奏（区别于 era 制）, V7.2 是首例 — 触发(用户反馈) → 诊断(根因) → 修复+美化 → 截图重拍 → 部署 → tag 归档

## 2026-06-05 — V7.3 速效 QoL 包（a5a8582）

- C1 暂停（P, 同顿帧语义+BGM 不断, 实测 3s 冻结）+ C2 倍速（F, dt×2 公平加速）+ D1 Frost 蓝 tint（ColorScale 乘数）+ D2 飘字并 J（渲染层过滤零逻辑变化）+ E1 term '4'（phase 分支保 menu 语义）+ F2 badges×3
- 销账: V4 盲点"飘字无开关" + V5 未竟项 ×2（Frost tint / term 4 键）
- 验收教训: 小图压缩偏色误判 cost 红绿 — 放大取色核实为全绿, 无 bug; 截图验收必要时裁原图放大
- 已部署线上; 待用户确认后 tag v7.3 归档

## 2026-06-05 — V7.4 工程包三项完成（a0b2954 / e8260ab / ee75ea0）

- A1 save 合并: save_core.go 唯一定义, 5 轮 family edit 根因消除, 126 测试零改动回归
- B1 BGM fetch: wasm 15.4→12.9MB; 踩坑 ×2 入档 — net/http 在 wasm +8MB（改 syscall/js 直调）; "call to released function" 基线对比定案为 oto 既有无害行为（线上 v7.3 同报）非本改动
- C5 平衡仿真器: 129 tests; 首跑产出 3 个平衡发现 — Normal 太松（全 3★）/ endless 曲线缓（wave 75）/ Hard 倒挂（前期失守后期碾压, gold 补偿过度）— 数据入档调整留用户
- 插曲: cwd 漂移混入 web/bgm_assets_wasm.go 旧版草稿入库, 发现即清（c0aafc3）— cat > 写文件前确认 cwd

## 2026-06-05 — V7.5 平衡调整（两轮数据驱动迭代）

- 第一轮误判即撤: gold 降档对 Normal 无效（钱花不完）反砸死 Lv11 — **改判: Normal 全 3★ 是合理休闲基准不调**（首跑结论修正, 仿真器防错价值二次兑现）
- 第二轮固化: endless HPScale 10%/wave（主杠杆是 HP 不是预算, 75→42 波）+ 预算 n²/6 保留 + Hard lives -1（救前期, 17/20）
- 方法论确认: 改→仿真→修订循环 0.2s/轮, 平衡迭代成本趋零; 测试预期同步 ×3（设计变更类）
- 边际收益判断: Hard Lv8/11 留作难点关, 不再内推 — 下一轮调整应由真实手感数据驱动

## 2026-06-05 — V7.5 收尾（tag v7.5）

- roadmap 归档含"规划假设被数据证伪"显式留档（gold 降档 → 撤销 + 改判）— 诚实记录比漂亮记录重要
- tag 全景: v1.0~v7.0 七 era + v7.2~v7.5 四个反馈/工程批次
- 项目状态: 无挂账, 129 tests, 线上 12.9MB wasm, 平衡有仿真回归门禁

## 2026-06-05 — 会话收尾（今天到这里）

**单日总账（2026-06-04~05 连续两日）：**
- 06-04: V4 收尾 → V5 启动+5 phase → V5 收尾 → V6 启动+4 phase → V6 收尾 → V7 启动+P1
- 06-05: V7 P2/P3+收尾（公开上线）→ Node24 → P-004 回传(#32) → 界面两批 → V7.3 QoL → V7.4 工程包 → V7.5 平衡 → 四批次全收尾
- 终态: tags v1.0~v7.5 / 129 tests / 线上 12.9MB / 零挂账 / 工作区干净

**下次会话入口（按优先级）：**
1. 用户实玩新平衡曲线（endless 42 波 / Hard lives -1）的手感反馈 → 下一轮平衡批次
2. standard#32 维护方响应跟进（bcorraychen 账号）
3. 新方向候选: itch.io 分发 / 新机制探索（多入口 path / 英雄单位 / 塔技能树）

## 2026-06-08 — V8 英雄单位 era（新机制，5 phase 全交付，待手感确认 tag）

- 用户从 gameplay 深度三候选拍板「英雄单位」(对等评估: 多入口 path 收益高但要重做 20 关数据风险最大 / 塔技能树最稳收益中 / 英雄补"主动操作"维度且加法式不重写单路径) → LMP L2
- 决策 A 控制 = 复用光标 + H 键设 rally (键盘网格唯一不新增移动输入、两端一致的方案); 决策 B 阻挡纳入并隔离 P3 (最高风险, 改敌移动状态机)
- **P1**(cdcff18) hero.go 实体+rally 移动纯逻辑+测试 / **P2**(442ec2a) damageEnemy 攻击+敌 meleeCD 反击+复活 / **P3**(3a7ad60) 贴身阻挡 / **P4**(6baa258) 两端渲染+H键+playwright 手感验 / **P5** 仿真接入+平衡+收尾
- 测试 129→148; 每 phase 三 build + 全量绿
- **fix-pattern-scan 实战 ×2**: 英雄生成 path 中点 → 干扰塔战斗精确 HP 断言 (P2 splashFixture/targeting 3 处 + P3 LosePushes 1 处), 一级修复 + 二级扫描确认其余鲁棒 (敌在英雄射程外/飞行/cell 5)
- **平衡仿真**: 英雄被动计入 (beginRun 生成), HeroNet 守护测试证 Hard 有英雄 19/无英雄 17 = 纯增量 +2 且无英雄基线 = V7.5 零回归; Normal 20/20 不变(满命天花板); endless 42 不变(HP 缩放主导)
- **手感验证局限**: playwright 本地 wasm 实测渲染/战斗/阻挡/阵亡复活全中, 但 headless 方向键投递不稳 (字母数字键可靠) → rally 引导线没拍到干净图; "好不好玩"仿真测不出 → 留用户实玩
- **已部署 (06-09)**: 用户选「直接部署线上玩着确认」→ push master (Actions 27178141887 绿, test gate 148 过 + WASM 13.3MB), 线上冒烟确认英雄上线 https://corray.github.io/kingdom-rush/
- **已收尾 (06-09)**: 用户实玩确认手感满意 → tag v8.0 @ 601f7c1; roadmap 版本史表补 V8 行 + 收尾记录, README 版本史 + 英雄 feature/H 键/截图
- 英雄 per-run 不入存档; 未竟: 解锁/等级/XP、专用 sprite、多英雄/技能

## 2026-06-09 — V8 收尾 + 八 era 总账

- tag 全景: v1.0~v8.0 八 era + v7.2~v7.5 四批次; 33→148 tests; terminal→桌面→WASM→公开上线→英雄机制
- V8 方法论复盘: LMP L2 全程兑现 (方向拍板 → 5 phase 规划 → 逐 phase 三 build+测试 → 收尾); 风险隔离 (阻挡独立 P3) + fix-pattern-scan 二级扫描 (英雄干扰塔战斗断言, 4 处修+其余鲁棒确认) + 仿真守护 (HeroNet 净非负) 三纪律兑现
- 项目状态: 零挂账, 148 tests, 线上 13.3MB wasm, 英雄机制上线
- **下次会话入口**: 1) 用户继续实玩 V8 手感 → 可能 V8.1 数值微调 2) 新方向候选 (英雄解锁/等级/技能树 = 英雄深化 / itch.io 分发 / 多入口 path) 3) standard#32 维护方响应跟进

## 2026-06-09 — V9 英雄成长 era (4 phase 全交付, 已部署待手感确认)

- 用户从英雄深化四子项 (解锁/等级XP / 专用sprite / 多英雄 / 技能树) 拍板「英雄成长」(对等评估: sprite 受素材阻塞 Kenney 无 hero / 多英雄+技能树留 V10) → LMP L2
- 决策 A 等级 per-run (不入存档, 持久化会让满级碾压固定难度早关) / B XP 来源 / C 主动技能 = AoE 横扫 G 键照搬 Meteor
- **P1**(02b6fce) Hero Level/XP + GainXP 升级回血/封顶 + Damage()/AttackRange() 按级缩放 / **P2**(32804b1) cleave (lvl3 解锁/8s冷却/半径2/Damage×3, 经 damageEnemy) / **P3**(f0d42f6) 两端 HUD 等级/XP条/技能冷却 + G 键 / **P4**(a3c8be3) 仿真+平衡
- **决策 B 两轮证伪改判** (同 V7.5 gold 反转性质): 仿真证「仅自身击杀」XP 被塔抢尾刀饿死 15/20 关停 L1 → 改「附近助攻」仍 12/20 停 L1 → 终定「被动」(在场每个击杀给 XP, killEnemy 唯一点统一给) → Lv1→L3 cleave 即解锁, 关内 L1→L5 成长弧
- **平衡零破坏**: 英雄升满级+仿真用 cleave 上界, Normal 20/20 / Hard 19/20 (同 V8 难点) / endless 44; per-run 后段加载自限 (难点波在英雄弱时发生)
- **test-infra**(579f7a6): 本会话环境无窗口服务器 → ebiten GLFW init panic, 把 2 个 ebiten-only 测试剥到 //go:build !term, 逻辑测试可 `go test -tags term` headless; 测试执行需 dangerouslyDisableSandbox 或 -tags term
- 测试 148→159; cwd 漂移坑复现 (cd web 起服务后残留) → go/git 用 -C 绑定根目录
- **已部署 (06-09)**: 用户选「直接部署线上确认」→ push (Actions 27190988283 绿), 线上冒烟确认 Hero L1 + G:Cleave@L3 上线
- **待办**: 用户线上实玩成长/cleave 手感 → 满意 tag v9.0; 需调 V9.1。cleave AoE 视觉未拍到干净图 (无塔 smoke 升不到 L3, 逻辑全单测)

## 2026-06-10 — V9 收尾 + 九 era 总账

- 用户线上实玩确认成长/cleave 手感满意 → tag v9.0 @ a3c8be3; roadmap 版本史补 V9 行 + 收尾记录定稿 (决策 B 两轮证伪改判 + per-run 自限假设被仿真证实, 双数据驱动留档), README 版本史补行
- tag 全景: v1.0~v9.0 九 era + v7.2~v7.5 四批次; 33→159 tests; terminal→桌面→WASM→公开上线→英雄→英雄成长
- V9 方法论复盘: 仿真器第三次反转规划预设 (V7.5 gold / V8 无 / V9 XP 来源 ×2 轮) — "改→仿真→修订"循环已是平衡类决策的默认路径; test-infra headless 解耦是环境约束逼出的正资产
- 项目状态: 零挂账, 159 tests, 线上 13.32MB wasm, 工作区干净
- **下次会话入口**: 1) V10 方向拍板 — 英雄深化剩余 (专用 sprite[素材阻塞] / 多英雄 / 技能树 / 跨局持久化[需配按关缩放]) 或 itch.io 分发 / 多入口 path 2) standard#32 维护方响应跟进 (bcorraychen)

## 2026-06-11 — V10 多英雄选择 era (4 phase 全交付, 已部署待手感确认)

- 用户拍板英雄深化 → 范围对比矩阵后定「多英雄」(技能树+持久化捆绑留 V11 meta 成长 — 需按关缩放配套, 平衡重做级; 真 sprite 素材阻塞继续留档); 阵容三人 Knight+Archer+Rogue (用户拍板)
- 决策 B 阻挡 per-class gating (Archer 不阻挡 = 风筝输出换隘口控制) / C 选择入口沿 difficulty 先例 (menu H 键 + Save.HeroChoice 零值兼容) / D 技能统一自身 AoE 只调参
- **P1**(94968f2) heroSpec 单例 → HeroClass 参数化, 测试纯机械符号替换零断言变化 / **P2**(3ae2a18) Archer/Rogue + gating + 8 测 / **P3**(45c276b) 选择UI+存档+程序化配色两端 (金/绿/紫, K/A/R 字形) / **P4**(f9dcd9c) 仿真矩阵 + 校准
- **盲点 3 实证反向**: 担心的是"某职业净负", 实际首跑 Archer/Rogue 净 +3 强于 Knight +2, Hard 20/20 抹平 V7.5 难点关 Lv11 → 一轮调参 (Archer dmg 12→9 / Rogue 9→7) 三职业 Hard 全 19/20 净 +2, 职业差异落在风格 (射程/速度/阻挡/复活/技能形状) 不落数值强度
- playwright 冒烟: menu H 循环 + 职业色 + Archer 绿圆 A 字本体 + 决策 B 行为差异实证 (Archer 无塔不阻挡 → 敌全漏 GAME OVER, vs V9 Knight 单英雄独守 wave1 — 职业选择有真实策略后果)
- 测试 159→173 (headless 171); 三 build + vet 全程绿; 每 phase 一 commit
- **已部署 (06-11)**: push master 触发 Actions (test gate + WASM + deploy-pages)
- **待办**: 用户线上实玩三职业手感 → 满意 tag v10.0; 需调 V10.1 (重点观察: Archer 不阻挡存在感 / Rogue 与 Knight 手感重叠度)

## 2026-06-11 — V10 收尾 + 十 era 总账

- 用户线上实玩确认三职业手感满意 → tag v10.0 @ f9dcd9c; roadmap 版本史补 V10 行 + 收尾记录 (盲点 3 反向实证 + 一轮校准收敛留档), README 版本史补行
- tag 全景: v1.0~v10.0 十 era + v7.2~v7.5 四批次; 33→173 tests
- V10 方法论复盘: 重构先行 (P1 零断言变化机器证明) + 仿真矩阵当天发现当天校准 (改→仿真→调一轮收敛) + playwright 行为差异实证 (决策 B 非换皮)
- 项目状态: 零挂账, 173 tests, 线上 13.33MB wasm, 三职业上线
- **V11 已拍板: Meta 成长** (跨局持久 + 技能树) — 启动方案设计中

## 2026-06-11 — V11 Meta 成长 era (4 phase 全交付, 已部署待手感确认)

- 用户拍板 V11 = Meta 成长 + 两哲学拍点: 满树允许碾压 Hard (earned power) / per-class 4 节点树
- **核心改判**: 「perk 预算」替代「按关缩放」— 持久层存技能树点不存等级, per-run 等级 (V9 资产) 原样保留, 平衡冲击从"重做 20 关"缩到"校准一张 perk 表" — 这个改判让 meta 成长从最重候选变成一天交付
- **P1**(7db4abe) skilltree.go 持久层: 星预算 (赚60/三树总价90=1.5×取舍) + BuyTreeNode gating + Save.TreeNodes 零值兼容 / **P2** HeroBonus 七字段聚合 + Hero 派生数值方法化 + beginRun 快照 + 零bonus=V10基线守护测试 / **P3** PhaseSkillTree 新 Phase + 菜单 T 入口 + 两端列表 UI / **P4** autoPlayTree 上下界矩阵
- **盲点 3「太弱白做」实证**: 首跑 Knight 满树 Hard +0 (30星无感, Archer/Rogue 已 +1) → 校准 Sharpened Blade +4→+6 / Undying -4s→-6s → 三职业满树统一 20/20 (+1) 征服 Lv11, 无树基线 19/20 零回归
- playwright 冒烟: T 进树/三列职业色/选中金边/0星买拒绝 "Need 3*"/localStorage HeroChoice 跨会话恢复顺带验证
- 测试 173→188 (headless 186); 三 build + vet 全程绿
- **已部署 (06-11)**: push master 触发 Actions
- **待办**: 用户线上实玩 (赚星→买 perk→实战验证增益) → 满意 tag v11.0; 重点观察: perk 手感 (买了有感吗) / 树 UI 可读性 / 老玩家爆买体验

## 2026-06-16 — V11 收尾 + 十一 era 总账

- 用户线上实玩确认技能树手感满意 → tag v11.0 @ cda346d; roadmap 版本史补 V11 行 + 收尾记录 (perk 预算改判 + 盲点 3 校准留档), README 版本史补行
- tag 全景: v1.0~v11.0 十一 era + v7.2~v7.5 四批次; 33→188 tests
- **英雄深化四子项收官**: V9 成长 / V10 多英雄 / V11 技能树+持久化 全交付; 仅"专用 sprite"因 Kenney 无 hero 素材悬留 (V8 起留档至今)
- V11 方法论复盘: 「perk 预算」改判是胜负手 (持久层存树点不存等级 → 平衡冲击从重做20关缩到一张perk表 → 最重候选变单日交付); 仿真上下界守护双向 (无树零回归下界 + 满树 earned power 上界); 盲点 3「太弱白做」一轮校准收敛 (Knight 满树+0→+1)
- 项目状态: 零挂账, 188 tests, 线上 13.38MB wasm, 工作区干净
- **下次会话入口**: 1) V12 方向拍板 — 英雄深化已收官, 候选回到 itch.io 分发 / 多入口 path (gameplay 广度) / 音画打磨批次 (找 CC0 hero sprite + 截图债) 2) standard#32 维护方响应跟进 (bcorraychen)

## 2026-06-22 — V12 收尾 + V13 新敌人类型 (两 era 同日交付)

- **V12 P3+P4 合并交付** (`4494364`): 20 关各添加 cps2 (第二入口路径, 末尾汇合同终点 cell); `TestLevels_DualPathIntegrity` 校验双路连通+汇流+地图范围; `rankedTowerSpots` 修复仅覆盖 `Paths[0]` → 覆盖所有路径; Normal 20/20, Hard 19/20 (难点移 Lv9), 无英雄 15/20 (双路削弱英雄阻挡 = 预期)
- **V12 盲点兑现**: 盲点 1 (平衡重校) = auto-player bug 根因, 零轮手调; 盲点 2 (汇流纪律) = 测试兜底 20/20 过; 盲点 3 (英雄只守一路) = 设计意图实现; 盲点 5 (漏改) = grep 清零 + 零回归
- tag `v12.0` @ `4494364`, 195 tests
- **V13 一次性交付** (`21cc588`): 敌型 5→8 — EShield (`d`, 护甲 4 减免 min 1, 克制 Archer 速射) / ERegen (`r`, 3hp/s 回血, 需集火) / EHealer (`h`, 2s 周期治疗半径 2 内最低血量盟友 +5hp); damageEnemy 加 armor 减免 / Update 加 regen+healer AI / wave DSL 扩展 d/r/h / L8 起渐进混入 20 关 / endless 按波次解锁 (4 Shield / 6 Regen / 7 Healer)
- 仿真 202 tests 全过 (Normal 20/20, Hard 19/20, 三职业净非负 — 新敌人替换旧敌人不增总量, auto-player 自然应对, 零手调)
- tag `v13.0` @ `21cc588`, 202 tests
- Roadmap V12 收尾记录 + V13 完整段 + 版本史表补两行 + 变更记录; spool 补本段; tag v12.0/v13.0 打标
- tag 全景: v1.0~v13.0 十三 era + v7.2~v7.5 四批次; 33→202 tests; gameplay 从单路5敌 → 双路8敌
- 项目状态: 零挂账, 202 tests, 已 push 待 CI 验证, 工作区干净
- **下次会话入口**: 1) CI 绿确认 + 线上实玩双路+新敌人手感 → 满意则归档 2) V14 候选: 真 sprite 替换 (素材阻塞长悬) / Boss 机制深化 / 截图重拍

## 2026-06-29 ~ 07-06 — V14/V15/V16 三批次 + audit + V16 收尾

- **V14 视觉+体验** (06-29): sprite 全映射 (3 新敌人 tile265/291/292 + 3 英雄职业 tile266/269/253 + tint, drawTileRotTint 新增) / wave 预览 (HUD 下波敌型图标+数量) / 地图 4 主题 (Forest L1-7 原色 / Desert L8-13 黄棕 / Snow L14-17 蓝白 / Lava L18-20 暗红, terrain tint + 路径配色) / Boss 行为三分区 (L8-13 护盾半血 3s 无敌 / L14-17 冲锋 8s 周期 ×3 速 2s / L18-20 召唤 10s 周期 2 ENormal, 渲染白圈/橙圈)
- **audit code-quality V12-V14** (06-29): 1 MEDIUM (endless Boss 无行为 — endlessLv.ID=0 走 default) + 3 LOW; 合规面全绿 (难度缩放统一路径 / term 渲染自动生效 / 零 TODO)
- **audit 修复** (06-30, `3c71545`): bossTier() 纯函数统一关卡+endless 分区 (endless wave 8+/12+/20+ 对应护盾/冲锋/召唤) + Boss spawn 3s 初始 bossCD (不入场即冲)
- **V15 Meta+UX** (07-01): 成就 16 项 ×5 维度 (Save 7 新字段零值兼容 / killEnemy 累积 / 胜利检查 / recordBestWave 触发) + 菜单 A 查看屏 + 结算屏 (半透明面板: 星级/击杀/金币/命数/英雄等级, KillsThisRun/GoldEarned 新计数) + 截图重拍 ×4 (菜单/双路战斗/成就/技能树, 替换 V7 旧图)
- **V16 塔系扩张** (07-01): Tesla (键5, 链电弹射 2-3 目标 50% 伤害, 打飞行) + Sniper (键6, 射程 6-8 / 伤 40-120 / 慢速, 不打飞行反 Boss) + 4 塔升级分支 (V 键 L1 切换 L2 锁定: Marksman/Rapidfire / Mortar/Gatling / Archmage/Enchanter / DeepFreeze/Hailstorm) — 10 条 build 路线
- **V16 收尾** (07-06, `6e5cada`): simStrategy 重构蓝图参数化 (simBlueprint, legacy [A,A,A,F] 逐位零回归) + TestBalance_NewTowerMatrix 5 蓝图矩阵。**首跑拦截 sniper-mix 16/20 坍塌** (单发溢出 + cost 120 铺塔慢 = 陷阱选项) → cost 120/100/150→100/90/140 + cd 3.0/2.5/2.0→2.6/2.2/1.8 一轮收敛 18/20。**仿真器第 4 次拦截拍脑袋数值** (V7.5 gold / V9 XP / V11 Knight 树 / V16 Sniper)
- tag: v14.0 @ 733493c / v15.0 @ 49bac79 / v16.0 @ 6e5cada; 测试 202→218; roadmap 补 V14-V16 归档 (版本史 +3 行 + 合并批次段 + 变更记录 +6 行)
- 项目状态: 218 tests, 十六 era, 6 塔 ×分支 / 8 敌 / 3 英雄 / 技能树 / 成就 / 双路 / 4 主题
- **下次会话入口**: 1) 线上实玩 V14-V16 手感 (新塔/分支/成就/结算屏) 2) 候选: itch.io 分发 (悬留最久) / 关卡 21-30 第二章 / 移动端触屏

## 2026-07-06 — itch.io 脚手架 + V17 第二章 (关卡 21-30)

- **itch.io 脚手架** (`e16ec17`): 用户拍板做分发 → docs/itch-page.md (配置/embed 1000×640/正文/标签) + make itch-zip (zip 根含 index.html, 6.3MB) + web/index.html 键位表 V7→V16 全量更新 + butler v15.27.0 装 ~/bin。**上传中断**: 用户选"装 butler 我来推"后改主意 → 等 `! ~/bin/butler login` + username, 随时可续
- **V17 拍板**: 决策 A 星星经济 = 树加第 5 节点 (30 关 90 星 = V11 三树总价 90, 取舍失效 → capstone ×3 价 10, 总价 120 = 1.33× 恢复) / 决策 B 难度 = 接续爬升
- **P1 树 5 节点**: Warlord (+20HP&+4dmg) / Windrunner (+0.5rng&+2dmg) / Phantom (+0.7spd&ability-2s); 前 4 节点不动零回归; V11 经济测试更新为 V17 语义 (90 星 → 两满树 80 剩 10, 第三树只能浅尝); 树屏 panelH 硬编码 4 → treeNodesPerClass 动态 (计划外发现: 第 5 节点与帮助行重叠)
- **P2 关卡 21-30** (Twilight 章): Swarm Nexus (Tesla 秀场 n50 海) / Iron Column (Sniper 秀场 d16+b4) / Unequal Roads (不对称双路: 超长蛇形 vs 直道) / Field Hospital (r16 h8 恢复大队) / Velvet Rush (f46) / Broodmoon (s14 h7 增殖) / Night Raid (g38) / Vanguard (b5 车轮) / Ragnarok (终局 6 波); 终局关 lives 2 (卫哨 ≥3 放宽 ≥2)
- **P3 菜单三列**: menuCols 动态 (ceil(n/10)), hitbox 与渲染同式; stats 压缩 `w%d g%d` (首拍右缘裁剪实锤盲点 2); 成就 +2 (Second Dawn 30 关 / Star Master 90 星, 既有条件不动防打破已解锁)
- **P4 仿真扩 30**: **盲点 1 反向 — Normal 首跑 30/30 零校准** (V16 蓝图矩阵基线立功); Hard 28/30 (Lv9 老难点 + Lv23 新装甲难点); 蓝图克制实证 (sniper-mix 过 L23 装甲关 / tesla 蓝图 L29 Boss 车轮失守 — 正反都对); **Rogue 满树 +0 判结构性** (失守关在英雄弱期/蓝图盲区, perk 救不动, V9 per-run 自限同因) → Phantom 加强一档收手 (V7.5 边际收益同判)
- playwright 冒烟: 三列 30 关 + L21 hitbox/unlock 链 + 树 5 节点两轮 (首轮抓 2 个 UI 溢出 → 修后复拍干净)
- tag v17.0 @ efb24ce; 218 tests; README 版本史补 v12-v17 + 30 关/5 节点/18 成就文案; .playwright-mcp/ 入 gitignore
- 项目状态: 十七 era, 30 关 × 5 主题 / 6 塔 × 分支 / 8 敌 / 3 英雄 × 5 节点树 / 18 成就
- **下次会话入口**: 1) 线上实玩第二章 (L21-30 手感 + capstone) 2) itch.io 续传 (butler login + username 即推) 3) 候选: 移动端触屏 / 在线排行榜 / 关卡编辑器
