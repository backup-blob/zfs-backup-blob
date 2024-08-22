package repository_test

import (
	"bytes"
	"go.uber.org/mock/gomock"
	"testing"

	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/repository"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRender(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRenderDriver := mocks.NewMockRenderDriver(ctrl)
	mockSizer := func(size *int64) string {
		return "1MB"
	}

	render := repository.NewRender(mockRenderDriver, mockSizer)

	Convey("Given RenderVolumeTable is called", t, func() {
		Convey("When everything works as expected", func() {
			Convey("Then the driver should be called ", func() {
				writer := &bytes.Buffer{}
				volumes := []*domain.ZfsVolume{
					{Name: "vol1", GroupName: "group1"},
					{Name: "vol2", GroupName: "group2"},
				}

				expectedHeaderRow := []interface{}{"Volume", "Group"}
				expectedRows := [][]interface{}{
					{"vol1", "group1"},
					{"vol2", "group2"},
				}

				mockRenderDriver.EXPECT().RenderTable(writer, expectedHeaderRow, expectedRows)

				render.RenderVolumeTable(writer, volumes)
			})
		})

		Convey("Given RenderBackupTable is called", func() {
			Convey("When everything works as expected", func() {
				Convey("Then the driver should be called ", func() {
					writer := &bytes.Buffer{}
					var size int64 = 100
					backups := []domain.BackupRecordWithKey{
						{Key: "key1", BackupRecord: domain.BackupRecord{Type: domain.Full, Size: &size}},
						{Key: "key2", BackupRecord: domain.BackupRecord{Type: domain.Incremental, Size: &size}},
					}

					expectedHeaderRow := []interface{}{"Key", "Type", "Size"}
					expectedRows := [][]interface{}{
						{"key1", "full", "1MB"},
						{"key2", "incremental", "1MB"},
					}

					mockRenderDriver.EXPECT().RenderTable(writer, expectedHeaderRow, expectedRows)

					render.RenderBackupTable(writer, backups)
				})
			})
		})
	})
}
