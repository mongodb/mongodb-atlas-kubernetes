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

version=${VERSION:?}

FILES=(
  "helm-charts/atlas-operator-crds/Chart.yaml"
  "helm-charts/atlas-operator/Chart.yaml"
)

for FILE in "${FILES[@]}"; do
  if [[ -f "$FILE" ]]; then
    echo "Updating $FILE..."
    awk -v version="$version" '
      BEGIN {
        updated_version = 0; updated_appVersion = 0; updated_img_version = 0;
      }
      /^version:/ && updated_version == 0 {
        $0 = "version: " version; updated_version = 1
      }
      /^appVersion:/ && updated_appVersion == 0 {
        $0 = "appVersion: " version; updated_appVersion = 1
      }
      /^    version:/ && updated_img_version == 0 {
        $0 = "    version: \"" version "\""; updated_img_version = 1
      }
      { print }
    ' "$FILE" > "${FILE}.tmp" && mv "${FILE}.tmp" "$FILE"
    rm -f "${FILE}.bak"
  else
    echo "Warning: File $FILE does not exist. Skipping."
  fi
done

echo "Version bump completed to ${version}."
