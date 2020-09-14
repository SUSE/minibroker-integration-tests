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

: "${NAMESPACE:=mits}"
: "${RELEASE_NAME:=mits}"
: "${CHART_TARBALL:=""}"
: "${CF_ADMIN_USERNAME:=admin}"
: "${CF_ADMIN_PASSWORD:=""}"
: "${CF_API_ENDPOINT:=""}"
: "${SET_OVERRIDE_PARAMS:=""}"

if ! kubectl --version 1> /dev/null 2> /dev/null; then
  >&2 echo "ERROR: Missing kubectl binary"
  exit 1
fi

if ! helm --version 1> /dev/null 2> /dev/null; then
  >&2 echo "ERROR: Missing helm binary"
  exit 1
fi

kubectl create namespace "${NAMESPACE}"
helm install "${RELEASE_NAME}" "${CHART_TARBALL}" \
  --namespace "${NAMESPACE}" \
  ${SET_OVERRIDE_PARAMS:+--set "config.minibroker.provisioning.override_params.enabled=true"} \
  --set "config.cf.admin.username=${CF_ADMIN_USERNAME}" \
  --set "config.cf.admin.password="${CF_ADMIN_PASSWORD}"" \
  --set "config.cf.api.endpoint=${CF_API_ENDPOINT}"
