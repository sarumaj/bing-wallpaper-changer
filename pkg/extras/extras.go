/*
Package extras contains extra functions that are not part of the main package.
Currently, custom watermarks are implemented here.
*/
package extras

import (
	"bytes"
	"embed"
	"io"
	"path/filepath"
	"slices"
	"strings"
)

const DefaultFontName = "unifont.ttf"
const DefaultWatermarkName = "sarumaj.png"

//go:embed fonts/*.ttf
var fonts embed.FS

//go:embed watermarks/*.png
var watermarks embed.FS

// EmbeddedFonts returns a map of available fonts.
var EmbeddedFonts = getEmbedded(fonts, "fonts")

// EmbeddedWatermarks returns a map of registered watermarks.
var EmbeddedWatermarks = getEmbedded(watermarks, "watermarks")

// Embedded is a map of embedded files.
type Embedded map[string]io.Reader

// Keys returns the keys of the embedded map.
func (e Embedded) Keys() []string {
	names := make([]string, 0, len(e))
	for k := range e {
		names = append(names, k)
	}

	slices.Sort(names)
	return names
}

// String returns the keys of the embedded map as a string.
func (e Embedded) String() string {
	return strings.Join(e.Keys(), ", ")
}

// getEmbedded returns a map of embedded files.
func getEmbedded(fsys embed.FS, path string) Embedded {
	m := make(Embedded)
	files, _ := fsys.ReadDir(path)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		r, _ := fsys.Open(filepath.Join(path, file.Name()))
		defer r.Close()

		buffer := bytes.NewBuffer(nil)
		_, _ = io.Copy(buffer, r)

		m[file.Name()] = buffer
	}

	return m
}
