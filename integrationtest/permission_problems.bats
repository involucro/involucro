#!/usr/bin/env bats

load find_inv

teardown() {
  $INV -e "inv.task('x').using('busybox').run('/bin/sh', '-c', 'rm -f /source/test/only_root')" x
  docker rmi 'inttest/wrap_root' || true
}

@test "able to wrap dirs containing root-only readable files" {
  $INV -e "inv.task('x').using('busybox').run('/bin/sh', '-c', 'echo FLAG > /source/test/only_root && chmod 0400 /source/test/only_root')" x
  test -e test/only_root
  test ! -r test/only_root

  $INV -e "inv.task('w').wrap('test').inImage('busybox').at('/data').as('inttest/wrap_root')" w

  docker run -it --rm inttest/wrap_root cat /data/only_root | grep FLAG
}
