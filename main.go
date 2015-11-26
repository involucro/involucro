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

	ctx.duk.EvalString(`inv.task('test').using('busybox').run('echo', 'Hello, Inxmail')`)

	for _, element := range (arguments["<task>"]).([]string) {
		for _, step := range ctx.Tasks[element] {
			step.DryRun()
			step.WithDockerClient(client)
		}
	}
}
