package core

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"

	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
	"github.com/tidwall/gjson"
	"golang.org/x/image/webp"
)

// bindURL is the URL to fetch the Bing wallpaper.
const bingURL = "https://www.bing.com/HPImageArchive.aspx?format=js&idx=%d&n=1&mkt=%s"

// DownloadAndDecode fetches the Bing wallpaper and decodes it.
func DownloadAndDecode(day types.Day, region types.Region, resolution types.Resolution) (*Image, error) {
	jsonRaw, err := fetch(fmt.Sprintf(bingURL, day, region))
	if err != nil {
		return nil, err
	}

	fmt.Println(string(jsonRaw))

	path := gjson.GetBytes(jsonRaw, "images.0.url").String()
	path = regexp.MustCompile(`_\d+x\d+`).ReplaceAllString(path, "_"+resolution.String())
	uri := "https://www.bing.com" + path

	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	decoder, err := getDecoder(parsed.Query().Get("id"))
	if err != nil {
		return nil, err
	}

	content, err := fetch(parsed.String())
	if err != nil {
		return nil, err
	}

	img, err := decoder(bytes.NewReader(content))

	imgBounds := img.Bounds()
	if imgBounds.Dx() != resolution.Width || imgBounds.Dy() != resolution.Height {
		return nil, fmt.Errorf("expected resolution: %s, got: %s", resolution, imgBounds.Size())
	}

	return &Image{
		Description: fmt.Sprintf(
			"%q, %s",
			gjson.GetBytes(jsonRaw, "images.0.title").String(),
			gjson.GetBytes(jsonRaw, "images.0.copyright").String(),
		),
		Image:       img,
		DownloadURL: uri,
		SearchURL:   gjson.GetBytes(jsonRaw, "images.0.copyrightlink").String(),
	}, err
}

// fetch fetches the content from the given URI.
func fetch(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
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
