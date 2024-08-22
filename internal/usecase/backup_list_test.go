package usecase_test

import (
	"bytes"
	"context"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/usecase"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSpecBackupList(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBackupState := mocks.NewMockBackupStateRepo(ctrl)
	mockNamestrategy := mocks.NewMockSnapshotNamestrategy(ctrl)
	mockRenderRepo := mocks.NewMockRenderRepository(ctrl)
	backupList := usecase.NewBackupList(mockBackupState, mockNamestrategy, mockRenderRepo)
	ctx := context.Background()

	Convey("Given i call the List function", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It should not error", func() {
				buf := bytes.NewBuffer([]byte{})
				backupState := domain.BackupState{Backups: map[string]domain.BackupRecord{
					"fff/fff/backup_blob_1300000000": {Type: domain.Incremental},
					"fff/fff/backup_blob_1400000000": {Type: domain.Full},
				}}
				mockBackupState.EXPECT().Download(gomock.Any(), gomock.Any()).Return(&backupState, nil)
				mockNamestrategy.EXPECT().IsGreater(gomock.Any(), gomock.Any()).DoAndReturn(func(i, j string) bool {
					return i > j
				})
				mockRenderRepo.EXPECT().RenderBackupTable(gomock.Any(), gomock.Any()).AnyTimes()

				err := backupList.List(ctx, "pool/vol1", buf)

				So(err, ShouldBeNil)
			})
		})
	})
}
