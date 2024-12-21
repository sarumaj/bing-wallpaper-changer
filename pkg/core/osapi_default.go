//go:build (darwin || linux) && !cgo

package core

// Run executes the given function with the given configuration.
func Run(execute func(*Config) *Image, cfg *Config) {
	// just execute
	_ = execute(cfg)
}
