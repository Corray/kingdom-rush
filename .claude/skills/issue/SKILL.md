---
name: issue
description: GitHub Issue standard handling flow — fetch, analyze, execute with doc-first approach
disable-model-invocation: true
installed-from: agent-dev-standard@c8fe883
installed-on: 2026-05-15
---

# /issue — GitHub Issue 标准处理流程

按 issue 处理协议（参考 `protocols/issue-process.md`）处理指定 issue。

---

## 输入

- `$ARGUMENTS`：issue 编号（如 `123`）或 `list`

---

## 项目配置

执行前读取项目 CLAUDE.md 中的 `## Issue 配置` 段，获取以下信息：

| 项 | 含义 |
|---|---|
| `issue_repo` | Issue 所在的 GitHub 仓库（owner/repo）|
| `doc_repo` | 共享文档仓库本地路径（用于文档同步）|
| `adr_path` | ADR 文件目录 |
| `code_path` | 代码根目录 |
| `compile_cmd` | 项目编译命令（如 `mvn compile -q` / `npm run build` 等）|
| `role` | 当前角色（默认 BE，可为 FE / QA）|

如果配置不存在，提示用户先运行 `/install` 或手动指定。

---

## 执行逻辑

### 当参数为 `list` 时

1. 拉取 issue_repo 的全部 open issues
2. 按优先级排序展示（high > medium > 无；pm-reviewed > raised）
3. 等待用户选择，不执行任何操作

### 当参数为 issue 编号时

#### 阶段一：展示 + 意图回读 + 建议（不做任何修改）

1. 拉取 issue 原文 + 所有 comments
2. 完整展示 issue 内容
3. **S/M/L 分级（在填自检表前先判断，决定后续每步深度）：**

   | 级别 | 判定特征 |
   |------|---------|
   | **S** | 单文件改动 / 无接口变更 / 无状态机影响 / 无歧义，可直接修 |
   | **M** | 单模块改动 / 可能有接口变更 / 状态机简单涉及 / 基本明确 |
   | **L** | 跨模块 / 接口或数据模型变更 / 需求有歧义 / 架构决策或 Gap |

   输出格式：`分级：S / M / L，理由：<一句话>`

4. **问题定性自检表（硬门禁，不得跳过）**

   **执行硬约束：**
   - **必须**先在 Issue 贴 comment 含完整自检表 4 维（下列），用户确认后才能进入意图回读 + 后续步骤
   - **不得**只在对话回复中展示——必须落到 Issue comment 形成永久追溯
   - **不得**隐式判断 / 直接动手——AI 跳过自检表 = 违规
   - **规模豁免**满足全部条件时（跨文件 ≤ 1 / 不涉接口 / 不涉数据模型 / 不涉安全 / 改动类型限定 typo / log 级别 / comment / 私有作用域重命名 / 格式化），架构师 3 维段标"⚪ 规模豁免（XXX 原因）"，**仍必须发 comment**——豁免是答案不是省略

   **Comment template（depth 按 S/M/L 调整内容详尽度）：**

   ````markdown
   ## Step 0 — 问题定性自检表

   **分级：** S / M / L
   **理由：** <一句话>

   ### 需求前提
   - [ ] 共识文档有明确定义 → 引用章节号
   - [ ] 共识文档有提及但模糊 → 引用 + 标注模糊点
   - [ ] 共识文档未提及 → 标注 **"需求空白"**

   ### 影响范围
   - [ ] 只影响当前接口 → 列出接口
   - [ ] 影响同模块其他接口 → 列出关联接口
   - [ ] 跨模块 → 列出受影响模块

   ### 修复前提
   - [ ] 不需要需求决策，代码逻辑明确有错 → 可直接修
   - [ ] 存在 ≥ 2 种合理实现 → 列出选项，等用户拍板
   - [ ] 需求本身未定义 → 标注 **"需 PM 确认"**，先建 Gap Issue

   ### 架构师视角 3 维

   **架构影响：**
   - [ ] 无影响（纯实现细节）
   - [ ] 消除约束（改善）：<一句话>
   - [ ] 新增约束（需评估）：<一句话>

   **技术债维度：**
   - [ ] 消除债：<具体哪条>
   - [ ] 无关
   - [ ] 新增债（需标注）：<什么债 / 为什么接受>

   **长期演化：**
   - [ ] 让未来简单（复用性增）
   - [ ] 无影响
   - [ ] 让未来复杂：<为什么接受 / 是否走 LMP>

   ### LMP 升级判定
   任一维度勾"新增约束 / 新增债 / 让未来复杂" → **强制升级 LMP**（即使代码改动小）
   ````

   **执行步骤：**
   1. AI 用 `gh issue comment <N> -F <tmpfile>` 贴上述 template comment（占位符填好）
   2. 在对话同步展示完整内容给用户
   3. 等用户在 Issue 下追评 / 在对话显式确认 → 进入下一步

5. **意图回读** — 用自己的话复述理解：
   ```
   【我理解的目标】…
   【我理解的约束】…
   【我不确定的地方】…
   ```
6. 判断当前状态，给出执行层处理建议：
   - 是否需要 PM 回复？是否可直接行动？
   - 技术影响范围
   - 推荐方案（如有多个选择逐一列出，触发 Large Module Protocol 时按 L1/L2/L3 呈现）
   - 需确认的设计决策（如有）
   - 补充验收条件（issue 中未明确的，主动提出）
7. **硬性停止 — 展示以上全部内容后，必须停下来等待用户确认。禁止自行进入阶段二。用户未回复 = 未确认。**

#### 阶段二：执行（用户确认后）

根据 issue 状态分支：

