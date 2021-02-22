#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

curDir=$(cd "$(dirname "$0")" && pwd)
cd "${curDir}" || exit 1
. ./env.sh

cd "${BUILD_SOURCE_HOME}" || exit 1
./${BUILD_PATH}/${BINARY} -t test/test.tmpl -i ignore -o test/output.md -d test/test-folder -v
