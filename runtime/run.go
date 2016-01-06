package runtime

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/runtime/translator"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
)

// executeImage executes the given config and host config, similar to "docker
// run"
type executeImage struct {
	Config                docker.Config
	HostConfig            docker.HostConfig
	ExpectedCode          int
	ExpectedStdoutMatcher *regexp.Regexp
	ExpectedStderrMatcher *regexp.Regexp
	ActualCode            int
}

func (img executeImage) WithRemoteDockerClient(c *docker.Client, remoteWorkDir string) error {
	if !path.IsAbs(remoteWorkDir) {
		remoteWorkDir = path.Join("/", remoteWorkDir)
	}
	return img.withAbsolutizedWorkDir(c, remoteWorkDir)
}

func (img executeImage) WithDockerClient(c *docker.Client, remoteWorkDir string) error {
	if !path.IsAbs(remoteWorkDir) {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		remoteWorkDir = path.Join(cwd, remoteWorkDir)
	}
	return img.withAbsolutizedWorkDir(c, remoteWorkDir)
}

func (img executeImage) withAbsolutizedWorkDir(c *docker.Client, remoteWorkDir string) error {
	var err error
	img.HostConfig, err = absolutizeBinds(img.HostConfig, remoteWorkDir)
	if err != nil {
		return err
	}

	container, err := img.createContainer(c)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{"ID": container.ID}).Debug("Container created, starting it")

	if err = c.StartContainer(container.ID, nil); err != nil {
		log.WithFields(log.Fields{"ID": container.ID, "err": err}).Warn("Container not started and not removed")
		return err
	}
	log.WithFields(log.Fields{"ID": container.ID}).Debug("Container started, await completion")

	err = img.loadAndProcessLogs(c, container.ID)
	if err != nil {
		return err
	}

	img.ActualCode, err = c.WaitContainer(container.ID)

	if img.ActualCode != img.ExpectedCode {
		log.WithFields(log.Fields{"ID": container.ID, "expected": img.ExpectedCode, "actual": img.ActualCode}).Error("Unexpected exit code, container not removed")
		return errors.New("Unexpected exit code")
	}

	log.WithFields(log.Fields{"Status": img.ActualCode, "ID": container.ID}).Debug("Execution complete")

	if err == nil && img.ActualCode == 0 {
		err := c.RemoveContainer(docker.RemoveContainerOptions{
			ID:    container.ID,
			Force: true,
		})
		if err != nil {
			log.WithFields(log.Fields{"ID": container.ID, "err": err}).Warn("Container not removed")
		} else {
			log.WithFields(log.Fields{"ID": container.ID}).Debug("Container removed")
		}
	} else {
		log.Debug("There was an error in execution or creation, container not removed")
	}

	return err
}

func absolutizeBinds(h docker.HostConfig, workDir string) (docker.HostConfig, error) {
	for ind, el := range h.Binds {
		parts := strings.Split(el, ":")
		if len(parts) != 2 {
			log.WithFields(log.Fields{"bind": el}).Error("Invalid bind, has to be of the form: source:dest")
			return h, errors.New("Invalid bind specification")
		}

		if !path.IsAbs(parts[0]) {
			h.Binds[ind] = path.Join(workDir, parts[0]) + ":" + parts[1]
		}
	}
	return h, nil
}

func (img executeImage) ShowStartInfo() {
	log.WithFields(log.Fields{"Image": img.Config.Image, "Cmd": img.Config.Cmd}).Info("run image")
}

