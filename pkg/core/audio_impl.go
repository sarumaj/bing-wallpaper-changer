//go:build !(darwin || (linux && arm64) || (linux && arm))

package core

import (
	"io"
	"sync"

	"github.com/hajimehoshi/oto"
)

// thread-safe audio context.
type audioContext struct {
	context *oto.Context
	mux     sync.Mutex
}

func (a *audioContext) Initialize(sampleRate, channels, bitDepth, bufferSize int) error {
	a.mux.Lock()
	defer a.mux.Unlock()

	if a.context == nil {
		ctx, err := oto.NewContext(sampleRate, channels, bitDepth, bufferSize)
		if err != nil {
			return err
		}

		a.context = ctx
	}

	return nil
}

func (a *audioContext) NewPlayer() io.WriteCloser {
	a.mux.Lock()
	defer a.mux.Unlock()

	return a.context.NewPlayer()
}

func init() {
	audioCtx = &audioContext{}
}
