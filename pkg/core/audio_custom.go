//go:build !linux || cgo

package core

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

var once sync.Once
var initialSampleRate atomic.Int32

// Play plays the audio.
func (a *Audio) Play() error {
	var err error
	once.Do(func() {
		initialSampleRate.Store(a.SampleRate)
		err = speaker.Init(beep.SampleRate(a.SampleRate), beep.SampleRate(a.SampleRate).N(time.Second/10))
	})
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	var stream beep.Streamer
	switch a.Encoding {
	case texttospeechpb.AudioEncoding_MP3.String():
		stream, _, err = mp3.Decode(io.NopCloser(io.TeeReader(a.Source, &buffer)))
		if err != nil {
			return err
		}
		defer stream.(beep.StreamSeekCloser).Close()

		if oldSampleRate := initialSampleRate.Load(); a.SampleRate != oldSampleRate {
			stream = beep.Resample(4, beep.SampleRate(oldSampleRate), beep.SampleRate(a.SampleRate), stream)
		}

	default:
		return fmt.Errorf("unsupported audio encoding: %s", a.Encoding)

	}

	speaker.PlayAndWait(stream)
	a.Source = &buffer
	return nil
}
