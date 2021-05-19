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
source_path=$(cd "$(dirname "$0")" && cd .. && pwd)
pushd "${source_path}" >/dev/null

./bin/"${repo}" -t test/test.tmpl -i ignore -o test/output.md -d test/test-folder

popd >/dev/null
