package types

import (
	"fmt"
	"strings"
)

var (
	LowDefinition       = Resolution{1366, 768, "SD"}
	HighDefinition      = Resolution{1920, 1080, "HD"}
	UltraHighDefinition = Resolution{3840, 2160, "UHD"}
)

var AllowedResolutions = Resolutions{LowDefinition, HighDefinition, UltraHighDefinition}

// Resolution is a struct for screen resolutions.
type Resolution struct {
	Width, Height int
	Alias         string
}

// BingFormat returns the Bing format of the Resolution.
func (r Resolution) BingFormat() string {
	switch r {
	case UltraHighDefinition:
		return r.Alias

	default:
		return fmt.Sprintf("%dx%d", r.Width, r.Height)

	}
}

// String returns the string representation of the Resolution.
func (r Resolution) String() string {
	return fmt.Sprintf("%dx%d", r.Width, r.Height)
}

// Resolutions is a slice of Resolution.
type Resolutions []Resolution

// String returns the string representation of the Resolutions.
func (r Resolutions) String() string {
	var s []string
	for _, v := range r {
		s = append(s, v.String()+" ("+v.Alias+")")
	}
	return strings.Join(s, ", ")
}
