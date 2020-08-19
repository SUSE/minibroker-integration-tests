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

: "${REPOSITORY:=SUSE/minibroker-integration-tests}"

latest_release=$(curl \
  --silent \
  --request GET \
  --header "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/${REPOSITORY}/releases/latest")
latest_release_semver=$(awk 'match($0, /^\s\s"tag_name":\s"(.*)",/, version){ print version[1] }' <<<"${latest_release}")

if [ -z "${latest_release_semver}" ]; then
  latest_release_semver="0.0.0"
fi
latest_release_semver="${latest_release_semver/#v/}"
latest_release_semver_major="${latest_release_semver%%.*}"
latest_release_semver_minor="${latest_release_semver#*.}"
latest_release_semver_minor="${latest_release_semver_minor%%.*}"
latest_release_semver_patch="${latest_release_semver##*.}"

next_release_semver_major="${latest_release_semver_major}"
next_release_semver_minor=$(("${latest_release_semver_minor}" + 1))
next_release_semver_patch="${latest_release_semver_patch}"
next_release_semver="${next_release_semver_major}.${next_release_semver_minor}.${next_release_semver_patch}"

echo "${next_release_semver}"
