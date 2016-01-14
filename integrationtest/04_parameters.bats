#!/usr/bin/env bats

INV=$(pwd)/../involucro

@test "parameters: accepts inline script" {
  $INV -e "inv.task('touch').using('busybox').run('touch', '/source/__asd')" touch
  test -f "__asd"
}

@test "parameters: rejects both inline script AND filename" {
  run $INV -e "inv.task('touch').using('busybox').run('touch', '/source/__asd')" -f invfile.js touch
  [ "$status" -ne 0 ]
}

teardown() {
  rm -f __asd
}
