package types

import "testing"

func TestParseLocale(t *testing.T) {
	for _, tt := range []struct {
		name    string
		locale  string
		want    Region
		wantErr bool
	}{
		{"test#1", "en-US", RegionUnitedStates, false},
		{"test#2", "en-GB", RegionUnitedKingdom, false},
		{"test#3", "de-DE", RegionGermany, false},
		{"test#5", "ja-JP", RegionJapan, false},
		{"test#6", "invalid", Region{}, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLocale(tt.locale)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLocale() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ParseLocale() = %v, want %v", got, tt.want)
			}
		})
	}
}
