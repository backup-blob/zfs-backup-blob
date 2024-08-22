package throttle_test

import (
	"bytes"
	"github.com/backup-blob/zfs-backup-blob/pkg/throttle"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestSpecThrottle(t *testing.T) {
	Convey("Given throttle a writer", t, func() {
		Convey("When everything works as expected", func() {
			Convey("Speed of writing should be limited", func() {
				buf := bytes.NewBuffer([]byte{})

				writer, err := throttle.SpeedlimitWriter(1)(buf)
				So(err, ShouldBeNil)

				now := time.Now()
				n, err := writer.Write([]byte("h"))
				So(err, ShouldBeNil)
				timePassed := time.Now().Sub(now)

				So(n, ShouldEqual, 1)
				So(timePassed, ShouldBeGreaterThanOrEqualTo, 1*time.Second)
			})
		})
	})
	Convey("Given throttle a reader", t, func() {
		Convey("When everything works as expected", func() {
			Convey("Speed of the reader should be limited", func() {
				reader := bytes.NewReader([]byte("h"))
				readerThrottled, err := throttle.SpeedlimitReader(1)(reader)
				So(err, ShouldBeNil)

				now := time.Now()
				m := make([]byte, 1)
				n, err := readerThrottled.Read(m)
				So(err, ShouldBeNil)
				timePassed := time.Now().Sub(now)

				So(n, ShouldEqual, 1)
				So(timePassed, ShouldBeGreaterThanOrEqualTo, 1*time.Second)
			})
		})
	})
}
