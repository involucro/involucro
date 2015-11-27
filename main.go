package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	file "github.com/thriqon/involucro/file"
	wrap "github.com/thriqon/involucro/steps/wrap"
	"os"
	"path/filepath"
)

func main() {

	arguments := parseArguments()
	fmt.Println(arguments)

	client, _ := docker.NewClient(arguments["--host"].(string))
	err := client.Ping()
	if err != nil {
		log.Error("Docker not reachable")
	}

	log.SetLevel(log.DebugLevel)
	if arguments["--wrap"] != nil {
		conf := wrap.AsImage{
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

	ctx := file.InstantiateRuntimeEnv(workingDir)

	fmt.Println("#!/bin/sh")
	ctx.RunFile(arguments["-f"].(string))

	for _, element := range (arguments["<task>"]).([]string) {
		steps := ctx.Tasks[element]
		if len(steps) == 0 {
			log.WithFields(log.Fields{"task": element}).Warn("no steps defined for task")
		}
		for _, step := range steps {
			if arguments["-n"].(bool) {
				step.DryRun()
			} else if arguments["-s"].(bool) {
				step.AsShellCommandOn(os.Stdout)
			} else {
				step.WithDockerClient(client)
			}
		}
	}
}
