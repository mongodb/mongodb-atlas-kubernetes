#!/bin/bash
# Copyright 2025 MongoDB Inc
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

# This script is made to make sure the test environment is clean after running the upgrade test.

expect_success_silent() {
  local cmd=$1
  if $cmd; then
    return 0
  fi
  return 1
}

oc_login() {
  echo "Logging in to the cluster..."
  expect_success_silent "oc login --token=${OC_TOKEN} --server=${CLUSTER_API_URL}"
  echo "Logged in to the cluster"
  oc get clusterversion/version
}

main() {
  echo "AKO openshift upgrade test cleaner"
  oc_login
  echo "Cleaning up the test leftovers..."
  namespaces=$(oc get namespaces | grep "ako-upgrade-test-" | awk '{print $1}')
  if [ -z "$namespaces" ]; then
    echo "No namespaces found for cleanup"
  else
    for ns in $namespaces; do
      echo "Deleting namespace: $ns"
      oc delete namespace --force "$ns" --wait=true --timeout=60s || true
    done
  fi
}

# Entrypoint
main

