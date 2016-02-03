package runtime

// **Wrap** provides functionality to encapsulate local or remote directories
// as descendants of a base image, producing a new image with a name and
// optional configuration.
//
// Since Docker provides no native way to do this, we have to develop our
// own solution. The format to upload a finished image is documented
// in the remote API specification.

import (
	"archive/tar"
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/ilog"
	"github.com/thriqon/involucro/runtime/translator"
)

// ## The Builder
//
// The wrap builder is pretty complex, as it provides many methods to control
// the resulting image. It contains a definition for the wrap step, and the usual
// information for building.

type wrapBuilderState struct {
	asImage
	upper        fm
	registerStep func(Step)
}

func newWrapSubBuilder(upper fm, register func(Step)) lua.Function {
	wbs := wrapBuilderState{
		upper:        upper,
		registerStep: register,
	}
	return wbs.wrap
}

// `wrap` serves as introductory method for the wrap step.  It takes the
// parameter and force-coerces into the source directory string.
func (wbs wrapBuilderState) wrap(l *lua.State) int {
	wbs.SourceDir = lua.CheckString(l, -1)
	return wbs.wrapTable(l)
}

// inImage sets the base image for the new image.
func (wbs wrapBuilderState) inImage(l *lua.State) int {
	wbs.ParentImage = lua.CheckString(l, -1)
	return wbs.wrapTable(l)
}

// at specifies the location of the directory in the new image.
func (wbs wrapBuilderState) at(l *lua.State) int {
	wbs.TargetDir = lua.CheckString(l, -1)
	return wbs.wrapTable(l)
}

// as specifies the name of the new image.
func (wbs wrapBuilderState) as(l *lua.State) int {
	wbs.NewRepositoryName = lua.CheckString(l, -1)

	wbs.registerStep(wbs.asImage)
	return wbs.wrapTable(l)
}

// withConfig sets the config for the new image.
func (wbs wrapBuilderState) withConfig(l *lua.State) int {
	wbs.Config = translator.ParseImageConfigFromLuaTable(l)
	return wbs.wrapTable(l)
}

// wrapTable builds a Lua table containing the methods for the wrap step.
func (wbs wrapBuilderState) wrapTable(l *lua.State) int {
	return tableWith(l, wbs.upper, fm{
		"inImage":    wbs.inImage,
		"as":         wbs.as,
		"at":         wbs.at,
		"withConfig": wbs.withConfig,
	})
}

// ## Wrapping a Directory

// asImage represents packing up a directory into
// an image derived from another.
type asImage struct {
	SourceDir         string
	TargetDir         string
	ParentImage       string
	NewRepositoryName string
	Config            docker.Config
	ForbidRemoteRetry bool
}

func (img asImage) Take(i *Runtime) error {
	if i.isUsingRemoteInstance() {
		return img.executeRemotelyIn(i)
	}

	var taker func(*Runtime) error

	switch img.ParentImage {
	case "":
		taker = img.wrapWithoutBaseImageLocally
	default:
		taker = img.wrapWithBaseImageLocally
	}

	if err := taker(i); err != nil {
		if img.ForbidRemoteRetry {
			return err
		}
		// Retry using "remote" execution this fixes some permission problems, but
		// is generally less efficient
		ilog.Warn.Logf("Local execution errorred with error [%s], retrying with remote execution", err)
		return img.executeRemotelyIn(i)
	}
	return nil
}

func packInto(sourceDir, prefix string) (io.Reader, chan error) {
	uploadReader, uploadWriter := io.Pipe()
	errchan := make(chan error, 1)

	go func() {
		defer uploadWriter.Close()
		defer close(errchan)
		if err := packItUp(sourceDir, uploadWriter, prefix); err != nil {
			ilog.Debug.Logf("Packing failed due to %s", err)
			errchan <- err
		} else {
			ilog.Debug.Logf("Packing succeeded")
		}
	}()

	return uploadReader, errchan
}

