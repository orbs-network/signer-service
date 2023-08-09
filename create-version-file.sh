#!/bin/bash -e

LATEST_SEMVER=$(git describe --tags --abbrev=0 --always)
SHORT_COMMIT=$(git rev-parse HEAD | cut -c1-8)
echo "$LATEST_SEMVER-$SHORT_COMMIT" > .version
