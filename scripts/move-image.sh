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

# Required environment variables
: "${IMAGE_SRC_REPO:?Missing IMAGE_SRC_REPO}"
: "${IMAGE_SRC_TAG:?Missing IMAGE_SRC_TAG}"
: "${IMAGE_DEST_REPO:?Missing IMAGE_DEST_REPO}"
: "${IMAGE_DEST_TAG:?Missing IMAGE_DEST_TAG}"

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

echo "Creating temporary buildx builder..."
docker buildx create --name "${BUILDER_NAME}" --use > /dev/null
docker buildx imagetools create "${image_src_url}" --tag "${image_dest_url}"
docker buildx rm "${BUILDER_NAME}" > /dev/null
echo "Successfully moved ${image_src_url} -> ${image_dest_url}"
