package driver

import (
	"io"
)

type FakeReader struct {
	r io.Reader
}

func (r *FakeReader) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}
