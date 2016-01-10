#!/bin/bash

REPO=involucro/tool

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

docker login -e="$DOCKER_EMAIL" -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"

if [[ $TAG != "latest" ]]; then
	docker tag $REPO:latest $REPO:$TAG
fi

docker push $REPO:$TAG
