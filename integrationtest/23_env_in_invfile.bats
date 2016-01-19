#!/usr/bin/env bats

load find_inv

@test "environment: available in invfile" {
  MESSAGE=inv_message $INV -e "inv.task('test').using('busybox').withExpectation({stdout = \"inv_message\"}).run('/bin/echo', ENV.MESSAGE)" test
}

