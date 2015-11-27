package wrap

import (
	"archive/tar"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func packItUp(sourceDirectory string, tarfile io.Writer, prefix string) error {
	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	_, err := os.Stat(sourceDirectory)
	if err != nil {
		return err
	}

	return filepath.Walk(sourceDirectory, func(os_path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		withoutPathPrefix := strings.TrimPrefix(info.Name(), sourceDirectory)
		asSlashPath := filepath.ToSlash(withoutPathPrefix)
		prefixWithoutLeadingSlash := strings.TrimPrefix(prefix, "/")
		withNewPrefix := path.Join(prefixWithoutLeadingSlash, asSlashPath)

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = withNewPrefix
		if err := tarball.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(os_path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(tarball, file)
		return err
	})
}
