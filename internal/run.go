package runtime

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/ilog"
	"github.com/thriqon/involucro/internal/translator"
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
	shortenedContainerID := container.ID[0:12]

	ilog.Debug.Logf("Created container [%s %s], starting it", shortenedContainerID, container.Name)

	if err = c.StartContainer(container.ID, nil); err != nil {
		ilog.Warn.Logf("Failed container [%s %s] not removed", shortenedContainerID, container.Name)
		return err
	}
	ilog.Debug.Logf("Container [%s %s] started, waiting for completion", shortenedContainerID, container.Name)

	if err := img.loadAndProcessLogs(c, container.ID); err != nil {
		return err
	}

	img.ActualCode, err = c.WaitContainer(container.ID)

	if err != nil {
		return err
	}

	if img.ActualCode != img.ExpectedCode {
		return fmt.Errorf("Unexpected exit code [%v] of container [%s %s], container preserved", img.ActualCode, shortenedContainerID, container.Name)
	}

	ilog.Debug.Logf("Container [%s %s] completed with exit code [%v] as expected", shortenedContainerID, container.Name, img.ActualCode)

	if err := c.RemoveContainer(docker.RemoveContainerOptions{
		ID:    container.ID,
		Force: true,
	}); err != nil {
		return err
	}
	ilog.Debug.Logf("Container [%s %s] removed", shortenedContainerID, container.Name)

	return nil
}

func absolutizeBinds(h docker.HostConfig, workDir string) (docker.HostConfig, error) {
	for ind, el := range h.Binds {
		parts := strings.Split(el, ":")
		if len(parts) != 2 {
			return h, fmt.Errorf("Invalid bind specification [%s], has to be of the form: source:dest", el)
		}

		if !path.IsAbs(parts[0]) {
			h.Binds[ind] = path.Join(workDir, parts[0]) + ":" + parts[1]
		}
	}
	return h, nil
}

func (img executeImage) ShowStartInfo() {
	logTask.Logf("Run image [%s] with command [%s]", img.Config.Image, img.Config.Cmd)
}

func (img executeImage) createContainer(c *docker.Client) (container *docker.Container, err error) {
	return createContainer(c, img.Config, img.HostConfig)
}

func createContainer(c *docker.Client, config docker.Config, hostConfig docker.HostConfig) (*docker.Container, error) {
	containerName := "step-" + randomIdentifier()

	opts := docker.CreateContainerOptions{
		Name:       containerName,
		Config:     &config,
		HostConfig: &hostConfig,
	}

	ilog.Debug.Logf("Creating container [%s]", containerName)

	container, err := c.CreateContainer(opts)
	if err == nil {
		return container, nil
	}

	if err != docker.ErrNoSuchImage {
		return nil, err
	}

	ilog.Debug.Logf("Image [%s] not present, pulling it", config.Image)
	if err := pull(c, config.Image); err != nil {
		return nil, err
	}

	return c.CreateContainer(opts)
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
		ilog.Warn.Logf("Overwriting the image set in .using() in .withConfig() is discouraged")
	} else {
		ubs.Config.Image = oldImageID
	}
	return ubs.usingTable(l)
}

func (ubs usingBuilderState) withHostConfig(l *lua.State) int {
	ubs.HostConfig = translator.ParseHostConfigFromLuaTable(l, ubs.HostConfig)
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
	input  chan string
	logger ilog.Logger
}

func (lw logWriter) run(wg *sync.WaitGroup) {
	defer wg.Done()

	for line := range lw.input {
		lw.logger.Logf("%s", line)
	}
}

func newLogWriter(logger ilog.Logger, wg *sync.WaitGroup) chan string {
	lw := logWriter{
		input:  make(chan string),
		logger: logger,
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
			ilog.Debug.Logf("Regex matched for %s", rcw.channel)
			return
		}
	}
	rcw.err = fmt.Errorf("Regex [%s] for %s not mached", rcw.re.String(), rcw.channel)
}

func (img *executeImage) loadAndProcessLogs(c dockerLogsProvider, containerID string) error {

	var wg sync.WaitGroup

	wg.Add(2)

	stdoutChans := []chan string{newLogWriter(logStdout, &wg)}
	stderrChans := []chan string{newLogWriter(logStderr, &wg)}

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

	return nil
}
