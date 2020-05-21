#!/bin/bash

docker login -u $DOCKER_HUB_LOGIN -p $DOCKER_HUB_PASSWORD

./create-version-file.sh
export VERSION=$(cat .version)

docker tag orbsnetwork/signer:$VERSION orbsnetworkstaging/signer:$VERSION
docker push orbsnetworkstaging/signer:$VERSION
