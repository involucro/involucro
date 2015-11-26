package wrap

import (
	"archive/tar"
	"path"
)

func versionFile(imageId string) (versionHeader tar.Header, contents []byte) {
	contents = []byte("1.0")

	versionHeader = tar.Header{
		Name:     path.Join(imageId, "VERSION"),
		Typeflag: tar.TypeReg,
		Size:     int64(len(contents)),
	}
	return
}

func imageDir(imageId string) tar.Header {
	return tar.Header{
		Name:     imageId + "/",
		Typeflag: tar.TypeDir,
	}
}