func (img asImage) wrapWithoutBaseImageLocally(i *Runtime) error {
	c := i.client
	intermediateImageRepo := "image-" + randomIdentifier()

	sourceDir := img.SourceDir
	if !path.IsAbs(sourceDir) {
		sourceDir = path.Join(i.workDir, sourceDir)
	}

	packStream, errChan := packInto(sourceDir, img.TargetDir)

	importErr := c.ImportImage(docker.ImportImageOptions{
		Repository:  intermediateImageRepo,
		Tag:         "latest",
		Source:      "-",
		InputStream: packStream,
	})
	defer c.RemoveImage(intermediateImageRepo)

	// If there is an error during packing, report that as main error.  A packing
	// error always causes an import error as well.
	for el := range errChan {
		return el
	}
	if importErr != nil {
		return importErr
	}

	container, err := createContainer(c, docker.Config{Image: intermediateImageRepo, Cmd: []string{"/bin/sh"}}, docker.HostConfig{})
	if err != nil {
		return err
	}

	var opts docker.CommitContainerOptions
	opts.Repository, _, opts.Tag = repoNameAndTagFrom(img.NewRepositoryName)
	opts.Container = container.ID
	opts.Message = "Created with Involucro"
	opts.Run = &img.Config

	if _, err := c.CommitContainer(opts); err != nil {
		return err
	}
	return c.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
}

func (img asImage) wrapWithBaseImageLocally(i *Runtime) error {
	c := i.client

	container, err := createContainer(c, docker.Config{Image: img.ParentImage, Cmd: []string{"/bin/sh"}}, docker.HostConfig{})
	if err != nil {
		return err
	}

	sourceDir := img.SourceDir
	if !path.IsAbs(sourceDir) {
		sourceDir = path.Join(i.workDir, sourceDir)
	}

	packStream, errChan := packInto(sourceDir, img.TargetDir)

	uploadErr := c.UploadToContainer(container.ID, docker.UploadToContainerOptions{
		InputStream:          packStream,
		Path:                 "/",
		NoOverwriteDirNonDir: false,
	})
	// Treat pack error as 'main' error (a packing error always triggers an upload
	// error)
	for el := range errChan {
		return el
	}
	if uploadErr != nil {
		return uploadErr
	}

	var opts docker.CommitContainerOptions
	opts.Repository, _, opts.Tag = repoNameAndTagFrom(img.NewRepositoryName)
	opts.Container = container.ID
	opts.Message = "Created with Involucro"
	opts.Run = &img.Config

	if _, err := c.CommitContainer(opts); err != nil {
		return err
	}
	return c.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
}

func (img asImage) ShowStartInfo() {
	logTask.Logf("Wrap [%s] as [%s]", img.SourceDir, img.NewRepositoryName)
}

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
	sourceDir = filepath.ToSlash(sourceDir)

	return rebaseFilename(sourceDir, prefixWithoutLeadingSlash, slashed)
}

func rebaseFilename(oldprefix, newprefix string, filename string) string {
	withoutOld := strings.TrimPrefix(filename, oldprefix)
	if withoutOld == filename {
		return filename
	}

	return path.Join(newprefix, withoutOld)
}

func (img asImage) forRemoteExecution() Step {
	dockerSocket := "/var/run/docker.sock"

	origSourceDir := img.SourceDir
	img.SourceDir = "/source"

	encoded, err := json.Marshal(img)
	if err != nil {
		panic(err)
	}

	return executeImage{
		Config: docker.Config{
			Image: "involucro/tool:latest",
			Cmd:   []string{"--wrap", string(encoded)},
		},
		HostConfig: docker.HostConfig{
			Binds: []string{
				dockerSocket + ":/sock",
				origSourceDir + ":/source",
			},
		},
	}
}

// DecodeWrapStep unmarshals the wrap step encoded in in (in JSON), and gives
// it back. This function panics if there are errors during decoding.
func DecodeWrapStep(in string) Step {
	img := asImage{}
	if err := json.Unmarshal([]byte(in), &img); err != nil {
		panic(err)
	}
	return img
}

func (img asImage) executeRemotelyIn(i *Runtime) error {
	img.ForbidRemoteRetry = true
	return img.forRemoteExecution().Take(i)
}
