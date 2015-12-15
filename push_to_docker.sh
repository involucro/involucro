#!/bin/bash

docker login -e="$DOCKER_EMAIL" -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"

if [[ $TRAVIS_PULL_REQUEST != "false" ]]; then
  exit 0
fi

if [[ $TRAVIS_TAG != "" ]]; then
  TAG=$TRAVIS_TAG
else
  TAG=$TRAVIS_BRANCH
fi
if [[ $TAG == "master" ]]; then
  TAG="latest"
fi

REPO=involucro/tool

set -e
mkdir -p dist/
cp involucro dist/
./involucro -e "inv.task('wrap-yourself').wrap('dist').at('/').withConfig({entrypoint = {'/involucro'}}).as('$REPO:$TAG')" wrap-yourself
rm -rf dist/

docker push $REPO:$TAG
