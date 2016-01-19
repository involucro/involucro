#!/usr/bin/env bats

load find_inv

@test "hook: example 'dynamic_invfile'" {
  $INV -f $(dirname $INV)/examples/dynamic_invfile/invfile.lua prep
}

@test "hook: using io.lines" {
SCRIPT=$(cat <<'EOS'
inv.task('write')
  .using('busybox')
  .run('/bin/sh', '-c', 'echo A > 37_lines.txt && echo B >> 37_lines.txt && echo C >> 37_lines.txt')

inv.task('modify')
  .hook(function ()
      for line in io.lines("37_lines.txt") do
        inv.task('do' .. line).using('busybox').run('echo', line)
      end
    end)
EOS
)
  $INV -e "$SCRIPT" -v write modify doB 2>&1| grep "stdout: B"
}

teardown() {
  rm -f 37_lines.txt
}
