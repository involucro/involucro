#!/usr/bin/env bats

load find_inv

teardown() {
  docker rmi inttest/14
}

@test "wrap options: set entrypoint" {
  $INV -e "inv.task('wrap').wrap('.').inImage('busybox').at('/data').withConfig({Entrypoint = {'/bin/echo', 'Hello_Options'}}).as('inttest/14')" wrap

  docker run -it --rm inttest/14 | grep "Hello_Options"
}

@test "wrap options: set entrypoint without base image" {
  $INV -e "inv.task('wrap').wrap('.').at('/data').withConfig({Entrypoint = {'/bin/echo', 'Hello_Options'}}).as('inttest/14')" wrap

  docker inspect -f "{{.Config.Entrypoint}}" inttest/14 | grep '/bin/echo Hello_Options'
}

@test "wrap: wrap in image without command set" {
  $INV -e "inv.task('wrap').wrap('.').inImage('alpine').at('/data').withConfig({Entrypoint = {'/bin/echo', 'Hello_Options'}}).as('inttest/14')" wrap

  docker inspect -f "{{.Config.Entrypoint}}" inttest/14 | grep '/bin/echo Hello_Options'
}
