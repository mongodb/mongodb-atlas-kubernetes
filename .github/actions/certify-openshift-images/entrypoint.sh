#!/bin/sh

set -eou pipefail

docker login -u mongodb+mongodb_atlas_kubernetes -p "${QUAY_PASSWORD}" quay.io

DIGESTS=$(docker manifest inspect "${REPOSITORY}:${VERSION}" | jq -r .manifests[].digest)

for DIGEST in $DIGESTS; do
    echo "Check and Submit result to RedHat Connect"
    # Send results to RedHat if preflight finished wthout errors
    preflight check container "quay.io/${REPOSITORY}@${DIGEST}" \
      --artifacts "${DIGEST}" \
      --pyxis-api-token="${RHCC_TOKEN}" \
      --certification-project-id="${RHCC_PROJECT}" \
      --docker-config="${HOME}/.docker/config.json" \
      --submit
done
