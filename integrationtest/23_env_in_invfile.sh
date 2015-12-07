#!/bin/bash

INV=$(pwd)/../involucro

set -e
MESSAGE=inv_message $INV -e "inv.task('test').using('busybox').withExpectation({stdout = \"inv_message\"}).run('/bin/echo', ENV.MESSAGE)" test


