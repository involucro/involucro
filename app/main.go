package app

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/ilog"
	"github.com/thriqon/involucro/runtime"
)

// Main represents the usual main method of the
// whole program. It is moved to its own package
// to testing using go utils.
func Main(args []string) error {
	initializeFlagSet()
	flags.Init("involucro", flag.ContinueOnError)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	if relativeWorkDir == "." {
		if val := os.Getenv("INXMAIL_WORKDIR"); val != "" {
			relativeWorkDir = val
		}
	}

	if remoteWrapTask != "" {
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

	ilog.StdLog.SetMinPrintLevel(logLevel)

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

	if controlScript != "" {
		if err := ctx.RunString(controlScript); err != nil {
			return err
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
				return err
			}
		} else {
			if err := ctx.RunFile(filename); err != nil {
				return err
			}
		}
	}

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
