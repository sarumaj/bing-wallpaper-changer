package core

import (
	_ "unsafe"

	"github.com/sarumaj/go-wallpaper"
)

// Mode represents the wallpaper setting mode.
type Mode = wallpaper.Mode

const (
	ModeCenter  = wallpaper.Center
	ModeCrop    = wallpaper.Crop
	ModeFit     = wallpaper.Fit
	ModeSpan    = wallpaper.Span
	ModeStretch = wallpaper.Stretch
	ModeTile    = wallpaper.Tile
)

// GetWallpaper returns the path to the current wallpaper.
func GetWallpaper() (string, error) {
	return wallpaper.Get()
}

// SetWallpaper sets the wallpaper from the given path.
func SetWallpaper(path string, mode Mode) error {
	if err := wallpaper.SetMode(mode); err != nil {
		return err
	}

	return wallpaper.SetFromFile(path)
}
