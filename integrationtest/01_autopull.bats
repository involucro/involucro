#!/usr/bin/env bats

load find_inv

setup() {
  if docker inspect tianon/true ; then
    docker rmi tianon/true
  fi
}

@test "autopull: when using absent image" {
  $INV -e "inv.task('test').using('tianon/true').run()" test
}
