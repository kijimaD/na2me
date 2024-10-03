package embeds

import "embed"

// 坊っちゃん
//
//go:embed scenario/bochan.sce
var Bochan []byte

// 坊っちゃん
//
//go:embed scenario/wagahai_ha_neko_dearu.sce
var WagahaiHaNekoDearu []byte

// 三四郎
//
//go:embed scenario/sanshirou.sce
var Sanshirou []byte

//go:embed *
var FS embed.FS
