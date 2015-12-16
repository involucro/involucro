#!/bin/bash

INV=$(pwd)/../involucro

TASK="inv.task('wrap').wrap('fixture_09').at('/blah').inImage('busybox').as('inttest/9')"


mkdir -p fixture_09
echo blahblubb > fixture_09/asd

chmod a-r fixture_09/asd

socat TCP-LISTEN:4243,fork UNIX-CONNECT:/var/run/docker.sock &
SOCAT_PID=$!

function finish() {
  set +e
  kill $SOCAT_PID

  rm -rf fixture_09
  docker rmi inttest/9
}
trap finish EXIT

set -e

# Demonstrate non-readability for this user
set +e
cat fixture_09/asd 2>/dev/null && exit 1
set -e

echo "Doesn't work when using it directly"
set +e
$INV -e $TASK wrap && exit 2
docker rmi inttest/9
set -e

echo "Does work when using it via TCP"

$INV -w $(pwd) -H tcp://127.0.0.1:4243 -e $TASK  wrap
docker inspect -f '{{.Id}}' inttest/9 > /dev/null

docker run -it --rm inttest/9 /bin/cat "/blah/asd" | grep blahblubb > /dev/null
