# FB Index — 启动前扫描索引

**定位：** 团队级 / 个人级 feedback（FB）的结构化元数据索引，支持 `/fb-scan` skill 按 skill / module / phase / category 筛选。

- **数据源：** `feedback/*.md`（按日期或主题分组的源文件）
- **维护规则：** 新 FB 录入时必须同步在此索引追加条目
- **字段定义：** 见本文件末尾 schema

---

## 编号规范

`FB-NNN` 连续递增，不按批次重置。

## 状态枚举

| 状态 | 含义 |
|------|------|
| `candidate` | 候选，未达 ≥ 2 例阈值 |
| `observing` | 观察期，已达阈值待累积更多实证 |
| `applied` | 已 applied 到规则文件 / SOP |
| `verified` | applied 后实际生效（产出 ≥ 1 次拦截真实问题）|
| `dismissed` | 排除（噪声 / 重复 / 已被 别 FB 覆盖）|

---

## FB 条目（示例骨架，删除后追加真实条目）

## FB-001 — <标题>
- **date**: YYYY-MM-DD
- **file**: <feedback/source-file.md>
- **category**: audit / process / design / implement / meta
- **skills**: <relevant skills>
- **modules**: <relevant modules or "(all)">
- **phases**: <relevant phases or "—">
- **severity**: low / medium / high / critical
- **status**: candidate
- **occurrences**: 1
- **guidance**: <guideline 一句话>
- **scan_when**: <什么场景启动前扫描本条 FB>
- **related**: <related FB IDs>

---

## 统计

| 维度 | 数量 |
|------|------|
| 总计 | 0 |
| critical | 0 |
| high | 0 |
| medium | 0 |
| low | 0 |
| applied 状态 | 0 |
| observing 状态 | 0 |

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
