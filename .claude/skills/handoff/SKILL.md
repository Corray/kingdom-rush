---
name: handoff
description: Check pending handoff tasks from PM/SA, execute by priority
disable-model-invocation: true
installed-from: agent-dev-standard@c8fe883
installed-on: 2026-05-15
---

# /handoff — 跨角色交接任务消费

从 `<handoff_dir>/pending/` 目录消费上游（PM / SA）写入的待办任务。这是上游 → 执行层的主通道，消除"用户逐条转达"的动作。

> **背景：** CLAUDE.md 描述式约定（"会话启动必扫描"）对 AI 无效。本 skill 是显式调用层兜底——用户主动 `/handoff` 触发扫描。

---

## 输入

- 无参数 — 扫描 `pending/`，按优先级 + 日期展示待办，等用户选择
- `<filename>` — 直接处理指定 handoff 文件
- `--list` — 只列出 pending，不进入执行
- `--status` — 统计 pending / in-progress / completed 数量

---

## 项目配置

| 项 | 值 |
|---|---|
| `handoff_dir` | `<project>/docs/handoff/`（通常无需调整）|
| `README` | `<handoff_dir>/README.md`（协议详细规范）|

> install 时确认路径。通常就是 `docs/handoff/`。

---

## 执行流程

### Step 1 — 前置检查

```bash
test -d <handoff_dir>/pending || {
  echo "❌ handoff 目录未初始化。请运行 /install handoff 或手动创建。"
  exit
}
```

### Step 2 — 扫描并排序

1. 读取 `<handoff_dir>/pending/*.md`
2. 每个文件解析 YAML header 的 `priority` + `date` + 一级标题
3. 排序：priority（HIGH > MEDIUM > LOW）→ date（早 → 晚）

### Step 3 — 展示清单

```
📋 待处理 handoff：<总数>

1. [HIGH]   YYYY-MM-DD  <topic-a>  (文件名.md)
2. [MEDIUM] YYYY-MM-DD  <topic-b>
3. [LOW]    YYYY-MM-DD  <topic-c>

请选择要处理的编号，或输入 N 跳过。
```

如果 pending 为空：

```
📋 pending/ 为空，无待办任务。
```

**不自行处理**，等用户选。

### Step 4 — 开始某任务

用户选 N 后：

1. 读取对应文件完整内容
2. **移入 in-progress**：`mv pending/<file>.md in-progress/<file>.md`
3. 更新 header：
   - `status: in-progress`
   - 追加 `started: YYYY-MM-DD HH:MM`
4. 按文件内容执行 —— handoff 任务本身的内容决定后续流程，可能触发 `/issue`、`/audit`、`/fix` 等其他 skill

### Step 5 — 完成后归档

任务执行完成，按 handoff 文件里定义的"完成记录"要求：

1. 更新 header：
   - `status: completed`
   - 追加 `completed: YYYY-MM-DD HH:MM`
2. 文件末尾追加"完成记录"（commit hash / Issue 编号 / 实际产出）
3. **移入 completed/**：`mv in-progress/<file>.md completed/YYYY-MM/<file>.md`
   - 目录按完成月份（非启动月份）分
4. 同步更新 problem-registry（如 handoff 涉及 P-xxx 条目）

---

## 错误处理

| 异常 | 处理 |
|------|------|
| `pending/` 不存在 | 提示运行 /install 或创建目录 |
| 文件 header 格式错误 | 跳过并警告，不阻塞其他文件 |
| 用户选的文件在处理前已被别的会话移走 | 重新扫描并提示 |
| `in-progress/` 里已有文件（上次中断）| 在清单里以 `[IN-PROGRESS]` 标注，允许继续 |

---

## 约束

- **扫描是只读，移动和修改必须用户触发** — 不自动开始执行任务
- **一次只处理一个** — 不要批量扫描所有 pending 一起跑
- **header 格式错误不阻塞** — 跳过但要警告，防止一个格式问题挡住其他待办
- **跨会话状态一致** — 会话中断后重新 /handoff，要能在 in-progress/ 发现"上次未完"

---

## 与其他机制的关系

| 机制 | 职责 |
|------|------|
| 本 skill | 执行层消费 handoff 的入口 |
| 项目 CLAUDE.md | 提供"会话启动建议 /handoff 检查"的提示（但不依赖 AI 自扫）|
| `docs/handoff/README.md` | handoff 协议的详细规范 |
| `/wrap-up`（上游 SA 侧）| 上游产出 handoff 后会话收尾（对称的另一端）|
| `templates/handoff.md.template` | handoff 文件模板 |

---

## 与协议层的关系

参考 `protocols/role-taxonomy.md` —— `/handoff` 是 SA / handoff agent → 执行层（FE / BE / QA）的协议节点。

---

## 试点观察

新装项目首次部署本 skill 时建议：
- 先跑 3~5 次再考虑 `--recover` 参数（恢复上次中断的 in-progress 任务）
- 是否要支持"多个并行在 in-progress"（当前设计是串行）
- 是否要和 `/install` 集成（install 时自动建目录）
