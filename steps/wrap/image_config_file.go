package wrap

import (
	"archive/tar"
	"encoding/json"
	"github.com/fsouza/go-dockerclient"
	"path"
	"time"
)

func imageConfigFile(parentID, imageID string) (tar.Header, []byte) {
	imageConfig, err := json.Marshal(docker.Image{
		ID:      imageID,
		Parent:  parentID,
		Comment: "Create with involucro 0.1",
		Created: time.Now(),
	})
	if err != nil {
		panic(err)
	}

	imageConfigHeader := tar.Header{
		Name:     path.Join(imageID, "json"),
		Typeflag: tar.TypeReg,
		Size:     int64(len(imageConfig)),
	}

	return imageConfigHeader, imageConfig
}
