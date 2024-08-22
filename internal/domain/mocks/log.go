package mocks

import "github.com/backup-blob/zfs-backup-blob/internal/domain"

type mockLogger struct {
}

func NewMockLogger() domain.LogDriver {
	return &mockLogger{}
}

func (l *mockLogger) Debugf(format string, v ...interface{}) {
}

func (l *mockLogger) Infof(format string, v ...interface{}) {
}

type mockLoggerRepo struct {
}

func NewMockLoggerRepo() domain.LogRepository {
	return &mockLoggerRepo{}
}

func (l *mockLoggerRepo) Debugf(format string, v ...interface{}) {
}

func (l *mockLoggerRepo) Infof(format string, v ...interface{}) {
}
