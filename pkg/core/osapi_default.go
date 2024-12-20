//go:build (darwin || linux) && !cgo

package core

// ShowTray shows the tray icon and menu
func ShowTray(execute func(*Config) *Image, cfg *Config) {
	// just execute
	_ = execute(cfg)
}
