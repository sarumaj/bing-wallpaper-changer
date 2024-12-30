package wallpaper

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type Mode int

const (
	Center Mode = iota
	Crop
	Fit
	Span
	Stretch
	Tile
)

// Desktop contains the current desktop environment on Linux.
// Empty string on all other operating systems.
var Desktop = os.Getenv("XDG_CURRENT_DESKTOP")

// DesktopSession is used by LXDE on Linux.
var DesktopSession = os.Getenv("DESKTOP_SESSION")

// ErrUnsupportedDE is thrown when Desktop is not a supported desktop environment.
var ErrUnsupportedDE = errors.New("your desktop environment is not supported")

func downloadImage(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("failed to download image: %s", res.Status)
	}

	cacheDir, err := getCacheDir()
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(filepath.Join(cacheDir, "wallpaper"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func execCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SetFromURL downloads the image to a cache directory and calls SetFromFile.
func SetFromURL(url string, desktop ...int) error {
	file, err := downloadImage(url)
	if err != nil {
		return err
	}

	return SetFromFile(file, desktop...)
}
