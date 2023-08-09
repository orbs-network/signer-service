#!/bin/bash

./create-version-file.sh
docker build -t orbsnetworkstaging/signer:$(cat .version) .