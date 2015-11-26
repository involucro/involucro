package lib

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRandomIds(t *testing.T) {
	Convey("Given a random identifier", t, func() {
		s := RandomIdentifier()
		Convey("Then it is a string", func() {
			So(s, ShouldHaveSameTypeAs, "asd")
		})
		Convey("When I generate another identifier", func() {
			t := RandomIdentifier()
			Convey("Then it is different from the original", func() {
				So(s, ShouldNotResemble, t)
			})
		})
	})
}

func TestRandomIdsWithLength(t *testing.T) {
	Convey("Given a random identifier of length 64", t, func() {
		s := RandomIdentifierOfLength(64)
		Convey("Then it has length 64", func() {
			So(s, ShouldHaveLength, 64)
		})
	})
}
