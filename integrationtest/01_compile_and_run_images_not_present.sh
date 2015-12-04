#!/bin/bash

docker rmi sojournlabs/gcc
docker rmi busybox

set -e

INV=$(pwd)/../involucro

CONTAINERS_BEFORE=$(docker ps -a)

cd ../examples/compile_and_run
rm -f dist/*
$INV compile package run

OUTPUT=$(docker run -i --rm test/showaddition:v1 /usr/local/bin/add)

echo "Assert correctness of output"
test "x$OUTPUT" = "x5 + 10 = 15"

echo "Assert that no containers where leaked"
CONTAINERS_AFTER=$(docker ps -a)
test "x$CONTAINERS_BEFORE" = "x$CONTAINERS_AFTER"

exit 0
