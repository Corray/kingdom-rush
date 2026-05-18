# Audit Report — issue-process (kingdom-rush)

**日期:** 2026-05-15
**phase:** issue-process
**窗口:** 30 天（2026-04-15 ~ 2026-05-15）
**范围:** 全局
**审查 Issues:** #1
**合规率:** 22/27 = 81%（首次审查,无趋势对比）

---

## 输入

- **Spec 侧:** `protocols/issue-process.md` § 审查维度（5 维度 27 项）
- **实现侧:**
  - `gh issue list --repo Corray/kingdom-rush --state all`: 1 issue (#1)
  - `gh api repos/Corray/kingdom-rush/issues/1/events`: 7 events
  - `docs/release-history.md`: 1 mock build (#1, head=`0859892`)
  - `git log`: 3 commits (`8c413bf` / `c87ea6b` / `0859892`)

## Issue #1 — events 时序

| 时间 (UTC) | 事件 | 来源 |
|---|---|---|
| 09:59:54 | issue created | gh issue create + labels: pm-reviewed, feature |
| 10:03:14 | +be-in-progress | /issue Step 1 |
| 10:04:36 | ❗ **closed** | commit `0859892` push 触发 `closes #1` |
| 10:10:18 | reopened | /release P-A 准备时人工 reopen |
| 10:10:51 | +be-confirmed | /release P-A Step 6.3 |
| 10:10:52 | -be-in-progress | /release P-A Step 6.3 |

---

## Step 3 正向比对（spec → 实际,27 项逐查）

### 维度 1 — Label 状态机合规性（5 项）

| # | spec 项 | 实际 | 状态 |
|---|---|---|---|
| 1 | 状态流转合法 | 流转含 closed → reopened 异常路径 | 🔴 **违规** → IPR-001 |
| 2 | 跳过中间态 | 实际经过 be-in-progress,无跳过 | ✅ |
| 3 | 未授权回退 | reopened 是 closed→open(合法),无 confirmed→in-progress 之类回退 | ✅ |
| 4 | closed 在 [role]-confirmed 之后 | closed (10:04:36) **早于** be-confirmed (10:10:51) | 🔴 **违规** → IPR-001 |
| 5a | be-in-progress 在方案 approve 后打 | pm-reviewed 后 3 min 才打,中间有自检表 + 用户确认 | ✅ |
| 5b | be-confirmed 由 /release 切 | P-A Step 6.3 切换 | ✅ |
| 5c | 已发版的 in-progress 及时切 confirmed | 已切（merge-base exit=0 + label 切换 1 sec 内） | ✅ |

### 维度 2 — 收尾 comment 完整性（5 项）

| # | spec 项 | 实际 | 状态 |
|---|---|---|---|
| 6 | commit hash 齐 | Step 6 comment 含 `commit: 0859892` × 2 处（自检 grep 命中） | ✅ |
| 7 | 文档链接指向共享文档仓库 | 单仓库流派,doc_repo == code_path | N/A |
| 8 | 状态注明合规 | "工作完成 — 等 /release 发版" + "已发布到 dev — 可关闭" | ✅ |
| 9 | /release 切 label 后含「可关闭」 | P-A Step 6.3 comment 含"**可关闭**" | ✅ |
| 10 | 文档同步每份有 commit hash | 三选一勾"无需更新"(豁免依据 #4458847614) | ✅（豁免）|

### 维度 3 — /issue Step 触发完整性（7 项）

| # | spec 项 | 实际 | 状态 |
|---|---|---|---|
| 11 | Step 0 自检表贴出 | comment #4458832324 | ✅ |
| 12 | 架构师 3 维填写 | 3 维全填（架构无影响 / 技术债无关 / 演化无影响） | ✅ |
| 13 | LMP 方案确认 | S 级不触发 LMP（自检判定） | ✅（豁免）|
| 14 | 文档先行 commit 在代码 commit 之前 | **无文档先行 commit** | 🟡 **偏差** → IPR-002 |
| 15 | 编译门禁 retry 超上限 | 编译一次过,无 retry | ✅ |
| 16 | 测试 — 收尾提及测试通过 / 新增测试用例 | 仅 `go run .` manual 验证,**无 unit test** | 🟡 **偏差** → IPR-003 |
| 17 | 6a 文档同步 / 6b Issue comment / be-in-progress label | 与维度 1/2 交叉,合规 | ✅ |

### 维度 4 — 跨规则交叉验证（5 项）

| # | spec 项 | 实际 | 状态 |
|---|---|---|---|
| 18 | be-confirmed 对应 dev build | mock build #1 + Step 6.3 comment 明示 build #1 | ✅ |
| 19 | 修复 commit 在 build 号之前 | commit `0859892` 时间 < mock build #1 时间 | ✅ |
| 20 | fix-pattern-scan 适用类型的家族扫描 | feature 类型不触发 fix-pattern-scan | N/A |
| 21 | ADR 关联（如有） | 本次无 ADR | N/A |
| 22 | tech-debt 新 TODO 同步 problem-registry | 改动无新 TODO 引入 | ✅ |

### 维度 5 — 趋势 / 模式（5 项）

| # | spec 项 | 实际 | 状态 |
|---|---|---|---|
| 23 | 违规率 | 5/27 项违规或偏差 = 81% 合规率 | 📊 基线 |
| 24 | 高频违规 step | 维度 1（label 状态机 closed/confirmed 时序） | 📊 基线 |
| 25 | 时间趋势 | prerequisite 累计 ≥ 3 次审查,本次首次 | N/A |
| 26 | 类型分布差异 | 仅 1 feature,无对比数据 | N/A |
| 27 | 责任人/来源差异 | actor=Corray + Claude（演练） | N/A |

---

## Step 4 反向比对（实际 → spec 应该如何）

实际中没有 spec 未覆盖的额外现象。所有事件 + comment + label 切换都可追溯到 5 维度某项检查。

**反向无新发现。**

---

## Step 4.5 安全语义升级

**安全基线扫描（3 项固定检查）:**

| # | 检查 | 结果 |
|---|---|---|
| S1 裸接口扫描 | grep Controller 方法权限注解 | N/A — 项目无 Controller / 无 HTTP 接口 |
| S2 凭证泄露扫描 | grep key/secret/token/password/credential 硬编码 | ✅ 无命中 |
| S3 OAuth 完整性 | state/nonce/token 失效机制 | N/A — 项目无 OAuth |

**Step 3/4 发现的安全语义升级:** Issue #1 涉及字面量改动,无认证/授权/加密/凭证语义,不升级。

---

## Step 4.6 Family-scan 合规

**触发判据:** Issue #1 类型 = feature（字面量改动）,不命中 fix-pattern-scan 触发场景（状态转换 / 状态校验 / 参数校验 / 异常处理 / 同名方法 / 第三方 API 参数 / 身份隔离）。

**结论:** N/A — 无 family-scan 触发场景,跳过两层扫描。

---

## Findings 汇总

| 编号 | 严重度 | 标题 | 维度命中 | 根因 | 推荐处置 |
|---|---|---|---|---|---|
| **IPR-001** | **HIGH** | closed 事件早于 be-confirmed | 维度 1 #1 + #4 | KR-FB-005（`closes #N` 触发 auto-close） | 升级 KR-FB-005 severity → high；standard 仓库 SKILL.md issue Step 6 加 commit message 关键词禁用列表 |
| **IPR-002** | **MEDIUM** | 无"文档先行" commit 在代码 commit 之前 | 维度 3 #14 | KR-FB-003（S 级 + 项目无 spec 豁免规则未明示） | 已有 KR-FB-003 覆盖；建议 audit skill 增加"识别豁免 comment"逻辑（avoid false positive） |
| **IPR-003** | **LOW** | 无 unit test,仅 manual `go run .` 验证 | 维度 3 #16 | 项目 bootstrap 阶段无 test 文件 + 字面量改动测试性低 | 项目 active phase 时补 `main_test.go` |

---

## Finding 详情

### IPR-001 — closed 事件早于 be-confirmed [HIGH]

**违规事实:**
- closed at 10:04:36 (UTC)
- be-confirmed at 10:10:51 (UTC)
- be-confirmed 晚于 closed 约 6 min

**spec 期望:** protocols/issue-process.md 维度 1 #4 — "closed 是否在 [role]-confirmed 之后"

**根因:** commit `0859892` message 含 `closes #1` → GitHub 在 push 到 default branch 时 auto-close issue。SKILL.md issue L211 明示 "closed 只由人工触发,agent 不主动 close",但 commit message 阶段无禁用列表。

**实证链:**
- commit `0859892` message body 含 `closes #1`
- events: 10:04:36 closed event,actor=Corray
- 后续不得不 reopen（10:10:18）+ /release 切 confirmed（10:10:51）

**关联:** KR-FB-005 / P-001（problem-registry 待写入）

**推荐处置:** 升级 KR-FB-005 severity → high；回传 standard 仓库时建议优先级最高。

---

### IPR-002 — 无"文档先行" commit 在代码 commit 之前 [MEDIUM]

**违规事实:** 三次 commit (`8c413bf` bootstrap, `c87ea6b` role fix, `0859892` v0.0→v0.1) 中没有"文档 commit 早于代码 commit"模式。

**spec 期望:** protocols/issue-process.md 维度 3 #14 — "文档先行: 是否有标'待实现'的文档 commit 在代码 commit 之前（结构化）"

**根因:** KR-FB-003（S 级 + 项目 bootstrap 无 spec/api 文档可改 → 文档先行豁免依据贴在 issue comment #4458847614,但 audit 机械化结构检查不识别该豁免 comment）。

**实证链:** comment #4458847614 含完整豁免依据（S 级 + bootstrap + Issue SSOT + 不触发 ADR）。

**关联:** KR-FB-003 / P-002

**推荐处置:** 已有 KR-FB-003 覆盖；建议 audit skill 自身改进 — 维度 3 #14 检查时识别 "Step 2 豁免" comment（grep 特定关键词如 "文档先行豁免" / "S 级豁免"）→ 不计违规。

---

### IPR-003 — 无 unit test [LOW]

**违规事实:** 项目无 `*_test.go` 文件,Issue #1 收尾仅提及 `go run .` manual 验证,无 unit test。

**spec 期望:** protocols/issue-process.md 维度 3 #16 — "测试: 收尾是否提及测试通过 / 新增测试用例（语义层）"

**根因:** 项目 bootstrap 阶段无 test infrastructure,字面量改动测试性低。

**关联:** P-003

**推荐处置:** 项目 active phase 时补 `main_test.go`；本次演练 LOW 不阻塞。

---

## 系统性建议（FB 候选）

1. **audit skill 自身改进 — 豁免 comment 识别（FB 候选）:** 维度 3 #14 / #16 等"结构化检查"应支持识别 "Step 2 文档先行豁免" / "S 级测试豁免" 等显式 comment,避免 false positive。
   - 跨项目普适: ✓（任何 bootstrap 阶段项目都会触发）
   - 根因在规范层: ✓（豁免规则未在 spec 中明示）
   - 改进规则后能降低问题: ✓（audit + /issue 双方都改进）

2. **issue-process protocol 维度 1 #4 应明示 reopen 场景:** 当前 spec 只检查 "closed 在 confirmed 之后",但**没明示 reopen 是否允许** + reopen 后的 confirmed 是否补救。本次 issue #1 是 closed→reopened→confirmed 链,标合规还是违规需要 spec 明确。
   - 跨项目普适: ✓（任何 auto-close 误触都会触发）
   - 根因在规范层: ✓
   - 改进规则后能降低问题: ✓

---

## 24h SLA Action Items（留给 SA）

按 SKILL.md L323 + `rules/extension/audit-fix-dispatch.md`,SA 应在 24h 内对每条 finding 做 fix dispatch 5 选 1 决策:

| Finding | 推荐决策 |
|---|---|
| IPR-001 | **escalate** — 升级 severity (medium → high) + 合入 KR-FB-005,优先回传 standard |
| IPR-002 | **dismiss** — 实质是 KR-FB-003 已覆盖,本条是 audit skill false positive 表现；不建独立 fix,留作 audit 改进 FB |
| IPR-003 | **defer** — 项目 active phase 时再补 unit test |

---

## 勘误段（review adjudication, 2026-05-18）

> audit 报告本身 immutable;本段为 review 拍板后的勘误归档,verdict 决策链见 `docs/handoff/completed/2026-05/2026-05-18-adjudication-issue-process-2026-05-15.md`。

### Verdict 落地

| Finding | 原 status | 新 status | 行动项执行情况 |
|---------|----------|----------|--------------|
| IPR-001 / P-001 | proposed | **confirmed** | KR-FB-005 已是 high（无需升级 severity）;衍生:occurrences 1 → 2,新增 IPR-001 作为第 2 例实证;standard 仓库 SKILL.md Step 6 加 commit message 关键词禁用列表 → 留待 standard 仓库回传时执行（本项目侧无法直接改 standard）|
| IPR-002 / P-002 | proposed | **dismissed** | 衍生 FB-006 候选已落 `docs/feedback/2026-05-18-audit-skill-exemption-detection-fb.md`（audit skill 维度 3 豁免 comment 识别能力）|
| IPR-003 / P-003 | proposed | **deferred** | 触发条件:active phase 引入第一个有逻辑分支的函数时改 confirmed → fixing |

### 报告自身瑕疵（本次不产出独立勘误文件,记此供下次 audit 改进）

1. **L84 合规率算式表述歧义** —— "5/27 项违规或偏差 = 81% 合规率",其中 5 是违规数、81% 是合规率,混读易误解。建议下次 audit 改为 "合规 22 / 违规 5 / 总 27 = 81%"。
2. **L8 头部 vs log L34 合规率分母不一致** —— 报告头写 `22/27=81%`,log 同时给出 `19/27=79%`/`23/27=85%`/`22/27=81%` 三种算法。审查报告对外口径应单一,内部 log 可保留三种以便比对。

### audit 自荐 fix dispatch vs review verdict 的差异

audit 自荐 L206-208 对 IPR-001 推荐 `escalate`（升 severity）,review 改判 `confirmed` —— 理由:KR-FB-005 已是 high（feedback 文件 L240 实证）,severity 维度无需 escalate;finding 自身分类正确即 confirmed,severity 升级是衍生动作不是 finding 状态变更。standard 对 "FB 实例 vs 独立 finding" 的边界未明示,本判定保留 confirmed,可能盲点已在 handoff 标注。
