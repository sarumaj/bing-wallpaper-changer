/*
Package extras contains extra functions that are not part of the main package.
Currently, custom watermarks are implemented here.
*/
package extras

import (
	"compress/gzip"
	"embed"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sarumaj/bing-wallpaper-changer/pkg/logger"
)

// Embedded is a map of embedded files.
type Embedded map[string]io.ReadCloser

func (e Embedded) ToFiles(name string) (string, error) {
	dir := filepath.Join(os.TempDir(), name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	for k, v := range e {
		path := filepath.Join(dir, k)
		f, err := os.Create(path)
		if err != nil {
			return "", err
		}

		if _, err := io.Copy(f, v); err != nil {
			return "", err
		}

		_ = f.Close()
	}

	return dir, nil
}

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

// GetEmbedded returns a map of embedded files.
func getEmbedded(fsys embed.FS, path string) Embedded {
	m := make(Embedded)
	files, _ := fsys.ReadDir(path)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		var r io.ReadCloser
		var err error
		r, err = fsys.Open(filepath.ToSlash(filepath.Join(path, file.Name())))
		if err != nil {
			logger.ErrLogger.Panicln(err)
		}

		defer r.Close()

		if filepath.Ext(file.Name()) == ".gz" {
			r, err = gzip.NewReader(r)
			if err != nil {
				logger.ErrLogger.Panicln(err)
			}

			defer r.Close()
		}

		raw, err := io.ReadAll(r)
		if err != nil {
			logger.ErrLogger.Panicln(err)
		}

		m[strings.TrimSuffix(file.Name(), ".gz")] = &multiReadReader{data: raw}
	}

	return m
}
