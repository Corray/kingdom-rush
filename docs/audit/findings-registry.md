# Findings Registry

所有审查发现的统一索引。每条发现有唯一编号,追踪从发现到处置的全生命周期。

**首次初始化**: 2026-05-15（kingdom-rush dogfood,首次 audit issue-process 演练）

---

## 编号空间约定

| Phase | 编号前缀 | 例 |
|-------|--------|----|
| Spec | `GAP-NNN` | GAP-001 |
| API | `API-NNN` 或 `API-<MODULE>-NNN` | API-001 / API-LLM-001 |
| Behavior | `BHV-NNN` | BHV-001 |
| Architecture | `ARCH-NNN` | ARCH-001 |
| Integration | `F-INT-NNN` | F-INT-001 |
| Data Model | `F-DM-NNN` | F-DM-001 |
| **Issue Process** | **`IPR-NNN`** 或 `IPR-T-NNN` (trend) | IPR-001 / IPR-T-001 |
| Code-Doc Gap | `GAP-CDG-NNN` | GAP-CDG-001 |
| Rule Coverage | `RC-NNN` | RC-001 |

---

## 状态枚举

| 状态 | 含义 |
|------|------|
| `proposed` | 新发现,待 review 确认 |
| `confirmed` | 已确认有效,进入处理 |
| `fixing` | 处理中（有对应 handoff / Issue / commit）|
| `resolved` | 已解决,有证据（commit / artifact）|
| `dismissed` | 排除（误报 / 测试噪声 / 范围外,需 reason）|
| `deferred` | 已确认但推迟（需触发条件）|
| `merged` | 合并到另一条（需 merged-into 引用）|
| `escalated` | 升级（严重度 / 重新分类 / 改编号）|

---

## 维护规则

- 每次 audit 产出新 finding → 在此追加条目（status=proposed）
- 状态变更时同步本文件 + 原 audit 报告末尾追加勘误（audit 报告 immutable）
- 编号一旦分配,不可重用、不可重号

---

## 条目（按 phase 分组）

### Issue Process 审查

| 编号 | 首次发现 | 当前状态 | 说明 | 关联 |
|------|---------|---------|------|------|
| IPR-001 | 2026-05-15 | **resolved** (2026-06-03) | closed 事件早于 be-confirmed,违反 protocol 维度 1 #4（HIGH）；review 认可,进 fixing；衍生动作:KR-FB-005 occurrences +1；2026-06-03 resolved:issue #1 显式闭环（4 字段 close comment）+ 约束沉淀 standard git-workflow §5.2 + 上报 standard#11 | audit/2026-05-15-issue-process.md / KR-FB-005 / P-001 / handoff/completed/2026-05/2026-05-18-adjudication-* |
| IPR-002 | 2026-05-15 | **dismissed** (2026-05-18) | 实质 KR-FB-003 已覆盖 + 豁免依据 comment #4458847614 已贴,本条为 audit skill false positive；衍生 FB-006（audit 豁免 comment 识别能力）| audit/2026-05-15-issue-process.md / KR-FB-003 / P-002 / feedback/2026-05-18-audit-skill-exemption-detection-fb.md |
| IPR-003 | 2026-05-15 | **resolved** (2026-06-03) | bootstrap 阶段 main.go 仅含字面量,无业务逻辑可测；触发条件:active phase 引入第一个有逻辑分支的函数时改 confirmed → fixing；2026-06-03 resolved:前提不成立——V1.6 起测试随版本累积,现 game_test.go 39 tests 全过,跳过 fixing 直接 resolved | audit/2026-05-15-issue-process.md / P-003 |

---

## 变更记录

| 日期 | 变更 |
|------|------|
| 2026-05-15 | 初建,首次 audit issue-process 演练产出 IPR-001/002/003 |
| 2026-05-18 | review adjudication 落地：IPR-001 → confirmed / IPR-002 → dismissed / IPR-003 → deferred（依据 handoff/completed/2026-05/2026-05-18-adjudication-issue-process-2026-05-15.md）|
| 2026-06-03 | 用户授权批处理：IPR-001 → resolved（issue #1 显式闭环 + git-workflow §5.2 沉淀 + standard#11 上报）；IPR-003 → resolved（前提不成立,39 tests 已存在）。与 problem-registry P-001/P-003 联动,原报告末尾追加勘误段 |
