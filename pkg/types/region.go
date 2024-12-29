package types

import (
	"fmt"
	"strings"
)

var (
	RegionBrazil        = Region{"BR", "pt"}
	RegionCanadaEnglish = Region{"CA", "en"}
	RegionCanadaFrench  = Region{"CA", "fr"}
	RegionChina         = Region{"CN", "zh"}
	RegionFrance        = Region{"FR", "fr"}
	RegionGermany       = Region{"DE", "de"}
	RegionItaly         = Region{"IT", "it"}
	RegionIndia         = Region{"IN", "hi"}
	RegionJapan         = Region{"JP", "ja"}
	RegionNewZealand    = Region{"NZ", "en"}
	RegionSpain         = Region{"ES", "es"}
	RegionOther         = Region{"ROW", "en"}
	RegionUnitedKingdom = Region{"GB", "en"}
	RegionUnitedStates  = Region{"US", "en"}
)

var AllowedRegions = Regions{
	RegionBrazil,
	RegionCanadaEnglish,
	RegionCanadaFrench,
	RegionChina,
	RegionFrance,
	RegionGermany,
	RegionItaly,
	RegionIndia,
	RegionJapan,
	RegionNewZealand,
	RegionSpain,
	RegionOther,
	RegionUnitedKingdom,
	RegionUnitedStates,
}

var EnglishRegions = AllowedRegions.Filter(func(r Region) bool { return r.LanguageCode == "en" })
var NonEnglishRegions = AllowedRegions.Filter(func(r Region) bool { return r.LanguageCode != "en" })

// Region is an enum type for locales.
type Region struct {
	Country      string
	LanguageCode string
}

// IsAny returns true if the Region is any of the provided Regions.
func (r Region) IsAny(o ...Region) bool {
	for _, v := range o {
		if r == v {
			return true
		}
	}

	return false
}

// String returns the code representation of the Region.
func (r Region) String() string {
	return r.LanguageCode + "-" + r.Country
}

// Regions is a slice of Region.
type Regions []Region

// Filter returns a new Regions with the Regions that satisfy the predicate.
func (r Regions) Filter(fn func(Region) bool) Regions {
	var s Regions
	for _, v := range r {
		if fn(v) {
			s = append(s, v)
		}
	}

	return s
}

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

	return Region{}, fmt.Errorf("unsupported region: %s, expected any of: %s", s, AllowedRegions)
}
