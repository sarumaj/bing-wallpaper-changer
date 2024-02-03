package extras

import "embed"

const DefaultFontName = "unifont.ttf"

//go:embed fonts/*.ttf.gz
var fonts embed.FS

// EmbeddedFonts returns a map of available fonts.
var EmbeddedFonts = getEmbedded(fonts, "fonts")
