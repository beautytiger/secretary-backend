#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DIR=$(cd $(dirname "${BASH_SOURCE[0]}") >/dev/null 2>&1 && pwd)
cd $DIR

# 摄像头获取截图的地址
export CAMURL="http://10.3.141.5:8080/?action=snapshot"

./main

# 登录地址
# http://localhost:8080/front/ admin:admin
