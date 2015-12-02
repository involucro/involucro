package run

import (
	"bufio"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"io"
	"regexp"
)

type dockerLogsProvider interface {
	Logs(docker.LogsOptions) error
}

func readAndMatchAgainst(r io.Reader, re *regexp.Regexp, val chan error, channel string) {
	runeReader := bufio.NewReader(r)
	if re.MatchReader(runeReader) {
		log.WithFields(log.Fields{"regex": re.String(), "channel": channel}).Debug("Regex matched")
		val <- nil
	} else {
		log.WithFields(log.Fields{"regex": re.String(), "channel": channel}).Warn("Regex not matched")
		val <- errors.New("Regex not matched")
	}
}

func outputLogLines(r io.Reader, val chan error, channel, container string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		log.WithFields(log.Fields{"container": container}).Debug(channel + ": " + scanner.Text())
	}
	val <- scanner.Err()
}

func setupMatching(re *regexp.Regexp, ch chan error, channel string) *io.PipeWriter {
	pR, pW := io.Pipe()
	go readAndMatchAgainst(pR, re, ch, channel)
	return pW
}

func setupDump(ch chan error, channel, container string) *io.PipeWriter {
	pR, pW := io.Pipe()
	go outputLogLines(pR, ch, channel, container)
	return pW
}

func (img *ExecuteImage) loadAndProcessLogs(c dockerLogsProvider, containerID string) error {
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
