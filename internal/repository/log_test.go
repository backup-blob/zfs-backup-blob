package repository_test

import (
	"github.com/backup-blob/zfs-backup-blob/internal/domain/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/repository"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLogger(t *testing.T) {
	mockDriver := mocks.NewMockLogger()
	logger := repository.NewLog(mockDriver)

	Convey("Debugf", t, func() {
		logger.Debugf("Debug message: %s", "test")
	})

	Convey("Infof", t, func() {
		logger.Infof("Info message: %s", "test")
	})
}
