package driver_test

import (
	"bytes"
	"github.com/backup-blob/zfs-backup-blob/internal/driver"
	"testing"
)
import (
	. "github.com/smartystreets/goconvey/convey"
)

func TestSpecRender(t *testing.T) {
	render := driver.NewRender()
	Convey("Given i call RenderBackupTable", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It should render table", func() {
				buf := bytes.NewBuffer([]byte{})

				render.RenderTable(buf, []interface{}{"Path", "Type", "Size"}, [][]interface{}{{"key", "full", "1MB"}})

				So(buf.String(), ShouldEqual, `+------+------+------+
| PATH | TYPE | SIZE |
+------+------+------+
| key  | full | 1MB  |
+------+------+------+
`)
			})
		})
	})
}
