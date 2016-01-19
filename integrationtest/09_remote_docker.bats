#!/usr/bin/env bats

INV=$(pwd)/../involucro

TASK="inv.task('wrap').wrap('fixture_09').at('/blah').inImage('busybox').as('inttest/9')"

setup() {
  mkdir -p fixture_09
  echo blahblubb > fixture_09/asd

  chmod a-r fixture_09/asd


  socat TCP-LISTEN:4243,fork UNIX-CONNECT:/var/run/docker.sock &
  SOCAT_PID=$!
}

@test "remote: unreadable for this user" {
  run cat fixture_09/asd
  [ "$status" -ne 0 ]
}

@test "remote: not working when connecting directly" {
  run $INV -e $TASK wrap
  [ "$status" -eq 1 ]
}

@test "remote: working when connecting via TCP" {
  $INV -v -w $(pwd) -H tcp://127.0.0.1:4243 -e $TASK  wrap
  docker inspect -f '{{.Id}}' inttest/9

  docker run -it --rm inttest/9 /bin/cat "/blah/asd" | grep blahblubb
}

teardown() {
  kill $SOCAT_PID

  rm -rf fixture_09
  if docker inspect -f '{{.Id}}' inttest/9; then
    docker rmi inttest/9
  fi
}
