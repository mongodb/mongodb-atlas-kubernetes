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

# This script retrives the git commit sha given an image sha

set -euo pipefail

registry="$1"
repo="$2"
image_sha="$3"

if [[ "$image_sha" == "latest" ]]; then
  tag="promoted-latest"
else
  tag="promoted-${image_sha}"
fi

full_image="${registry}/${repo}:${tag}"

sha=$(skopeo inspect "docker://${full_image}" | jq -r '.Labels["org.opencontainers.image.revision"]')

if [[ -z "$sha" || "$sha" == "null" ]]; then
  echo "Error: Could not extract commit SHA from $full_image" >&2
  exit 1
fi

echo "$sha"
