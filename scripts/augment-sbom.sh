#!/bin/bash
# Copyright 2025 MongoDB Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


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
registry=artifactory.corp.mongodb.com/release-tools-container-registry-public-local
silkbomb_img="${registry}/silkbomb:2.0"
docker_platform="linux/amd64"

# Arguments
sbom_lite_json=$1
[ -z "${sbom_lite_json}" ] && echo "Missing SBOM lite JSON path parameter" && exit 1
target_dir=$2
[ -z "${target_dir}" ] && echo "Missing target directory parameter #2" && exit 1

# Environment inputs
kondukto_token="${KONDUKTO_TOKEN:?KONDUKTO_TOKEN must be set}"
kondukto_repo="${KONDUKTO_REPO:?KONDUKTO_REPO must be set}"
kondukto_branch_prefix="${KONDUKTO_BRANCH_PREFIX:?KONDUKTO_BRANCH_PREFIX must be set}"

# Computed values
arch=$(jq -r '.components[0].properties[] | select( .name == "syft:metadata:architecture" ) | .value' <"${sbom_lite_json}")
kondukto_branch="${kondukto_branch_prefix}-linux-${arch}"
target="${target_dir}/linux-${arch}.sbom.json"

echo "Computed Kondukto branch: ${kondukto_branch}"

# Download
mkdir -p "${target_dir}"
docker run --platform="${docker_platform}" -it --rm -v "${PWD}":/pwd \
  -e KONDUKTO_TOKEN="${kondukto_token}" \
  "${silkbomb_img}" augment --sbom-in "/pwd/${sbom_lite_json}" \
  --repo "${kondukto_repo}" --branch "${kondukto_branch}" --sbom-out "/pwd/${target}"

echo "${target} augmented with Kondukto scan results"
