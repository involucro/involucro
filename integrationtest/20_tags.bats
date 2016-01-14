#!/usr/bin/env bats

INV=$(pwd)/../involucro

@test "tag: wraps with correct name" {
  $INV -e "inv.task('package').wrap('.').inImage('busybox').at('/data').as('inttest/20:v1')" package
  docker inspect -f "{{.ID}}" inttest/20:v1

  $INV -e "inv.task('run').using('inttest/20:v1').run('true')" run
}

teardown() {
  docker rmi inttest/20
  docker rmi inttest/20:v1
}
