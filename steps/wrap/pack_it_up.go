package wrap

import (
	"archive/tar"
	utils "github.com/thriqon/involucro/lib"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func random_tarball_file_name() string {
	dir := os.TempDir()
	tarid := utils.RandomIdentifier()
	return filepath.Join(dir, "involucro-volume-"+tarid+".tar")
}

type TemporaryFile struct {
	Filename string
	Close func()
}

func pack_it_up(source_directory string, prefix string) (string, error) {
	tarfile_name := random_tarball_file_name()
	tarfile, err := os.Create(tarfile_name)
	if err != nil {
		return "", err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source_directory)
	if err != nil {
		return "", err
	}

	err := filepath.Walk(source_directory, func(os_path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		without_path_prefix := strings.TrimPrefix(source_directory, info.Name())
		as_slash_path := filepath.ToSlash(without_path_prefix)
		with_new_prefix := path.Join(prefix, as_slash_path)

		header, err := tar.FileInfoHeader(info, with_new_prefix)
		if err != nil {
			return err
		}
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
	if err != nil {
		return 
}
