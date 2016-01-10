#!/bin/bash

INV=$(pwd)/../involucro
FILE=$(mktemp)
DIR=$(mktemp -d)

cat > $FILE <<EOF
# Testing

See, it supports markdown: **emphasis**

We can also define tasks:

> inv.task('blubb')

Altough, interrupted!

> .using('busybox').run('/bin/echo', 'Test OK')

That's all
EOF

cp $FILE $DIR/invfile.lua.md

function finish {
	rm -f $FILE
	rm -rf $DIR
}

trap finish EXIT

set -e
echo "filename given"
$INV -vv -f $FILE blubb 2>&1| grep "stdout: Test OK"

echo "folder name"
cd $DIR
$INV -vv blubb 2>&1| grep "stdout: Test OK"
