package wrap

import (
	"archive/tar"
	"encoding/json"
	"strings"
)

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
