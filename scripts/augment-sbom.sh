#!/bin/bash

set -euo pipefail

###
# This script is responsible for uploading an SBOM and augmenting it using the SBOM scan results from Kondukto
#
# See: https://docs.devprod.prod.corp.mongodb.com/mms/python/src/sbom/silkbomb/docs/commands/AUGMENT
#
# Usage:
#  KONDUKTO_BRANCH_PREFIX=... augment-sbom ${SBOM_JSON_LITE_PATH} ${TARGET_DIR}
# Where:
#   KONDUKTO_BRANCH_PREFIX is the environment variable with the Kondukto branch common prefix
#   SBOM_JSON_LITE_PATH is the path to the SBOM lite json file to upload to Kondukto
#   TARGET_DIR is the local directory in where to place the augmented SBOMs
###

# Constants
registry=artifactory.corp.mongodb.com/release-tools-container-registry-local
silkbomb_img="${registry}/silkbomb:2.0"
docker_platform="linux/amd64"

# Arguments
sbom_lite_json=$1
[ -z "${sbom_lite_json}" ] && echo "Missing SBOM lite JSON path parameter" && exit 1
target_dir=$2
[ -z "${target_dir}" ] && echo "Missing target directory parameter #2" && exit 1

# Environment inputs
artifactory_usr="${ARTIFACTORY_USERNAME}"
artifactory_pwd="${ARTIFACTORY_PASSWORD}"
kondukto_token="${KONDUKTO_TOKEN}"
kondukto_repo="${KONDUKTO_REPO}"
kondukto_branch_prefix="${KONDUKTO_BRANCH_PREFIX}"

# Computed values
arch=$(jq -r < "${sbom_lite_json}" '.components[0].properties[] | select( .name == "syft:metadata:architecture" ) | .value')
kondukto_branch="${kondukto_branch_prefix}-linux-${arch}"
target="${target_dir}/linux-${arch}.sbom.json"

echo "Computed Kondukto branch: ${kondukto_branch}"

# Login to docker registry
echo "${artifactory_pwd}" |docker login "${registry}" -u "${artifactory_usr}" --password-stdin

# Download
mkdir -p "${target_dir}"
docker run --platform="${docker_platform}" -it --rm -v "${PWD}":/pwd \
  -e KONDUKTO_TOKEN="${kondukto_token}" \
  "${silkbomb_img}" augment --sbom-in "/pwd/${sbom_lite_json}" \
  --repo "${kondukto_repo}" --branch "${kondukto_branch}" -sbom-out "/pwd/${target}" 

echo "${target} augmented with Kondukto scan results"
