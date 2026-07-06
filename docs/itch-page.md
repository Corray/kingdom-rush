# itch.io 页面配置（Gopher Defense）

> 分发 artifact：itch.io 项目创建/更新时照抄本文件。上传包 = 仓库根 `gopher-defense-itch.zip`（`make wasm` 后由 web/ 打包，index.html 在 zip 根目录）。

## 基本信息

| 字段 | 值 |
|------|----|
| Title | Gopher Defense |
| Project URL | gopher-defense（建议 slug） |
| Classification | Games |
| Kind of project | HTML |
| Pricing | No payments（免费） |
| Genre | Strategy |
| Tags | tower-defense, pixel-art, golang, wasm, singleplayer, strategy, 2d |

## Embed 设置

| 设置 | 值 |
|------|----|
| Viewport dimensions | 1000 × 640 |
| Fullscreen button | 开 |
| Mobile friendly | 关（键鼠为主，触屏未适配） |
| Automatically start on page load | 关（用户点击启动 → 浏览器 autoplay 政策下音频正常解锁） |

上传 zip 后勾选 **"This file will be played in the browser"**。

## Short description / tagline

> A tower defense game written in Go, running in your browser. 20 dual-path levels, 6 towers with upgrade branches, 3 hero classes, skill trees, achievements, endless mode.

## Description（页面正文，粘贴为 markdown）

```
Defend the kingdom across **20 handcrafted dual-path levels** — every level has two entry
points converging on a shared exit, so you can never turtle on a single chokepoint.

**Play in browser. No install. Progress saves locally.**

### Features

- **6 tower types** — Archer, Cannon, Magic, Frost, Tesla (chain lightning), Sniper
  (long-range anti-boss) — each with 3 levels and **branching upgrades**: pick Marksman
  or Rapidfire, Mortar or Gatling, Archmage or Enchanter, Deep Freeze or Hailstorm
- **8 enemy types** — runners, fliers, armored shields, regenerators, healers, spawners,
  and bosses with zone-specific abilities (shield / charge / summon)
- **3 hero classes** — Knight (tanky blocker), Archer (long-range kiter), Rogue (fast
  skirmisher). Your hero levels up each battle and unlocks a cleave ability
- **Skill trees** — spend stars earned from level ratings on permanent per-class perks;
  the three trees cost more than you can ever earn, so pick your build
- **16 achievements**, 4 themed zones (Forest / Desert / Snow / Lava), 3 difficulties,
  3-star ratings, and a budget-scaling **endless mode**
- Full SFX + chiptune BGM (all CC0 assets)

### Controls

Arrows / mouse to aim · 1-6 select tower · Space / click build & upgrade · V switch
upgrade branch · X sell · T targeting · H rally hero · G cleave · R meteor strike ·
P pause · F 2x speed · full key list on the title screen

---

Written in Go, compiled to WebAssembly. Open source (MIT) — code on
[GitHub](https://github.com/Corray/kingdom-rush). Also playable on
[GitHub Pages](https://corray.github.io/kingdom-rush/).

Art & audio CC0 (Kenney, Juhani Junkala). Inspired by Kingdom Rush — not affiliated
with Ironhide Game Studio.
```

## 截图（Screenshots 上传，仓库路径）

1. `docs/screenshots/battle.png` — 双路战斗
2. `docs/screenshots/menu.png` — 关卡选择
3. `docs/screenshots/skilltree.png` — 技能树
4. `docs/screenshots/achievements.png` — 成就

Cover image（630×500）：可从 battle.png 裁切，或后续单独做。

## butler 推送（可选，装好后一条命令更新）

```sh
# 安装（macOS）: https://itch.io/docs/butler/installing.html
butler login
butler push gopher-defense-itch.zip <username>/gopher-defense:html5
```

## 更新记录

| 日期 | 动作 |
|------|------|
| 2026-07-06 | 初版文案 + zip 打包脚手架（V16 / 218 tests 构建） |
