package testing

import (
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
	"testing"
)

func shouldAccept(actual interface{}, expected ...interface{}) string {
	regex := actual.(*regexp.Regexp)
	if regex == nil {
		return "Regular expression is nil"
	}
	for _, x := range expected {
		s := x.(string)
		if !regex.MatchString(s) {
			return regex.String() + " did not accept " + s + ", but it should!"
		}
	}
	return ""
}

func shouldNotAccept(actual interface{}, expected ...interface{}) string {
	if s := shouldAccept(actual, expected...); s == "" {
		return actual.(*regexp.Regexp).String() + " accepted, but it shouldnt't!"
	} else {
		return ""
	}
}

func TestAcceptAssertion(t *testing.T) {
	Convey("When I use the regex /ttt.../", t, func() {
		regex := regexp.MustCompile("ttt.[0-9].")
		Convey("Then the assertion accepts the example strings", func() {
			So(shouldAccept(regex, "ttta8a", "ttt.2."), ShouldResemble, "")
		})
		Convey("Then the assertion rejects the example strings", func() {
			So(shouldAccept(regex, "ttta8a", "ttt.2"), ShouldNotResemble, "")
			So(regex, shouldNotAccept, "ttta8a", "ttt.2")
		})
		Convey("Then the assertion rejects a nil regex", func() {
			var empty *regexp.Regexp
			So(shouldAccept(empty, "ttta8a", "ttt.2."), ShouldNotResemble, "")
			So(empty, shouldNotAccept, "ttta8a")
		})
	})
}
