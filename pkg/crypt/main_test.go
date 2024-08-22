package crypt_test

import (
	"bytes"
	"crypto/cipher"
	"errors"
	"github.com/backup-blob/zfs-backup-blob/pkg/crypt"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"math/rand"
	"strings"
	"testing"
)

func FuzzEncrypt(f *testing.F) {
	f.Add(1024, 1024, 1024)
	f.Fuzz(func(t *testing.T, _payloadSize int, _blockSize int, _bufferSize int) {
		Convey("It should decrypt it", t, func() {
			payloadSize := minMaxInt(_payloadSize, 1, 1024*1024*100)       // 1 bye - 100mb
			blockSize := minMaxInt(_blockSize, 1024*1024*1, 1024*1024*100) // 1mb - 100mb
			bufferSize := minMaxInt(_bufferSize, 1024, 1024*1024*100)      // 1kb - 100mb
			setupF := NewSetup(&payloadSize)
			crypter := crypt.NewCrypter(blockSize)

			encryptedReader, errR := crypter.Encrypt(setupF.C, setupF.InputReader)
			So(errR, ShouldBeNil)

			copyChunked(setupF.EncryptedBuf, encryptedReader, bufferSize)

			decryptedWriter, errD := crypter.Decrypt(setupF.C, setupF.DecryptedBuf)
			So(errD, ShouldBeNil)

			copyChunked(decryptedWriter, setupF.EncryptedBuf, bufferSize)

			So(setupF.DecryptedBuf.String(), ShouldEqual, setupF.Input)
		})
	})
}

func TestSpecCrypt(t *testing.T) {
	Convey("Should error if key is too short", t, func() {
		res, err := crypt.NewAESGCM(crypt.Key("fff"))

		So(err, ShouldBeError, "aes.NewCipher failed: crypto/aes: invalid key size 3")
		So(res, ShouldBeNil)
	})
	Convey("Encrypt should propagate reader errors", t, func() {
		setupF := NewSetup(nil)
		crypter := crypt.NewCrypter(1024)

		encryptedReader, errR := crypter.Encrypt(setupF.C, FailedReader{})
		So(errR, ShouldBeNil)

		_, err := io.Copy(setupF.EncryptedBuf, encryptedReader)
		So(err, ShouldBeError, "simulated read error")
	})
	Convey("Decrypt should fail if payload is invalid", t, func() {
		setupF := NewSetup(nil)
		crypter := crypt.NewCrypter(1024)

		decryptedWriter, errD := crypter.Decrypt(setupF.C, setupF.DecryptedBuf)
		So(errD, ShouldBeNil)

		_, err := io.Copy(decryptedWriter, strings.NewReader("something-not-encrypted<EOF>"))
		So(err, ShouldBeError, "cipher: message authentication failed")
	})
	Convey("Decrypt should fail if writer fails", t, func() {
		setupF := NewSetup(nil)
		crypter := crypt.NewCrypter(1024)

		encryptedReader, errR := crypter.Encrypt(setupF.C, setupF.InputReader)
		So(errR, ShouldBeNil)

		io.Copy(setupF.EncryptedBuf, encryptedReader)

		decryptedWriter, errD := crypter.Decrypt(setupF.C, FailingWriter{})
		So(errD, ShouldBeNil)

		_, err := io.Copy(decryptedWriter, setupF.EncryptedBuf)

		So(err, ShouldBeError, "simulated write error")
	})
	Convey("Should handle buffer-size < block-size", t, func() {
		setupF := NewSetup(nil)
		crypter := crypt.NewCrypter(1024)
		bufferSize := 900

		encryptedReader, errR := crypter.Encrypt(setupF.C, setupF.InputReader)
		So(errR, ShouldBeNil)

		copyChunked(setupF.EncryptedBuf, encryptedReader, bufferSize)

		decryptedWriter, errD := crypter.Decrypt(setupF.C, setupF.DecryptedBuf)
		So(errD, ShouldBeNil)

		copyChunked(decryptedWriter, setupF.EncryptedBuf, bufferSize)

		So(setupF.DecryptedBuf.String(), ShouldEqual, setupF.Input)
	})
	Convey("Should handle buffer-size > block-size", t, func() {
		setupF := NewSetup(nil)
		crypter := crypt.NewCrypter(1024)
		bufferSize := 1200

		encryptedReader, errR := crypter.Encrypt(setupF.C, setupF.InputReader)
		So(errR, ShouldBeNil)

		copyChunked(setupF.EncryptedBuf, encryptedReader, bufferSize)

		decryptedWriter, errD := crypter.Decrypt(setupF.C, setupF.DecryptedBuf)
		So(errD, ShouldBeNil)

		copyChunked(decryptedWriter, setupF.EncryptedBuf, bufferSize)

		So(setupF.DecryptedBuf.String(), ShouldEqual, setupF.Input)
	})
}

type FailingWriter struct{}

func (w FailingWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("simulated write error")
}

type FailedReader struct{}

func (r FailedReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

type setup struct {
	C            cipher.AEAD
	Input        string
	InputReader  io.Reader
	EncryptedBuf *bytes.Buffer
	DecryptedBuf *bytes.Buffer
}

func NewSetup(payloadSize *int) *setup {
	size := 5000
	if payloadSize != nil {
		size = *payloadSize
	}
	c, err := crypt.NewAESGCM(crypt.NewKey("hello"))
	if err != nil {
		panic(err)
	}
	str := RandStringRunes(size)
	input := strings.NewReader(str)
	encrypted := new(bytes.Buffer)
	decrypted := new(bytes.Buffer)

	return &setup{
		C:            c,
		Input:        str,
		InputReader:  input,
		EncryptedBuf: encrypted,
		DecryptedBuf: decrypted,
	}
}

func minMaxInt(i int, min int, max int) int {
	if i > max {
		return max
	}
	if i < min {
		return min
	}
	return i
}

func copyChunked(target io.Writer, source io.Reader, bufferSize int) {
	// Create a buffer with the defined size
	buffer := make([]byte, bufferSize)

	// Loop to read from source and write to destination in chunks
	for {
		// Read from the source file into the buffer
		bytesRead, err := source.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}

		// Write the read bytes to the destination file
		_, errF := target.Write(buffer[:bytesRead])
		if errF != nil {
			panic(errF)
		}

		// If we've reached the end of the file, break out of the loop
		if err == io.EOF {
			break
		}
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
