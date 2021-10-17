#!/usr/bin/env bash

set -euo pipefail

pushd /tmp
git clone https://github.com/shoenig/donutdns
cd donutdns
TAG=`git describe --tags $(git rev-list --tags --max-count=1)`
popd
rm -rf /tmp/donutdns

echo "building image for $TAG"

docker build --no-cache -t shoenig/donutdns:${TAG} .
docker push shoenig/donutdns:${TAG}
