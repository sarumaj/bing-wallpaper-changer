package extras

import "io"

// Assert that multiReadReader implements the io.ReadCloser interface.
var _ io.ReadCloser = &multiReadReader{}

// multiReadReader is an io.Reader that allows for reading the same data multiple times.
type multiReadReader struct {
	data   []byte
	offset int
}

// Read reads the data from the multiReadReader.
func (r *multiReadReader) Read(p []byte) (int, error) {
	if r.offset >= len(r.data) {
		r.offset = 0     // reset offset
		return 0, io.EOF // end of file reached
	}

	n := copy(p, r.data[r.offset:])
	r.offset += n
	return n, nil
}

// Close closes the multiReadReader.
func (r *multiReadReader) Close() error {
	return nil
}
