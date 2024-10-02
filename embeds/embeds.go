package embeds

import "embed"

// 坊っちゃん
//
//go:embed scenario/bochan.sce
var Bochan []byte

//go:embed *
var FS embed.FS
