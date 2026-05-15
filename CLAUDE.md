# kingdom-rush — Project Context

> **本文件 = 项目级 agent 上下文。** 会话启动时与 `~/.claude/CLAUDE.md`（全局）同时加载，给 agent 提供项目特定信息。
>
> **维护方：** PM（项目元数据 + 当前阶段 + 业务定位）/ EL（项目特定 rules，如有）
>
> **跨机器约定（重要）：**
> 本文件 commit 到 docs 仓库（schema 共享），但 §Standard 路径 段是**用户本地化字段**——clone 项目后必须按本机 standard 仓库绝对路径调整。建议本地化变更**不入 commit**（`git update-index --skip-worktree CLAUDE.md` 或保持 working tree 修改不 stage）。

---

## §Standard 路径（用户本地化字段，PM init 时填）

```yaml
standard_path: /Users/chat/backend-ai-workflow/agent-dev-standard
# 例: /Users/gsr/claude/agent-dev-standard
```

**作用：**
- 项目级 standard 引用源
- 全局 `~/.claude/CLAUDE.md` 已 `@import` 9 条核心 rules（个人 default 体验）
- 本字段允许项目级**扩展引用** standard 文档——如 `init-flow.md` / `project-env-spec.md` / `concepts/*` 等非全局加载的内容
- 不同用户机器 standard 路径不同（团队默认版 / 个人 fork / 实验分支），所以本字段必须本地化

**用法：** 引用 standard 文档时使用绝对路径 `<standard_path>/docs/concepts/<doc>.md`（agent 读到本字段后展开）。

---

## 项目元数据

| 字段 | 值 |
|------|----|
| **名称** | `kingdom-rush` |
| **流派** | `github` |
| **代码仓库（本地目录名）** | `kingdom-rush` |
| **共享文档仓库（本地目录名）** | `kingdom-rush`（单仓库，code 与 docs 同目录）|
| **env.yaml 路径** | `kingdom-rush/env.yaml` |

---

## Issue 配置

> `/issue` skill 执行前读取本段（字段来源 standard `skills/core/issue/SKILL.md` L21-32）。

- **issue_repo**: `Corray/kingdom-rush`
- **doc_repo**: `/Users/chat/Desktop/games/kingdom-rush` (单仓库流派，code 与 docs 同目录)
- **adr_path**: `docs/adr/` (待首次 ADR 时创建)
- **code_path**: `/Users/chat/Desktop/games/kingdom-rush`
- **compile_cmd**: `go build ./...`
- **role**: `be`

---

## Standard 引用（基础规则）

本项目采用 **agent-dev-standard** 工作规范（ADR-004 决策 1：standard 是规范源，不 install 到本项目）。

**全局已加载（无需重复引用）：** `~/.claude/CLAUDE.md` 已 `@import` 9+ 条核心 rules（spec-to-code-flow / problem-handling-pattern / artifact-based-handoff / task-lifecycle / fix-pattern-scan / research-first / architecture-constraints / security-review / tech-debt / 等）。

**项目级扩展（按需）：**
- handoff 协议: 详见 `<standard_path>/docs/concepts/project-init-flow.md` + standard rule `artifact-based-handoff.md`
- env.yaml 规范: 详见 `<standard_path>/docs/concepts/project-env-spec.md`
- 流派差异: 同上 §一 末"流派项目背景差异"对比表

---

## Flow 工作模式

### 通用（不分流派）

- **agent 启动加载顺序：** `~/.claude/CLAUDE.md`（全局）→ 本文件（项目级）→ `docs/env.yaml`（运行时按需读）
- **commit + push 硬门禁：** handoff 收尾未 push = 未完成（Obs-7 教训）
- **角色独立性：** PM / BE 角色边界清晰，参数（ticket id 等）通过 TAPD / GitHub Issue 原生载体流转，不通过 handoff 文件硬塞

### TAPD 流派专属（流派 = tapd 时启用）

- **BE 任务收尾 4+1 件套：** commit / push / comment / status 流转 / **工时**（缺一不算闭环）
- **完整 ID 格式：** TAPD API 接受 `<workspace_id>001<short_id>`，UI 短 ID 不能直接调 API
- **群消息：** curl webhook 直发（MCP `send_qiwei_message` 当前有缺口）

### GitHub 流派专属（流派 = github 时启用）

- **BE 任务收尾 3 件套：** commit / push / Issue comment
- **工时管理：** 弱约束，daily worklog 周度 / 月度汇总即可

---

## 项目特定 context（PM 维护）

> PM 按项目实际情况填以下段，可标 TBD 后续补。

### 业务定位
- Kingdom Rush 游戏相关项目（克隆 / MOD / 工具方向，具体方向待 PM 细化）

### 关键决策点
- TBD

### 当前 active 阶段
- 参见 `docs/env.yaml`（TAPD 流派 `active_stories` / GitHub 流派 `active_milestones`）

### 项目特殊约束（如有）
- TBD

---

## 项目特定 rules（非必填）

> 本段仅当项目有 standard 未覆盖的特殊约束时填，**不重复定义 standard rules 已有内容**。

- 无（简单项目通常无项目特定 rules）

---

## 变更记录

| 日期 | 变更 |
|------|------|
| 2026-05-15 | 初建（agent-dev-standard install + PM init 元数据填充）|
