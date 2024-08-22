package repository_test

import (
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/repository"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"os/exec"
	"strings"
	"testing"
)

func TestSpec(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockZfsDriver(ctrl)
	namestrategy := mocks.NewMockSnapshotNamestrategy(ctrl)
	repo := repository.NewSnapshot(m, namestrategy)

	Convey("Given the Delete function is called", t, func() {
		Convey("When everthing works as expected", func() {
			Convey("It should not return an error", func() {
				m.EXPECT().Destroy("vol/path@snap")

				err := repo.Delete(&domain.ZfsSnapshot{Name: "snap", VolumeName: "vol/path"})

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given a snapshot is tagged with a backuptype", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It should not error", func() {
				m.EXPECT().Snapshot("test@snapname").Times(1).Return(fakeCommand(""))
				m.EXPECT().SetField("backup_blob::type", "full", "test@snapname").Times(1).Return(nil)
				namestrategy.EXPECT().GetName().Return("snapname")

				name, err := repo.CreateWithType(&domain.ZfsVolume{Name: "test"}, domain.Full)

				So(err, ShouldBeNil)
				So(name, ShouldEqual, "test@snapname")
			})
		})
	})

	Convey("Given a snapshot has a backuptype tag", t, func() {
		Convey("When i fetch the type with GetType and no error happens", func() {
			Convey("It should return the type", func() {
				m.EXPECT().GetField("backup_blob::type", "volume1").Times(1).Return(fakeCommand("full \n"))
				bType, err := repo.GetType("volume1")

				So(err, ShouldBeNil)
				So(bType, ShouldEqual, domain.Full)
			})
		})
	})

	Convey("Given the snapshot is instantiated", t, func() {
		Convey("When i call the Create function with all parameters", func() {
			Convey("It should not return an error", func() {
				m.EXPECT().Snapshot("test@snapname").Times(1).Return(fakeCommand(""))
				namestrategy.EXPECT().GetName().Return("snapname")

				name, err := repository.NewSnapshot(m, namestrategy).Create(&domain.ZfsVolume{Name: "test"})

				So(err, ShouldBeNil)
				So(name, ShouldEqual, "test@snapname")
			})
		})
		Convey("When i call the List function", func() {
			Convey("It should not return an error", func() {
				m.EXPECT().List(gomock.Any()).Times(1).Return(fakeCommand(`NAME
zfs-pool@1
zfs-pool/33@1`))
				snaps, err := repository.NewSnapshot(m, namestrategy).List()

				So(snaps[0].Name, ShouldEqual, "1")
				So(snaps[0].VolumeName, ShouldEqual, "zfs-pool")
				So(err, ShouldBeNil)
			})
		})
		Convey("When i call the ListByVolume function", func() {
			Convey("It should return a filtered list of snapshots by volume name", func() {
				m.EXPECT().List(gomock.Any()).Times(1).Return(fakeCommand(`NAME
zfs-pool@1
zfs-pool/33@1`))
				filter := domain.FilterCriteria{VolumeName: "zfs-pool/33"}
				snaps, err := repository.NewSnapshot(m, namestrategy).ListFilter(&filter)

				So(len(snaps), ShouldEqual, 1)
				So(snaps[0].Name, ShouldEqual, "1")
				So(snaps[0].VolumeName, ShouldEqual, "zfs-pool/33")
				So(err, ShouldBeNil)
			})
			Convey("It should return a filtered list of snapshots with valid names", func() {
				m.EXPECT().List(gomock.Any()).Times(1).Return(fakeCommand(`NAME
zfs-pool@snap1
zfs-pool/33@snap2`))
				namestrategy.EXPECT().IsMatching("snap1").Return(false)
				namestrategy.EXPECT().IsMatching("snap2").Return(true)
				filter := domain.FilterCriteria{IgnoreInvalidSnapshotNames: true}
				snaps, err := repository.NewSnapshot(m, namestrategy).ListFilter(&filter)

				So(len(snaps), ShouldEqual, 1)
				So(snaps[0].Name, ShouldEqual, "snap2")
				So(snaps[0].VolumeName, ShouldEqual, "zfs-pool/33")
				So(err, ShouldBeNil)
			})
		})
	})
}

func fakeCommand(output string) *exec.Cmd {
	cmd := exec.Command("cat", "-")
	reader := strings.NewReader(output)
	cmd.Stdin = reader
	return cmd
}
