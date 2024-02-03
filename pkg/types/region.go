package types

import (
	"fmt"
	"strings"
)

const (
	Canada Region = iota
	China
	Germany
	Japan
	NewZealand
	UnitedKingdom
	UnitedStates
)

var AllowedRegions = Regions{Canada, China, Germany, Japan, NewZealand, UnitedKingdom, UnitedStates}

// Region is an enum type for locales.
type Region int

// String returns the code representation of the Region.
func (r Region) String() string {
	switch r {
	case Canada:
		return "en-CA"

	case China:
		return "zh-CN"

	case Germany:
		return "de-DE"

	case Japan:
		return "ja-JP"

	case NewZealand:
		return "en-NZ"

	case UnitedKingdom:
		return "en-GB"

	case UnitedStates:
		return "en-US"

	default:
		return "Unknown"

	}
}

// Regions is a slice of Region.
type Regions []Region

// String returns the string representation of the Regions.
func (r Regions) String() string {
	var s []string
	for _, v := range r {
		s = append(s, v.String())
	}

	return strings.Join(s, ", ")
}

// ParseLocale parses a string and returns a Region.
func ParseLocale(s string) (Region, error) {
	for _, allowed := range AllowedRegions {
		if s == allowed.String() {
			return allowed, nil
		}
	}

	return -1, fmt.Errorf("unsupported region: %s, expected any of: %s", s, AllowedRegions)
}
