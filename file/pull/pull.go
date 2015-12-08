package pull

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"io"
	"sync"
)

type Pullable interface {
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

// Pull pulls the image with the given identifier from
// the repository
func Pull(c Pullable, repositoryName string) error {
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
			if err := dec.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				if t, ok := err.(*json.UnmarshalTypeError); ok {
					log.WithFields(log.Fields{"value": t.Value, "notAssignableTo": t.Type}).Warn("Decode log message error")
				} else {
					log.WithFields(log.Fields{"error": err}).Warn("Decode log message error")
				}
			} else {
				log.WithFields(log.Fields{"message": m}).Debug("Progress")
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
