package app

import (
	"errors"
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

	arguments, err := parseArguments(argv, exit)
	if err != nil {
		return err
	}
	if add := arguments["-v"].(int); add == 0 {
		log.SetLevel(log.WarnLevel)
	} else if add == 1 {
		log.SetLevel(log.InfoLevel)
	} else if add == 2 {
		log.SetLevel(log.DebugLevel)
	}

	if arguments["--encoded-state"].(bool) != false {
		if steps, err := runtime.DecodeState(os.Getenv("STATE")); err != nil {
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

	ctx := runtime.New(additionalArguments)

	if arguments["-e"] != nil {
		if err := ctx.RunString(arguments["-e"].(string)); err != nil {
			log.WithFields(log.Fields{"error": err}).Fatal("Failed executing script")
		}
	} else {
		filename := arguments["-f"].(string)

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			if _, err := os.Stat(filename + ".md"); err == nil {
				filename += ".md"
			}
		}

		if err := ctx.RunFile(filename); err != nil {
			if err := ctx.RunLiterateFile(filename); err != nil {
				log.WithFields(log.Fields{"error": err}).Fatal("Failed executing file")
			}
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
