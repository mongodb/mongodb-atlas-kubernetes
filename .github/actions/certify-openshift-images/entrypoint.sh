#!/bin/bash

set -eou pipefail

docker login -u mongodb+mongodb_atlas_kubernetes -p "${QUAY_PASSWORD}" quay.io

DIGESTS=$(docker manifest inspect "quay.io/${REPOSITORY}:${VERSION}" | jq -r '.manifests[] | select(.platform.os!="unknown") | .digest')
mapfile -t PLATFORMS < <(docker manifest inspect "quay.io/${REPOSITORY}:${VERSION}" | jq -r '.manifests[] | select(.platform.os!="unknown") | .platform.architecture')

INDEX=0
for DIGEST in $DIGESTS; do
    echo "Check and Submit result to RedHat Connect"
    # Send results to RedHat if preflight finished wthout errors
    preflight check container "quay.io/${REPOSITORY}@${DIGEST}" \
      --artifacts "${DIGEST}" \
      --platform "${PLATFORMS[$INDEX]}" \
      --pyxis-api-token="${RHCC_TOKEN}" \
      --certification-project-id="${RHCC_PROJECT}" \
      --docker-config="${HOME}/.docker/config.json" \
      --submit

  (( INDEX++ )) || true
done
