#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DIR=$(cd $(dirname "${BASH_SOURCE[0]}") >/dev/null 2>&1 && pwd)
cd $DIR

image="10.33.132.23:5000/ai/secretary:front"
docker build -t ${image} .
docker push ${image}
