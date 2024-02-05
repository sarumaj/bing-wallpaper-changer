/*
Package extras contains extra functions that are not part of the main package.
Currently, custom watermarks are implemented here.
*/
package extras

import (
	"compress/gzip"
	"embed"
	"io"
	"path/filepath"
	"slices"
	"strings"
)

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

		var r io.ReadCloser
		var err error
		r, err = fsys.Open(filepath.Join(path, file.Name()))
		if err != nil {
			panic(err)
		}

		defer r.Close()

		if filepath.Ext(file.Name()) == ".gz" {
			r, err = gzip.NewReader(r)
			if err != nil {
				panic(err)
			}

			defer r.Close()
		}

		raw, err := io.ReadAll(r)
		if err != nil {
			panic(err)
		}

		m[strings.TrimSuffix(file.Name(), ".gz")] = &multiReadReader{data: raw}
	}

	return m
}
