package usecase_test

import (
	"context"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/usecase"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSpecBackupSync(t *testing.T) {
	var (
		mockBackupUsecase *mocks.MockBackupUsecase
		mockSnapRepo      *mocks.MockSnapshotRepository
		mockStateRepo     *mocks.MockBackupStateRepo
		mockNaming        *mocks.MockSnapshotNamestrategy
		backupAutoUsecase domain.BackupSyncUsecase
		mockVolumeRepo    *mocks.MockVolumeRepository
	)
	prepare := func() {
		ctrl := gomock.NewController(t)
		mockBackupUsecase = mocks.NewMockBackupUsecase(ctrl)
		mockSnapRepo = mocks.NewMockSnapshotRepository(ctrl)
		mockStateRepo = mocks.NewMockBackupStateRepo(ctrl)
		mockNaming = mocks.NewMockSnapshotNamestrategy(ctrl)
		mockVolumeRepo = mocks.NewMockVolumeRepository(ctrl)
		backupAutoUsecase = usecase.NewBackupSync(mockBackupUsecase, mockSnapRepo, mockStateRepo, mockNaming, mockVolumeRepo)
		mockNaming.EXPECT().IsGreater(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(a, b string) bool {
			return a > b
		})
	}

	ctx := context.Background()

	Convey("Given i call the backup function", t, func() {
		Convey("When the groupName is empty", func() {
			Convey("It should set default groupName", func() {
				prepare()
				vols := []*domain.ZfsVolume{}
				mockVolumeRepo.EXPECT().ListVolumesByGroup(domain.DefaultVolumeGroup).Times(1).Return(vols, nil)

				err := backupAutoUsecase.Backup(ctx, "")

				So(err, ShouldBeNil)
			})
		})
		Convey("When everything works as expected", func() {
			Convey("It should return requests", func() {
				prepare()
				vols := []*domain.ZfsVolume{{Name: "vol1"}}
				snaps := []*domain.ZfsSnapshot{{Name: "snap1", VolumeName: "vol1"}, {Name: "snap2", VolumeName: "vol1"}}
				mockVolumeRepo.EXPECT().ListVolumesByGroup("group1").Times(1).Return(vols, nil)
				mockSnapRepo.EXPECT().ListFilter(gomock.Any()).DoAndReturn(func(filter *domain.FilterCriteria) ([]*domain.ZfsSnapshot, error) {
					So(filter.VolumeName, ShouldEqual, vols[0].Name)
					So(filter.IgnoreInvalidSnapshotNames, ShouldBeTrue)
					return snaps, nil
				})
				mockStateRepo.EXPECT().DownloadOrDefault(gomock.Any(), "vol1/.backupstate.yaml").Return(&domain.BackupState{}, nil)
				mockSnapRepo.EXPECT().GetType(gomock.Any()).Times(2).Return(domain.Full, nil)
				mockBackupUsecase.EXPECT().BackupFull(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

				err := backupAutoUsecase.Backup(ctx, "group1")

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given i call the CalcSnapsToBackup function", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It should return backup requests", func() {
				prepare()
				state := domain.BackupState{}
				snaps := []*domain.ZfsSnapshot{{Name: "snap1", VolumeName: "vol1"}, {Name: "snap2", VolumeName: "vol1"}}
				mockSnapRepo.EXPECT().GetType(gomock.Any()).Times(2).Return(domain.Full, nil)

				reqs, err := backupAutoUsecase.CalcSnapsToBackup(&state, snaps)

				So(err, ShouldBeNil)
				So(len(reqs), ShouldEqual, 2)
			})
		})
		Convey("When the remote head and the local head matches", func() {
			Convey("It should return nil reqs", func() {
				prepare()
				state := domain.BackupState{
					Head: "/pool1/fs/snap1",
				}
				snaps := []*domain.ZfsSnapshot{{Name: "snap1", VolumeName: "vol1"}, {Name: "snap2", VolumeName: "vol1"}}

				reqs, err := backupAutoUsecase.CalcSnapsToBackup(&state, snaps)

				So(err, ShouldBeNil)
				So(reqs, ShouldBeNil)
			})
		})
		Convey("When the remote head greater then the local head", func() {
			Convey("It should return error", func() {
				prepare()
				state := domain.BackupState{
					Head: "/pool1/fs/snap99999",
				}
				snaps := []*domain.ZfsSnapshot{{Name: "snap1", VolumeName: "vol1"}, {Name: "snap2", VolumeName: "vol1"}}

				reqs, err := backupAutoUsecase.CalcSnapsToBackup(&state, snaps)

				So(err.Error(), ShouldContainSubstring, "remote head is newer then local head snap99999 > snap1")
				So(reqs, ShouldBeNil)
			})
		})
		Convey("When the remote head cannot be found locally", func() {
			Convey("It should return error", func() {
				prepare()
				state := domain.BackupState{
					Head: "/pool1/fs/snap0",
				}
				snaps := []*domain.ZfsSnapshot{{Name: "snap1", VolumeName: "vol1"}, {Name: "snap2", VolumeName: "vol1"}}
				mockSnapRepo.EXPECT().GetType(gomock.Any()).Times(2).Return(domain.Full, nil)

				reqs, err := backupAutoUsecase.CalcSnapsToBackup(&state, snaps)

				So(err.Error(), ShouldContainSubstring, "remote head 'snap0' could not be found locally")
				So(reqs, ShouldBeNil)
			})
		})
		Convey("When a backup type is unknown", func() {
			Convey("It should skip those snapshots", func() {
				prepare()
				state := domain.BackupState{
					Head: "/pool1/fs/snap1",
				}
				snaps := []*domain.ZfsSnapshot{{Name: "snap3", VolumeName: "vol1"}, {Name: "snap2", VolumeName: "vol1"}, {Name: "snap1", VolumeName: "vol1"}}
				mockSnapRepo.EXPECT().GetType(gomock.Any()).Times(1).Return(domain.Full, nil)
				mockSnapRepo.EXPECT().GetType(gomock.Any()).Times(1).Return(domain.Unknown, nil)
				mockSnapRepo.EXPECT().GetType(gomock.Any()).Times(1).Return(domain.Full, nil)

				reqs, err := backupAutoUsecase.CalcSnapsToBackup(&state, snaps)

				So(err, ShouldBeNil)
				So(reqs[0].Snapshot.Name, ShouldEqual, "snap1")
				So(reqs[1].Snapshot.Name, ShouldEqual, "snap3")
				So(len(reqs), ShouldEqual, 2)
			})
		})
		Convey("When the first backup is not a full backup", func() {
			Convey("It should error", func() {
				prepare()
				state := domain.BackupState{}
				snaps := []*domain.ZfsSnapshot{{Name: "snap3", VolumeName: "vol1"}}
				mockSnapRepo.EXPECT().GetType(gomock.Any()).Times(1).Return(domain.Incremental, nil)

				reqs, err := backupAutoUsecase.CalcSnapsToBackup(&state, snaps)

				So(err.Error(), ShouldContainSubstring, "first backup needs to be a full backup")
				So(reqs, ShouldBeNil)
			})
		})
	})
}
