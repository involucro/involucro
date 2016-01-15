package runtime

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"io"
	"sync"
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
				log.WithFields(log.Fields{"error": err}).Warn("Decode log message error")
			case log.GetLevel() == log.DebugLevel:
				fmt.Printf("Pull Progress: %#v\r", m)
			}
		}
	}()

	pio := docker.PullImageOptions{
		Repository:    repositoryName,
		OutputStream:  pipeWriter,
		RawJSONStream: true,
	}
	log.WithFields(log.Fields{"repository": repositoryName}).Debug("Pull Image")
	return c.PullImage(pio, docker.AuthConfiguration{})
}
