#!/usr/bin/env bats

INV=$(pwd)/../involucro

setup() {
  if docker inspect tianon/true ; then
    docker rmi tianon/true
  fi
}

@test "autopull: when using absent image" {
  $INV -e "inv.task('test').using('tianon/true').run()" test
}
