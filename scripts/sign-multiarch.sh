#!/bin/bash

set -euo pipefail

REPO=${IMG_REPO:-docker.io/mongodb/mongodb-atlas-kubernetes-operator-prerelease}
img=${IMG:-$REPO:$VERSION}
action=${1:-sign}

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

docker pull "${img}"
MULTIARCH_IMG_SHA=$(docker inspect "${img}" |jq -rc '.[0].Id')
IMG_PLATFORMS_SHAS=$(docker manifest inspect "${img}" | \
  jq -rc '.manifests[] | select(.platform.os != "unknown" and .platform.architecture != "unknown") | .digest')

echo "${action} parent multiarch image ${img}@${MULTIARCH_IMG_SHA}..."
IMG="${img}@${MULTIARCH_IMG_SHA}" "${SCRIPT_DIR}/${action}.sh"

for platform_sha in ${IMG_PLATFORMS_SHAS}; do
  echo "${action} platform image ${img}@${platform_sha}..."
  IMG="${img}@${platform_sha}" "${SCRIPT_DIR}/${action}.sh"
done
