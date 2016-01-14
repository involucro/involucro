#!/usr/bin/env bats

INV=$(pwd)/../involucro

@test "from scratch: wrap without parent image" {
  $INV -e "inv.task('package').wrap('.').at('/').as('inttest/16')" package

  PARENT=$(docker inspect -f "{{.Parent}}" inttest/16)
  test "x$PARENT" = "x"
}

teardown() {
  docker rmi inttest/16
}

