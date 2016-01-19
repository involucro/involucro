#!/usr/bin/env bats

load find_inv

setup() {
  DIR=$(mktemp -d)
}

teardown() {
  rm -rf $DIR
  docker rmi test/int34:v1 || true
  docker rmi test/int34:v2 || true
  docker rmi test/int34:latest || true
}

@test "tag image: both images have the same ID" {
  $INV -e "inv.task('wrap').wrap('$DIR').at('/').as('test/int34:v1').tag('test/int34:v1').as('test/int34:v2')" wrap
  V1=$(docker inspect -f '{{.Id}}' test/int34:v1)
  V2=$(docker inspect -f '{{.Id}}' test/int34:v2)

  [ $V1 = $V2 ]
}

@test "tag image: can tag v1 to latest (without specifying the tag name)" {
  $INV -e "inv.task('wrap').wrap('$DIR').at('/').as('test/int34:v1').tag('test/int34:v1').as('test/int34')" wrap
  V1=$(docker inspect -f '{{.Id}}' test/int34:v1)
  LATEST=$(docker inspect -f '{{.Id}}' test/int34:latest)

  [ $V1 = $LATEST ]
}

@test "tag image: can tag latest to v1" {
  $INV -e "inv.task('wrap').wrap('$DIR').at('/').as('test/int34:latest').tag('test/int34:latest').as('test/int34:v1')" wrap
  V1=$(docker inspect -f '{{.Id}}' test/int34:v1)
  LATEST=$(docker inspect -f '{{.Id}}' test/int34:latest)

  [ $V1 = $LATEST ]
}

@test "tag image: can tag an image ID to v1" {
  $INV -e "inv.task('wrap').wrap('$DIR').at('/').as('test/int34:latest')" wrap
  LATEST=$(docker inspect -f '{{.Id}}' test/int34:latest)

  $INV -e "inv.task('wrap2').tag('$LATEST').as('test/int34:v1')" wrap2
  V1=$(docker inspect -f '{{.Id}}' test/int34:v1)

  [ $V1 = $LATEST ]
}