func (img executeImage) createContainer(c *docker.Client) (container *docker.Container, err error) {
	containerName := "step-" + randomIdentifier()

	opts := docker.CreateContainerOptions{
		Name:       containerName,
		Config:     &img.Config,
		HostConfig: &img.HostConfig,
	}

	log.WithFields(log.Fields{"containerName": containerName}).Debug("Create Container")
	container, err = c.CreateContainer(opts)

	if err == docker.ErrNoSuchImage {
		if err = pull(c, img.Config.Image); err != nil {
			log.WithFields(log.Fields{"err": err}).Warn("pull failed")
			return
		}

		log.WithFields(log.Fields{"containerName": containerName}).Debug("Retry: Create Container")
		container, err = c.CreateContainer(opts)
	}

	if err != nil {
		log.WithFields(log.Fields{"image": img.Config.Image, "err": err}).Warn("create container failed")
	}
	return
}

type usingBuilderState struct {
	executeImage
	upper        fm
	registerStep func(Step)
}

func newRunSubBuilder(upper fm, register func(Step)) lua.Function {
	ubs := usingBuilderState{
		registerStep: register,
		upper:        upper,
		executeImage: executeImage{
			HostConfig: docker.HostConfig{
				Binds: []string{
					"./:/source",
				},
			},
		},
	}

	return ubs.using
}

func (ubs usingBuilderState) usingTable(l *lua.State) int {
	return tableWith(l, ubs.upper, fm{
		"run":             ubs.usingRun,
		"withExpectation": ubs.usingWithExpectation,
		"withConfig":      ubs.withConfig,
		"withHostConfig":  ubs.withHostConfig,
	})
}

func (ubs usingBuilderState) using(l *lua.State) int {
	ubs.Config.Image = lua.CheckString(l, -1)
	return ubs.usingTable(l)
}

func (ubs usingBuilderState) usingRun(l *lua.State) int {
	ubs.Config.Cmd = argumentsToStringArray(l)
	if ubs.Config.WorkingDir == "" {
		ubs.Config.WorkingDir = "/source"
	}

	ubs.registerStep(ubs.executeImage)
	return ubs.usingTable(l)
}

func (ubs usingBuilderState) usingWithExpectation(l *lua.State) int {
	if l.Top() != 1 {
		lua.Errorf(l, "expected exactly one argument to 'withExpectation'")
		panic("unreachable")
	}
	lua.ArgumentCheck(l, l.IsTable(-1), 1, "Expected table as argument")

	l.Field(-1, "code")
	if !l.IsNil(-1) {
		ubs.ExpectedCode = lua.CheckInteger(l, -1)
		log.WithFields(log.Fields{"code": ubs.ExpectedCode}).Debug("Expecting code")
	}
	l.Pop(1)

	l.Field(-1, "stdout")
	if !l.IsNil(-1) {
		str := lua.CheckString(l, -1)
		if regex, err := regexp.Compile(str); err != nil {
			lua.ArgumentError(l, 1, "invalid regular expression in stdout: "+err.Error())
			panic("unreachable")
		} else {
			ubs.ExpectedStdoutMatcher = regex
		}
	}
	l.Pop(1)

	l.Field(-1, "stderr")
	if !l.IsNil(-1) {
		str := lua.CheckString(l, -1)
		if regex, err := regexp.Compile(str); err != nil {
			lua.ArgumentError(l, 1, "invalid regular expression in stderr: "+err.Error())
			panic("unreachable")
		} else {
			ubs.ExpectedStderrMatcher = regex
		}
	}
	l.Pop(1)

	return ubs.usingTable(l)
}

func (ubs usingBuilderState) withConfig(l *lua.State) int {
	oldImageID := ubs.Config.Image
	ubs.Config = translator.ParseImageConfigFromLuaTable(l)
	if ubs.Config.Image != "" {
		log.Warn("Overwriting the used image in withConfig is discouraged")
	} else {
		ubs.Config.Image = oldImageID
	}
	return ubs.usingTable(l)
}

func (ubs usingBuilderState) withHostConfig(l *lua.State) int {
	ubs.HostConfig = translator.ParseHostConfigFromLuaTable(l)
	return ubs.usingTable(l)
}

func argumentsToStringArray(l *lua.State) (args []string) {
	top := l.Top()
	args = make([]string, top)
	for i := 1; i <= top; i++ {
		args[i-1] = lua.CheckString(l, i)
	}
	return
}

