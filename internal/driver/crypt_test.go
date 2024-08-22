package driver_test

import (
	"bytes"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
	"github.com/backup-blob/zfs-backup-blob/internal/driver"
	"io"
	"strings"
	"testing"
)
import (
	. "github.com/smartystreets/goconvey/convey"
)

func TestSpecCrypt(t *testing.T) {
	conf := config.CryptConfig{Password: "hello"}
	Convey("Given i call Read", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It should encrypt data", func() {
				readerIn := strings.NewReader("hello")

				readerOut, err := driver.NewCrypt(&conf).Read(readerIn)
				encrypted, errR := io.ReadAll(readerOut)

				So(err, ShouldBeNil)
				So(errR, ShouldBeNil)
				So(encrypted, ShouldNotBeNil)
			})
		})
	})
	Convey("Given i call Write", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It should decrypt data", func() {
				data := "helloooo"
				readerIn := strings.NewReader(data)
				writerIn := bytes.NewBuffer([]byte{})

				readerOut, err := driver.NewCrypt(&conf).Read(readerIn)
				So(err, ShouldBeNil)
				writerOut, errW := driver.NewCrypt(&conf).Write(writerIn)
				So(errW, ShouldBeNil)
				_, errC := io.Copy(writerOut, readerOut)
				So(errC, ShouldBeNil)
				So(writerIn.String(), ShouldEqual, data)
			})
		})
	})
}
