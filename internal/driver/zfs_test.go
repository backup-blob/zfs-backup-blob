package driver_test

import (
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/driver"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestSpec(t *testing.T) {

	Convey("Given zfs is instantiated by a config", t, func() {
		Convey("when everything works as expected", func() {
			Convey("it should return a instance", func() {
				zfsConfig := config.ZfsConfig{ZfsPath: "zfs2"}
				result := driver.NewZfsFromConfig(&zfsConfig, mocks.NewMockLogger())

				So(result, ShouldNotBeNil)
			})
		})
	})

	Convey("Given the zfs is instantiated", t, func() {
		var calls []string
		ctrl := driver.NewZfs("zfs", fakeExecCommand(&calls), mocks.NewMockLogger())
		Convey("When i call the send function with all parameters", func() {
			Convey("The zfs cli should be called with all parameters", func() {
				cmd := ctrl.Send(&domain.SendParameters{
					WithParameters:       true,
					PreviousSnapshotName: "snap-previous",
					SnapshotName:         "snap",
				})

				_, err := cmd.Output()

				So(calls[0], ShouldEqual, "zfs send --raw -p -I snap-previous snap")
				So(err, ShouldBeNil)
			})
		})
		Convey("When i call the setField function", func() {
			Convey("It should set a field", func() {
				err := ctrl.SetField("field1", "value", "volume1/deep2")

				So(err, ShouldBeNil)
				So(calls[0], ShouldEqual, "zfs set field1=value volume1/deep2")
			})
		})
		Convey("When i call the Destroy function", func() {
			Convey("It should set a field", func() {
				err := ctrl.Destroy("volume1/deep2")

				So(err, ShouldBeNil)
				So(calls[0], ShouldEqual, "zfs destroy volume1/deep2")
			})
		})
		Convey("When i call the getField function", func() {
			Convey("It should get a field", func() {
				cmd := ctrl.GetField("field1", "volume1/deep2")

				_, err := cmd.Output()

				So(err, ShouldBeNil)
				So(calls[0], ShouldEqual, "zfs get -o value -H field1 volume1/deep2")
			})
		})
		Convey("When i call the receive function with all parameters", func() {
			Convey("The zfs cli should be called with all parameters", func() {
				cmd := ctrl.Receive(&domain.ReceiveParameters{
					TargetName: "target",
				})

				_, err := cmd.Output()

				So(calls[0], ShouldEqual, "zfs receive -u target")
				So(err, ShouldBeNil)
			})
		})
		Convey("When i call the snapshot function with all parameters", func() {
			Convey("The zfs cli should be called with all parameters", func() {
				cmd := ctrl.Snapshot("snap-name")

				_, err := cmd.Output()

				So(calls[0], ShouldEqual, "zfs snapshot snap-name")
				So(err, ShouldBeNil)
			})
		})
		Convey("When i call the list function with all parameters", func() {
			Convey("The zfs cli should be called with all parameters", func() {
				cmd := ctrl.List(&domain.ListParameters{
					Type:   []string{"snapshot", "volume"},
					Fields: []string{"name"},
				})

				_, err := cmd.Output()

				So(calls[0], ShouldEqual, "zfs list -t snapshot,volume -o name")
				So(err, ShouldBeNil)
			})
		})
	})
}

func fakeExecCommand(calls *[]string) func(string, ...string) *exec.Cmd {
	return func(command string, args ...string) *exec.Cmd {
		*calls = append(*calls, command+" "+strings.Join(args, " "))
		cs := []string{"-test.run=TestHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}
