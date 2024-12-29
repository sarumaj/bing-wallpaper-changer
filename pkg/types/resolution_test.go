package types

import "testing"

func TestParseResolution(t *testing.T) {

	for _, tt := range []struct {
		name    string
		res     string
		want    Resolution
		wantErr bool
	}{
		{"test#1", "1920x1080", HighDefinition, false},
		{"test#2", "1920x1200", Resolution{}, true},
		{"test#3", "3840x2160", UltraHighDefinition, false},
		{"test#4", "3840x2400", Resolution{}, true},
		{"test#5", "1366x768", LowDefinition, false},
		{"test#6", "invalid", Resolution{}, true},
		{"test#7", "HD", HighDefinition, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseResolution(tt.res)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseResolution() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ParseResolution() = %v, want %v", got, tt.want)
			}
		})
	}
}
