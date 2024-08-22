package driver

import (
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
	"github.com/backup-blob/zfs-backup-blob/pkg/throttle"
	"io"
)

type Throttle struct {
	Conf *config.ThrottleConfig
}

func NewThrottle(conf *config.ThrottleConfig) domain.Middleware {
	return &Throttle{Conf: conf}
}

func (t *Throttle) Write(w io.Writer) (wp io.Writer, err error) {
	if t.Conf.WriteSpeed == 0 {
		return w, nil
	}

	return throttle.SpeedlimitWriter(t.Conf.WriteSpeed)(w)
}

func (t *Throttle) Read(r io.Reader) (rp io.Reader, err error) {
	if t.Conf.ReadSpeed == 0 {
		return r, nil
	}

	return throttle.SpeedlimitReader(t.Conf.ReadSpeed)(r)
}
