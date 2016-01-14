#!/usr/bin/env bats

INV=$(pwd)/../involucro

@test "output task list when using -T and directly supplying script" {
  TASKS=$($INV -e "inv.task('a').using('busybox').run('x'); inv.task('b').using('busybox').run('z')" -T | sort | base64)
  [ "$TASKS" = "YQpiCg==" ]
}

@test "output task list when using --tasks and directly supplying script" {
  TASKS=$($INV -e "inv.task('a').using('busybox').run('x'); inv.task('b').using('busybox').run('z')" --tasks | sort | base64)
  [ "$TASKS" = "YQpiCg==" ]
}

