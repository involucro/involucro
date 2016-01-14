#!/usr/bin/env bats

INV=$(pwd)/../involucro

teardown() {
  docker rmi test/inttesting
}

@test "timing: creation date within seconds of now" {
  $INV -e "inv.task('package').wrap('../integrationtest').inImage('busybox').at('/data').as('test/inttesting')" package
  docker images | grep test/inttesting
  docker images | grep test/inttesting | grep "second"
}
