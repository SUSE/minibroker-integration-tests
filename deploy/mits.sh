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

if ! kubectl version 1> /dev/null 2> /dev/null; then
  >&2 echo "ERROR: Missing kubectl binary"
  exit 1
fi

if ! helm version 1> /dev/null 2> /dev/null; then
  >&2 echo "ERROR: Missing helm binary"
  exit 1
fi

>&2 echo "Deploying MITS..."

if [ -z "$(kubectl get namespace "${NAMESPACE}" --output name)" ]; then
  kubectl create namespace "${NAMESPACE}"
fi

helm install "${RELEASE_NAME}" "${CHART_TARBALL}" \
  --namespace "${NAMESPACE}" \
  ${SET_OVERRIDE_PARAMS:+--set "config.minibroker.provisioning.override_params.enabled=true"} \
  --set "config.cf.admin.username=${CF_ADMIN_USERNAME}" \
  --set "config.cf.admin.password="${CF_ADMIN_PASSWORD}"" \
  --set "config.cf.api.endpoint=${CF_API_ENDPOINT}"

function on_exit() {
  helm delete "${RELEASE_NAME}" --namespace "${NAMESPACE}"
  kubectl delete namespace "${NAMESPACE}"
}

trap on_exit EXIT

>&2 echo "Waiting for MITS to be ready..."
until kubectl get pod \
  --namespace "${NAMESPACE}" \
  --selector "job-name=${RELEASE_NAME}-mits" \
  --output name \
  2> /dev/null \
  | wc -l \
  | awk '$0 == 0 { exit 1 }'; do
    sleep 1
done

pod_name=$(kubectl get pod \
  --namespace "${NAMESPACE}" \
  --selector "job-name=${RELEASE_NAME}-mits" \
  --output name)
pod_name="${pod_name/#pod\//}"

kubectl wait pod "${pod_name}" \
  --namespace "${NAMESPACE}" \
  --for condition=ready \
  --timeout 3m

kubectl logs "${pod_name}" \
  --follow \
  --timestamps \
  --namespace "${NAMESPACE}"

# Wait for the container to terminate and then exit the script with the container's exit code.
jsonpath='{.status.containerStatuses[?(@.name == "mits")].state.terminated.exitCode}'
while true; do
  exit_code=$(kubectl get pod "${pod_name}" \
    --namespace "${NAMESPACE}" \
    --output "jsonpath=${jsonpath}")
  if [[ -n "${exit_code}" ]]; then
    exit "${exit_code}"
  fi
  sleep 1
done
