#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

curDir=$(cd "$(dirname "$0")" && pwd)
cd "${curDir}" || exit 1
. ./env.sh

DATE=$(date "+%Y%m%d-%H:%M:%S")
VERSION=$(git describe --tags "$(git rev-list --tags --max-count=1)")
REVISION=$(git rev-parse --short HEAD)
LDFLAGS="-X ${PKG}/pkg/version.version=${VERSION} -X ${PKG}/pkg/version.revision=${REVISION} -X ${PKG}/pkg/version.buildDate=${DATE}"

cd "${BUILD_SOURCE_HOME}" || exit 1
go build -o "${BUILD_SOURCE_HOME}/${BUILD_PATH}/${BINARY}"  -ldflags "${LDFLAGS}"
