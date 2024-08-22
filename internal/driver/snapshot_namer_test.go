package driver

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type mockNower struct {
	add int
}

func (n *mockNower) Now() time.Time {
	return time.Unix(1712939306, 10).Add(time.Second * time.Duration(n.add))
}

func TestSpecNamer(t *testing.T) {

	Convey("Given i call the default naming strategy is instantiated", t, func() {
		strategy := NewDefaultNamer(&mockNower{})
		Convey("When i call GetName", func() {
			Convey("I expect to receive a snapshot name", func() {
				name := NewDefaultNamer(&mockNower{}).GetName()

				So(name, ShouldEqual, "backup_blob_2024-04-12T16-28-26Z")
			})
		})
		Convey("When i call IsMatching and the name is matching", func() {
			Convey("I expect to receive true", func() {
				result := strategy.IsMatching(strategy.GetName())

				So(result, ShouldEqual, true)
			})
		})
		Convey("When i call IsMatching and the name is NOT matching", func() {
			Convey("I expect to receive false", func() {
				result := strategy.IsMatching("bacp_blob_44444")

				So(result, ShouldEqual, false)
			})
		})
		Convey("Given i call IsGreater", func() {
			Convey("When i have two snapshots where the first one is lex greater then the second", func() {
				Convey("I expect to receive true", func() {
					strategyCustom := NewDefaultNamer(&mockNower{add: 10})

					result := strategy.IsGreater(strategyCustom.GetName(), strategy.GetName())

					So(result, ShouldEqual, true)
				})
			})
		})
	})
}
