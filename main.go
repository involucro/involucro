package main

import "fmt"
import "github.com/fsouza/go-dockerclient"
import log "github.com/Sirupsen/logrus"

func main() {


	arguments := parse()
	fmt.Println(arguments)
	log.Info("Starting")

	ctx := InstantiateRuntimeEnv()
	ctx.EvalString(`2 + 3`)
	result := ctx.GetNumber(-1)
	ctx.Pop()
	log.WithFields(log.Fields{ "result": result, }).Info("result received")

	ctx.EvalString(`inv.task('blah')`)
	ctx.EvalString(`inv.task(5)`)

	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	client.Ping()
	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
	}
}
