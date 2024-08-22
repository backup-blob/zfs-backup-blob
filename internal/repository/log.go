package repository

import "github.com/backup-blob/zfs-backup-blob/internal/domain"

type logger struct {
	logDriver domain.LogDriver
}

func NewLog(logDriver domain.LogDriver) domain.LogRepository {
	return &logger{logDriver: logDriver}
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.logDriver.Debugf(format, v...)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.logDriver.Infof(format, v...)
}
