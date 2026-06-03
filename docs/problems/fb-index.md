# FB Index — 启动前扫描索引

**定位：** 团队级 / 个人级 feedback（FB）的结构化元数据索引，支持 `/fb-scan` skill 按 skill / module / phase / category 筛选。

- **数据源：** `feedback/*.md`（按日期或主题分组的源文件）
- **维护规则：** 新 FB 录入时必须同步在此索引追加条目
- **字段定义：** 见本文件末尾 schema

---

## 编号规范

- **既有条目（2026-05 dogfood 批次）：** `KR-FB-001~005`（kingdom-rush dogfood 临时前缀，源文件约定合入 standard 时由维护方重号）+ `FB-006`（IPR-002 adjudication 衍生）。**编号保留不重号** —— 已被 problem-registry / audit findings-registry / standard#11 交叉引用，重号断追溯链。
- **新条目（2026-05-25 起）：** 按 standard ADR-008 用 `FB-YYYYMMDD-<4-char-hash>` 格式，防多会话撞号。

## 状态枚举

| 状态 | 含义 |
|------|------|
| `candidate` | 候选，未达 ≥ 2 例阈值 |
| `observing` | 观察期，已达阈值待累积更多实证 |
| `applied` | 已 applied 到规则文件 / SOP |
| `verified` | applied 后实际生效（产出 ≥ 1 次拦截真实问题）|
| `dismissed` | 排除（噪声 / 重复 / 已被 别 FB 覆盖）|

---

## FB 条目

## KR-FB-001 — `.gitignore` 模板不区分 hub vs business 项目
- **date**: 2026-05-15
- **file**: feedback/2026-05-15-skill-dogfood-fb.md
- **category**: implement
- **skills**: install
- **modules**: install/modules/06-templates.sh + templates/hub-gitignore.template
- **phases**: install Phase 6
- **severity**: medium
- **status**: candidate
- **occurrences**: 1
- **guidance**: 06-templates.sh 渲染 .gitignore 时必须按项目类型（hub vs business）选模板；统一套 hub 模板导致业务项目 `CLAUDE.md` / `.claude/` 被静默忽略
- **scan_when**: 新项目 install 完后第一次 git add / commit / push 前
- **related**: KR-FB-002（同源于 install 模板未规范化）

## KR-FB-002 — `role` 字段大小写未规范化
- **date**: 2026-05-15
- **file**: feedback/2026-05-15-skill-dogfood-fb.md
- **category**: implement
- **skills**: install + issue
- **modules**: templates/CLAUDE.md.template + skills/core/issue/SKILL.md
- **phases**: install + issue execution
- **severity**: medium
- **status**: candidate
- **occurrences**: 1
- **guidance**: labels.yml 全 lowercase（be-reviewed 等），CLAUDE.md `role` 字段必须 lowercase，否则 `[role]-in-progress` 拼接出不存在的 label
- **scan_when**: install 完跑 `/issue` 时第一次 `gh issue edit --add-label [role]-*`
- **related**: KR-FB-001（同源于 install 模板未规范化）

## KR-FB-003 — S 级 + 项目无 spec 的文档先行豁免规则未明示
- **date**: 2026-05-15
- **file**: feedback/2026-05-15-skill-dogfood-fb.md
- **category**: design
- **skills**: issue
- **modules**: skills/core/issue/SKILL.md
- **phases**: issue Step 2（文档先行）
- **severity**: low
- **status**: candidate
- **occurrences**: 1
- **guidance**: S 级改动 + bootstrap 阶段（无 spec / api / 共识文档）时文档先行应可豁免；SKILL.md Step 2 应明示豁免段（类比 Step 0 规模豁免），定义触发条件 + 必填字段
- **scan_when**: 第一次跑 `/issue` 时项目尚无 spec 文档
- **related**: FB-006（豁免 comment 的 audit 识别侧）

