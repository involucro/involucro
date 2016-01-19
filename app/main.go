package app

import (
	"errors"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/runtime"
	"os"
	"strings"
)

// Main represents the usual main method of the
// whole program. It is moved to its own package
// to testing using go utils.
func Main(argv []string, exit bool) error {
	flag.Parse()

	switch {
	case silent:
		log.SetLevel(log.WarnLevel)
	case verbose:
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	if remoteWrapTask != "" {
		step := runtime.DecodeWrapStep(remoteWrapTask)
		client, err := docker.NewClient("unix:///sock")
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Fatal("Unable to connect to Docker")
		}

		if err := step.WithDockerClient(client, "/"); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Unable to run step")
			return err
		}
		return nil
	}

	client, isremote, err := connectToDocker(dockerUrl)
	if err != nil {
		log.Fatal("Unable to create Docker client")
	}

	if err := client.Ping(); err != nil {
		log.Fatal("Docker not reachable")
	}

	ctx := runtime.New(variables)

	if controlScript != "" && isControlFileOverriden() {
		log.Fatal("Specified both -e and -f")
	}

	if controlScript != "" {
		if err := ctx.RunString(controlScript); err != nil {
			log.WithFields(log.Fields{"error": err}).Fatal("Failed executing script")
		}
	} else {
		filename := controlFile

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			if _, err := os.Stat(filename + ".md"); err == nil {
				filename += ".md"
			}
		}

		if strings.HasSuffix(filename, ".md") {
			if err := ctx.RunLiterateFile(filename); err != nil {
				log.WithFields(log.Fields{"error": err}).Fatal("Failed executing file")
			}
		} else {
			if err := ctx.RunFile(filename); err != nil {
				log.WithFields(log.Fields{"error": err}).Fatal("Failed executing file")
			}
		}
	}

	if showTasks {
		for _, id := range ctx.TaskIDList() {
			fmt.Println(id)
		}
		return nil
	}

	taskrunner := ctx.RunLocallyTaskWith
	if isremote {
		taskrunner = ctx.RunTaskOnRemoteSystemWith
	}
	for _, element := range flag.Args() {
		if ctx.HasTask(element) {
			if err := taskrunner(element, client, relativeWorkDir); err != nil {
				log.WithFields(log.Fields{"error": err}).Fatal("Error during task processing")
			}
		} else {
			log.WithFields(log.Fields{"task": element}).Warn("no steps defined for task")
		}
	}
	return nil
}

func parseAdditionalArguments(in []string) (map[string]string, error) {
	answer := make(map[string]string)
	for _, x := range in {
		parts := strings.SplitN(x, "=", 2)
		if len(parts) < 2 {
			return nil, errors.New("Invalid parameter usage, should be KEY=VALUE")
		}
		answer[parts[0]] = parts[1]
	}
	return answer, nil
}

func connectToDocker(url string) (client *docker.Client, isremote bool, err error) {
	isremote = strings.HasPrefix(url, "tcp:")
	client, err = docker.NewClient(url)
	return
}
