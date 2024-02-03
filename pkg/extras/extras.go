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

const DefaultWatermarkName = "sarumaj.png"

//go:embed watermarks/*.png
var watermarks embed.FS

var RegisteredWatermarks = func() map[string]io.Reader {
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
