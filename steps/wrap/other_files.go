package wrap

import (
	"archive/tar"
	"path"
)

func versionFile(imageID string) (versionHeader tar.Header, contents []byte) {
	contents = []byte("1.0")

	versionHeader = tar.Header{
		Name:     path.Join(imageID, "VERSION"),
		Typeflag: tar.TypeReg,
		Size:     int64(len(contents)),
	}
	return
}

func imageDir(imageID string) tar.Header {
	return tar.Header{
		Name:     imageID + "/",
		Typeflag: tar.TypeDir,
	}
}
