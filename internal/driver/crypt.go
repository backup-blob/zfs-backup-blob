package driver

import (
	"crypto/cipher"
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
	"github.com/backup-blob/zfs-backup-blob/pkg/crypt"
	"io"
)

var DefaultBlockSize = 1024 * 1024 * 4

type CryptDriver struct {
	Conf    *config.CryptConfig
	Crypter crypt.EncryptDecrypter
}

func NewCrypt(conf *config.CryptConfig) domain.Middleware {
	return &CryptDriver{Conf: conf, Crypter: crypt.NewCrypter(DefaultBlockSize)}
}

func (cd *CryptDriver) loadCipher() (cipher.AEAD, error) {
	key := crypt.NewKey(cd.Conf.Password)
	return crypt.NewAESGCM(key)
}

func (cd *CryptDriver) Write(w io.Writer) (wp io.Writer, err error) {
	c, err := cd.loadCipher()
	if err != nil {
		return nil, fmt.Errorf("cd.loadCipher failed %w", err)
	}

	return cd.Crypter.Decrypt(c, w)
}

func (cd *CryptDriver) Read(r io.Reader) (rp io.Reader, err error) {
	c, err := cd.loadCipher()
	if err != nil {
		return nil, fmt.Errorf("cd.loadCipher failed %w", err)
	}

	return cd.Crypter.Encrypt(c, r)
}
