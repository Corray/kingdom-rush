# FB-006 (候选) — audit skill 维度 3 / #14/#16 缺豁免 comment 识别能力

**来源:** kingdom-rush dogfood / IPR-002 adjudication 衍生（2026-05-18）

---

## FB Summary

- **date:** 2026-05-18
- **file:** docs/feedback/2026-05-18-audit-skill-exemption-detection-fb.md
- **category:** audit
- **skills:** audit（issue-process phase 维度 3 #14 文档先行 / #16 测试覆盖）
- **modules:** standard `skills/core/audit/SKILL.md` Step 3 正向比对 + protocols/issue-process.md 维度 3
- **phases:** audit issue-process
- **severity:** medium
- **status:** candidate（首例,未达 ≥ 2 例阈值）
- **occurrences:** 1
- **guidance:** audit 维度 3 #14（文档先行）/ #16（测试覆盖）等结构化检查应识别 issue 内的 "Step 2 豁免 comment"（grep 特定关键词如 "文档先行豁免" / "S 级豁免" / "测试豁免"）→ 不计违规,避免 false positive
- **scan_when:** 跑 `/audit issue-process` Step 3 维度 3 检查时（特别是 S 级或 bootstrap 阶段项目）
- **related:** KR-FB-003（豁免规则未明示）/ IPR-002（首例实证）

---

## 背景

### 触发事件

2026-05-15 kingdom-rush 首次 audit issue-process 演练:
- Issue #1 触发了"文档先行" commit 检查（维度 3 #14）
- 实际:项目 bootstrap 阶段无 spec/api 文档可改,Step 2 自检判定豁免,comment `#4458847614` 已贴完整豁免依据
- audit skill 仍机械化结构检查 → 标违规 IPR-002（MEDIUM）
- review 判 dismissed（false positive）

### 根因

`skills/core/audit/SKILL.md` Step 3 正向比对当前是**纯结构化检查**（"有没有文档先行 commit"）,**不识别 issue 内的语义豁免 comment**。即使 `/issue` Step 2 已贴明确豁免依据,audit 仍会标违规。

跨项目普适性:✓ —— 任何 bootstrap 阶段项目 / S 级 issue / 无 spec 项目都会触发同类 false positive。

### 改进建议

**短期（audit skill 侧）:**
1. Step 3 维度 3 #14 检查时,先 `gh issue view <N> --comments | grep -E "文档先行豁免|S 级豁免|bootstrap 豁免"`,命中则标 ✅（豁免）而非 🟡 偏差
2. 同理 #16（测试覆盖）应识别 "测试豁免" 关键词

**中期（spec 侧）:**
3. `protocols/issue-process.md` 维度 3 增段 "豁免识别规则",明示哪些 comment 关键词等价于豁免依据
4. `skills/core/issue/SKILL.md` Step 2 自检表增"豁免 comment 落标规则",指引如何贴出 audit 可识别的豁免 comment

**长期:**
5. 把"豁免 comment 识别"提取为 audit skill 的通用能力,跨维度（#13 LMP / #14 文档先行 / #16 测试 / 等）共用同一识别引擎

---

## 实证

| 编号 | 日期 | 项目 | 维度 | 现象 |
|------|------|------|------|------|
| 1 | 2026-05-15 | kingdom-rush | 维度 3 #14 | IPR-002 / P-002,豁免依据 comment #4458847614 已贴,audit 仍标 false positive |

---

## 达阈值（≥ 2 例）后的形式化路径

按 `formalization-timing.md` 类型 A（异象探索）:
- 候选 status → 累积第 2 例 → observing
- observing 期内验证 + 考验样本各 ≥ 1（不同结构家族）→ applied 到 standard audit SKILL.md
- applied 后实际生效（拦截 ≥ 1 次真实 false positive）→ verified

**本条暂为 candidate（1 例）,先观察 kingdom-rush 后续 audit 是否复现;若仅 1 例不复现,3 个月后归 dismissed。**

---

## fb-index 同步

> `docs/problems/fb-index.md` 当前是骨架（仅 FB-001 示例条目）,本 FB 暂不录入 fb-index,等 fb-index 正式启用 KR-FB-XXX 录入流程时一并加入。
