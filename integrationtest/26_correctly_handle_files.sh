#!/bin/bash

INV=$(pwd)/../involucro

docker rmi inttest/26
rm -rf /tmp/involucro_test_26/
mkdir -p /tmp/involucro_test_26/p
ln -s ../p /tmp/involucro_test_26/p/cur

set -e
$INV -e "inv.task('wrap').wrap('/tmp/involucro_test_26').at('/data').inImage('busybox').as('inttest/26')" wrap

docker run -it --rm inttest/26 ls /data/p/cur/cur/cur

set +e

docker rmi inttest/26
rm -rf /tmp/involucro_test_26/
