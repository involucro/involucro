#!/bin/bash

INV=$(pwd)/../involucro

docker rmi inttest/20
docker rmi inttest/20:v1

set -e
$INV -e "inv.task('package').wrap('.').at('/data').as('inttest/20:v1')" package

docker inspect -f "{{.ID}}" inttest/20:v1

set +e
docker rmi inttest/20
docker rmi inttest/20:v1
