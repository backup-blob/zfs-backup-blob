package driver

import (
	"fmt"
	"io"

	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
	"github.com/klauspost/compress/zstd"
)

// speed thresholds for compression levels.
// These map user-friendly levels to zstd encoder speeds.
const (
	fastestThreshold = 2
	defaultThreshold = 5
	betterThreshold  = 9
)

// CompressDriver implements compression middleware using zstd.
type CompressDriver struct {
	Conf *config.CompressConfig
}

// NewCompress creates a new compression middleware.
func NewCompress(conf *config.CompressConfig) domain.Middleware {
	return &CompressDriver{Conf: conf}
}

// encoderLevel maps user level (1-22) to zstd encoder levels.
func (cd *CompressDriver) encoderLevel() zstd.EncoderLevel {
	switch {
	case cd.Conf.Level <= fastestThreshold:
		return zstd.SpeedFastest
	case cd.Conf.Level <= defaultThreshold:
		return zstd.SpeedDefault
	case cd.Conf.Level <= betterThreshold:
		return zstd.SpeedBetterCompression
	default:
		return zstd.SpeedBestCompression
	}
}

// Read wraps the reader with compression (used during backup).
// When backing up: zfs send output -> Read(compress) -> upload.
func (cd *CompressDriver) Read(r io.Reader) (rp io.Reader, err error) {
	pr, pw := io.Pipe()

	go func() {
		encoder, err := zstd.NewWriter(pw, zstd.WithEncoderLevel(cd.encoderLevel()))
		if err != nil {
			pw.CloseWithError(fmt.Errorf("failed to create zstd encoder: %w", err))
			return
		}

		_, copyErr := io.Copy(encoder, r)
		closeErr := encoder.Close()

		switch {
		case copyErr != nil:
			pw.CloseWithError(copyErr)
		case closeErr != nil:
			pw.CloseWithError(closeErr)
		default:
			pw.Close()
		}
	}()

	return pr, nil
}

// decompressWriter wraps the pipe writer and waits for decompression to complete.
type decompressWriter struct {
	pw     *io.PipeWriter
	done   chan error
	closed bool
}

// Write implements io.Writer.
func (d *decompressWriter) Write(p []byte) (n int, err error) {
	return d.pw.Write(p)
}

// Close implements io.Closer and waits for decompression to finish.
func (d *decompressWriter) Close() error {
	if d.closed {
		return nil
	}

	d.closed = true

	err := d.pw.Close()
	if err != nil {
		return err
	}

	// Wait for decompression to complete.
	return <-d.done
}

// Write wraps the writer with decompression (used during restore).
// When restoring: download -> Write(decompress) -> zfs receive.
func (cd *CompressDriver) Write(writer io.Writer) (wp io.Writer, err error) {
	pr, pw := io.Pipe()
	done := make(chan error, 1)

	go func() {
		decoder, err := zstd.NewReader(pr)
		if err != nil {
			pr.CloseWithError(fmt.Errorf("failed to create zstd decoder: %w", err))
			done <- err

			return
		}

		defer decoder.Close()

		_, copyErr := io.Copy(writer, decoder)
		pr.CloseWithError(copyErr)
		done <- copyErr
	}()

	return &decompressWriter{pw: pw, done: done}, nil
}
