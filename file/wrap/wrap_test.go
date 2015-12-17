package wrap

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRandomFileName(t *testing.T) {
	filename := randomTarballFileName()
	if !strings.Contains(filename, "involucro") {
		t.Errorf("Didn't contain involucro: %s", filename)
	}
	otherFilename := randomTarballFileName()
	if otherFilename == filename {
		t.Errorf("Other filename is not different from the original: %s == %s", otherFilename, filename)
	}

	if _, err := os.Stat(filename); err == nil {
		t.Errorf("Stat succeeded, file shouldn't exist: %s", filename)
	}

	if info, err := os.Stat(filepath.Dir(filename)); err != nil {
		t.Errorf("Parent failed to stat: %s", err)
	} else {
		if !info.IsDir() {
			t.Errorf("Parent should be a directory")
		}
	}
}
