# Skill Dogfood Feedback — kingdom-rush bootstrap (2026-05-15)

> **产出方:** kingdom-rush 项目 dogfood session（首次 `/issue` + `/release` 演练）
> **目标仓库:** agent-dev-standard
> **格式:** 按 `templates/registry-schemas/fb-index.md.template` schema
> **编号约定:** 本文档用临时前缀 `KR-FB-NNN`（kingdom-rush dogfood），合入 standard `fb-index.md` 时由维护方按 `FB-NNN` 连续递增重号

---

## 背景

kingdom-rush 项目刚跑完 `agent-dev-standard install`，作为首次 dogfood 演练:

1. 从零 bootstrap (git init + GitHub remote + Go 骨架 + GHA labels 同步)
2. `/issue` 演练 — 建 issue #1（main 函数 v0.0 → v0.1）+ 完整 PDCA 11 step
3. `/release` 演练（P-A: Mock + 只跑 Step 6.3）— 算法走通 + label 切换

演练过程发现 5 条规则级 FB（与项目本身无关，属 standard 仓库需改进项）。所有实证均可追溯到具体 file:line + git commit + issue comment 链接。

---

## 总览

| 编号 | 标题 | category | severity | 影响 skill/file |
|---|---|---|---|---|
| KR-FB-001 | `.gitignore` 模板不区分 hub vs business 项目 | implement | medium | install/modules/06-templates.sh + templates/hub-gitignore.template |
| KR-FB-002 | `role` 字段大小写未规范化 | implement | medium | templates/CLAUDE.md.template + skills/core/issue/SKILL.md |
| KR-FB-003 | S 级 + 项目无 spec 的文档先行豁免规则未明示 | design | low | skills/core/issue/SKILL.md Step 2 |
| KR-FB-004 | `gh issue list --label X` 在新建 label 后立刻跑会扑空 | implement | medium | skills/core/release/SKILL.md Step 6.3 |
| KR-FB-005 | commit message `closes #N` 触发 GitHub auto-close，违反 SKILL.md 约束 | implement | **high** | skills/core/issue/SKILL.md Step 6 |

**severity 分布:** high × 1 / medium × 3 / low × 1

---

## KR-FB-001 — `.gitignore` 模板不区分 hub vs business 项目

- **date:** 2026-05-15
- **file:** docs/feedback/2026-05-15-skill-dogfood-fb.md
- **category:** implement
- **skills:** install
- **modules:** install/modules/06-templates.sh + templates/hub-gitignore.template
- **phases:** install Phase 6
- **severity:** medium
- **status:** candidate
- **occurrences:** 1
- **guidance:** 06-templates.sh 渲染 .gitignore 时必须根据项目类型 (hub vs business) 选不同模板；当前统一套 `hub-gitignore.template` 导致业务项目 `CLAUDE.md` 和 `.claude/` 被静默忽略
- **scan_when:** 新项目 install 完后第一次 git add / commit / push 前
- **related:** —

### 背景与影响

`hub-gitignore.template` 自己头部注释 (L1-2) 明示 "适用：standard hub repo 及类似分发型仓库 / 避免 self-install dogfood 把'给别人 install 后才该产生的'产物误提交到主仓库"——但 `06-templates.sh` L244-258 不区分项目类型，对**所有项目**都套同一份模板。

业务项目跑完 install 后:
- `CLAUDE.md`（项目级 agent 上下文）被忽略 → 无法 commit → 跨机器协作的基础缺失
- `.claude/`（per-project skills + protocols）被忽略 → 团队共享 skill 配置失败

`06-templates.sh` L244 注释承认这是 "v0.3 Finding 5"——已发现但未完成修复。

### 实证

- **standard 文件:** `install/modules/06-templates.sh:244-258`（cp `$hub-gitignore.template` `$project/.gitignore`，无 project_type 区分）
- **standard 文件:** `templates/hub-gitignore.template:39` `/CLAUDE.md`，`:42` `/.claude/`
- **kingdom-rush 实测:**
  - 第一次 `git status` 显示 CLAUDE.md / .claude/ 都不在 untracked 列表（被 ignore）
  - 手动删 .gitignore 两行后才正常 stage
  - 见 commit `8c413bf` message 第 5 项

### 修复建议

1. 新增 `templates/business-gitignore.template` —— 不含 `/CLAUDE.md` 和 `/.claude/` 忽略规则，其他保留
2. `06-templates.sh` 加 `project_type` 参数（`hub` | `business`，默认 `business`），按类型选模板
3. `hub-gitignore.template` 头部注释强化警告："仅给 hub repo self-install dogfood 用，业务项目勿套"
4. install.sh 交互流程加一个问询：项目类型（hub / business），默认 business

