package extras

import "embed"

const DefaultWatermarkName = "sarumaj.png"

//go:embed watermarks/*.png.gz
var watermarks embed.FS

// EmbeddedWatermarks returns a map of registered watermarks.
var EmbeddedWatermarks = getEmbedded(watermarks, "watermarks")
