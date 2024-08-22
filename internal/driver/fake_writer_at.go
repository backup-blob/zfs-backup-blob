package driver

import (
	"errors"
	"io"
	"sync"
)

type FakeWriterAt struct {
	w              io.Writer
	expectedOffset int64
	mutex          sync.Mutex
}

func (fw *FakeWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	if offset < fw.expectedOffset {
		return 0, errors.New("out of order write")
	}

	// Move the expected offset to the end of the current write.
	fw.expectedOffset = offset + int64(len(p))

	return fw.w.Write(p)
}
