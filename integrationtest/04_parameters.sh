#!/bin/bash


set -e

rm -f __asd
../involucro -e "inv.task('touch').using('busybox').run('touch', '/source/__asd')" touch
test -f "__asd"
rm -f __asd

set +e
../involucro -e "inv.task('touch').using('busybox').run('touch', '/source/__asd')" -f invfile.js touch && { echo "Accepted -e ... -f" ; exit 1; }
rm -f __asd
set -e

