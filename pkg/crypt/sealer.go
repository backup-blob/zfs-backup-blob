package crypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

type sealer struct {
	c         cipher.AEAD
	reader    io.Reader
	encrypted bytes.Buffer
	eof       bool
	blockSize int
}

func NewSealer(c cipher.AEAD, r io.Reader, blockSize int) io.Reader {
	return &sealer{
		c:         c,
		reader:    r,
		blockSize: blockSize,
	}
}

func (s *sealer) Read(payload []byte) (int, error) { //
	if s.encrypted.Len() == 0 && !s.eof {
		var buffer bytes.Buffer

		errF := s.fillBuffer(&buffer)
		if errF != nil {
			return 0, errF
		}

		errE := s.encrypt(&buffer)
		if errE != nil {
			return 0, errE
		}
	}

	// forward
	chunk := make([]byte, len(payload))

	n, errR := s.encrypted.Read(chunk)
	if errR != nil && errR != io.EOF {
		return 0, fmt.Errorf("buffer read failed %w", errR)
	}

	nW := copy(payload, chunk[0:n])

	// did we reach end?
	if nW == 0 && s.eof {
		return copy(payload, EOFMarker), io.EOF
	}

	return nW, nil
}

func (s *sealer) encrypt(buffer *bytes.Buffer) error {
	nonce := make([]byte, s.c.NonceSize())

	_, err := rand.Read(nonce)
	if err != nil {
		return fmt.Errorf("generating nonce failed: %w", err)
	}

	encryptedBuf := s.c.Seal(nonce, nonce, buffer.Bytes(), nil)

	s.encrypted.Write(encryptedBuf)

	return nil
}

func (s *sealer) fillBuffer(buffer *bytes.Buffer) error {
	blockSizeMinusOverhead := s.blockSize - (s.c.Overhead() + s.c.NonceSize())

	_, errC := io.CopyN(buffer, s.reader, int64(blockSizeMinusOverhead))
	if errC == io.EOF {
		s.eof = true
	}

	if errC != nil && errC != io.EOF {
		return errC
	}

	return nil
}
