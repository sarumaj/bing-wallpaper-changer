package core

import (
	"testing"

	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
)

func TestDownloadAndDecode(t *testing.T) {
	SkipOAT(t)

	type args struct {
		day        types.Day
		region     types.Region
		resolution types.Resolution
	}

	for _, tt := range []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test#1", args{types.Today, types.Germany, types.HighDefinition}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DownloadAndDecode(tt.args.day, tt.args.region, tt.args.resolution)
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadAndDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Fetched wallpaper: %#v", got)
		})
	}
}
