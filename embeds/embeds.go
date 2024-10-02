package embeds

import "embed"

//go:embed input.sce
var Input []byte

//go:embed *
var FS embed.FS
