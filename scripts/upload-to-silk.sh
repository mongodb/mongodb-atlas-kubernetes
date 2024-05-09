#!/bin/bash

set -euo pipefail

# Constants
registry=artifactory.corp.mongodb.com/release-tools-container-registry-local
silkbomb_img="${registry}/silkbomb:1.0"
platform=linux/amd64

# Arguments
sbom_lite_json=$1

# Environment inputs
artifactory_usr="${ARTIFACTORY_USERNAME}"
artifactory_pwd="${ARTIFACTORY_PASSWORD}"
client_id="${SILK_CLIENT_ID}"
client_secret="${SILK_CLIENT_SECRET}"
asset_group="${SILK_ASSET_GROUP}"

# Compute paths to map within silkbomb docker container
relative_dir=$(dirname docs/releases/v2.2.2/linux-amd64.sbom.json)
dir=$(cd "${relative_dir}" && pwd)

# Login to docker registry
echo "${artifactory_pwd}" |docker login "${registry}" -u "${artifactory_usr}" --password-stdin

# Upload
docker run --platform="${platform}" -it --rm -v "${PWD}":/pwd -v "${dir}":"/${relative_dir}" \
  -e SILK_CLIENT_ID="${client_id}" -e SILK_CLIENT_SECRET="${client_secret}" \
  "${silkbomb_img}" upload --silk_asset_group "${asset_group}" --sbom_in "/${sbom_lite_json}"
echo "${sbom_lite_json} uploaded to Silk"
