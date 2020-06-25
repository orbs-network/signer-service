#!/bin/bash

./create-version-file.sh
docker build -t orbsnetwork/signer:$(cat .version) .