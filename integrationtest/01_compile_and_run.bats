#!/usr/bin/env bats

INV=$(pwd)/../involucro

@test "examples/compile_and_run" {
  cd ../examples/compile_and_run
  rm -f dist/*
  $INV -vv compile package run

  OUTPUT=$(docker run -i --rm test/showaddition:v1 /usr/local/bin/add)

  test "x$OUTPUT" = "x5 + 10 = 15"
}

teardown() {
  docker rmi test/showaddition:v1
}