##### `raised`（产品经理尚未回复）
1. 写 BE/FE/QA 技术分析 comment（含影响范围、推荐方案、补充的验收条件）
2. 更新 label → `[role]-reviewed`
3. 如果 label 含 `needs-pm` → 同步录入 `<project>/docs/problems/needs-pm-queue.md` 的 Open 段
4. 停止，等待产品经理回复

##### `pm-reviewed`（执行层可实现）

按以下步骤顺序执行：

**Step 0 — 在 Issue 贴执行清单 comment（硬性前置，不得跳过）**

用户确认方案后，**第一个动作**是在 Issue 贴一条清单 comment，作为本次执行的唯一追踪源：

```markdown
## 执行清单（[role]）

- [ ] 标 `[role]-in-progress` + comment "开始实现"
- [ ] 文档先行（标记待实现，同步共享仓库）
- [ ] 实现
- [ ] 编译门禁（项目 CLAUDE.md 指定的 compile_cmd）
- [ ] 测试
- [ ] 6a — 模块文档 / api-spec / 追溯链更新，同步共享仓库
- [ ] 6b — problem-registry 同步
- [ ] 6b — Issue comment（commit hash + 文档链接 + 摘要）
- [ ] 6b — label **保持 `[role]-in-progress`**（不在收尾时切），注明 "工作完成 — 等 /release 发版"
- [ ] 6b — 等 `/release` 发版成功后由 release skill 批量切 `[role]-in-progress` → `[role]-confirmed` + comment 嵌入「可关闭」
```

贴出后**每完成一步立即更新对应 checkbox**（`[ ]` → `[x]`）。清单是执行的唯一追踪源，不靠记忆。

1. **标 `[role]-in-progress`（硬门禁）**
   - 0a. 执行 `gh issue view <N> --json labels` 验证含 `[role]-in-progress` label
   - 0b. 无则立即补：`gh issue edit <N> --add-label [role]-in-progress` + comment "开始实现"（不填 ETA）
   - 0c. **未通过禁止进入 Step 2** —— 这是机器化门禁，不是建议
   - 背景：纯文字约束在违规率高的环境下需机器化校验
2. **文档先行** — 判断是否触发 ADR（架构调整 / 技术选型 / 明显取舍）；更新 / 创建相关文档，标记"待实现"，同步共享仓库并推送
3. **实现** — 确认分支状态（工作区干净、基准分支正确）；实现代码 / 文档改动；遵循 `rules/core/research-first.md` 和 `rules/core/incremental-verification.md`
4. **编译门禁（硬门禁）** — 跑 `<compile_cmd>`（项目 CLAUDE.md 指定），失败 → 修 → retry（max 2 次），仍失败 → 停下来上报用户
5. **测试** — 按改动类型执行最低测试要求
6. **收尾 6a** —
   - 更新文档状态为"已实现"，同步共享仓库并推送
   - 更新模块清单追溯链：补充本次新增 / 变更的 API、数据模型、ADR、Issue 关联
7. **收尾 6b** —
   - 同步 problem-registry（如有对应 P-xxx 条目，更新状态为 resolved）
   - **needs-pm-queue 状态同步**（如 Issue 曾入 Open 段）
   - 写 Issue comment（commit hash + 共享仓库文档链接 + 摘要）
   - **文档同步声明强制三选一**（防止 50% 缺失率）：
     - [ ] 列出每份已更新文档的 commit hash
     - [ ] 显式标 "无需更新（原因：<具体原因>）"
     - [ ] 标 "延后处理（追加 Issue ref：#NNN）"
     - **三选一缺失 = 收尾未完成**（不能切 confirmed 状态）
   - **commit hash 格式**：用 `` commit: `<hash>` `` 格式（前缀 `commit:` + 反引号包裹），机器可解析，供 `/release` 反查是否已部署
     - ✅ 正确：`` commit: `9818485` ``
     - ❌ 错误：heredoc 内 `commit: \\\`9818485\\\``（反斜杠转义后落库变字面量，regex 不命中）
     - 写法守则：heredoc / shell 字符串内**不要给反引号加反斜杠**。如有 shell 解析顾虑，用 `gh issue comment <N> --body-file <path>` 改走外部文件
   - **comment 落库后自检（强制）**：写 comment 后立即 grep 验证可识别：
     ```bash
     gh issue view <N> --comments | grep -oE 'commit:\s*`[a-f0-9]{7,40}`'
     ```
     命中 = 合规（≥ 1 行输出）；空输出 = 格式不合规，**立即重写 comment**
   - **label 保持 `[role]-in-progress` 不变**（不在收尾时切；由 `/release` 发版成功后批量切 `[role]-confirmed`）
   - 注明 "工作完成 — 等 /release 发版"（6a 全部完成 + 三选一已填后；「可关闭」由 /release Step 6b.1 切 label 时同步嵌入）

##### 纯文档类
1. 更新文档 + 同步共享仓库
2. 更新模块清单追溯链（如影响模块关联产物）
3. 同步 problem-registry（如有对应 P-xxx 条目）
4. 写 comment + 更新 label → `[role]-confirmed` + comment 嵌入「可关闭」

---

## 约束

- 每一步完成后再进入下一步，不跳步
- Issue comment 中文档链接指向**共享文档仓库**，不使用代码仓库路径
- `closed` 只由人工触发，agent 不主动 close，完成后 comment 中注明「可关闭」
- 遵循 `rules/core/large-module.md`：涉及较大改动时先呈现方案，等用户确认
- 遵循 `rules/extension/adr-discipline.md`（按需启用）：架构级决策在编码前生成 ADR

---

## 与协议层的关系

参考 `protocols/issue-process.md` —— 本 skill 是 issue-process 协议的执行层 skill。状态机定义 / 5 维度审查由 protocol 描述，本 skill 实施 PDCA。
