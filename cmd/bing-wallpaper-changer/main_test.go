package main

import (
	"testing"

	"github.com/sarumaj/bing-wallpaper-changer/pkg/core"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/extras"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
)

func TestFlow(t *testing.T) {
	img, err := core.DownloadAndDecode(types.Yesterday, types.Germany, types.HighDefinition)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Fetched wallpaper: %#v", img)

	if err := img.DrawWatermark("red-dot.png"); err != nil {
		t.Error(err)
		return
	}

	t.Logf("Watermark drawn: %#v", img)

	if err := img.DrawDescription(types.BottomCenter, extras.DefaultFontName); err != nil {
		t.Error(err)
		return
	}

	t.Logf("Description drawn: %#v", img)

	if err := img.DrawQRCode(types.HighDefinition, types.TopRight); err != nil {
		t.Error(err)
		return
	}

	t.Logf("QR code drawn: %#v", img)

	path, err := img.EncodeAndDump(t.TempDir())
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Wallpaper saved to: %s", path)
}
