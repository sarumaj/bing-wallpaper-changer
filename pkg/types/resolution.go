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
		s = append(s, v.String())
	}
	return strings.Join(s, ", ")
}

// ParseResolution parses a string and returns a Resolution.
func ParseResolution(s string) (Resolution, error) {
	defaultErr := fmt.Errorf("invalid resolution: %s, allowed values arr: %s", s, AllowedResolutions)

	var r Resolution
	_, err := fmt.Sscanf(s, "%dx%d", &r.Width, &r.Height)
	for _, allowed := range AllowedResolutions {
		switch {
		case strings.EqualFold(s, allowed.Alias), r.Width == allowed.Width && r.Height == allowed.Height:
			return Resolution{
				Width:  allowed.Width,
				Height: allowed.Height,
				Alias:  allowed.Alias,
			}, nil
		}
	}

	if err != nil {
		return r, defaultErr
	}

	return Resolution{}, defaultErr
}
