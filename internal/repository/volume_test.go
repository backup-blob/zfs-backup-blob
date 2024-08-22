package repository_test

import (
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/repository"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSpecVol(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockZfsDriver(ctrl)
	repo := repository.NewVolume(m)

	Convey("Given the TagVolumeWithGroup is called", t, func() {
		Convey("When everything goes positive", func() {
			Convey("It should not return an error", func() {
				m.EXPECT().SetField(gomock.Any(), "group1", "pool1/path").Return(nil)

				err := repo.TagVolumeWithGroup(&domain.ZfsVolume{GroupName: "group1", Name: "pool1/path"})

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given i call the ListVolumes function", t, func() {
		Convey("When everything works as expected", func() {

			Convey("It should not return an error", func() {
				m.EXPECT().List(gomock.Any()).Times(1).Return(fakeCommand(`NAME         BACKUP_BLOB::GROUP
zfs-pool     -
zfs-pool/33  group1
zfs-pool2    group2`))
				volumes, err := repo.ListVolumes()

				So(volumes[0].Name, ShouldEqual, "zfs-pool")
				So(volumes[0].GroupName, ShouldEqual, "-")
				So(volumes[1].Name, ShouldEqual, "zfs-pool/33")
				So(volumes[1].GroupName, ShouldEqual, "group1")
				So(volumes[2].Name, ShouldEqual, "zfs-pool2")
				So(volumes[2].GroupName, ShouldEqual, "group2")
				So(err, ShouldBeNil)
			})
		})
	})
	Convey("Given i call ListVolumesByGroup ", t, func() {
		Convey("When everything works as exptect", func() {
			Convey("It should return volumes filtered by group", func() {
				m.EXPECT().List(gomock.Any()).Times(1).Return(fakeCommand(`NAME         BACKUP_BLOB::GROUP
zfs-pool     -
zfs-pool/33  group1
zfs-pool2    group2`))
				volumes, err := repo.ListVolumesByGroup("group1")

				So(err, ShouldBeNil)
				So(len(volumes), ShouldEqual, 1)
				So(volumes[0].Name, ShouldEqual, "zfs-pool/33")
			})
		})
	})
}
