#!/bin/bash

set -eou pipefail

docker login -u mongodb+mongodb_atlas_kubernetes -p "${REGISTRY_PASSWORD}" "${REGISTRY}"

DIGESTS=$(docker manifest inspect "${REGISTRY}/${REPOSITORY}:${VERSION}" | jq -r '.manifests[] | select(.platform.os!="unknown") | .digest')
mapfile -t PLATFORMS < <(docker manifest inspect "${REGISTRY}/${REPOSITORY}:${VERSION}" | jq -r '.manifests[] | select(.platform.os!="unknown") | .platform.architecture')

submit_flag=--submit
if [ "${SUBMIT}" == "false" ]; then
  submit_flag=
fi

INDEX=0
for DIGEST in $DIGESTS; do
    echo "Check and Submit result to RedHat Connect"
    # Send results to RedHat if preflight finished wthout errors
    preflight check container "${REGISTRY}/${REPOSITORY}@${DIGEST}" \
      --artifacts "${DIGEST}" \
      --platform "${PLATFORMS[$INDEX]}" \
      --pyxis-api-token="${RHCC_TOKEN}" \
      --certification-project-id="${RHCC_PROJECT}" \
      --docker-config="${HOME}/.docker/config.json" \
      ${submit_flag}

  (( INDEX++ )) || true
done
