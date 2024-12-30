//go:build linux
// +build linux

package wallpaper

import (
	"bufio"
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

func (mode Mode) getKDEString() string {
	str, ok := map[Mode]string{
		Center:  "6",
		Crop:    "2",
		Fit:     "1",
		Span:    "2",
		Stretch: "0",
		Tile:    "3",
	}[mode]
	if !ok {
		panic("invalid wallpaper mode")
	}
	return str
}

func evalKDE(script string) error {
	return execCmd("qdbus", "org.kde.plasmashell", "/PlasmaShell", "org.kde.PlasmaShell.evaluateScript", script)
}

func getKDE() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	filename := filepath.Join(usr.HomeDir, ".config", "plasma-org.kde.plasma.desktop-appletsrc")
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := scanner.Text(); strings.HasPrefix(line, "Image=") {
			return strings.TrimSpace(removeProtocol(strings.TrimPrefix(line, "Image="))), nil
		}
	}
	if scanner.Err() != nil {
		return "", scanner.Err()
	}

	return "", errors.New("kde image not found")
}

func setKDE(path string) error {
	return evalKDE(`
		for (const desktop of desktops()) {
			desktop.currentConfigGroup = ["Wallpaper", "org.kde.image", "General"]
			desktop.writeConfig("Image", ` + strconv.Quote("file://"+path) + `)
		}
	`)
}

func setKDEMode(mode Mode) error {
	return evalKDE(`
		for (const desktop of desktops()) {
			desktop.currentConfigGroup = ["Wallpaper", "org.kde.image", "General"]
			desktop.writeConfig("FillMode", ` + mode.getKDEString() + `)
		}
	`)
}
