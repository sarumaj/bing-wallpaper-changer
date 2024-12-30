package wallpaper

import (
	"os/exec"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

func (mode Mode) getGNOMEString() string {
	str, ok := map[Mode]string{
		Center:  "centered",
		Crop:    "zoom",
		Fit:     "scaled",
		Span:    "spanned",
		Stretch: "stretched",
		Tile:    "wallpaper",
	}[mode]
	if !ok {
		panic("invalid wallpaper mode")
	}
	return str
}

func getGNOME() (string, error) {
	style, err := parseDconf("dconf", "read", "/org/gnome/desktop/interface/color-scheme")
	if err != nil {
		return "", err
	}

	if style == "prefer-dark" {
		return parseDconf("dconf", "read", "/org/gnome/desktop/background/picture-uri-dark")
	}

	return parseDconf("dconf", "read", "/org/gnome/desktop/background/picture-uri")
}

func isGNOMECompliant() bool {
	return strings.Contains(Desktop, "GNOME") || Desktop == "Unity" || Desktop == "Pantheon"
}

func parseDconf(command string, args ...string) (string, error) {
	output, err := exec.Command(command, args...).Output()
	if err != nil {
		return "", err
	}

	// unquote string
	var unquoted string
	// the output is quoted with single quotes, which cannot be unquoted using strconv.Unquote, but it is valid yaml
	err = yaml.Unmarshal(output, &unquoted)
	if err != nil {
		return unquoted, err
	}

	return removeProtocol(unquoted), nil
}

func removeProtocol(input string) string {
	return strings.TrimPrefix(input, "file://")
}

func setGNOME(path string) error {
	if err := execCmd("dconf", "write", "/org/gnome/desktop/background/picture-uri", strconv.Quote("file://"+path)); err != nil {
		return err
	}

	return execCmd("dconf", "write", "/org/gnome/desktop/background/picture-uri-dark", strconv.Quote("file://"+path))
}

func setGNOMEMode(mode Mode) error {
	return execCmd("dconf", "write", "/org/gnome/desktop/background/picture-options", strconv.Quote(mode.getGNOMEString()))
}
