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

# Get manifest and validate it's JSON before parsing
MANIFEST_OUTPUT=$(docker manifest inspect "${img}" 2>&1)
if ! echo "${MANIFEST_OUTPUT}" | jq empty >/dev/null 2>&1; then
    echo "Error: docker manifest inspect returned invalid JSON for ${img}" >&2
    echo "Output (first 500 chars): ${MANIFEST_OUTPUT:0:500}" >&2
    exit 1
fi

IMG_PLATFORMS_SHAS=$(echo "${MANIFEST_OUTPUT}" | \
  jq -rc '.manifests[] | select(.platform.os != "unknown" and .platform.architecture != "unknown") | .digest')

echo "${action} parent multiarch image ${img}@${MULTIARCH_IMG_SHA}..."
IMG="${img}@${MULTIARCH_IMG_SHA}" "${SCRIPT_DIR}/${action}.sh"

for platform_sha in ${IMG_PLATFORMS_SHAS}; do
  echo "${action} platform image ${img}@${platform_sha}..."
  IMG="${img}@${platform_sha}" "${SCRIPT_DIR}/${action}.sh"
done

msg="All signed"
if [ "${action}" == "verify" ]; then
  msg="All verified OK"
fi
echo "✅ ${msg}"
