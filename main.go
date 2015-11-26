package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

func main() {

	arguments := parse()
	fmt.Println(arguments)
	log.Info("Starting")

	ctx := InstantiateRuntimeEnv()
	ctx.EvalString(`2 + 3`)
	result := ctx.GetNumber(-1)
	ctx.Pop()
	log.WithFields(log.Fields{"result": result}).Info("result received")

	ctx.EvalString(`inv.task('blah')`)
	defer func() {
		if r := recover(); r != nil {
			log.Info("recovered, something went wrong")
		}
	}()
	ctx.EvalString(`inv.task(5)`)

	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	client.Ping()
	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
	for _, img := range imgs {
		log.WithFields(log.Fields{"id": img.ID}).Info("have docker image")
	}
}
