package usecase_test

import (
	"context"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	mocks_config "github.com/backup-blob/zfs-backup-blob/internal/domain/config/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/usecase"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSpecBackup(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBackupRepo := mocks.NewMockBackupRepository(ctrl)
	mockBackupState := mocks.NewMockBackupStateRepo(ctrl)
	mockConfigRepo := mocks_config.NewMockConfigRepo(ctrl)
	backupUsecase := usecase.NewBackup(mockBackupRepo, mockBackupState, mocks.NewMockLoggerRepo(), mockConfigRepo)
	ctx := context.Background()
	backupKey := "folder1/folder2/backup1"
	backupParentKey := "folder1/folder2/backup2"
	backupStateKey := "folder1/folder2/.backupstate.yaml"
	uploadRes := domain.UploadResponse{Size: 1}
	mockConfigRepo.EXPECT().GetMiddlewares().AnyTimes().Return([]domain.Middleware{})

	Convey("Given i call the Restore function", t, func() {
		Convey("When the backup does not exists on the remote", func() {
			Convey("I should error", func() {
				mockBackupState.EXPECT().DownloadOrDefault(gomock.Any(), gomock.Any()).Return(&domain.BackupState{}, nil)

				err := backupUsecase.Restore(ctx, &domain.RestoreParams{
					TargetZfsLocation: "folder1",
					BlobKey:           backupKey,
				})

				So(err, ShouldNotBeNil)
			})
		})
		Convey("When a full backup is restored", func() {
			Convey("I should restore the backup", func() {
				state := domain.BackupState{Backups: map[string]domain.BackupRecord{backupKey: {
					Type: domain.Full,
				}}}
				mockBackupState.EXPECT().DownloadOrDefault(gomock.Any(), backupStateKey).Return(&state, nil)
				mockBackupRepo.EXPECT().Restore(gomock.Any(), gomock.Any()).Return(nil)

				err := backupUsecase.Restore(ctx, &domain.RestoreParams{
					TargetZfsLocation: "folder1",
					BlobKey:           backupKey,
				})

				So(err, ShouldBeNil)
			})
		})
		Convey("When a incremental backup is restored", func() {
			Convey("I should restore the backup", func() {
				state := domain.BackupState{Backups: map[string]domain.BackupRecord{backupKey: {
					Type: domain.Incremental,
				}}}
				mockBackupState.EXPECT().DownloadOrDefault(gomock.Any(), backupStateKey).Return(&state, nil)
				mockBackupRepo.EXPECT().Restore(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, br *domain.BackupRestore) error {
					So(br.BlobKey, ShouldEqual, backupKey)
					So(br.TargetZfsLocation, ShouldEqual, "folder1")
					return nil
				})

				err := backupUsecase.Restore(ctx, &domain.RestoreParams{
					TargetZfsLocation: "folder1",
					BlobKey:           backupKey,
				})

				So(err, ShouldBeNil)
			})
		})
		Convey("When a incremental backup is restored until the full", func() {
			Convey("I should restore the backup", func() {
				state := domain.BackupState{Backups: map[string]domain.BackupRecord{
					backupKey: {
						ParentBackupKey: backupParentKey,
						Type:            domain.Incremental,
					},
					backupParentKey: {
						Type: domain.Full,
					},
				}}
				mockBackupState.EXPECT().DownloadOrDefault(gomock.Any(), backupStateKey).Return(&state, nil)
				mockBackupRepo.EXPECT().Restore(gomock.Any(), gomock.Any()).Times(2).DoAndReturn(func(_ context.Context, br *domain.BackupRestore) error {
					So(br.BlobKey, ShouldBeIn, []string{backupKey, backupParentKey})
					So(br.TargetZfsLocation, ShouldEqual, "folder1")
					return nil
				})

				err := backupUsecase.Restore(ctx, &domain.RestoreParams{
					TargetZfsLocation: "folder1",
					BlobKey:           backupKey,
					RestoreAll:        true,
				})

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given i call the BackupFull function", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It should create a full backup", func() {
				mockBackupRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&uploadRes, nil)
				mockBackupState.EXPECT().UpdateState(gomock.Any(), "folder1/folder2/.backupstate.yaml", gomock.Any()).DoAndReturn(func(_ context.Context, _ string, f func(state *domain.BackupState) error) error {
					s := domain.BackupState{Backups: make(map[string]domain.BackupRecord)}
					f(&s)
					entry := s.Backups["folder1/folder2/snap"]
					So(entry.Type, ShouldEqual, domain.Full)
					So(entry.Size, ShouldEqual, &uploadRes.Size)
					So(entry.ParentBackupKey, ShouldEqual, "")
					So(s.Head, ShouldEqual, "folder1/folder2/snap")
					return nil
				})

				err := backupUsecase.BackupFull(ctx, "folder1/folder2@snap", true)

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given i call the BackupIncremental function", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It should create a incremental backup", func() {
				mockBackupRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&uploadRes, nil)
				mockBackupState.EXPECT().UpdateState(gomock.Any(), "folder1/folder2/.backupstate.yaml", gomock.Any()).DoAndReturn(func(_ context.Context, _ string, f func(state *domain.BackupState) error) error {
					s := domain.BackupState{Backups: make(map[string]domain.BackupRecord)}
					f(&s)
					entry := s.Backups["folder1/folder2/snap2"]
					So(entry.Type, ShouldEqual, domain.Incremental)
					So(entry.Size, ShouldEqual, &uploadRes.Size)
					So(entry.ParentBackupKey, ShouldEqual, "folder1/folder2/snap1")
					So(s.Head, ShouldEqual, "folder1/folder2/snap2")
					return nil
				})

				err := backupUsecase.BackupIncremental(ctx, "folder1/folder2@snap1", "folder1/folder2@snap2", true)

				So(err, ShouldBeNil)
			})
		})
	})
}
