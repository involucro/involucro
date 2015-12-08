#!/bin/bash

set -e

INV=$(pwd)/../involucro

cd ../examples/compile_and_run
rm -f dist/*
$INV compile package run

OUTPUT=$(docker run -i --rm test/showaddition:v1 /usr/local/bin/add)

echo "Assert correctness of output"
test "x$OUTPUT" = "x5 + 10 = 15"

exit 0
