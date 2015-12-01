package pull

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"io"
)

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
func Pull(c *docker.Client, repositoryName string) error {
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		dec := json.NewDecoder(pipeReader)
		for dec.More() {
			var m progress
			if err := dec.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				t := err.(*json.UnmarshalTypeError)
				log.WithFields(log.Fields{"value": t.Value, "notAssignableTo": t.Type}).Warn("Decode log message error")
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
	err := c.PullImage(pio, docker.AuthConfiguration{})

	return err
}
