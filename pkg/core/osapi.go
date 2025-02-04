package core

import (
	"fmt"
	_ "unsafe"

	"github.com/sarumaj/go-wallpaper"
)

var AllowedModes = Modes{ModeCenter, ModeCrop, ModeFit, ModeSpan, ModeStretch, ModeTile}

const (
	ModeCenter  = Mode(wallpaper.Center)
	ModeCrop    = Mode(wallpaper.Crop)
	ModeFit     = Mode(wallpaper.Fit)
	ModeSpan    = Mode(wallpaper.Span)
	ModeStretch = Mode(wallpaper.Stretch)
	ModeTile    = Mode(wallpaper.Tile)
)

// Mode represents the wallpaper mode.
type Mode wallpaper.Mode

// Modes represents a list of wallpaper modes.
type Modes []Mode

// Contains returns true if the mode is in the list of modes.
func (ms Modes) Contains(m Mode) bool {
	for _, mode := range ms {
		if mode == m {
			return true
		}
	}
	return false
}

// String returns the string representation of the mode.
func (m Mode) String() string {
	s, ok := map[Mode]string{
		ModeCenter:  "center",
		ModeCrop:    "crop",
		ModeFit:     "fit",
		ModeSpan:    "span",
		ModeStretch: "stretch",
		ModeTile:    "tile",
	}[m]
	if !ok {
		return "Unknown"
	}
	return s
}

// GetWallpaper returns the path to the current wallpaper.
func GetWallpaper() (string, error) {
	return wallpaper.Get()
}

// SetWallpaper sets the wallpaper from the given path.
func SetWallpaper(path string, mode Mode) error {
	if !AllowedModes.Contains(mode) {
		return fmt.Errorf("invalid mode: %s", mode)
	}

	if err := wallpaper.SetMode(wallpaper.Mode(mode)); err != nil {
		return err
	}

	return wallpaper.SetFromFile(path)
}
