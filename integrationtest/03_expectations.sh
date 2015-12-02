#!/bin/bash

set -e
INV=$(pwd)/../involucro

$INV -e "inv.task('test').using('busybox').withExpectation({stdout = 'Hello, World!\n'}).run('echo', 'Hello, World!')" test

set +e
$INV -e "inv.task('test').using('busybox').withExpectation({stdout = 'Hello, World!\n'}).run('echo', 'Hello, Moon')" test && { echo "Accepted output..."; exit 1; }
set -e

exit 0


