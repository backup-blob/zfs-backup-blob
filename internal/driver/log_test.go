package driver_test

import (
	"bytes"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/driver"
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
	"testing"
)

func TestSpecLog(t *testing.T) {
	re := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(.*)")
	Convey("Given i call the Debugf function", t, func() {
		Convey("When the debug level is debug", func() {
			Convey("It should log", func() {
				buf := bytes.NewBufferString("")
				driver.NewLog(buf, domain.DebugLevel).Debugf("%s", "hello")

				So(re.Match(buf.Bytes()), ShouldBeTrue)
			})
		})
		Convey("When the debug level is info", func() {
			Convey("It should not log", func() {
				buf := bytes.NewBufferString("")
				driver.NewLog(buf, domain.InfoLevel).Debugf("%s", "hello")

				So(re.Match(buf.Bytes()), ShouldBeFalse)
			})
		})
	})
	Convey("Given i call the Infof function", t, func() {
		Convey("When the debug level is info", func() {
			Convey("It should log", func() {
				buf := bytes.NewBufferString("")
				driver.NewLog(buf, domain.InfoLevel).Infof("%s", "hello")

				So(re.Match(buf.Bytes()), ShouldBeTrue)
			})
		})
	})
}
