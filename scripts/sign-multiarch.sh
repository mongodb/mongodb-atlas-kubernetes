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


set -euo pipefail

REPO=${IMG_REPO:-docker.io/mongodb/mongodb-atlas-kubernetes-operator-prerelease}
img=${IMG:-$REPO:$VERSION}
action=${1:-sign}

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

"${SCRIPT_DIR}"/retry.sh docker pull "${img}"
MULTIARCH_IMG_SHA=$(docker inspect --format='{{index .RepoDigests 0}}' "${img}" |awk -F@ '{print $2}')
IMG_PLATFORMS_SHAS=$(docker manifest inspect "${img}" | \
  jq -rc '.manifests[] | select(.platform.os != "unknown" and .platform.architecture != "unknown") | .digest')

echo "${action} parent multiarch image ${img}@${MULTIARCH_IMG_SHA}..."
IMG="${img}@${MULTIARCH_IMG_SHA}" "${SCRIPT_DIR}/${action}.sh"

for platform_sha in ${IMG_PLATFORMS_SHAS}; do
  echo "${action} platform image ${img}@${platform_sha}..."
  IMG="${img}@${platform_sha}" "${SCRIPT_DIR}/${action}.sh"
done
