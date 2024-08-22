package crypt

import (
	"crypto/cipher"
	"io"
)

var EOFMarker = "<EOF>"

type Key []byte

type Encrypter interface {
	Encrypt(c cipher.AEAD, r io.Reader) (io.Reader, error)
}

type Decrypter interface {
	Decrypt(c cipher.AEAD, w io.Writer) (io.Writer, error)
}

type EncryptDecrypter interface {
	Encrypter
	Decrypter
}
