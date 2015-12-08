package wrap

import (
	"archive/tar"
	"encoding/json"
	"github.com/fsouza/go-dockerclient"
	"path"
	"strings"
	"time"
)

func imageConfigFile(parentID, imageID string, containerConfig docker.Config) (tar.Header, []byte) {
	imageConfig, err := json.Marshal(docker.Image{
		ID:      imageID,
		Parent:  parentID,
		Comment: "Create with involucro 0.1",
		Created: time.Now(),
		Config:  &containerConfig,
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

func repositoriesFile(newRepositoryName, id string) (tar.Header, []byte) {
	topMap := make(map[string]map[string]string)

	repo, tag := repoNameAndTagFrom(newRepositoryName)

	topMap[repo] = make(map[string]string)
	topMap[repo][tag] = id
	val, _ := json.Marshal(topMap)

	repositoriesFileHeader := tar.Header{
		Name:     "repositories",
		Typeflag: tar.TypeReg,
		Size:     int64(len(val)),
	}

	return repositoriesFileHeader, val
}

func repoNameAndTagFrom(name string) (repo, tag string) {
	parts := strings.Split(name, ":")
	switch len(parts) {
	case 1:
		repo = parts[0]
		tag = "latest"
	case 2:
		repo = parts[0]
		tag = parts[1]
	default:
		panic("Invalid repository name")
	}
	return
}

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
