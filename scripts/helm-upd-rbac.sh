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


set -eou pipefail

echo "Working dir: $(pwd)"

if [[ -z "${HELM_RBAC_FILE}" ]]; then
  echo "HELM_RBAC_FILE is not set"
  exit 1
fi

if [ ! -f "${HELM_RBAC_FILE}" ]; then
  echo "File ${HELM_RBAC_FILE} does not exist. Skipping RBAC validation"
  exit 0
fi

yq '.spec.install.spec.clusterPermissions[0].rules' ./bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml > rbac.yaml

echo "Comparing RBAC for CSV to RBAC in AKO helm chart"
if ! diff rbac.yaml "$HELM_RBAC_FILE"; then
  echo "Copying RBAC"
  cp rbac.yaml "$HELM_RBAC_FILE"
else
  echo "No changes detected"
fi
