#!/bin/bash

set -eou pipefail

docker login -u mongodb+mongodb_atlas_kubernetes -p "${REGISTRY_PASSWORD}" "${REGISTRY}"

submit_flag=--submit
if [ "${SUBMIT}" == "false" ]; then
  submit_flag=
fi

echo "Check and Submit result to RedHat Connect"
# Send results to RedHat if preflight finished wthout errors
preflight check container "${REGISTRY}/${REPOSITORY}:${VERSION}" \
  --pyxis-api-token="${RHCC_TOKEN}" \
  --certification-project-id="${RHCC_PROJECT}" \
  --docker-config="${HOME}/.docker/config.json" \
  ${submit_flag}
