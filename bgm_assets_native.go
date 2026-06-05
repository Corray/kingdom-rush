//go:build !js && !term

// V7.4 B1: BGM 资产提供层 — native build 保持 embed (二进制自包含)。
package main

import _ "embed"

//go:embed assets/bgm/menu.ogg
var oggBgmMenu []byte

//go:embed assets/bgm/battle.ogg
var oggBgmBattle []byte

func loadBGMAssets() {
	queueBGMData(bgmMenu, oggBgmMenu)
	queueBGMData(bgmBattle, oggBgmBattle)
}
