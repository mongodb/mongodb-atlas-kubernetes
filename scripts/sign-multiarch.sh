#!/bin/bash

set -euo pipefail

REPO=${IMG_REPO:-docker.io/mongodb/mongodb-atlas-kubernetes-operator-prerelease}
img_to_sign=${IMG_TO_SIGN:-$REPO:$VERSION}

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

docker pull "${img_to_sign}"
MULTIARCH_IMG_SHA=$(docker inspect "${img_to_sign}" |jq -rc '.[0].Id')
IMG_PLATFORMS_SHAS=$(docker manifest inspect "${img_to_sign}" | \
  jq -rc '.manifests[] | select(.platform.os != "unknown" and .platform.architecture != "unknown") | .digest')

echo "Signing parent multiarch image ${img_to_sign}@${MULTIARCH_IMG_SHA}..."
IMG_TO_SIGN="${img_to_sign}@${MULTIARCH_IMG_SHA}" "${SCRIPT_DIR}"/sign.sh

for platform_sha in ${IMG_PLATFORMS_SHAS}; do
  echo "Signing platform image ${img_to_sign}@${platform_sha}..."
  IMG_TO_SIGN="${img_to_sign}@${platform_sha}" "${SCRIPT_DIR}"/sign.sh
done
