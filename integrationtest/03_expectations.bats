#!/usr/bin/env bats

load find_inv

@test "expectation: successful matching of stdout" {
  $INV -e "inv.task('test').using('busybox').withExpectation({stdout = 'Hello, World!'}).run('echo', 'Hello, World!')" test
}

@test "expectation: failed matching of stdout" {
  run $INV -e "inv.task('test').using('busybox').withExpectation({stdout = 'Hello, World!'}).run('echo', 'Hello, Moon')" test
  [ "$status" -eq 1 ]
}

@test "expectation: successful matching of stderr" {
  $INV -e "inv.task('test').using('busybox').withExpectation({stderr = 'Hello, World'}).run('/bin/sh', '-c', 'echo Hello, World 1>&2')" test
}

@test "expectation: failed matching of stderr" {
  run $INV -e "inv.task('test').using('busybox').withExpectation({stderr = 'Hello, World'}).run('/bin/sh', '-c', 'echo Hello, Moon 1>&2')" test
  [ "$status" -eq 1 ]
}

@test "expectation: match exit code 1" {
  $INV -e "inv.task('test').using('busybox').withExpectation({code = 1}).run('false')" test
}

@test "expectation: failed match for exit code" {
  run $INV -e "inv.task('test').using('busybox').withExpectation({code = 1}).run('true')" test
  [ "$status" -eq 1 ]
}



