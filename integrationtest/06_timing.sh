#!/bin/bash


docker rmi test/inttesting

set -e
../involucro -e "inv.task('package').wrap('../integrationtest').inImage('busybox').at('/data').as('test/inttesting')" package


docker images | grep test/inttesting

docker images | grep test/inttesting | grep "second"

set +e
docker rmi test/inttesting
