#!/usr/bin/env bats

load find_inv

@test "variables: with -s k=asd" {
  $INV -v -e "inv.task('test').using('busybox').run('/bin/echo', VAR['k'])" -s k=asd test 2>&1| grep "stdout: asd"
}

@test "variables: with --set k=asd" {
  $INV -v -e "inv.task('test').using('busybox').run('/bin/echo', VAR['k'])" --set k=asd test 2>&1| grep "stdout: asd"
}

@test "variables: with --set k=asd=6" {
  $INV -v -e "inv.task('test').using('busybox').run('/bin/echo', VAR['k'])" --set k=asd=6 test 2>&1| grep "stdout: asd=6"
}

