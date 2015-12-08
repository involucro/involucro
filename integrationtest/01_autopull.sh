#!/bin/bash

INV=$(pwd)/../involucro

docker rmi tianon/true

set -e
$INV -e "inv.task('test').using('tianon/true').run()" test
set +e


