package app

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/involucro/involucro/ilog"
	runtime "github.com/involucro/involucro/internal"
)

// Main represents the usual main method of the
// whole program. It is moved to its own package
// to testing using go utils.
func Main(args []string) error {
	initializeFlagSet()
	flags.Init("involucro "+version, flag.ContinueOnError)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	if showVersion {
		fmt.Printf("involucro %s\n", versionNotice)
		return nil
	}

	if relativeWorkDir == "." {
		if val := os.Getenv("INVOLUCRO_WORKDIR"); val != "" {
			relativeWorkDir = val
		}
	}

	if remoteWrapTask != "" {
		return runRemoteWrapTask()
	}

	ilog.StdLog.SetMinPrintLevel(-verbosity)

	var client *docker.Client
	var err error
	if dockerURL != defaultDockerURL {
		client, err = docker.NewClient(dockerURL)
	} else {
		client, err = docker.NewClientFromEnv()
	}
	if err != nil {
		return fmt.Errorf("Unable to create Docker client: %s", err)
	}

	if err := client.Ping(); err != nil {
		return fmt.Errorf("Docker not reachable: %s", err)
	}

	ctx := runtime.New(variables, client, relativeWorkDir)

	if controlScript != "" && isControlFileOverriden() {
		return fmt.Errorf("Specified both -e and -f")
	}

	runControlScriptOn(&ctx)

	if showTasks {
		tasks := ctx.TaskIDList()
		sort.Strings(tasks)
		for _, id := range tasks {
			fmt.Println(id)
		}
		return nil
	}

	for _, element := range flags.Args() {
		if err := ctx.RunTask(element); err != nil {
			return err
		}
	}
	return nil
}

func runControlScriptOn(ctx *runtime.Runtime) error {
	if controlScript != "" {
		return ctx.RunString(controlScript)
	}

	filename := controlFile
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		filename += ".md"
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("Control file %v (and %v.md) not found", filename, filename)
	}

	if strings.HasSuffix(filename, ".md") {
		return ctx.RunLiterateFile(filename)
	}

	return ctx.RunFile(filename)
}

func runRemoteWrapTask() error {
	ilog.StdLog.SetMinPrintLevel(-2)
	step := runtime.DecodeWrapStep(remoteWrapTask)
	client, err := docker.NewClient("unix:///sock")
	if err != nil {
		return err
	}

	ctx := runtime.New(make(map[string]string), client, "/")
	if err := step.Take(&ctx); err != nil {
		return err
	}
	return nil
}
