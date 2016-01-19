#!/usr/bin/env bats

load find_inv

setup() {
  mkdir -p int21
  touch int21/testfile
}

teardown() {
  rm -rf int21
}

@test "flexible mount: via relative path" {
  $INV -e "inv.task('p').using('busybox').withHostConfig({Binds = {'./int21:/ttt'}}).run('rm', '/ttt/testfile')" p
  test ! -f "int21/testfile"
}

@test "flexible mount: via absolute path" {
  $INV -e "inv.task('p').using('busybox').withHostConfig({Binds = {'$(pwd)/int21:/ttt'}}).run('rm', '-f', '/ttt/testfile')" p
  test ! -f "int21/testfile"
}


