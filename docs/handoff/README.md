---
date: YYYY-MM-DD
from: <发起角色，如 SA / TL / PM>
to: <接收角色，如 执行层 / EL>
priority: HIGH | MEDIUM | LOW
kind: <handoff 类型，如 audit / fix-dispatch / cleanup / rule-update / bug-fix>
status: pending
related:
  - <相关 audit 报告 / Issue / commit / 上游 handoff 路径>
---

# Handoff — <一句话主题，含日期>

## 背景

<2-3 段：触发原因 / 上游决策 / 当前状态>

## Why（为什么需要做）

<列出 1-3 条核心理由：实证 / 风险 / 上游依赖 / 其他>

## 建议动作（接收方决策范围）

**核心动作：**
- <具体可执行 step 1>
- <具体可执行 step 2>
- ...

**预期产出：**
- <commit / 文档更新 / Issue / 其他 artifact>

## 不预设（接收方自决）

接收方有权基于上下文决定以下，发起方不强加：
1. <决策点 1：A vs B 选择>
2. <决策点 2：是否启动派生 handoff>
3. <决策点 3：实施粒度 / 顺序>
4. <决策点 4：是否反哺 observations>

## 验证场景（Check 必过后才能进收尾）

> 涉及代码 / 接口 / 数据模型变更的 handoff 必含此段。纯文档 / 元层 handoff 可省。

| # | 操作 | 预期结果 |
|---|----|----|
| C-1 | <验证步骤 1> | <预期 1> |
| C-2 | <验证步骤 2> | <预期 2> |

**验证通过条件：** 所有场景全部步骤通过 → 进入收尾。
任何步骤失败 → 修复后重跑该步骤，不跳过。

## 收尾

完成后：
1. mv 本 handoff 从 `pending/` 到 `completed/YYYY-MM/`
2. 在接收方 spool 留底（一句话总结）
3. 如有派生发现，开新 handoff 或起 Issue
4. status 变更：pending → completed

---

## 历史变更（可选）

| 日期 | 状态 | 变更说明 |
|------|----|----|
| YYYY-MM-DD | pending | 起草 |
| YYYY-MM-DD | in-progress | 接收方启动 |
| YYYY-MM-DD | completed | 收尾 |
