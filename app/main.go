package app

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	file "github.com/thriqon/involucro/file"
	wrap "github.com/thriqon/involucro/file/wrap"
	"path/filepath"
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

	client, _ := docker.NewClient(arguments["--host"].(string))

	if err := client.Ping(); err != nil {
		log.Fatal("Docker not reachable")
	}

	if arguments["--wrap"] != nil {
		conf := wrap.AsImage{
			SourceDir:         arguments["--wrap"].(string),
			TargetDir:         arguments["--at"].(string),
			ParentImage:       arguments["--into-image"].(string),
			NewRepositoryName: arguments["--as"].(string),
		}
		log.WithFields(log.Fields{"conf": conf}).Debug("Starting wrap")

		if err := conf.WithDockerClient(client); err != nil {
			log.WithFields(log.Fields{"error": err}).Panic("Failed wrapping")
		}
		return nil
	}

	relativeWorkDir := arguments["-w"].(string)
	workingDir, _ := filepath.Abs(relativeWorkDir)
	log.WithFields(log.Fields{"workdir": workingDir}).Info("Start")

	additionalArguments, err := parseAdditionalArguments(arguments["--set"].([]string))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Unable to parse arguments")
	}

	ctx := file.InstantiateRuntimeEnv(workingDir, additionalArguments)

	if arguments["-e"] != nil {
		if err := ctx.RunString(arguments["-e"].(string)); err != nil {
			log.WithFields(log.Fields{"error": err}).Fatal("Failed executing script")
		}
	} else {
		if err := ctx.RunFile(arguments["-f"].(string)); err != nil {
			log.WithFields(log.Fields{"error": err}).Fatal("Failed executing file")
		}
	}

	for _, element := range (arguments["<task>"]).([]string) {
		if ctx.HasTask(element) {
			if err := ctx.RunTaskWith(element, client); err != nil {
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
