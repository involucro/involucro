#!/usr/bin/env bats

INV=$(pwd)/../involucro

@test "examples/hello_world" {
  cd ../examples/hello_world
  $INV greet
}
