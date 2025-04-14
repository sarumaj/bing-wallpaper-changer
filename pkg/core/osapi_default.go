//go:build (darwin || linux) && !cgo

package core

import (
	"net/http"

	"github.com/sarumaj/bing-wallpaper-changer/pkg/logger"
)

// Controller is the controller of the application.
type Controller struct {
	cfg     *Config
	img     *Image
	execute func(*Config) *Image
}

// OnReady is called when the application is ready.
func (c *Controller) OnReady() {
	c.img.Update(c.execute(c.cfg))
}

// OnExit is called when the application is closed.
func (c *Controller) OnExit() {}

// Run executes the given function with the given configuration.
func Run(execute func(*Config) *Image, cfg *Config) {
	img := &Image{}
	if !cfg.Daemon {
		img.Update(execute(cfg))
		return
	}

	server := NewServer(cfg, &Controller{img: img, cfg: cfg, execute: execute})
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Fatalf("Failed to start API server: %v", err)
	}
}
