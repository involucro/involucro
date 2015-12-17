package app

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/file"
	"github.com/thriqon/involucro/file/utils"
	"os"
	"strings"
)

// Main represents the usual main method of the
// whole program. It is moved to its own package
// to testing using go utils.
func Main(argv []string, exit bool) error {

	arguments, err := parseArguments(argv, exit)
	if err != nil {
		return err
	}
	log.SetLevel(log.DebugLevel)

	if arguments["--encoded-state"].(bool) != false {
		if steps, err := utils.DecodeState(os.Getenv("STATE")); err != nil {
			log.WithFields(log.Fields{"error": err}).Fatal("Unable to parse state")
			return err
		} else {
			client, err := docker.NewClient("unix://" + arguments["--socket"].(string))
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Fatal("Unable to connect to Docker")
			}
			for _, step := range steps {
				if err := step.WithDockerClient(client, "/"); err != nil {
					log.WithFields(log.Fields{"error": err}).Error("Unable to run step")
					return err
				}
			}
			return nil
		}
	}

	client, isremote, err := connectToDocker(arguments)
	if err != nil {
		log.Fatal("Unable to create Docker client")
	}

	if err := client.Ping(); err != nil {
		log.Fatal("Docker not reachable")
	}

	relativeWorkDir := arguments["-w"].(string)

	additionalArguments, err := parseAdditionalArguments(arguments["--set"].([]string))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Unable to parse arguments")
	}

	ctx := file.InstantiateRuntimeEnv(additionalArguments)

	if arguments["-e"] != nil {
		if err := ctx.RunString(arguments["-e"].(string)); err != nil {
			log.WithFields(log.Fields{"error": err}).Fatal("Failed executing script")
		}
	} else {
		if err := ctx.RunFile(arguments["-f"].(string)); err != nil {
			log.WithFields(log.Fields{"error": err}).Fatal("Failed executing file")
		}
	}

	taskrunner := ctx.RunLocallyTaskWith
	if isremote {
		taskrunner = ctx.RunTaskOnRemoteSystemWith
	}
	for _, element := range (arguments["<task>"]).([]string) {
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

func connectToDocker(args argumentsMap) (client *docker.Client, isremote bool, err error) {
	var url string
	var ok bool
	if url, ok = args["--host"].(string); !ok {
		url, err = docker.DefaultDockerHost()
		if err != nil {
			return
		}
	}
	client, err = docker.NewClient(url)
	isremote = strings.HasPrefix(url, "tcp:")
	return
}
