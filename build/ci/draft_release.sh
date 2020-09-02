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

git_root="$(git rev-parse --show-toplevel)"

################################################################################
# VERIFY ASSETS
################################################################################

chart_file=$(find "${git_root}/output/" -name 'mits*')

if [ -z "${chart_file}" ]; then
  >&2 echo "Failed to publish: no chart found"
  exit 1
fi

if [ "$(wc -l <<<"${chart_file}")" -gt 1 ]; then
  >&2 echo "Failed to publish: found more than one chart candidate"
  exit 1
fi

# Push the tag before creating the release, which would trigger the creation of
# the tag automatically.
git push origin "refs/tags/${GIT_TAG}"

# Construct the release body as a draft first. We remove the draft after the
# chart asset was uploaded.
release_data=$(cat <<EOF
{
  "name": "${GIT_TAG}",
  "tag_name": "${GIT_TAG}",
  "body": "A MITS release.",
  "draft": true,
  "prerelease": false
}
EOF
)

>&2 echo "Creating draft release"

# Create the release as a draft and get its ID.
release_id=$(curl \
  --fail \
  --request POST \
  --header "Authorization: Bearer ${GITHUB_TOKEN}" \
  --header "Content-Type: application/json" \
  --header "Accept: application/vnd.github.v3+json" \
  --data "${release_data}" \
  "https://api.github.com/repos/${REPOSITORY}/releases" \
  | awk 'match($0, /^\s\s"id":\s(.*),/, id){ print id[1] }')

>&2 echo "Uploading chart asset ${chart_file} to draft release"

# Upload the chart asset.
curl \
  --silent \
  --fail \
  --request POST \
  --header "Authorization: Bearer ${GITHUB_TOKEN}" \
  --header "Content-Type: $(file --brief --mime-type "${chart_file}")" \
  --header "Accept: application/vnd.github.v3+json" \
  --data-binary "@${chart_file}" \
  "https://uploads.github.com/repos/${REPOSITORY}/releases/${release_id}/assets?name=$(basename "${chart_file}")"
