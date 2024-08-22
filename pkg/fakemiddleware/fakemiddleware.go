package fakemiddleware

import (
	"fmt"
	"io"
)

type fakeMiddleware struct {
}

func NewFakeMiddleware() *fakeMiddleware {
	return &fakeMiddleware{}
}

type FakeReader struct {
	r     io.Reader
	total int
}

func (r *FakeReader) Read(p []byte) (int, error) {
	result, err := r.r.Read(p)
	r.total += result
	fmt.Printf("read %d\n", r.total)
	return result, err
}

type FakeWriter struct {
	r     io.Writer
	total int
}

func (r *FakeWriter) Write(p []byte) (n int, err error) {
	result, err := r.r.Write(p)
	r.total += result
	fmt.Printf("writtern %d\n", r.total)
	return result, err
}

func (f fakeMiddleware) Write(w io.Writer) (wp io.Writer, err error) {
	ww := &FakeWriter{w, 0}
	return ww, err
}

func (f fakeMiddleware) Read(r io.Reader) (rp io.Reader, err error) {
	ff := &FakeReader{r, 0}
	return ff, nil
}
