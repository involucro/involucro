package wrap

import (
	"archive/tar"
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestPackItUp(t *testing.T) {
	Convey("Given I have a prepared temp dir", t, func() {
		dir, err := ioutil.TempDir("", "involucro-test-wrap-packitup")
		So(err, ShouldBeNil)
		defer os.RemoveAll(dir)

		os.MkdirAll(filepath.Join(dir, "a", "b", "c", "d"), 0777)
		ioutil.WriteFile(filepath.Join(dir, "a", "b", "asd"), []byte("ASD"), 0777)

		Convey("When I pack that up into a buffer with prefix 'blubb'", func() {
			var buf bytes.Buffer
			err := packItUp(dir, &buf, "blubb")
			Convey("Then no error occurred", func() {
				So(err, ShouldBeNil)
			})

			Convey("When I read that with archive/tar", func() {
				tarReader := tar.NewReader(&buf)
				Convey("Then I find all the directories, e.g. blubb/a, blubb/a/b...", func() {
					seen := make(map[string]struct{})
					for {
						header, err := tarReader.Next()
						if err == io.EOF {
							break
						}
						if header.Typeflag == tar.TypeDir {
							seen[header.Name] = struct{}{}
						}
					}

					So(seen, ShouldContainKey, "blubb")
					So(seen, ShouldContainKey, "blubb/a")
					So(seen, ShouldContainKey, "blubb/a/b")
					So(seen, ShouldContainKey, "blubb/a/b/c")
					So(seen, ShouldContainKey, "blubb/a/b/c/d")
				})
				Convey("Then I can find and read blubb/a/b/asd", func() {
					var contents []byte
					for {
						header, err := tarReader.Next()
						if err == io.EOF {
							break
						}
						if header.Name == "blubb/a/b/asd" {
							var err error
							contents, err = ioutil.ReadAll(tarReader)
							So(err, ShouldBeNil)
						}
					}
					So(contents, ShouldResemble, []byte("ASD"))
				})
			})
		})
	})
	Convey("When I give it an not existing directory", t, func() {
		var buf bytes.Buffer
		err := packItUp("/not_existing", &buf, "blubb")
		Convey("Then it returns an error", func() {
			So(err, ShouldNotBeNil)
		})
	})
}