---

## KR-FB-002 — `role` 字段大小写未规范化

- **date:** 2026-05-15
- **file:** docs/feedback/2026-05-15-skill-dogfood-fb.md
- **category:** implement
- **skills:** install + issue
- **modules:** templates/CLAUDE.md.template + skills/core/issue/SKILL.md
- **phases:** install + issue execution
- **severity:** medium
- **status:** candidate
- **occurrences:** 1
- **guidance:** `labels.yml.template` 用 lowercase (be-reviewed / be-in-progress)，CLAUDE.md `## Issue 配置` 段 `role` 字段也必须 lowercase；install 时强制 normalize 或 SKILL.md issue 明示约束
- **scan_when:** install 完跑 `/issue` 时第一次 `gh issue edit --add-label [role]-*`
- **related:** KR-FB-001（同源于 install 模板未规范化）

### 背景与影响

`templates/labels.yml.template` 定义 18 state labels 全部 lowercase (`be-reviewed` / `be-in-progress` / `be-confirmed` 等)。`skills/core/issue/SKILL.md` 使用 `[role]-in-progress` 占位符直接字符串拼接。

如果 `CLAUDE.md` 的 `role: BE`（大写），拼出 `BE-in-progress`，`gh issue edit --add-label BE-in-progress` 失败（label 不存在）。

CLAUDE.md.template 应该明示 `role` lowercase 约束，或 install 时强制 normalize。

### 实证

- **standard 文件:** `templates/labels.yml.template:20-83` —— 18 state labels 全部 lowercase
- **standard 文件:** `skills/core/issue/SKILL.md:166-170` —— `gh issue edit <N> --add-label [role]-in-progress` 直接拼接
- **kingdom-rush 实测:**
  - 第一次写 `## Issue 配置` 段 role 字段时按习惯填 `BE`（大写）
  - 演练前发现问题，提前修正 → commit `c87ea6b` "fix: role 字段 BE → be (对齐 labels.yml 小写约定)"

### 修复建议

1. `templates/CLAUDE.md.template` `## Issue 配置` 段示例改为 `role: be`（lowercase）+ 加注释 "**必须 lowercase**，对齐 labels.yml"
2. `skills/core/issue/SKILL.md` Step 0 加前置 normalize：`role=$(echo $role | tr 'A-Z' 'a-z')`
3. install.sh 在询问 role 时显式 lowercase prompt

---

## KR-FB-003 — S 级 + 项目无 spec 的文档先行豁免规则未明示

- **date:** 2026-05-15
- **file:** docs/feedback/2026-05-15-skill-dogfood-fb.md
- **category:** design
- **skills:** issue
- **modules:** skills/core/issue/SKILL.md
- **phases:** issue Step 2 (文档先行)
- **severity:** low
- **status:** candidate
- **occurrences:** 1
- **guidance:** S 级改动 + 项目 bootstrap 阶段（无 spec / api / 共识文档）时，文档先行应可豁免；SKILL.md Step 2 应明示豁免段（类似 Step 0 的"规模豁免"），定义触发条件 + 必填字段
- **scan_when:** 第一次跑 `/issue` 时项目尚无 spec 文档
- **related:** —

### 背景与影响

SKILL.md issue L171 写: "**文档先行** — 判断是否触发 ADR（架构调整 / 技术选型 / 明显取舍）；更新 / 创建相关文档，标记'待实现'，同步共享仓库并推送"

但**没明示**:
- 项目 bootstrap 阶段（尚无 spec / api / 共识文档）时怎么办？
- S 级改动（单文件单行字面量）是否需要文档先行？
- 如果豁免，理由应该怎么写？

对比 Step 0 自检表已经定义了"规模豁免"（L66）的明确触发条件，Step 2 没有对等的豁免段。

### 实证

- **standard 文件:** `skills/core/issue/SKILL.md:171` —— 未涵盖项目无 spec 场景
- **standard 文件:** `skills/core/issue/SKILL.md:66` —— Step 0 已有"规模豁免"段可作类比
- **kingdom-rush 实测:** 演练 #1 跑到 Step 2 时只能"自行推理"豁免依据 + 写 comment 解释（见 issue #1 [comment#4458847614](https://github.com/Corray/kingdom-rush/issues/1#issuecomment-4458847614)）

### 修复建议

在 SKILL.md issue Step 2 增加豁免段:

