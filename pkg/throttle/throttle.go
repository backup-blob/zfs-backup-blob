package throttle

import (
	"github.com/fujiwara/shapeio"
	"io"
)

type WriterChain = func(writer io.Writer) (io.Writer, error)
type ReaderChain = func(writer io.Reader) (io.Reader, error)

// SpeedlimitWriter limits the writer to the maxSpeed of bytes per second.
func SpeedlimitWriter(maxSpeed int64) WriterChain {
	return func(writer io.Writer) (io.Writer, error) {
		throttledWriter := shapeio.NewWriter(writer)
		throttledWriter.SetRateLimit(float64(maxSpeed))

		return throttledWriter, nil
	}
}

// SpeedlimitReader limits the reader to the maxSpeed of bytes per second.
func SpeedlimitReader(maxSpeed int64) ReaderChain {
	return func(reader io.Reader) (io.Reader, error) {
		throttledReader := shapeio.NewReader(reader)
		throttledReader.SetRateLimit(float64(maxSpeed))

		return throttledReader, nil
	}
}
