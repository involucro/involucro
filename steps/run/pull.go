package run

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/cheggaaa/pb"
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

type pblogger struct {
	bar *pb.ProgressBar
}

func (pb pblogger) Write(p []byte) (int, error) {
	var entry progress
	err := json.Unmarshal(p, &entry)
	log.WithFields(log.Fields{"err": err, "entry": entry, "raw": string(p)}).Info("Progress")
	return len(p), nil
}

func pull(c *docker.Client, repositoryName string) error {
	total := 100

	pipeReader, pipeWriter := io.Pipe()
	bar := pb.New(total)
	bar.ShowTimeLeft = false
	bar.ShowFinalTime = false
	defer func() {
		bar.ShowCounters = false
		bar.Prefix(repositoryName + ": Pull complete  ")
		bar.Total = 100
		bar.Set(100)
		bar.Update()
	}()

	go func() {
		if !log.IsTerminal() {
			return
		}

		bar.ShowCounters = true
		bar.SetUnits(pb.U_BYTES)
		bar.Total = 1
		bar.Prefix(repositoryName + ": Pulling...  ")

		bar.Format("[=>.]")

		dec := json.NewDecoder(pipeReader)
		for dec.More() {
			var m progress
			if err := dec.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				t := err.(*json.UnmarshalTypeError)
				log.WithFields(log.Fields{"value": t.Value, "notAssignableTo": t.Type}).Warn("Decode log message error")
			}
			bar.Prefix(repositoryName + ": " + m.Status + "  ")
			bar.Total = m.ProgressDetail.Total
			bar.Set64(m.ProgressDetail.Current)
			bar.Update()
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
