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

# This script sets up the .netrc file with the GitHub Actions credentials so CI
# can perform git operations in the origin.

set -o errexit -o nounset -o pipefail

cat <<EOF > "${HOME}/.netrc"
  machine github.com
  login ${GITHUB_ACTOR}
  password ${GITHUB_TOKEN}

  machine api.github.com
  login ${GITHUB_ACTOR}
  password ${GITHUB_TOKEN}
EOF

chmod 600 "${HOME}/.netrc"
