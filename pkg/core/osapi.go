package core

import (
	"github.com/reujab/wallpaper"
)

// GetWallpaper returns the path to the current wallpaper.
func GetWallpaper() (string, error) {
	return wallpaper.Get()
}

// SetWallpaper sets the wallpaper from the given path.
func SetWallpaper(path string) error {
	wallpaper.SetMode(wallpaper.Fit)
	return wallpaper.SetFromFile(path)
}
