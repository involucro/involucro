package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"path/filepath"
)

func main() {

	arguments := parseArguments()
	fmt.Println(arguments)

	relativeWorkDir := arguments["-w"].(string)
	workingDir, _ := filepath.Abs(relativeWorkDir)
	log.SetLevel(log.DebugLevel)
	log.WithFields(log.Fields{"workdir": workingDir}).Info("Start")

	ctx := InstantiateRuntimeEnv(workingDir)

	ctx.duk.EvalString(`inv.task('test').using('busybox').run('echo', 'Hello, Inxmail')`)

	client, _ := docker.NewClient(arguments["--host"].(string))
	client.Ping()

	for _, element := range (arguments["<task>"]).([]string) {
		for _, step := range ctx.Tasks[element] {
			step.DryRun()
			step.WithDockerClient(client)
		}
	}
}
