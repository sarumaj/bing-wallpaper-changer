package core

import (
	"github.com/rawnly/go-wallpaper"
)

// GetWallpaper returns the path to the current wallpaper.
func GetWallpaper() (string, error) {
	return wallpaper.Get()
}

// SetWallpaper sets the wallpaper from the given path.
func SetWallpaper(path string) error {
	if err := wallpaper.SetMode(wallpaper.Fit); err != nil {
		return err
	}

	return wallpaper.SetFromFile(path)
}
