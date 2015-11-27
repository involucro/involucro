package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	wrap "github.com/thriqon/involucro/steps/wrap"
	"path/filepath"
)

func main() {

	arguments := parseArguments()
	fmt.Println(arguments)

	client, _ := docker.NewClient(arguments["--host"].(string))
	client.Ping()

	log.SetLevel(log.DebugLevel)
	if arguments["--wrap"] != nil {
		conf := wrap.WrapAsImage{
			SourceDir:         arguments["--wrap"].(string),
			TargetDir:         arguments["--at"].(string),
			ParentImage:       arguments["--into-image"].(string),
			NewRepositoryName: arguments["--as"].(string),
		}
		log.WithFields(log.Fields{"conf": conf}).Debug("Starting wrap")

		err := conf.WithDockerClient(client)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Panic("Failed wrapping")
		}
		return
	}

	relativeWorkDir := arguments["-w"].(string)
	workingDir, _ := filepath.Abs(relativeWorkDir)
	log.WithFields(log.Fields{"workdir": workingDir}).Info("Start")

	ctx := InstantiateRuntimeEnv(workingDir)

	ctx.duk.PevalFile(arguments["-f"].(string))

	for _, element := range (arguments["<task>"]).([]string) {
		steps := ctx.Tasks[element]
		if len(steps) == 0 {
			log.WithFields(log.Fields{"task": element}).Warn("no steps defined for task")
		}
		for _, step := range steps {
			if arguments["-n"].(bool) {
				step.DryRun()
			} else if arguments["-s"].(bool) {
				//TODO
			} else {
				step.WithDockerClient(client)
			}
		}
	}
}
