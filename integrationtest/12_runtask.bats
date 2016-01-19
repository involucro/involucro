#!/usr/bin/env bats

INV=$(pwd)/../involucro

@test "runtask: other task present" {
  $INV -v -e "inv.task('blah').runTask('test'); inv.task('test').using('busybox').run('echo', 'TEST8102')" blah 2>&1 | grep "stdout: TEST8102"
}

@test "runtask: other task not present" {
  run $INV -e "inv.task('test').runTask('udef')" test
  [ "$status" -ne 0 ]
}

@test "runtask: examples/nested_tasks" {
  cd ../examples/nested_tasks
  $INV all
}