## KR-FB-004 — `gh issue list --label X` 在新建 label 后立刻跑会扑空
- **date**: 2026-05-15
- **file**: feedback/2026-05-15-skill-dogfood-fb.md
- **category**: implement
- **skills**: release（可能影响 audit / 其他依赖 gh issue list 的 skill）
- **modules**: skills/core/release/SKILL.md
- **phases**: release Step 6.3
- **severity**: medium
- **status**: candidate
- **occurrences**: 1
- **guidance**: GitHub label-based search 索引有延迟（实测 sync 后几分钟内 `--label` 过滤返回空但 issue view 可见）；label 同步后立刻跑的查询需加 fallback 到客户端 jq filter
- **scan_when**: install 完立刻跑 `/release` 演练，或 labels 同步完几分钟内的任何 label-based 查询
- **related**: —

## KR-FB-005 — commit message `closes #N` 触发 GitHub auto-close，违反 SKILL.md 约束
- **date**: 2026-05-15
- **file**: feedback/2026-05-15-skill-dogfood-fb.md
- **category**: implement
- **skills**: issue（影响 release Step 6.3 算法前提）
- **modules**: skills/core/issue/SKILL.md
- **phases**: issue Step 6 commit
- **severity**: high
- **status**: candidate
- **occurrences**: 2（2026-05-18 更新：IPR-001 作为第 2 例实证，见 audit/findings-registry.md）
- **guidance**: commit message 禁用 auto-close 关键词（closes / fixes / resolves 全变体）指向本次 issue，用 `refs #N` / `related #N` / `see #N` 替代；违反则 GitHub 替 agent 间接 close，破坏 release Step 6.3 `--state open` 前提
- **scan_when**: 跑 `/issue` Step 6 写 commit message 时
- **related**: IPR-001 / P-001（第 2 例实证）

## FB-006 — audit skill 维度 3 / #14/#16 缺豁免 comment 识别能力
- **date**: 2026-05-18
- **file**: feedback/2026-05-18-audit-skill-exemption-detection-fb.md
- **category**: audit
- **skills**: audit（issue-process phase 维度 3 #14 文档先行 / #16 测试覆盖）
- **modules**: standard skills/core/audit/SKILL.md Step 3 + protocols/issue-process.md 维度 3
- **phases**: audit issue-process
- **severity**: medium
- **status**: candidate（首例，未达 ≥ 2 例阈值）
- **occurrences**: 1
- **guidance**: audit 维度 3 结构化检查应先识别 issue 内豁免 comment（grep "文档先行豁免" / "S 级豁免" / "测试豁免" 等关键词），命中标 ✅ 豁免而非违规，避免 false positive
- **scan_when**: 跑 `/audit issue-process` Step 3 维度 3 检查时（特别是 S 级或 bootstrap 阶段项目）
- **related**: KR-FB-003（豁免规则未明示）/ IPR-002 / P-002（首例实证）

---

## 统计

| 维度 | 数量 |
|------|------|
| 总计 | 6 |
| critical | 0 |
| high | 1 |
| medium | 4 |
| low | 1 |
| applied 状态 | 0 |
| observing 状态 | 0 |

> 注：KR-FB-005 occurrences=2 已达 ≥ 2 例阈值，但源文件 status 仍为 candidate——status 流转待拍板（见 problem-registry 拍板项），本索引镜像源文件不单边改。

---

## 变更记录

| 日期 | 变更 |
|------|------|
| 2026-06-03 | 初次补录：KR-FB-001~005 + FB-006（闭环 2026-05-18 spool 遗留项 (b)），删示例骨架，编号规范增既有保留 + ADR-008 新格式双轨说明 |

---

## Schema（字段定义）

| 字段 | 必填 | 类型 | 说明 |
|------|----|----|----|
| date | yes | date YYYY-MM-DD | 首次发现日期 |
| file | yes | path | feedback 详细内容文件 |
| category | yes | enum | audit / process / design / implement / meta |
| skills | yes | list | 关联 skills |
| modules | yes | list | 关联模块（"(all)" 表示通用）|
| phases | yes | list | 关联 phase（"—" 表示无）|
| severity | yes | enum | low / medium / high / critical |
| status | yes | enum | candidate / observing / applied / verified / dismissed |
| occurrences | no | int | 实证累计次数 |
| guidance | yes | string | 一句话指引（scan_when 触发时呈现）|
| scan_when | yes | string | 启动前扫描时机 |
| related | no | list | 相关 FB IDs |
