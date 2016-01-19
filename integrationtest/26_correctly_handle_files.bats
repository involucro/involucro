#!/usr/bin/env bats

load find_inv

setup() {
  mkdir -p /tmp/involucro_test_26/p
  ln -s ../p /tmp/involucro_test_26/p/cur
}

@test "file handling: properly handle symlinks" {
  $INV -e "inv.task('wrap').wrap('/tmp/involucro_test_26').at('/data').inImage('busybox').as('inttest/26')" wrap

  docker run -it --rm inttest/26 ls /data/p/cur/cur/cur
}


teardown() {
  docker rmi inttest/26
  rm -rf /tmp/involucro_test_26/
}
