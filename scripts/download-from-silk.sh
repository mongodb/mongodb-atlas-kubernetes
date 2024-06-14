#!/bin/bash

set -euo pipefail

###
# This script is responsible for downloading augmented SBOM assets from Silk
#
# See: https://docs.devprod.prod.corp.mongodb.com/mms/python/src/sbom/silkbomb/docs/commands/DOWNLOAD
#
# Usage:
#  SILK_ASSET_GROUP=... download-from-silk ${TARGET_ARCH} ${TARGET_DIR}
# Where:
#   SILK_ASSET_GROUP is the environment variable with the silk assert group common prefix
#   TARGET_ARCH is the architecture to download from Silk
#   TARGET_DIR is the local directory in where to place the Silk downloaded SBOMs
###

# Constants
registry=artifactory.corp.mongodb.com/release-tools-container-registry-local
silkbomb_img="${registry}/silkbomb:1.0"
docker_platform="linux/amd64"

# Arguments
arch=$1
[ -z "${arch}" ] && echo "Missing arch parameter #1" && exit 1
target_dir=$2
[ -z "${target_dir}" ] && echo "Missing target directory parameter #2" && exit 1

# Environment inputs
artifactory_usr="${ARTIFACTORY_USERNAME}"
artifactory_pwd="${ARTIFACTORY_PASSWORD}"
client_id="${SILK_CLIENT_ID}"
client_secret="${SILK_CLIENT_SECRET}"
asset_group_prefix="${SILK_ASSET_GROUP}"

# Computed values
asset_group="${asset_group_prefix}-linux-${arch}"
target="${target_dir}/linux-${arch}.sbom.json"

echo "Computed Silk Asset Group: ${asset_group}"

# Login to docker registry
echo "${artifactory_pwd}" |docker login "${registry}" -u "${artifactory_usr}" --password-stdin

# Download
mkdir -p "${target_dir}"
docker run --platform="${docker_platform}" -it --rm -v "${PWD}":/pwd \
  -e SILK_CLIENT_ID="${client_id}" -e SILK_CLIENT_SECRET="${client_secret}" \
  "${silkbomb_img}" download -o "/pwd/${target}" --silk_asset_group "${asset_group}"

echo "${target} downloaded from Silk"
