package main

import (
	"fmt"
	"testing"

	"github.com/sarumaj/bing-wallpaper-changer/pkg/core"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/extras"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
)

func TestFlow(t *testing.T) {
	tempDir := t.TempDir()

	if core.FromMock(t) {
		core.MockServers(t)
	}

	type args struct {
		resolution     types.Resolution
		region         types.Region
		day            types.Day
		qrcodePosition types.Position
		titlePosition  types.Position
	}

	for _, tt := range []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test#1", args{types.HighDefinition, types.RegionGermany, types.DayToday, types.PositionBottomRight, types.PositionTopCenter}, false},
		{"test#2", args{types.HighDefinition, types.RegionGermany, types.DayToday, types.PositionBottomRight, types.PositionBottomLeft}, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := func(t testing.TB) error {
				t.Helper()

				img, err := core.DownloadAndDecode(tt.args.day, tt.args.region, tt.args.resolution)
				if err != nil {
					return fmt.Errorf("DownloadAndDecode() failed: %w", err)
				}

				t.Logf("Fetched wallpaper: %#v", img)

				if err := img.DrawWatermark(extras.DefaultWatermarkName, false); err != nil {
					return fmt.Errorf("DrawWatermark() failed: %w", err)
				}

				t.Logf("Watermark drawn: %#v", img)

				if err := img.DrawDescription(tt.args.titlePosition, extras.DefaultFontName); err != nil {
					return fmt.Errorf("DrawDescription() failed: %w", err)
				}

				t.Logf("Description drawn: %#v", img)

				if err := img.DrawQRCode(tt.args.resolution, tt.args.qrcodePosition); err != nil {
					return fmt.Errorf("DrawQRCode() failed: %w", err)
				}

				t.Logf("QR code drawn: %#v", img)

				path, err := img.EncodeAndDump(tempDir)
				if err != nil {
					return fmt.Errorf("EncodeAndDump() failed: %w", err)
				}

				t.Logf("Wallpaper saved to: %s", path)
				return nil
			}(t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Flow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
