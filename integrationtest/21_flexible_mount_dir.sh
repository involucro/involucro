#!/bin/bash

INV=$(pwd)/../involucro

mkdir -p int21

set -e

$INV -e "inv.task('p').using('busybox').withHostConfig({Binds = {'./int21:/ttt'}}).run('touch', '/ttt/asd')" p
test -f "int21/asd"

$INV -e "inv.task('p').using('busybox').withHostConfig({Binds = {'$(pwd)/int21:/ttt'}}).run('rm', '-f', '/ttt/asd')" p
test ! -f "int21/asd"

set +e

rmdir int21
