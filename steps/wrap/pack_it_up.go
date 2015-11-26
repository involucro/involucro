package wrap

import (
	"archive/tar"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func pack_it_up(source_directory string, tarfile io.Writer, prefix string) error {
	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	_, err := os.Stat(source_directory)
	if err != nil {
		return err
	}

	return filepath.Walk(source_directory, func(os_path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		without_path_prefix := strings.TrimPrefix(info.Name(), source_directory)
		as_slash_path := filepath.ToSlash(without_path_prefix)
		prefix_without_leading_slash := strings.TrimPrefix(prefix, "/")
		with_new_prefix := path.Join(prefix_without_leading_slash, as_slash_path)

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = with_new_prefix
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
