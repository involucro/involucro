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
	"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/runtime/translator"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
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
}

func (img asImage) WithDockerClient(c *docker.Client, remoteWorkDir string) error {
	imageID := randomIdentifierOfLength(64)

	parentImageID, err := img.enforceParentImagePresenceAndGetId(c, img.ParentImage)
	if err != nil {
		return err
	}

	uploadReader, uploadWriter := io.Pipe()

	layerBallName := randomTarballFileName()

	layerFile, err := os.Create(layerBallName)
	if err != nil {
		return err
	}
	defer os.Remove(layerBallName)
	defer layerFile.Close()

	log.WithFields(log.Fields{"layerBallName": layerBallName}).Debug("Packing")

	sourceDir := img.SourceDir
	if !path.IsAbs(sourceDir) {
		sourceDir = path.Join(remoteWorkDir, sourceDir)
	}
	err = packItUp(sourceDir, layerFile, img.TargetDir)
	if err != nil {
		return err
	}
	layerFile.Close()

	var uploadErr error
	var wg sync.WaitGroup
	wg.Add(2)

	log.WithFields(log.Fields{"image_id": imageID, "repository": img.NewRepositoryName}).Debug("Wrapping up")

	go func() {
		defer wg.Done()
		defer uploadWriter.Close()
		writeUploadBallInto(uploadWriter, layerBallName, img.NewRepositoryName, parentImageID, imageID, img.Config)
		return
	}()

	go func(r io.Reader) {
		defer wg.Done()
		err := c.LoadImage(docker.LoadImageOptions{InputStream: r})
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("Failed loading image into Docker")
			uploadErr = err
		} else {
			log.Debug("Upload finished")
		}
	}(uploadReader)

	wg.Wait()

	return uploadErr
}

// enforceParentImagePresenceAndGetId gets the ID of the image with the given
// name, and pulls the image first, if neccessary.
func (img asImage) enforceParentImagePresenceAndGetId(c *docker.Client, parentName string) (string, error) {
	// If there is no image name present, we are done. This happens if the user
	// wants to pack the data into a parent-less image.
	if parentName == "" {
		return "", nil
	}

	// Fetch the config from Docker. If it works, we can return the ID.
	// If anything goes wrong, we go on.
	if parentImageConfig, err := c.InspectImage(parentName); err == nil {
		return parentImageConfig.ID, nil
	}

	log.WithFields(log.Fields{"image": parentName}).Debug("Parent image not found, pulling it")

	if err := pull(c, parentName); err != nil {
		log.WithFields(log.Fields{"image": parentName, "error": err}).Error("Pulling failed")
		return "", err
	}

	parentImageConfig, err := c.InspectImage(parentName)
	if err != nil {
		log.WithFields(log.Fields{"image": parentName, "error": err}).Error("Image still not found after pulling, bailing")
		return "", err
	}
	return parentImageConfig.ID, nil
}

func writeUploadBallInto(w io.Writer, layerBallName string, newRepositoryName string, parentImageID string, imageID string, config docker.Config) error {
	uploadBall := tar.NewWriter(w)
	defer uploadBall.Close()

	repositoriesFileHeader, repoFile := repositoriesFile(newRepositoryName, imageID)
	uploadBall.WriteHeader(&repositoriesFileHeader)
	uploadBall.Write(repoFile)

	imageDirHeader := imageDir(imageID)
	uploadBall.WriteHeader(&imageDirHeader)

	versionFileHeader, versionFileContents := versionFile(imageID)
	uploadBall.WriteHeader(&versionFileHeader)
	uploadBall.Write(versionFileContents)

	configFileHeader, configFileContents := imageConfigFile(parentImageID, imageID, config)
	uploadBall.WriteHeader(&configFileHeader)
	uploadBall.Write(configFileContents)

	info, err := os.Stat(layerBallName)
	if err != nil {
		return err
	}
	layerBallHeader := tar.Header{
		Name:     path.Join(imageID, "layer.tar"),
		Typeflag: tar.TypeReg,
		Size:     info.Size(),
	}
	uploadBall.WriteHeader(&layerBallHeader)
	layerBallFile, err := os.Open(layerBallName)
	if err != nil {
		return err
	}
	defer layerBallFile.Close()

	_, err = io.Copy(uploadBall, layerBallFile)

	log.Debug("Pipe finished")

	return err
}

func randomTarballFileName() string {
	dir := os.TempDir()
	tarid := randomIdentifier()
	return filepath.Join(dir, "involucro-volume-"+tarid+".tar")
}

func (img asImage) ShowStartInfo() {
	log.WithFields(log.Fields{"as": img.NewRepositoryName}).Info("wrap")
}

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

	repo, _, tag := repoNameAndTagFrom(newRepositoryName)

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

func init() {
	RegisterEncodeableType(asImage{})
}

func (img asImage) forRemoteExecution() Step {
	dockerSocket := "/var/run/docker.sock"

	origSourceDir := img.SourceDir
	img.SourceDir = "/source"

	steps := []Step{img}

	return executeImage{
		Config: docker.Config{
			Image: "involucro/tool:latest",
			Env:   []string{"STATE=" + EncodeState(steps)},
			Cmd:   []string{"--encoded-state", "--socket", "/sock"},
		},
		HostConfig: docker.HostConfig{
			Binds: []string{
				dockerSocket + ":/sock",
				origSourceDir + ":/source",
			},
		},
	}
}

func (img asImage) WithRemoteDockerClient(c *docker.Client, remoteWorkDir string) error {
	return img.forRemoteExecution().WithRemoteDockerClient(c, remoteWorkDir)
}
