#!/usr/bin/env bash

# Copyright 2020 SUSE
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit -o nounset -o pipefail

# If no tags exist, create the first one starting with v0.1.0.
if ! git describe --tags 1> /dev/null 2> /dev/null; then
  git tag v0.1.0
else
  git_root="$(git rev-parse --show-toplevel)"
  next_version=$("${git_root}/third-party/kubecf-tools/versioning/versioning.rb" --next minor)
  git tag "v${next_version}"
fi
