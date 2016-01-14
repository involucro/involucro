#!/usr/bin/env bats

INV=$(pwd)/../involucro

setup() {
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
}

@test "literate mode: using filename" {
  $INV -vv -f $FILE blubb 2>&1| grep "stdout: Test OK"
}

@test "literate mode: using automatic recognition in directory" {
  cd $DIR
  $INV -vv blubb 2>&1| grep "stdout: Test OK"
}

teardown() {
	rm -f $FILE
	rm -rf $DIR
}
