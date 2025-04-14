//go:build (darwin || linux) && !cgo

package core

import (
	"net/http"

	"github.com/sarumaj/bing-wallpaper-changer/pkg/logger"
)

// Run executes the given function with the given configuration.
func Run(execute func(*Config) *Image, cfg *Config) {
	img := &Image{}
	if !cfg.Daemon {
		_ = execute(cfg)
		return
	}

	server := NewServer(cfg, img, execute, nil)
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Fatalf("Failed to start API server: %v", err)
	}
}
