package core

import (
	"bytes"
	"io"
	"os"
)

type Audio struct {
	Encoding   string
	Source     io.Reader
	SampleRate int32
}

// Implements the io.Closer interface.
func (a *Audio) Close() error { return nil }

// Dump dumps the audio to the target path.
func (a *Audio) Dump(path string) error {
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
