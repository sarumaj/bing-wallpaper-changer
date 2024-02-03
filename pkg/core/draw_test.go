package core

import (
	"testing"

	"github.com/sarumaj/bing-wallpaper-changer/pkg/extras"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
)

func TestDrawDescription(t *testing.T) {
	img := SetupTestImage()

	type args struct {
		fontName string
		position types.Position
	}

	for _, tt := range []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test#1", args{extras.DefaultFontName, types.TopCenter}, false},
		{"test#2", args{extras.DefaultFontName, types.BottomCenter}, false},
		{"test#3", args{extras.DefaultFontName, types.Position(-1)}, true},
		{"test#4", args{"unknown", types.TopCenter}, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := SetupTestImage()

			err := got.DrawDescription(tt.args.position, tt.args.fontName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DrawDescription(%q, %q) error = %v, wantErr %t", tt.args.position, tt.args.fontName, err, tt.wantErr)
				return
			}

			if tt.wantErr != got.Equals(img) {
				t.Errorf("DrawDescription(%q, %q) = %v, want %v", tt.args.position, tt.args.fontName, got, img)
			}
		})
	}
}

func TestDrawQRCode(t *testing.T) {
	img := SetupTestImage()

	type args struct {
		resolution types.Resolution
		position   types.Position
	}

	for _, tt := range []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test#1", args{types.HighDefinition, types.TopLeft}, false},
		{"test#2", args{types.HighDefinition, types.BottomRight}, false},
		{"test#3", args{types.HighDefinition, types.BottomLeft}, false},
		{"test#4", args{types.HighDefinition, types.TopRight}, false},
		{"test#5", args{types.HighDefinition, types.Position(-1)}, true},
		{"test#6", args{types.Resolution{}, types.TopLeft}, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := SetupTestImage()

			err := got.DrawQRCode(tt.args.resolution, tt.args.position)
			if (err != nil) != tt.wantErr {
				t.Errorf("DrawQRCode(%q, %q) error = %v, wantErr %t", tt.args.resolution, tt.args.position, err, tt.wantErr)
				return
			}

			if tt.wantErr != got.Equals(img) {
				t.Errorf("DrawQRCode(%q, %q) = %v, want %v", tt.args.resolution, tt.args.position, got, img)
			}
		})
	}

}

func TestDrawWatermark(t *testing.T) {
	img := SetupTestImage()

	type args struct {
		watermarkFile          string
		rotateCounterClockwise bool
	}

	for _, tt := range []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test#1", args{extras.DefaultWatermarkName, false}, false},
		{"test#2", args{extras.DefaultWatermarkName, true}, false},
		{"test#3", args{"unknown", false}, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := SetupTestImage()

			err := got.DrawWatermark(tt.args.watermarkFile, tt.args.rotateCounterClockwise)
			if (err != nil) != tt.wantErr {
				t.Errorf("DrawWatermark(%q, %t) error = %v, wantErr %t", tt.args.watermarkFile, tt.args.rotateCounterClockwise, err, tt.wantErr)
				return
			}

			if tt.wantErr != got.Equals(img) {
				t.Errorf("DrawWatermark(%q, %t) = %v, want %v", tt.args.watermarkFile, tt.args.rotateCounterClockwise, got, img)
			}
		})
	}
}