```markdown
**文档先行豁免规则:**

满足全部条件可豁免，需在 comment 明示理由 + 字段:
- S 级改动（Step 0 自检表分级为 S）
- 项目处 bootstrap / 早期阶段，尚无 spec / api / 共识文档可对齐
- Issue body 自含 SSOT (背景 + 验收条件齐备)
- 不触发 ADR

豁免 comment 必填字段:
- 豁免理由（具体条件命中）
- 项目当前阶段
- 后置承诺（项目 active phase 后补建文档时的追溯锚点）
```

---

## KR-FB-004 — `gh issue list --label X` 在新建 label 后立刻跑会扑空

- **date:** 2026-05-15
- **file:** docs/feedback/2026-05-15-skill-dogfood-fb.md
- **category:** implement
- **skills:** release（可能影响 audit / 其他依赖 gh issue list 的 skill）
- **modules:** skills/core/release/SKILL.md
- **phases:** release Step 6.3
- **severity:** medium
- **status:** candidate
- **occurrences:** 1
- **guidance:** GitHub label-based search 索引有延迟（实测刚 sync 完几分钟内 `--label` 过滤返回空）；SKILL.md release Step 6.3 L159-161 直接用 `gh issue list --label X` 在 install 直跑场景下必踩；建议加 retry 或 fallback 到客户端 jq filter
- **scan_when:** install 完立刻跑 `/release` 演练，或 `09-github-labels.sh` 同步完几分钟内的任何 label-based 查询
- **related:** —

### 背景与影响

SKILL.md release L159-161:
```bash
gh issue list --repo <owner/repo> --label [role]-in-progress --state open \
  --json number,title,comments
```

这是 Step 6.3 V2 算法的 Algo Step 1，决定后续所有处理。

实测: 在 `09-github-labels.sh` 同步完 24 labels + 创建 issue + 加 label 后**几分钟内**，`gh issue list --label be-in-progress` 返回 `[]`，**但**:
- `gh issue view 1 --json labels` 显示 `be-in-progress` label 确实存在
- `gh api repos/Corray/kingdom-rush/issues/1` 直查 REST API 也确实显示 label

说明 issue **有该 label**，只是 GitHub 的 **label-based search index** 还没收录新 label。`gh issue list --label` 走 search API，所以扑空。

生产环境 install 一次后不易踩（label 早就稳定了），但**install 直跑 audit / release / 任何 label-based 查询会必踩**。

### 实证

- **standard 文件:** `skills/core/release/SKILL.md:159-161`
- **kingdom-rush 实测:**
  - 同步 24 labels: `bash 09-github-labels.sh`（约 17:45）
  - 创建 issue #1: 17:50
  - 加 be-in-progress label: 17:54
  - 跑 `gh issue list --label be-in-progress`（约 17:56）→ `[]`
  - 同时间 `gh issue view 1 --json labels` → `["pm-reviewed","be-in-progress","feature"]`
  - 改用 `gh issue list --state all --json` + `jq '[.[] | select(.labels[]?.name == "be-in-progress")]'` → 正确找到 #1

### 修复建议

1. **首选 fallback**: SKILL.md L159 增加备用命令 + 自动 fallback 逻辑:
   ```bash
   # 主路径
   issues=$(gh issue list --label [role]-in-progress --state open --json number,title,comments)
   # 主路径返回空且 label 同步晚于 label-search-stable 阈值 → 走 fallback
   if [ "$issues" = "[]" ]; then
       issues=$(gh issue list --state open --json number,title,comments,labels \
                | jq '[.[] | select(.labels[]?.name == "[role]-in-progress")]')
   fi
   ```
2. **明示警告**: SKILL.md L159 加注释 "**注意:** GitHub label-based search 索引有延迟（约 5-15 min）。`09-github-labels.sh` 同步完立刻跑可能扑空，必要时改用客户端 jq filter（见 fallback）"
3. **长远**: 考虑给 `09-github-labels.sh` 加可选 `--wait-for-index` 参数，sync 完后 sleep 一段或主动 poll search 直到索引就绪

---

## KR-FB-005 — commit message `closes #N` 触发 GitHub auto-close，违反 SKILL.md 约束

