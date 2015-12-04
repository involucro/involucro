#!/bin/bash

INV=$(pwd)/../involucro

docker rmi inttest/16

set -e
$INV -e "inv.task('package').wrap('.').at('/').as('inttest/16')" package

PARENT=$(docker inspect -f "{{.Parent}}" inttest/16)
test "x$PARENT" = "x"

set +e
docker rmi inttest/16

