package wallpaper

import (
	"os/user"
	"path/filepath"

	ini "gopkg.in/ini.v1"
)

func (mode Mode) getLXDEString() string {
	str, ok := map[Mode]string{
		Center:  "center",
		Crop:    "crop",
		Fit:     "fit",
		Span:    "screen",
		Stretch: "stretch",
		Tile:    "tile",
	}[mode]
	if !ok {
		panic("invalid wallpaper mode")
	}
	return str
}

func getLXDE() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	if DesktopSession == "" {
		DesktopSession = "LXDE"
	}

	cfg, err := ini.Load(filepath.Join(usr.HomeDir, ".config/pcmanfm/"+DesktopSession+"/desktop-items-0.conf"))
	if err != nil {
		return "", err
	}

	key, err := cfg.Section("*").GetKey("wallpaper")
	if err != nil {
		return "", err
	}

	return key.String(), err
}
