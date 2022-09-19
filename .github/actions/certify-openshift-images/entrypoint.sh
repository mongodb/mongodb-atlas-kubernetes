#!/bin/bash

set -eou pipefail

if [ -z "${IMAGE+x}" ]; then
  echo "IMAGE is not set"
  exit 1
fi

if [ -z "${VERSION+x}" ]; then
  echo "VERSION is not set"
  exit 1
fi

if [ -z "${RH_CERTIFICATION_OSPID+x}" ]; then
  echo "RH_CERTIFICATION_OSPID is not set"
  exit 1
fi

if [ -z "${RH_CERTIFICATION_TOKEN+x}" ]; then
  echo "RH_CERTIFICATION_TOKEN is not set"
  exit 1
fi

if [ -z "${RH_CERTIFICATION_PYXIS_API_TOKEN+x}" ]; then
  echo "RH_CERTIFICATION_PYXIS_API_TOKEN is not set"
  exit 1
fi

if [ -z "${CONTAINER_ENGINE+x}" ]; then
  echo "CONTAINER_ENGINE is not set, defaulting to podman"
  CONTAINER_ENGINE=podman
fi

preflight --version
${CONTAINER_ENGINE} --version

${CONTAINER_ENGINE} login -u unused -p "${RH_CERTIFICATION_TOKEN}" scan.connect.redhat.com --authfile ./authfile.json

IMG_SHA=$("${CONTAINER_ENGINE}" inspect --format='{{ index .RepoDigests 0}}' "${IMAGE}":"${VERSION}")

# Do the preflight check first
preflight check container "${IMG_SHA}" --docker-config=./authfile.json

# Send results to RedHat if preflight finished without errors
preflight check container "${IMG_SHA}" \
  --submit \
  --pyxis-api-token="${RH_CERTIFICATION_PYXIS_API_TOKEN}" \
  --certification-project-id="${RH_CERTIFICATION_OSPID}" \
  --docker-config=./authfile.json