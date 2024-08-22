package crypt

import (
	"crypto/cipher"
	"io"
)

type crypter struct {
	blockSize int
}

func NewCrypter(blockSize int) EncryptDecrypter {
	return &crypter{
		blockSize: blockSize,
	}
}

func (a *crypter) Encrypt(c cipher.AEAD, r io.Reader) (io.Reader, error) {
	return NewSealer(c, r, a.blockSize), nil
}

func (a *crypter) Decrypt(c cipher.AEAD, w io.Writer) (io.Writer, error) {
	return NewUnsealer(c, w, a.blockSize), nil
}
