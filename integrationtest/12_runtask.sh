#!/bin/bash

INV=$(pwd)/../involucro

set -e
$INV -e "inv.task('blah').runTask('test'); inv.task('test').using('busybox').run('echo', 'TEST8102')" blah 2>&1 | grep "stdout: TEST8102"

