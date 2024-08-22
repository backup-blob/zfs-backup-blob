package crypt

import (
	"bytes"
	"crypto/cipher"
	"fmt"
	"io"
)

type unsealer struct {
	writer    io.Writer
	c         cipher.AEAD
	buffer    bytes.Buffer
	chunkSize int
}

func NewUnsealer(c cipher.AEAD, w io.Writer, chunkSize int) io.Writer {
	return &unsealer{
		c:         c,
		writer:    w,
		chunkSize: chunkSize,
	}
}

func (u *unsealer) Write(p []byte) (int, error) {
	// fill buffer
	bytesWritten, _ := u.buffer.Write(p)
	eofLength := len(EOFMarker)

	// drain buffer + decrypt + write
	for u.buffer.Len() > 0 {
		var errRC error

		isLast := u.buffer.Len() >= eofLength && string(u.buffer.Bytes()[u.buffer.Len()-eofLength:u.buffer.Len()]) == EOFMarker

		if u.buffer.Len() >= u.chunkSize { //nolint:gocritic // if is ok
			errRC = u.readChunk(u.chunkSize)
		} else if isLast {
			errRC = u.readChunk(u.buffer.Len() - eofLength)
			u.buffer.Reset()
		} else {
			break
		}

		if errRC != nil {
			return 0, errRC
		}
	}

	return bytesWritten, nil
}

func (u *unsealer) readChunk(lenToRead int) error {
	// read nonce.
	nonce := make([]byte, u.c.NonceSize())
	_, errRN := u.buffer.Read(nonce)

	if errRN != nil {
		return fmt.Errorf("nonce read failed %w", errRN)
	}

	// read data.
	sliceToRead := make([]byte, lenToRead-u.c.NonceSize())

	_, errRB := u.buffer.Read(sliceToRead)
	if errRB != nil {
		return fmt.Errorf("buffer read failed %w", errRB)
	}

	// decrypt.
	decryptedBuf, errO := u.c.Open(nil, nonce, sliceToRead, nil)
	if errO != nil {
		return errO
	}

	_, errCW := io.Copy(u.writer, bytes.NewReader(decryptedBuf))
	if errCW != nil {
		return errCW
	}

	return nil
}
