package wallpaper

import (
	"os/exec"
	"path"
	"strings"
)

func (mode Mode) getXFCEString() string {
	str, ok := map[Mode]string{
		Center:  "1",
		Crop:    "5",
		Fit:     "4",
		Span:    "5",
		Stretch: "3",
		Tile:    "2",
	}[mode]
	if !ok {
		panic("invalid wallpaper mode")
	}

	return str
}

func getXFCE() (string, error) {
	desktops, err := getXFCEProps("last-image")
	if err != nil || len(desktops) == 0 {
		return "", err
	}

	output, err := exec.Command("xfconf-query", "--channel", "xfce4-desktop", "--property", desktops[0]).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func getXFCEProps(key string) ([]string, error) {
	output, err := exec.Command("xfconf-query", "--channel", "xfce4-desktop", "--list").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.Trim(string(output), "\n"), "\n")
	var desktops []string

	for _, line := range lines {
		if path.Base(line) == key {
			desktops = append(desktops, line)
		}
	}

	return desktops, nil
}

func setXFCE(file string) error {
	desktops, err := getXFCEProps("last-image")
	if err != nil {
		return err
	}

	for _, desktop := range desktops {
		if err := execCmd("xfconf-query", "--channel", "xfce4-desktop", "--property", desktop, "--set", file); err != nil {
			return err
		}
	}

	return nil
}

func setXFCEMode(mode Mode) error {
	styles, err := getXFCEProps("image-style")
	if err != nil {
		return err
	}

	for _, style := range styles {
		if err := execCmd("xfconf-query", "--channel", "xfce4-desktop", "--property", style, "--set", mode.getXFCEString()); err != nil {
			return err
		}
	}

	return nil
}
