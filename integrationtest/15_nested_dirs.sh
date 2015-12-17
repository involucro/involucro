#!/bin/bash

INV=$(pwd)/../involucro

set -e
mkdir -p /tmp/involucro_inttest15/asd/p/aaa
echo 123 > /tmp/involucro_inttest15/asd/p/aaa/a
echo 456 > /tmp/involucro_inttest15/asd/p/aaa/b

function finish() {
  set +e
  rm -rf /tmp/involucro_inttest15/
  docker rmi test/i15
}
trap finish EXIT

cd /tmp/

$INV -e "inv.task('wrap').wrap('involucro_inttest15').inImage('busybox').at('/data').as('test/i15')" wrap

docker run -it --rm test/i15 /bin/grep 123 /data/asd/p/aaa/a
docker run -it --rm test/i15 /bin/grep 456 /data/asd/p/aaa/b
