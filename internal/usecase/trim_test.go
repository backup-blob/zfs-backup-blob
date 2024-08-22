package usecase_test

import (
	"context"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
	mocks_config "github.com/backup-blob/zfs-backup-blob/internal/domain/config/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/usecase"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSpecTrim(t *testing.T) {
	var ctrl *gomock.Controller
	var mockLogger domain.LogDriver
	var mockVolumeRepo *mocks.MockVolumeRepository
	var mockStateRepo *mocks.MockBackupStateRepo
	var mockConfigUse *mocks_config.MockConfigRepo
	var trimUsecase domain.TrimUsecase
	var mockBackupRepo *mocks.MockBackupRepository
	var mockSnapRepo *mocks.MockSnapshotRepository

	ctx := context.Background()
	vols := []*domain.ZfsVolume{{Name: "vol/vol1", GroupName: "test"}}

	var setup = func() {
		ctrl = gomock.NewController(t)
		mockLogger = mocks.NewMockLogger()
		mockVolumeRepo = mocks.NewMockVolumeRepository(ctrl)
		mockStateRepo = mocks.NewMockBackupStateRepo(ctrl)
		mockConfigUse = mocks_config.NewMockConfigRepo(ctrl)
		mockBackupRepo = mocks.NewMockBackupRepository(ctrl)
		mockSnapRepo = mocks.NewMockSnapshotRepository(ctrl)

		trimUsecase = usecase.NewTrimUseCase(mockLogger, mockVolumeRepo, mockStateRepo, mockConfigUse, mockBackupRepo, mockSnapRepo)
	}

	Convey("Given the TrimLocal function is called", t, func() {
		Convey("When dry run is active", func() {
			Convey("It should not actually delete any backups", func() {
				setup()
				mockConfigUse.EXPECT().GetConfig().AnyTimes().Return(&config.Config{LocalTrimPolicy: "F"})
				mockVolumeRepo.EXPECT().ListVolumesByGroup("test").Return(vols, nil)
				snapsWithType := map[string]*domain.ZfsSnapshotWithType{
					"1": {ZfsSnapshot: domain.ZfsSnapshot{Name: "backup_blob_2011-03-13T07-06-40Z", VolumeName: "vol1"}, Type: domain.Full},
					"2": {ZfsSnapshot: domain.ZfsSnapshot{Name: "backup_blob_2011-03-14T07-06-40Z", VolumeName: "vol1"}, Type: domain.Full},
				}
				var snaps []*domain.ZfsSnapshot
				for _, v := range snapsWithType {
					snaps = append(snaps, &v.ZfsSnapshot)
					mockSnapRepo.EXPECT().GetType(v.ZfsSnapshot.FullName()).Return(v.Type, nil)
				}
				mockSnapRepo.EXPECT().ListFilter(gomock.Any()).Return(snaps, nil)

				err := trimUsecase.TrimLocal(ctx, &domain.TrimLocalParameters{GroupName: "test", DryRun: true})

				So(err, ShouldBeNil)
			})
		})
		Convey("When everything works as expected", func() {
			setup()
			mockConfigUse.EXPECT().GetConfig().AnyTimes().Return(&config.Config{LocalTrimPolicy: "F"})
			mockVolumeRepo.EXPECT().ListVolumesByGroup("test").Return(vols, nil)
			snapsWithType := map[string]*domain.ZfsSnapshotWithType{
				"1": {ZfsSnapshot: domain.ZfsSnapshot{Name: "backup_blob_2011-03-12T07-06-40Z", VolumeName: "vol1"}, Type: domain.Incremental},
				"2": {ZfsSnapshot: domain.ZfsSnapshot{Name: "backup_blob_2011-03-13T07-06-40Z", VolumeName: "vol1"}, Type: domain.Full},
				"3": {ZfsSnapshot: domain.ZfsSnapshot{Name: "backup_blob_2011-03-14T07-06-40Z", VolumeName: "vol1"}, Type: domain.Full},
			}
			var snaps []*domain.ZfsSnapshot
			for _, v := range snapsWithType {
				snaps = append(snaps, &v.ZfsSnapshot)
				mockSnapRepo.EXPECT().GetType(v.ZfsSnapshot.FullName()).Return(v.Type, nil)
			}
			mockSnapRepo.EXPECT().ListFilter(gomock.Any()).Return(snaps, nil)
			mockSnapRepo.EXPECT().Delete(&snapsWithType["1"].ZfsSnapshot).Return(nil)
			mockSnapRepo.EXPECT().Delete(&snapsWithType["2"].ZfsSnapshot).Return(nil)

			err := trimUsecase.TrimLocal(ctx, &domain.TrimLocalParameters{GroupName: "test"})

			So(err, ShouldBeNil)
		})
		Convey("When the remote trim policy is invalid", func() {
			Convey("It should return an error", func() {
				setup()
				mockConfigUse.EXPECT().GetConfig().AnyTimes().Return(&config.Config{LocalTrimPolicy: ""})
				err := trimUsecase.TrimLocal(ctx, &domain.TrimLocalParameters{GroupName: "test"})

				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given the TrimRemote function is called", t, func() {
		for _, testCase := range testCases {
			Convey("When test Case: "+testCase.name+" is executed", func() {
				setup()
				mockConfigUse.EXPECT().GetConfig().AnyTimes().Return(&config.Config{RemoteTrimPolicy: testCase.policy})
				mockVolumeRepo.EXPECT().ListVolumesByGroup("test").Return(vols, nil)
				mockStateRepo.EXPECT().Download(gomock.Any(), gomock.Any()).Return(&domain.BackupState{
					Backups: testCase.backupMap,
				}, nil)
				deletedKeys := []string{}
				mockBackupRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(ctx context.Context, bd *domain.BackupDelete) error {
					deletedKeys = append(deletedKeys, bd.BlobKey)
					return nil
				})
				mockStateRepo.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(ctx context.Context, key string, state *domain.BackupState) error {
					So(state.Backups, ShouldNotContainKey, "1")
					for _, bKey := range testCase.backupsRemaining {
						So(state.Backups, ShouldContainKey, bKey)
					}
					So(len(state.Backups), ShouldEqual, len(testCase.backupsRemaining))
					return nil
				})

				err := trimUsecase.TrimRemote(ctx, &domain.TrimRemoteParameters{GroupName: "test"})

				for _, bKey := range testCase.backupsDeleted {
					So(deletedKeys, ShouldContain, bKey)
				}
				So(len(deletedKeys), ShouldEqual, len(testCase.backupsDeleted))

				So(err, ShouldBeNil)
			})
		}
		Convey("When everything works as expected", func() {
			Convey("It should deleted trim-able backups", func() {
				setup()
				mockConfigUse.EXPECT().GetConfig().AnyTimes().Return(&config.Config{RemoteTrimPolicy: "IIIFF"})
				mockVolumeRepo.EXPECT().ListVolumesByGroup("test").Return(vols, nil)
				bs := domain.BackupState{
					Backups: map[string]domain.BackupRecord{
						"1": {Type: domain.Full},
						"2": {Type: domain.Full},
						"3": {Type: domain.Full},
					},
				}
				mockStateRepo.EXPECT().Download(gomock.Any(), gomock.Any()).Return(&bs, nil)
				mockBackupRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, bd *domain.BackupDelete) error {
					So(bd.BlobKey, ShouldEqual, "1")
					return nil
				})
				mockStateRepo.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, key string, state *domain.BackupState) error {
					So(state.Backups, ShouldNotContainKey, "1")
					So(state.Backups, ShouldContainKey, "2")
					So(state.Backups, ShouldContainKey, "3")
					return nil
				})

				err := trimUsecase.TrimRemote(ctx, &domain.TrimRemoteParameters{GroupName: "test"})

				So(err, ShouldBeNil)
			})
		})
		Convey("When dry run is active", func() {
			Convey("It should not actually delete any backups", func() {
				setup()
				mockConfigUse.EXPECT().GetConfig().AnyTimes().Return(&config.Config{RemoteTrimPolicy: "IIIFF"})
				mockVolumeRepo.EXPECT().ListVolumesByGroup("test").Return(vols, nil)
				bs := domain.BackupState{
					Backups: map[string]domain.BackupRecord{
						"1": {Type: domain.Full},
						"2": {Type: domain.Full},
						"3": {Type: domain.Full},
					},
				}
				mockStateRepo.EXPECT().Download(gomock.Any(), gomock.Any()).Return(&bs, nil)

				err := trimUsecase.TrimRemote(ctx, &domain.TrimRemoteParameters{GroupName: "test", DryRun: true})

				So(err, ShouldBeNil)
			})
		})
		Convey("When the remote trim policy is invalid", func() {
			Convey("It should return an error", func() {
				setup()
				mockConfigUse.EXPECT().GetConfig().AnyTimes().Return(&config.Config{RemoteTrimPolicy: ""})
				err := trimUsecase.TrimRemote(ctx, &domain.TrimRemoteParameters{GroupName: "test"})

				So(err, ShouldNotBeNil)
			})
		})
	})
}
