package wrap

import (
	"archive/tar"
	"encoding/json"
	"github.com/fsouza/go-dockerclient"
	"path"
)

func imageConfigFile(parentId, imageId string) (tar.Header, []byte) {
	imageConfig, err := json.Marshal(docker.Image{
		ID:      imageId,
		Parent:  parentId,
		Comment: "Create with involucro 0.1",
	})
	if err != nil {
		panic(err)
	}

	imageConfigHeader := tar.Header{
		Name:     path.Join(imageId, "json"),
		Typeflag: tar.TypeReg,
		Size:     int64(len(imageConfig)),
	}

	return imageConfigHeader, imageConfig
}
