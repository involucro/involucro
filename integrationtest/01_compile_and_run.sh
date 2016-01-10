#!/bin/bash
docker rmi test/showaddition:v1

set -e

INV=$(pwd)/../involucro

cd ../examples/compile_and_run
rm -f dist/*
$INV -vv compile package run

OUTPUT=$(docker run -i --rm test/showaddition:v1 /usr/local/bin/add)

echo "Assert correctness of output"
test "x$OUTPUT" = "x5 + 10 = 15"

set +e
docker rmi test/showaddition:v1
exit 0
