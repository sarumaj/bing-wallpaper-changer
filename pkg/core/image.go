package core

import (
	"fmt"
	"image"
	"image/png"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Image is a wrapper around the image.Image interface.
type Image struct {
	image.Image
	Description string
	SearchURL   string
	DownloadURL string
}

// Equals returns true if the given image is equal to the receiver.
func (i *Image) Equals(other *Image) bool {
	if i.Description != other.Description || i.SearchURL != other.SearchURL || i.DownloadURL != other.DownloadURL {
		return false
	}
	selfBounds, otherBounds := i.Bounds(), other.Bounds()
	if selfBounds.Dx() != otherBounds.Dx() || selfBounds.Dy() != otherBounds.Dy() {
		return false
	}

	for x := 0; x < selfBounds.Dx(); x++ {
		for y := 0; y < selfBounds.Dy(); y++ {
			if i.At(x, y) != other.At(x, y) {
				return false
			}
		}
	}

	return true
}

// EncodeAndDump encodes the image and dumps it to the target directory.
func (img *Image) EncodeAndDump(targetDir string) (string, error) {
	parsed, err := url.Parse(img.DownloadURL)
	if err != nil {
		return "", err
	}

	fileName := parsed.Query().Get("id")
	if fileName == "" {
		return "", fmt.Errorf("missing file name in URL: %s", img.DownloadURL)
	}

	_ = os.MkdirAll(targetDir, os.ModePerm)
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".png"
	filePath := filepath.Join(targetDir, fileName)
	target, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return "", err
	}

	return target.Name(), png.Encode(target, img)
}
