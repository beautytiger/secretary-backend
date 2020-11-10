#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DIR=$(cd $(dirname "${BASH_SOURCE[0]}") >/dev/null 2>&1 && pwd)
cd $DIR

platform='arm'
pushd ..
GOOS=linux GOARCH=${platform} go build -o build/backend-${platform}
popd

image="10.33.132.23:5000/ai/secretary:backend"
docker build -t ${image} .
docker push ${image}
