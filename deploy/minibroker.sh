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

set -o errexit -o nounset -o pipefail -o xtrace

: "${NAMESPACE:=minibroker}"
: "${RELEASE_NAME:=minibroker}"
: "${SET_OVERRIDE_PARAMS:=""}"

if ! kubectl version 1> /dev/null 2> /dev/null; then
  >&2 echo "ERROR: Missing kubectl binary"
  exit 1
fi

if ! helm version 1> /dev/null 2> /dev/null; then
  >&2 echo "ERROR: Missing helm binary"
  exit 1
fi

script_dir="$(dirname "$(realpath "${BASH_SOURCE}")")"

if [ -z "$(kubectl get namespace "${NAMESPACE}" --output name)" ]; then
  kubectl create namespace "${NAMESPACE}"
fi

helm install "${RELEASE_NAME}" "${CHART_TARBALL}" \
  --wait \
  --namespace "${NAMESPACE}" \
  ${SET_OVERRIDE_PARAMS:+--values "${script_dir}/minibroker/override_params_values.yaml"} \
  --set "deployServiceCatalog=false" \
  --set "defaultNamespace=${NAMESPACE}"
