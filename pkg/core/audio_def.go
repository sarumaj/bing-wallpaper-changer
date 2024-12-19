//go:build darwin || (linux && !cgo)

package core

import (
	"fmt"
	"io"
	"runtime"
)

// thread-safe audio context.
type audioContext struct{}

func (a *audioContext) Initialize(sampleRate, channels, bitDepth, bufferSize int) error {
	return fmt.Errorf("audio context not supported on %s-%s", runtime.GOOS, runtime.GOARCH)
}

func (a *audioContext) NewPlayer() io.WriteCloser {
	return nil
}

func init() {
	audioCtx = &audioContext{}
}
