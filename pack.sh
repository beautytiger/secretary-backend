#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DIR=$(cd $(dirname "${BASH_SOURCE[0]}") >/dev/null 2>&1 && pwd)
cd $DIR

pushd front/front-repo/
git pull
popd

cp -r front/front-repo/dist dist

platform=${1:-amd64}
GOOS=linux GOARCH=${platform} go build main.go

tar zcvf backend.tgz main dist meeting.jpg run.sh

#rm -rf dist main