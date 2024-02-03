/*
Package extras contains extra functions that are not part of the main package.
Currently, custom watermarks are implemented here.
*/
package extras

import (
	"bytes"
	"embed"
	"io"
)

const DefaultFontName = "unifont.ttf"
const DefaultWatermarkName = "sarumaj.png"

//go:embed fonts/*.ttf
var fonts embed.FS

//go:embed watermarks/*.png
var watermarks embed.FS

// EmbeddedFonts returns a map of available fonts.
var EmbeddedFonts = func() map[string]io.Reader {
	m := make(map[string]io.Reader)
	files, _ := fonts.ReadDir("fonts")

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		r, _ := fonts.Open("fonts/" + file.Name())
		defer r.Close()
		b, _ := io.ReadAll(r)
		m[file.Name()] = bytes.NewReader(b)
	}

	return m
}()

// EmbeddedWatermarks returns a map of registered watermarks.
var EmbeddedWatermarks = func() map[string]io.Reader {
	m := make(map[string]io.Reader)
	files, _ := watermarks.ReadDir("watermarks")

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		r, _ := watermarks.Open("watermarks/" + file.Name())
		defer r.Close()
		b, _ := io.ReadAll(r)
		m[file.Name()] = bytes.NewReader(b)
	}

	return m
}()
