package run

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"io"
	"regexp"
	"sync"
)

// ExecuteImage executes the given config and host config, similar to "docker
// run"
type ExecuteImage struct {
	Config                docker.Config
	HostConfig            docker.HostConfig
	ExpectedCode          int
	ExpectedStdoutMatcher *regexp.Regexp
	ExpectedStderrMatcher *regexp.Regexp
	ActualCode            int
}

// WithDockerClient executes the task on the given Docker instance
func (img ExecuteImage) WithDockerClient(c *docker.Client) error {

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

	img.loadAndProcessLogs(c, container)

	img.ActualCode, err = c.WaitContainer(container.ID)

	log.WithFields(log.Fields{"Status": img.ActualCode, "ID": container.ID}).Info("Execution complete")

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

func readAndMatchAgainst(r io.Reader, re *regexp.Regexp, wg *sync.WaitGroup, channel string) {
	defer wg.Done()
	runeReader := bufio.NewReader(r)
	if re.MatchReader(runeReader) {
		log.WithFields(log.Fields{"regex": re.String(), "channel": channel}).Debug("Regex matched")
	} else {
		log.WithFields(log.Fields{"regex": re.String(), "channel": channel}).Warn("Regex not matched")
	}
}

func outputLogLines(r io.Reader, wg *sync.WaitGroup, channel, container string) {
	defer wg.Done()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		log.WithFields(log.Fields{"container": container}).Debug(channel + ": " + scanner.Text())
	}
}

func (img *ExecuteImage) loadAndProcessLogs(c *docker.Client, container *docker.Container) error {
	var stdOutWriters, stdErrWriters []io.Writer
	var wg sync.WaitGroup

	defer wg.Wait()

	if img.ExpectedStdoutMatcher != nil {
		stdOutMatcherPipeReader, stdOutMatcherPipeWriter := io.Pipe()
		defer stdOutMatcherPipeWriter.Close()

		stdOutWriters = append(stdOutWriters, stdOutMatcherPipeWriter)

		wg.Add(1)
		go readAndMatchAgainst(stdOutMatcherPipeReader, img.ExpectedStdoutMatcher, &wg, "stdout")
	}
	if img.ExpectedStderrMatcher != nil {
		stdErrMatcherPipeReader, stdErrMatcherPipeWriter := io.Pipe()
		defer stdErrMatcherPipeWriter.Close()

		stdErrWriters = append(stdErrWriters, stdErrMatcherPipeWriter)

		wg.Add(1)
		go readAndMatchAgainst(stdErrMatcherPipeReader, img.ExpectedStderrMatcher, &wg, "stderr")
	}

	stdOutPrinterPipeReader, stdOutPrinterPipeWriter := io.Pipe()
	defer stdOutPrinterPipeWriter.Close()
	wg.Add(1)
	go outputLogLines(stdOutPrinterPipeReader, &wg, "stdout", container.ID)
	stdOutWriters = append(stdOutWriters, stdOutPrinterPipeWriter)

	stdErrPrinterPipeReader, stdErrPrinterPipeWriter := io.Pipe()
	defer stdErrPrinterPipeWriter.Close()
	wg.Add(1)
	go outputLogLines(stdErrPrinterPipeReader, &wg, "stderr", container.ID)
	stdErrWriters = append(stdErrWriters, stdErrPrinterPipeWriter)

	var lOpt docker.LogsOptions
	lOpt.Container = container.ID
	lOpt.Follow = true
	lOpt.Stdout = true
	lOpt.Stderr = true

	if len(stdOutWriters) > 0 {
		lOpt.OutputStream = io.MultiWriter(stdOutWriters...)
	}
	if len(stdErrWriters) > 0 {
		lOpt.ErrorStream = io.MultiWriter(stdErrWriters...)
	}
	if err := c.Logs(lOpt); err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("Log retrieval failed")
		return err
	} else {
		log.Debug("Logs processed.")
	}
	return nil
}