- **date:** 2026-05-15
- **file:** docs/feedback/2026-05-15-skill-dogfood-fb.md
- **category:** implement
- **skills:** issue（影响 release Step 6.3 算法前提）
- **modules:** skills/core/issue/SKILL.md
- **phases:** issue Step 6 commit
- **severity:** **high**
- **status:** candidate
- **occurrences:** 2 _(2026-05-18 更新: IPR-001 作为第 2 例实证,见 docs/audit/findings-registry.md)_
- **guidance:** commit message 不要用 git auto-close 关键词（closes / fixes / resolves / 等）指向本次 issue，会触发 GitHub 自动关闭违反 "closed 只由人工触发" 约束；用 `refs #N` / `related #N` / `see #N` 替代
- **scan_when:** 跑 `/issue` Step 6 写 commit message 时
- **related:** —

### 背景与影响

SKILL.md issue L211 明示:
> `closed` 只由人工触发，agent 不主动 close，完成后 comment 中注明「可关闭」

但**Step 6 commit 阶段没有任何指引**关于 commit message 怎么写。

GitHub 默认行为: 含 `closes #N` / `fixes #N` / `resolves #N` 等关键词的 commit push 到 default branch 时**自动关闭** issue #N。

agent 写出 `closes #1` 的 commit → push → GitHub 替 agent 关掉 issue → **等效于 agent 间接 close**，违反 L211 约束。

**真实影响:**
- `/release` Step 6.3 V2 算法 L159-161 假设 in-progress issue 是 `--state open`。如果 issue 因 `closes #N` 提前关掉，list 查不到，label 切换跳过 → 与 `/release` 的核心承诺（"`[role]-confirmed` 由 release 切"）矛盾
- 必须 reopen 才能继续走 Step 6.3

### 实证

- **standard 文件:** `skills/core/issue/SKILL.md:211` —— "closed 只由人工触发,agent 不主动 close"
- **kingdom-rush 实测:**
  - commit `0859892` message:
    ```
    fix: main 函数版本号 v0.0 → v0.1

    closes #1
    ```
  - `git push` 完后立即 `gh api repos/Corray/kingdom-rush/issues/1` → `state: "closed"`
  - 后续 `/release` Step 6.3 Algo Step 1 `gh issue list --state open` 扑空（state 不对）
  - 不得不 `gh issue reopen 1` 才能继续演练

### 修复建议

1. **强约束**: SKILL.md issue Step 6 commit message 段加禁用列表:
   ```
   ❌ 禁止 commit message 含: closes / close / closed / closing
                              fixes / fix / fixed / fixing
                              resolves / resolve / resolved / resolving
                              (case-insensitive)
   ✅ 推荐: refs #N / related #N / see #N / cf #N
   ```
2. **正向引导**: 给 commit message 提供模板示例:
   ```
   <type>: <subject>

   refs #<issue-number>
   ```
3. **机器化**: 长远在 commit-msg hook 中 grep 拦截 auto-close 关键词 + issue 引用组合（hook 由 install 模块 05-core-hooks 注入）

---

## 附录: 元观察

**Standard 仓库自身未启用 fb-index.md 实例。** `templates/registry-schemas/fb-index.md.template` 已就绪 + schema 完整，但 standard 仓库自身没有 `feedback/` 目录 + `fb-index.md` 实例文件。意味着 standard 的"规则级 FB 收集机制" template 完成但 dogfood 尚未启动——本批 FB 是潜在的首批 dogfood 样本。

**严格说不是缺陷**（template 就绪即可，由维护方决定何时启用 instance），但作为元观察记录。

---

## 回传 standard 仓库的方式

按 standard `README.md:33` 指引: **反馈 / 提议 → hub Issues at https://github.com/shangronggu-cyber/agent-dev-standard/issues**

| 方式 | 描述 | 工作量 |
|---|---|---|
| **A. 单 GitHub Issue** | 在 hub 提一个总览 issue，body 含完整 5 条 FB + 链接本文档 | 5 min |
| **B. 5 个独立 GitHub Issue** | 每条 FB 一个 issue，便于 standard 维护方独立 label / track / close | 15 min |
| **C. PR + 直接修补** | clone standard 仓库，按修复建议直接改文件，提 PR | 30+ min |
| **D. 只保存本地** | 本文档留 kingdom-rush 项目，不立即回传，PM 阶段性人工整理 | 0 min |

**说明:**
- standard 仓库 origin 在 Bitbucket (`chatly-biz-tool/agent-dev-standard`)，但反馈渠道是 GitHub (`shangronggu-cyber/agent-dev-standard`)——可能跨平台镜像
- A / B 都需要先验证 `Corray` 账号能否向 `shangronggu-cyber/agent-dev-standard` 提 issue (public + community / private + collaborator 视设置而定)
- C 是最直接的贡献方式，但需要 standard 维护方接受 PR 流程
- D 是最保守的，没回传节奏的话先攒着
