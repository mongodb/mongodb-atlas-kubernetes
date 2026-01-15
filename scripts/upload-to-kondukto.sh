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
# This script is responsible for uploading SBOM lite to Kondukto
#
# See: https://docs.devprod.prod.corp.mongodb.com/mms/python/src/sbom/silkbomb/docs/commands/UPLOAD
#
# Usage:
#  KONDUKTO_BRANCH_PREFIX=... store_ ${SBOM_JSON_LITE_PATH}
# Where:
#   KONDUKTO_BRANCH_PREFIX is the environment variable with the Kondukto branch common prefix
#   SBOM_JSON_LITE_PATH is the path to the SBOM lite json file to upload to Kondukto
###

# Constants
registry=artifactory.corp.mongodb.com/release-tools-container-registry-public-local
silkbomb_img="${registry}/silkbomb:2.0"
docker_platform="linux/amd64"

# Arguments
sbom_lite_json=$1
[ -z "${sbom_lite_json}" ] && echo "Missing SBOM lite JSON path parameter" && exit 1

# Environment inputs
kondukto_token="${KONDUKTO_TOKEN:?KONDUKTO_TOKEN must be set}"
kondukto_repo="${KONDUKTO_REPO:?KONDUKTO_REPO must be set}"
kondukto_branch_prefix="${KONDUKTO_BRANCH_PREFIX:?KONDUKTO_BRANCH_PREFIX must be set}"

# Computed values
arch=$(jq -r < "${sbom_lite_json}" '.components[0].properties[] | select( .name == "syft:metadata:architecture" ) | .value')
kondukto_branch="${kondukto_branch_prefix}-linux-${arch}"

echo "Computed Kondukto branch: ${kondukto_branch}"

# Upload
docker run --platform="${docker_platform}" --rm -v "${PWD}":/pwd \
  -e KONDUKTO_TOKEN="${kondukto_token}" \
  "${silkbomb_img}" upload --sbom-in "/pwd/${sbom_lite_json}" \
  --repo "${kondukto_repo}" --branch "${kondukto_branch}"

echo "${sbom_lite_json} uploaded to Kondukto"
