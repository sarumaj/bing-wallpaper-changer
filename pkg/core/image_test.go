package core

import (
	"testing"
)

func TestEncodeAndDump(t *testing.T) {
	img := SetupTestImage(t)

	if _, err := img.EncodeAndDump(t.TempDir()); err != nil {
		t.Errorf("EncodeAndDump() error = %v, wantErr %v", err, false)
	}
}
