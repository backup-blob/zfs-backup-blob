package usecase_test

import (
	"errors"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/usecase"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSpecSnapshot(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSnapRepo := mocks.NewMockSnapshotRepository(ctrl)
	mockVolumeRepo := mocks.NewMockVolumeRepository(ctrl)

	snapUsecase := usecase.NewSnapshot(mockSnapRepo, mockVolumeRepo, mocks.NewMockLoggerRepo())

	Convey("Given the CreateByVolume function is called", t, func() {
		Convey("When everything goes positive", func() {
			Convey("It should create snapshot", func() {
				backupType := domain.Full
				volume := &domain.ZfsVolume{Name: "pool/path", GroupName: "group"}
				mockSnapRepo.EXPECT().CreateWithType(volume, backupType)

				err := snapUsecase.CreateByVolume(volume, backupType)

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given the create function is called", t, func() {
		Convey("When no volumes are in the given group", func() {
			Convey("It should error", func() {
				mockVolumeRepo.EXPECT().ListVolumesByGroup("group1").Times(1).Return([]*domain.ZfsVolume{}, nil)

				err := snapUsecase.CreateByGroup("group1", "full")

				So(err, ShouldNotBeNil)
			})
		})
		Convey("When a empty group is supplied", func() {
			Convey("It should default to the default group", func() {
				vol1 := domain.ZfsVolume{
					Name:      "vol1",
					GroupName: domain.DefaultVolumeGroup,
				}
				mockVolumeRepo.EXPECT().ListVolumesByGroup(domain.DefaultVolumeGroup).Times(1).Return([]*domain.ZfsVolume{&vol1}, nil)
				mockSnapRepo.EXPECT().CreateWithType(&vol1, domain.Full).Times(1).Return("", nil)

				err := snapUsecase.CreateByGroup("", "full")

				So(err, ShouldBeNil)
			})
		})
		Convey("When everything works as expected", func() {
			Convey("It should create snapshots of the volumes in the group", func() {
				vol1 := domain.ZfsVolume{
					Name:      "vol1",
					GroupName: "group1",
				}
				mockVolumeRepo.EXPECT().ListVolumesByGroup("group1").Times(1).Return([]*domain.ZfsVolume{&vol1}, nil)
				mockSnapRepo.EXPECT().CreateWithType(&vol1, domain.Full).Times(1).Return("", nil)

				err := snapUsecase.CreateByGroup("group1", "full")

				So(err, ShouldBeNil)
			})
		})
		Convey("When the creation of the snapshot fails", func() {
			Convey("Then it should error", func() {
				vol1 := domain.ZfsVolume{
					Name:      "vol1",
					GroupName: "group1",
				}
				mockVolumeRepo.EXPECT().ListVolumesByGroup("group1").Times(1).Return([]*domain.ZfsVolume{&vol1}, nil)
				mockSnapRepo.EXPECT().CreateWithType(&vol1, domain.Full).Times(1).Return("", errors.New("error"))

				err := snapUsecase.CreateByGroup("group1", "full")

				So(err, ShouldNotBeNil)
			})
		})
		Convey("When listing of volumes fails", func() {
			Convey("Then it should error", func() {
				mockVolumeRepo.EXPECT().ListVolumesByGroup("group1").Times(1).Return([]*domain.ZfsVolume{}, errors.New("error"))

				err := snapUsecase.CreateByGroup("group1", "full")

				So(err, ShouldNotBeNil)
			})
		})
	})
}
