package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

var audioCtx audioContext

type (
	// thread-safe audio context.
	audioContext struct {
		context *oto.Context
		mux     sync.Mutex
	}

	Audio struct {
		Encoding   string
		Source     io.Reader
		SampleRate int
	}
)

func (a *audioContext) Initialize(sampleRate, channels, bitDepth, bufferSize int) error {
	a.mux.Lock()
	defer a.mux.Unlock()

	if audioCtx.context == nil {
		ctx, err := oto.NewContext(sampleRate, channels, bitDepth, bufferSize)
		if err != nil {
			return err
		}

		audioCtx.context = ctx
	}

	return nil
}

func (a *audioContext) NewPlayer() io.WriteCloser {
	a.mux.Lock()
	defer a.mux.Unlock()

	return a.context.NewPlayer()
}

// Implements the io.Closer interface.
func (a *Audio) Close() error { return nil }

// Dump dumps the audio to the target path.
func (a *Audio) Dump(path string) error {
	if a == nil {
		return nil
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()

	var buffer bytes.Buffer
	if _, err := io.ReadAll(io.TeeReader(a.Source, io.MultiWriter(&buffer, f))); err != nil {
		return err
	}

	a.Source = &buffer
	return nil
}

// Play the audio stream.
func (a *Audio) Play() error {
	if a == nil {
		return nil
	}

	if err := audioCtx.Initialize(a.SampleRate, 2, 2, 32); err != nil {
		return err
	}

	var buffer bytes.Buffer
	player := audioCtx.NewPlayer()
	defer player.Close()
	switch a.Encoding {
	case texttospeechpb.AudioEncoding_LINEAR16.String():
		if _, err := io.Copy(player, io.TeeReader(a.Source, &buffer)); err != nil {
			return err
		}

	case texttospeechpb.AudioEncoding_MP3.String():
		decoder, err := mp3.NewDecoder(io.TeeReader(a.Source, &buffer))
		if err != nil {
			return err
		}

		if _, err := io.Copy(player, decoder); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported audio encoding: %s", a.Encoding)

	}

	a.Source = &buffer
	return nil
}