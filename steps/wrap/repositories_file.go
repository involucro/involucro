package wrap

import (
	"archive/tar"
	"encoding/json"
)

func repositoriesFile(newRepositoryName, tag, id string) (tar.Header, []byte) {
	topMap := make(map[string]map[string]string)
	topMap[newRepositoryName] = make(map[string]string)
	topMap[newRepositoryName][tag] = id
	val, _ := json.Marshal(topMap)

	repositoriesFileHeader := tar.Header{
		Name:     "repositories",
		Typeflag: tar.TypeReg,
		Size:     int64(len(val)),
	}

	return repositoriesFileHeader, val
}
