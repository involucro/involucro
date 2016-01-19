#!/usr/bin/env bats

load find_inv

teardown() {
  docker rmi inttest/14
}

@test "wrap options: set entrypoint" {
  $INV -e "inv.task('wrap').wrap('.').inImage('busybox').at('/data').withConfig({Entrypoint = {'/bin/echo', 'Hello_Options'}}).as('inttest/14')" wrap

  docker run -it --rm inttest/14 | grep "Hello_Options"
}

