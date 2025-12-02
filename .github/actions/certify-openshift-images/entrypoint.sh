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

docker login -u mongodb+mongodb_atlas_kubernetes -p "${REGISTRY_PASSWORD}" "${REGISTRY}"

submit_flag=--submit
if [ "${SUBMIT}" == "false" ]; then
  submit_flag=
fi

echo "Check and Submit result to RedHat Connect"
# Send results to RedHat if preflight finished wthout errors
preflight check container "${REGISTRY}/${REPOSITORY}:${VERSION}" \
  --pyxis-api-token="${RHCC_TOKEN}" \
  --certification-component-id="${RHCC_PROJECT}" \
  --docker-config="${HOME}/.docker/config.json" \
  ${submit_flag}
