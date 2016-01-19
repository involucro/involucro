package runtime

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/runtime/translator"
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

func (img executeImage) Take(i *Runtime) error {
	remoteWorkDir := i.workDir

	if !path.IsAbs(remoteWorkDir) {
		base := "/"
		if !i.isUsingRemoteInstance() {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			base = cwd
		}

		remoteWorkDir = path.Join(base, remoteWorkDir)
	}

	c := i.client

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

	if err := img.loadAndProcessLogs(c, container.ID); err != nil {
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

type dataAcceptorLineSender struct {
	cs []chan string
	pr *io.PipeReader
	pw *io.PipeWriter
	br *bufio.Reader
}

func newDataAcceptorLineSender(cs []chan string) dataAcceptorLineSender {
	pr, pw := io.Pipe()
	dals := dataAcceptorLineSender{cs, pr, pw, bufio.NewReader(pr)}
	go dals.run()
	return dals
}

func (dals dataAcceptorLineSender) Write(p []byte) (int, error) {
	return dals.pw.Write(p)
}

func (dals dataAcceptorLineSender) Close() error {
	dals.pw.Close()
	return dals.pr.Close()
}

func (dals dataAcceptorLineSender) run() {
	defer func() {
		for _, c := range dals.cs {
			close(c)
		}
		dals.Close()
	}()

	for {
		line, err := dals.br.ReadString('\n')
		if err == nil {
			line = strings.TrimSpace(line)
		}
		if err == nil || err == io.EOF {
			for _, c := range dals.cs {
				c <- line
			}
		}
		if err != nil {
			return
		}
	}
}

type dockerLogsProvider interface {
	Logs(docker.LogsOptions) error
}

type logWriter struct {
	input     chan string
	logger    *log.Logger
	channel   string
	container string
}

func (lw logWriter) run(wg *sync.WaitGroup) {
	defer wg.Done()

	for line := range lw.input {
		lw.logger.WithFields(log.Fields{"container": lw.container}).Debug(lw.channel + ": " + line)
	}
}

func newLogWriter(channel, container string, logger *log.Logger, wg *sync.WaitGroup) chan string {
	lw := logWriter{
		input:     make(chan string),
		logger:    logger,
		channel:   channel,
		container: container,
	}
	go lw.run(wg)
	return lw.input
}

type regexpCheckWriter struct {
	input     chan string
	re        *regexp.Regexp
	channel   string
	container string
	err       error
}

func (rcw *regexpCheckWriter) run(wg *sync.WaitGroup) {
	defer wg.Done()

	for line := range rcw.input {
		if rcw.re.MatchString(line) {
			log.WithFields(log.Fields{"regex": rcw.re.String(), "channel": rcw.channel, "container": rcw.container}).Debug("Regex matched")
			return
		}
	}
	log.WithFields(log.Fields{"regex": rcw.re.String(), "channel": rcw.channel, "container": rcw.container}).Warn("Regex not matched")
	rcw.err = ErrRegexNotMatched{rcw.channel}
}

func (img *executeImage) loadAndProcessLogs(c dockerLogsProvider, containerID string) error {

	var wg sync.WaitGroup

	wg.Add(2)

	stdoutChans := []chan string{newLogWriter("stdout", containerID, log.StandardLogger(), &wg)}
	stderrChans := []chan string{newLogWriter("stderr", containerID, log.StandardLogger(), &wg)}

	stdoutMatcher := regexpCheckWriter{make(chan string), img.ExpectedStdoutMatcher, "stdout", containerID, nil}
	stderrMatcher := regexpCheckWriter{make(chan string), img.ExpectedStderrMatcher, "stderr", containerID, nil}

	if img.ExpectedStdoutMatcher != nil {
		wg.Add(1)
		stdoutChans = append(stdoutChans, stdoutMatcher.input)
		go stdoutMatcher.run(&wg)
	}

	if img.ExpectedStderrMatcher != nil {
		wg.Add(1)
		stderrChans = append(stderrChans, stderrMatcher.input)
		go stderrMatcher.run(&wg)
	}
	dalsOutput := newDataAcceptorLineSender(stdoutChans)
	dalsError := newDataAcceptorLineSender(stderrChans)

	lOpt := docker.LogsOptions{
		Container:    containerID,
		Follow:       true,
		Stdout:       true,
		Stderr:       true,
		OutputStream: dalsOutput,
		ErrorStream:  dalsError,
	}

	if err := c.Logs(lOpt); err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("Log retrieval failed")
		return err
	}

	dalsOutput.Close()
	dalsError.Close()

	wg.Wait()

	if stdoutMatcher.err != nil {
		return stdoutMatcher.err
	}
	if stderrMatcher.err != nil {
		return stderrMatcher.err
	}

	log.Debug("Logs processed.")
	return nil
}
