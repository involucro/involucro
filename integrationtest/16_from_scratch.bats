#!/usr/bin/env bats

load find_inv

@test "from scratch: wrap without parent image" {
  $INV -e "inv.task('package').wrap('.').at('/').as('inttest/16')" package

  # currently a limitation of the wrap algorithm
  test $(docker history --quiet inttest/16 | wc -l) = 2
}

teardown() {
  docker rmi inttest/16
}

