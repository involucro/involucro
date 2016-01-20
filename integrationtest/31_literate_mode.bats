#!/usr/bin/env bats

load find_inv

setup() {
  FILE=$(mktemp).md
  DIR=$(mktemp -d)

  cat > $FILE <<EOF
# Testing

See, it supports markdown: **emphasis**

We can also define tasks:

> inv.task('blubb')

Altough, interrupted!


I can also use four spaces:

    .using('busybox').run('/bin/echo', 'Test OK')

But it needs to be have an empty row before the code.
    This doesn't crash even though it's illegal Lua.

That's all
EOF

  cp $FILE $DIR/invfile.lua.md
}

@test "literate mode: using filename" {
  $INV -v -f $FILE blubb 2>&1| grep "stdout: Test OK"
}

@test "literate mode: using automatic recognition in directory" {
  cd $DIR
  $INV -v blubb 2>&1| grep "stdout: Test OK"
}

teardown() {
	rm -f $FILE
	rm -rf $DIR
}
