package wrap

import (
	utils "github.com/thriqon/involucro/lib"
	"os"
	"path/filepath"
)

func randomTarballFileName() string {
	dir := os.TempDir()
	tarid := utils.RandomIdentifier()
	return filepath.Join(dir, "involucro-volume-"+tarid+".tar")
}
