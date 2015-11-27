#!/bin/bash

set -e

INV=$(pwd)/../involucro

CONTAINERS_BEFORE=$(docker ps -a)

cd ../examples/hello_world
$INV greet

CONTAINERS_AFTER=$(docker ps -a)
test "x$CONTAINERS_BEFORE" = "x$CONTAINERS_AFTER"

exit 0
