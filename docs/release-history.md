# Release History

> 由 `/release` skill 维护,append-only,仅保留最近 10 条成功构建。
> `head_commit` 用于 Step 6.3 的 `git merge-base --is-ancestor` 判断"某 commit 是否已部署"。

| 构建号 | 日期 | 分支 | 耗时(秒) | head_commit | 结果 | 备注 |
|-------|------|------|---------|-------------|------|------|
| #1 | 2026-05-15 | master | 30 (mock) | `0859892` | SUCCESS | **P-A 演练 mock 记录** — 实际未触发 CI,本条人工写入用于演练 Step 6.3 label 切换算法 |
