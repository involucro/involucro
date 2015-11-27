#!/bin/bash


set -e

rm -f __asd
../involucro -e "inv.task('touch').using('busybox').run('touch', '/source/__asd')" touch
test -f "__asd"
rm -f __asd

../involucro -n -e "inv.task('touch').using('busybox').run('touch', '/source/__asd')" touch
test ! -f "__asd"

../involucro -s -e "inv.task('touch').using('busybox').run('touch', '/source/__asd')" touch
test ! -f "__asd"

set +e
../involucro -n -e "inv.task('touch').using('busybox').run('touch', '/source/__asd')" -f invfile.js touch && { echo "Accepted -e ... -f" ; exit 1; }
set -e
