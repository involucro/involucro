#!/usr/bin/env bats

INV="$(pwd)/../involucro -vv"

setup() {
  FILE=$(mktemp)
  echo ASD >> $FILE
  echo DSA >> $FILE

  SCRIPTFILE=$FILE.lua
  cat > $SCRIPTFILE <<EOS
for line in io.lines('$FILE') do
  inv.task('do' .. line)
    .using('busybox')
      .run('echo', 'TESTOK')
end
EOS

  RUN="$INV -f $SCRIPTFILE"
}

@test "io.lines: can run task ASD" {
  $RUN doASD 2>&1 | grep "TESTOK"
}

@test "io.lines: can run task DSA" {
  $RUN doDSA 2>&1| grep "TESTOK"
}

teardown() {
  if [ "x$FILE" != "x" ]; then
    rm -rf $FILE $SCRIPTFILE
  fi
}
