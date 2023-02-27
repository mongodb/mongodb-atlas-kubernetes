#!/bin/sh

set -eou pipefail

docker login -u mongodb+mongodb_atlas_kubernetes -p "${QUAY_PASSWORD}" quay.io

DIGESTS=$(docker manifest inspect "${REPOSITORY}:${VERSION}" | jq -r .manifests[].digest)

for DIGEST in $DIGESTS; do
  echo "Checking image $DIGEST"
  # Do the preflight check first
  preflight check container "${REPOSITORY}@${DIGEST}" --artifacts "${DIGEST}" --docker-config="${HOME}/.docker/config.json"

  if [ "$SUBMIT" = "true" ]; then
    rm -rf "${DIGEST}"
    echo "Submitting result to RedHat Connect"
    # Send results to RedHat if preflight finished wthout errors
    preflight check container "${REPOSITORY}@${DIGEST}" \
      --artifacts "${DIGEST}" \
      --pyxis-api-token="${RHCC_TOKEN}" \
      --certification-project-id="${RHCC_PROJECT}" \
      --docker-config="${HOME}/.docker/config.json" \
      --submit
  fi
done
