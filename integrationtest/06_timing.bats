#!/usr/bin/env bats

load find_inv

teardown() {
  docker rmi test/inttesting
}

@test "timing: creation date within seconds of now" {
  $INV -e "inv.task('package').wrap('"$(dirname $INV)"/integrationtest').inImage('busybox').at('/data').as('test/inttesting')" package
  docker images | grep test/inttesting
  docker images | grep test/inttesting | grep "second"
}
