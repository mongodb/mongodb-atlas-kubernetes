#!/bin/bash
# Copyright 2025 MongoDB Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# You may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script moves a multi-arch image from one registry to another using docker buildx.
set -euo pipefail

# Required env vars
: "${IMAGE_SRC_REPO:?Missing IMAGE_SRC_REPO}"
: "${IMAGE_SRC_TAG:?Missing IMAGE_SRC_TAG}"
: "${IMAGE_DEST_REPO:?Missing IMAGE_DEST_REPO}"
: "${IMAGE_DEST_TAG:?Missing IMAGE_DEST_TAG}"

# Optional env vars -> ALIAS TAG can be updated to a list later on
ALIAS_TAG="${ALIAS_TAG:-}"
ALIAS_ENABLED="${ALIAS_ENABLED:-false}"

image_src_url="${IMAGE_SRC_REPO}:${IMAGE_SRC_TAG}"
image_dest_url="${IMAGE_DEST_REPO}:${IMAGE_DEST_TAG}"

echo "Checking if ${image_dest_url} already exists remotely..."
if docker manifest inspect "${image_dest_url}" > /dev/null 2>&1; then
  echo "Image ${image_dest_url} already exists. Skipping transfer."
  exit 0
fi

echo "Transferring multi-arch image:"
echo "  From: ${image_src_url}"
echo "  To:   ${image_dest_url}"

BUILDER_NAME="tmpbuilder-move-image"
# Remove builder if it already exists (from previous failed run)
if docker buildx inspect "${BUILDER_NAME}" > /dev/null 2>&1; then
  docker buildx rm "${BUILDER_NAME}" > /dev/null 2>&1 || true
fi
docker buildx create --name "${BUILDER_NAME}" --use > /dev/null
docker buildx imagetools create "${image_src_url}" --tag "${image_dest_url}"

if [[ "${ALIAS_ENABLED}" == "true" && -n "${ALIAS_TAG}" ]]; then
  echo "Aliasing ${image_src_url} as ${IMAGE_DEST_REPO}:${ALIAS_TAG}"
  docker buildx imagetools create "${image_src_url}" --tag "${IMAGE_DEST_REPO}:${ALIAS_TAG}"
  echo "Successfully aliased as ${IMAGE_DEST_REPO}:${ALIAS_TAG}"
fi

docker buildx rm "${BUILDER_NAME}" > /dev/null 2>&1 || true
echo "Successfully moved ${image_src_url} -> ${image_dest_url}"
