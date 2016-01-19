#!/usr/bin/env bats

load find_inv

@test "runoptions: can set environment variables" {
  $INV -v -e "inv.task('a').using('busybox').withConfig({Env = { 'FOO=bar' }}).run('/bin/sh', '-c', 'echo \$FOO')" a 2>&1 | grep "stdout: bar"
}
