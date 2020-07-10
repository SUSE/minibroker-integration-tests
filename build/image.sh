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

git_root=$(git rev-parse --show-toplevel)

: "${CHART_VERSION:=$(awk '/^version: /{ print $2 }' < "${git_root}/chart/mits/Chart.yaml")}"
: "${IMAGE_TAG:=splatform/mits:${CHART_VERSION}}"

if [[ "${MINIKUBE:=''}" == "true" ]]; then
  >&2 echo "Building using Minikube's Docker daemon..."
  eval "$(minikube docker-env)"
fi

docker build \
  --tag "${IMAGE_TAG}" \
  --file "${git_root}/image/Dockerfile" \
  .
