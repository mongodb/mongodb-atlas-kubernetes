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

# Prefers GNU semantics for sed, grep, sort, etc

supported_releases() {
    local all_versions
    local last_three_minor_version
    all_versions=$(git tag -l | grep '^v[0-9]*\.[0-9]*\.[0-9]*$' | sort -ru)
    last_three_minor_version=$(git tag -l | grep '^v[0-9]*\.[0-9]*\.[0-9]*$' | awk -F. '{print $1 "." $2}' | sort -ru | head -3)
    echo "${all_versions}" | grep "${last_three_minor_version}" | sed 's/^v//'
}

supported_releases_json() {
    local releases_json
    releases_json=""
    for release in $(supported_releases); do
        releases_json="${releases_json}\"${release}\","
    done
    # shellcheck disable=SC2001
    releases_json=$(echo "${releases_json}" | sed 's/,$//')
    echo "[${releases_json}]"
}

supported_releases_json
