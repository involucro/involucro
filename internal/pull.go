package runtime

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/fsouza/go-dockerclient"
	"github.com/involucro/involucro/auth"
	"github.com/involucro/involucro/ilog"
)

type pullimager interface {
	PullImage(docker.PullImageOptions, docker.AuthConfiguration) error
}

type progress struct {
	Status         string `json:"status"`
	Progress       string `json:"-"`
	ProgressDetail struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	} `json:"progressDetail"`
	ErrorMessage string `json:"error,omitempty"`
}

// pull pulls the image with the given identifier from
// the repository
func pull(c pullimager, repositoryName string) error {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	pipeReader, pipeWriter := io.Pipe()
	defer pipeWriter.Close()

	go func() {
		defer wg.Done()
		dec := json.NewDecoder(pipeReader)
		for dec.More() {
			var m progress
			err := dec.Decode(&m)
			switch {
			case err == io.EOF:
				return
			case err != nil:
				ilog.Warn.Logf("Decode log message error: %s", err)
			default:
				logProgress.Logf("%v", m)
			}
		}
	}()

	ac, foundAuthentication, err := auth.ForServer(serverOfRepo(repositoryName))
	if err != nil {
		return err
	}

	pio := docker.PullImageOptions{
		Repository:    repositoryName,
		OutputStream:  pipeWriter,
		RawJSONStream: true,
	}
	ilog.Debug.Logf("Pull Image [%s]", repositoryName)
	if err := c.PullImage(pio, ac); err != nil {
		if !foundAuthentication {
			ilog.Warn.Logf("Pull may have failed due to missing authentication information in ~/.involucro")
		}
		return err
	}
	return nil
}
