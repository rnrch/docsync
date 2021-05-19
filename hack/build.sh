#!/bin/bash
# Copyright 2021 rnrch
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o nounset
set -o errexit
set -o pipefail

repo=$(basename -s .git "$(git config --get remote.origin.url)")
repo_path=github.com/rnrch/"${repo}"

source_path=$(cd "$(dirname "$0")" && cd .. && pwd)
pushd "$source_path" >/dev/null

date=$(date "+%Y%m%d-%H:%M:%S")
version=$(git describe --tags --dirty 2>/dev/null || echo 'unknown')
revision=$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')
branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'unknown')

ldflags="
  -X ${repo_path}/pkg/version.Version=${version}
  -X ${repo_path}/pkg/version.Revision=${revision}
  -X ${repo_path}/pkg/version.Branch=${branch}
  -X ${repo_path}/pkg/version.BuildDate=${date}"

echo "Building with -ldflags $ldflags"
go build -ldflags "${ldflags}" -o "${source_path}/bin/${repo}" "${source_path}/cmd"

popd >/dev/null
