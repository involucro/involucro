#!/bin/bash

INV=$(pwd)/../involucro


set -e
$INV -vv -e "inv.task('a').using('busybox').withConfig({Env = { 'FOO=bar' }}).run('/bin/sh', '-c', 'echo \$FOO')" a 2>&1 | grep "stdout: bar"
set +e
