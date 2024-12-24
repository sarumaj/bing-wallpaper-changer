package core

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"
)

// Image is a wrapper around the image.Image interface.
type Image struct {
	image.Image
	Audio       *Audio
	Description string
	SearchURL   string
	DownloadURL string
	Location    string
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
// If audio description is available, it will be dumped as well.
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

	defer target.Close()

	if img.Audio != nil {
		audioPath := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + "." + strings.ToLower(img.Audio.Encoding)
		if err := img.Audio.Dump(audioPath); err != nil {
			return "", err
		}
	}

	img.Location = filePath
	return target.Name(), png.Encode(target, img)
}

// getDecoder returns the decoder for the given file path.
func getDecoder(path string) (decoder func(io.Reader) (image.Image, error), err error) {
	switch ext := filepath.Ext(path); ext {
	case ".jpg", ".jpeg":
		decoder = jpeg.Decode

	case ".png":
		decoder = png.Decode

	case ".webp":
		decoder = webp.Decode

	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)

	}

	return decoder, nil
}