// ErrRegexNotMatched indicates that the given stream does not match the supplied regular expression.
type ErrRegexNotMatched struct {
	Channel string
}

func (e ErrRegexNotMatched) Error() string {
	return fmt.Sprintf("Regex not matched for channel %s", e.Channel)
}

func outputLogLines(r io.Reader, errCh chan error, channel, container string, logger *log.Logger) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		logger.WithFields(log.Fields{"container": container}).Debug(channel + ": " + scanner.Text())
	}
	errCh <- scanner.Err()
}

func readAndMatchAgainst(r io.Reader, re *regexp.Regexp, val chan error, channel string) {
	runeReader := bufio.NewReader(r)
	if re.MatchReader(runeReader) {
		log.WithFields(log.Fields{"regex": re.String(), "channel": channel}).Debug("Regex matched")
		val <- nil
	} else {
		log.WithFields(log.Fields{"regex": re.String(), "channel": channel}).Warn("Regex not matched")
		val <- ErrRegexNotMatched{channel}
	}
}

func setupMatching(re *regexp.Regexp, ch chan error, channel string) *io.PipeWriter {
	pR, pW := io.Pipe()
	go readAndMatchAgainst(pR, re, ch, channel)
	return pW
}

func setupDump(ch chan error, channel, container string) *io.PipeWriter {
	pR, pW := io.Pipe()
	go outputLogLines(pR, ch, channel, container, log.StandardLogger())
	return pW
}

type dockerLogsProvider interface {
	Logs(docker.LogsOptions) error
}

func (img *executeImage) loadAndProcessLogs(c dockerLogsProvider, containerID string) error {
	var stdOutWriters, stdErrWriters []io.Writer
	var closers []io.Closer

	stdoutErrorChan := make(chan error, 1)
	stderrErrorChan := make(chan error, 1)
	stdoutDumpChan := make(chan error, 1)
	stderrDumpChan := make(chan error, 1)

	if img.ExpectedStdoutMatcher != nil {
		w := setupMatching(img.ExpectedStdoutMatcher, stdoutErrorChan, "stdout")
		closers = append(closers, w)
		defer w.Close()
		stdOutWriters = append(stdOutWriters, w)
	} else {
		stdoutErrorChan <- nil
	}

	if img.ExpectedStderrMatcher != nil {
		w := setupMatching(img.ExpectedStderrMatcher, stderrErrorChan, "stderr")
		closers = append(closers, w)
		defer w.Close()
		stdErrWriters = append(stdErrWriters, w)
	} else {
		stderrErrorChan <- nil
	}

	{
		w := setupDump(stdoutDumpChan, "stdout", containerID)
		closers = append(closers, w)
		defer w.Close()
		stdOutWriters = append(stdOutWriters, w)
	}
	{
		w := setupDump(stderrDumpChan, "stderr", containerID)
		closers = append(closers, w)
		defer w.Close()
		stdErrWriters = append(stdErrWriters, w)
	}

	lOpt := docker.LogsOptions{
		Container: containerID,
		Follow:    true,
		Stdout:    true,
		Stderr:    true,
	}

	if len(stdOutWriters) > 0 {
		lOpt.OutputStream = io.MultiWriter(stdOutWriters...)
	}
	if len(stdErrWriters) > 0 {
		lOpt.ErrorStream = io.MultiWriter(stdErrWriters...)
	}
	if err := c.Logs(lOpt); err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("Log retrieval failed")
		return err
	}

	for _, w := range closers {
		w.Close()
	}

	chans := [...]chan error{stdoutErrorChan, stderrErrorChan, stdoutDumpChan, stderrDumpChan}
	for _, c := range chans {
		x := <-c
		if x != nil {
			log.WithFields(log.Fields{"error": x}).Warn("Log processing failed")
			return x
		}
	}
	log.Debug("Logs processed.")
	return nil
}
