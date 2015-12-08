#!/bin/bash

set -e
INV=$(pwd)/../involucro

$INV -e "inv.task('test').using('busybox').withExpectation({stdout = 'Hello, World!\n'}).run('echo', 'Hello, World!')" test

set +e
$INV -e "inv.task('test').using('busybox').withExpectation({stdout = 'Hello, World!\n'}).run('echo', 'Hello, Moon')" test && { echo "Accepted output..."; exit 1; }
set -e

$INV -e "inv.task('test').using('busybox').withExpectation({code = 1}).run('false')" test

set +e
$INV -e "inv.task('test').using('busybox').withExpectation({code = 1}).run('true')" test && { echo "Accepted output..."; exit 1; }
set -e

exit 0


