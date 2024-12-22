//go:build linux && !cgo

package core

import (
	"fmt"
	"runtime"
)

// Play plays the audio.
func (a *Audio) Play() error {
	return fmt.Errorf("unsupported platform: %s-%s", runtime.GOOS, runtime.GOARCH)
}
