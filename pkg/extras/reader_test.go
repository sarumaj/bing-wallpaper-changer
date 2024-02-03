package extras

import (
	"io"
	"testing"
)

func Test_multiReadReader(t *testing.T) {
	reader := &multiReadReader{data: []byte("test")}

	for i := 0; i < 2; i++ {
		got, err := io.ReadAll(reader)
		if err != nil {
			t.Errorf("ReadAll() error = %v, wantErr %v", err, false)
		}

		if string(got) != "test" {
			t.Errorf("ReadAll() = %v, want %v", got, "test")
		}
	}
}
