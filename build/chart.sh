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

: "${VERSION:="$("${git_root}/third-party/kubecf-tools/versioning/versioning.rb")"}"
: "${IMAGE_REPOSITORY:=minibroker-integration-tests}"
: "${IMAGE_TAG:="${VERSION}"}"
: "${IMAGE:="${IMAGE_REPOSITORY}:${IMAGE_TAG}"}"
: "${OUTPUT_CHARTS_DIR:="${git_root}/output"}"
: "${TMP_BUILD_DIR:="${git_root}/tmp"}"
: "${CHART_SRC:="${git_root}/chart/mits"}"

>&2 echo "Building chart to ${OUTPUT_CHARTS_DIR}"

mkdir -p "${TMP_BUILD_DIR}"
tmp_chart_build_dir="${TMP_BUILD_DIR}/mits"
rm -rf "${tmp_chart_build_dir}"
cp -R "${CHART_SRC}" "${tmp_chart_build_dir}"

sed -i "s#<%image%>#${IMAGE}#" "${tmp_chart_build_dir}/values.yaml"

helm package ${CHART_SIGN_KEY:+--sign --key "${CHART_SIGN_KEY}"} \
    --destination "${OUTPUT_CHARTS_DIR}" \
    --version "${VERSION}" \
    "${tmp_chart_build_dir}"

rm -rf "${tmp_chart_build_dir}"
