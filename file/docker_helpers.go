package runtime

import "github.com/fsouza/go-dockerclient"

func repoNameAndTagFrom(name string) (repo, tag, autotag string) {
	repo, tag = docker.ParseRepositoryTag(name)

	autotag = tag

	if autotag == "" {
		autotag = "latest"
	}

	return
}
