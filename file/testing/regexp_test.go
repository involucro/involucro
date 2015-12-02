package testing

import (
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
	"testing"
)

func TestExpectedRegexpBehavior(t *testing.T) {
	Convey("The regex /asd/ ", t, func() {
		regex := regexp.MustCompile("asd")
		Convey("matches the string tasda", func() {
			So(regex, shouldAccept, "tasda")
		})
		Convey("doesn't match the string tada", func() {
			So(regex, shouldNotAccept, "tada")
		})
	})
	Convey("The regex /(?i)asd/ matches the string tASDAt", t, func() {
		regex := regexp.MustCompile("(?i)asd")
		So(regex, shouldAccept, "tASDAt")
	})

	Convey("The regex /(?m)asd/", t, func() {
		regex := regexp.MustCompile("(?m)asd")
		Convey("matches tasdt", func() {
			So(regex, shouldAccept, "tasdt")
		})
	})
	Convey("The regex /^asd/", t, func() {
		regex := regexp.MustCompile("^asd")
		Convey("doesn't match the example multi line string", func() {
			So(regex, shouldNotAccept, "hiuq\nasdppp\nlkl")
		})
	})
	Convey("The regex /(?m)^asd/", t, func() {
		regex := regexp.MustCompile("(?m)^asd")
		Convey("matches the example multi line string", func() {
			So(regex, shouldAccept, "hiuq\nasdppp\nlkl")
		})
	})
}
