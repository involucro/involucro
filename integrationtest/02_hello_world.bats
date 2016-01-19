#!/usr/bin/env bats

load find_inv

@test "examples/hello_world" {
  cd $(dirname $INV)/examples/hello_world
  $INV greet
}
