#!/usr/bin/env bats

load find_inv

setup() {
  DIR=$(mktemp -d)
  mkdir -p $DIR/asd/p/aaa
  echo 123 > $DIR/asd/p/aaa/a
  echo 456 > $DIR/asd/p/aaa/b
}

teardown() {
  rm -rf $DIR/
  docker rmi test/i15
}

@test "nested dirs: wraps correctly" {
  cd $DIR/../

  $INV -e "inv.task('wrap').wrap('$(basename $DIR)').inImage('busybox').at('/data').as('test/i15')" wrap

  docker run -it --rm test/i15 /bin/grep 123 /data/asd/p/aaa/a
  docker run -it --rm test/i15 /bin/grep 456 /data/asd/p/aaa/b
}
