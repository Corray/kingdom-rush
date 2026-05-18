# Problem Registry

**初始化日期**: 2026-05-15
**数据来源**: findings-registry（audit 产出）+ GitHub Issues + 开发中发现

> 项目级质量模式（tendency）见 `project-patterns.md`（可选独立文件）— 记录此项目特别容易犯哪几类错（PP-NNN）。本项目首次 audit,暂无 PP 条目。

---

## 字段约定

| 字段 | 含义 |
|------|------|
| 编号 | `P-NNN` 项目内连续递增 |
| 来源 | `audit` / `issue` / `开发` / `用户反馈` 等 |
| 日期 | 首次记录日期 |
| 模块 | 业务模块名 |
| 标题 | 一句话描述 |
| 类型 | 缺失 / 偏差 / 风险 / 改进建议 |
| 层级 | 项目级 / 规则级（规则级会被上报全局）|
| 状态 | 见下文状态枚举 |
| 关联 | 关联 finding ID / Issue # / commit hash 等 |

---

## 状态枚举

与 findings-registry 一致: proposed / confirmed / fixing / resolved / dismissed / deferred / merged / escalated。

---

## 条目（按时间倒序,最新在顶部）

| 编号 | 来源 | 日期 | 模块 | 标题 | 类型 | 层级 | 状态 | 关联 |
|------|------|------|------|------|------|------|------|------|
| P-003 | audit | 2026-05-15 | (bootstrap) | 无 unit test,manual `go run .` only | 改进建议 | 项目级 | **deferred** (2026-05-18) | IPR-003 |
| P-002 | audit | 2026-05-15 | issue-process | 无文档先行 commit（实质 KR-FB-003 触发,audit skill 未识别豁免 comment）| 偏差 | 规则级 | **dismissed** (2026-05-18) | IPR-002 / KR-FB-003 / FB-006(衍生) |
| P-001 | audit | 2026-05-15 | issue-process | closed 事件早于 be-confirmed（KR-FB-005 实证）| 偏差 | 规则级 | **confirmed** (2026-05-18) | IPR-001 / KR-FB-005 |

---

## 变更记录

| 日期 | 变更 |
|------|------|
| 2026-05-15 | 初建,首次 audit issue-process 演练产出 P-001/002/003 |
| 2026-05-18 | review adjudication 落地：P-001 → confirmed / P-002 → dismissed / P-003 → deferred（依据 handoff/completed/2026-05/2026-05-18-adjudication-issue-process-2026-05-15.md）|
