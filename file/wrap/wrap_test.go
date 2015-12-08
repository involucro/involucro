package wrap

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"path/filepath"
	"testing"
)

func TestRandomFileName(t *testing.T) {
	Convey("Given a random file name", t, func() {
		filename := randomTarballFileName()
		Convey("Then it contains the label 'involucro'", func() {
			So(filename, ShouldContainSubstring, "involucro")
		})
		Convey("When I generate another file name", func() {
			otherFilename := randomTarballFileName()
			Convey("Then it is a different path", func() {
				So(otherFilename, ShouldNotResemble, filename)
			})
		})

		Convey("Then it doesn't exist yet", func() {
			_, err := os.Stat(filename)
			So(err, ShouldNotBeNil)
		})
		Convey("Then is has an existing parent directory", func() {
			info, err := os.Stat(filepath.Dir(filename))
			So(err, ShouldBeNil)
			So(info.IsDir(), ShouldBeTrue)
		})
	})
}
