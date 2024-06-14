#!/bin/bash

set -euo pipefail

###
# This script is responsible for uploading SBOM lite assets to Silk
#
# See: https://docs.devprod.prod.corp.mongodb.com/mms/python/src/sbom/silkbomb/docs/commands/UPLOAD
###

# Constants
registry=artifactory.corp.mongodb.com/release-tools-container-registry-local
silkbomb_img="${registry}/silkbomb:1.0"
docker_platform="linux/amd64"

# Arguments
sbom_lite_json=$1

# Environment inputs
artifactory_usr="${ARTIFACTORY_USERNAME}"
artifactory_pwd="${ARTIFACTORY_PASSWORD}"
client_id="${SILK_CLIENT_ID}"
client_secret="${SILK_CLIENT_SECRET}"
asset_group_prefix="${SILK_ASSET_GROUP}"

# Computed values
arch=$(jq -r < "${sbom_lite_json}" '.components[0].properties[] | select( .name == "syft:metadata:architecture" ) | .value')
asset_group="${asset_group_prefix}-linux-${arch}"

echo "Computed Silk Asset Group: ${asset_group}"

# Login to docker registry
echo "${artifactory_pwd}" |docker login "${registry}" -u "${artifactory_usr}" --password-stdin

# Upload
docker run --platform="${docker_platform}" -it --rm -v "${PWD}":/pwd \
  -e SILK_CLIENT_ID="${client_id}" -e SILK_CLIENT_SECRET="${client_secret}" \
  "${silkbomb_img}" upload --silk_asset_group "${asset_group}" --sbom_in "/pwd/${sbom_lite_json}"

echo "${sbom_lite_json} uploaded to Silk"
