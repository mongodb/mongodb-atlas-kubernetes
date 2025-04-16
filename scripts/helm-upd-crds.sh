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

if [[ -z "${HELM_CRDS_PATH}" ]]; then
  echo "HELM_CRDS_PATH is not set"
  exit 1
fi

filesToCopy=()
for filename in ./bundle/manifests/atlas.mongodb.com_*.yaml; do
  absName="$(basename "$filename")"
  echo "Verifying file: ${filename}"
  if ! diff "$filename" "${HELM_CRDS_PATH}"/"$absName"; then
    filesToCopy+=("$filename")
  fi
done

fLen=${#filesToCopy[@]}
if [ "$fLen" -eq 0 ]; then
  echo "No CRD changes detected"
  exit 0
fi

echo "The following CRD changes detected:"
for (( i=0; i < fLen; i++ )); do
  echo "${filesToCopy[$i]}"
done

for (( i=0; i < fLen; i++ )); do
  echo "Copying ${filesToCopy[$i]} to ${HELM_CRDS_PATH}/"
  cp "${filesToCopy[$i]}" "${HELM_CRDS_PATH}"/
done
