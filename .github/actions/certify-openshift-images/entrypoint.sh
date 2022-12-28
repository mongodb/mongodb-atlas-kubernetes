#!/bin/bash

set -eou pipefail

docker login -u unused -p "${QUAY_PASSWORD}" quay.io

DIGESTS=$(docker manifest inspect "${REPOSITORY}:${VERSION}" | jq -r .manifests[].digest)

for DIGEST in $DIGESTS; do
  echo "Checking image $DIGEST"
  # Do the preflight check first
  preflight check container "${DIGEST}" --docker-config="${HOME}/.docker/config.json"

  # Send results to RedHat if preflight finished without errors
  preflight check container "${DIGEST}" \
    --submit \
    --pyxis-api-token="${RHCC_TOKEN}" \
    --certification-project-id="${RHCC_PROJECT}" \
    --docker-config="${HOME}/.docker/config.json"
done
