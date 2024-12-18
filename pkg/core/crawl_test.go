package core

import (
	"testing"

	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
)

func Test_furiganizeGooLabsApi(t *testing.T) {
	if FromMock(t) {
		MockServers(t)
	}

	if cfg.furiganaApiAppId == "" {
		t.Skip("no Goo Labs API AppId provided")
	}

	for _, tt := range []struct {
		name string
		args string
		want string
	}{
		{"test#1", "今日はダーウィンの日, ガラパゴスゾウガメ", "今日[きょう]はダーウィンの日[ひ], ガラパゴスゾウガメ"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := furiganizeGooLabsApi(tt.args)
			if err != nil {
				t.Errorf("furiganizeGooLabsApi() error = %v, wantErr %v", err, false)
				return
			}

			if got != tt.want {
				t.Errorf("furiganizeGooLabsApi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDownloadAndDecode(t *testing.T) {
	if FromMock(t) {
		MockServers(t)
	}

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
