package types

import "testing"

func TestParseLocale(t *testing.T) {
	for _, tt := range []struct {
		name    string
		locale  string
		want    Region
		wantErr bool
	}{
		{"test#1", "en-US", UnitedStates, false},
		{"test#2", "en-GB", UnitedKingdom, false},
		{"test#3", "de-DE", Germany, false},
		{"test#5", "ja-JP", Japan, false},
		{"test#6", "invalid", Region(-1), true},
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
