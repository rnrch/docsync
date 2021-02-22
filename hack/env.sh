#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

export GOOS=$(go env GOOS)
export GOARCH=$(go env GOARCH)
export BUILD_PATH=bin/${GOOS}_${GOARCH}
export BINARY=docsync
export PKG=github.com/rnrch/docsync
export BUILD_SOURCE_HOME=$(cd ".." && pwd)
