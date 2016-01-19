package app

import (
	"flag"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/thriqon/involucro/runtime"
)

// Main represents the usual main method of the
// whole program. It is moved to its own package
// to testing using go utils.
func Main() error {
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

		ctx := runtime.New(make(map[string]string), client, "/")
		if err := step.Take(&ctx); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Unable to run step")
			return err
		}
		return nil
	}

	client, err := docker.NewClient(dockerUrl)
	if err != nil {
		log.Fatal("Unable to create Docker client")
	}

	if err := client.Ping(); err != nil {
		log.Fatal("Docker not reachable")
	}

	ctx := runtime.New(variables, client, relativeWorkDir)

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

	for _, element := range flag.Args() {
		if err := ctx.RunTask(element); err != nil {
			log.WithFields(log.Fields{"error": err}).Fatal("Error during task processing")
			return err
		}
	}
	return nil
}
