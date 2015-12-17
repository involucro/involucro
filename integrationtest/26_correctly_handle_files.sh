#!/bin/bash

INV=$(pwd)/../involucro

mkdir -p /tmp/involucro_test_26/p
ln -s ../p /tmp/involucro_test_26/p/cur

function finish() {
  set +e
  docker rmi inttest/26
  rm -rf /tmp/involucro_test_26/
}
trap finish EXIT

set -e
$INV -e "inv.task('wrap').wrap('/tmp/involucro_test_26').at('/data').inImage('busybox').as('inttest/26')" wrap

docker run -it --rm inttest/26 ls /data/p/cur/cur/cur

