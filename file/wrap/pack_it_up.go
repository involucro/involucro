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
		withNewPrefix := preparePathForTarHeader(os_path, sourceDirectory, prefix)

		var symlinkTarget string
		if info.Mode()&os.ModeSymlink > 0 {
			symlinkOsTarget, err := os.Readlink(os_path)
			if err != nil {
				return err
			}

			symlinkTarget = preparePathForTarHeader(symlinkOsTarget, sourceDirectory, prefix)
		}

		header, err := tar.FileInfoHeader(info, symlinkTarget)
		if err != nil {
			return err
		}
		header.Name = withNewPrefix
		if err := tarball.WriteHeader(header); err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
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

func preparePathForTarHeader(filename string, sourceDir, prefix string) string {
	prefixWithoutLeadingSlash := strings.TrimPrefix(prefix, "/")

	slashed := filepath.ToSlash(filename)

	return rebaseFilename(sourceDir, prefixWithoutLeadingSlash, slashed)
}

func rebaseFilename(oldprefix, newprefix string, filename string) string {
	withoutOld := strings.TrimPrefix(filename, oldprefix)
	if withoutOld == filename {
		return filename
	}

	return path.Join(newprefix, withoutOld)
}
