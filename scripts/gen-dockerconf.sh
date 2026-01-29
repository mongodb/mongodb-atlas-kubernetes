#!/bin/bash
# Copyright 2026 MongoDB Inc
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

# Usage: ./gen-dockerconf.sh <username> [registry=ghcr.io] [<token>=GITHUB_TOKEN]

user="$1"
registry="${2:-ghcr.io}"
token="${3:-${GITHUB_TOKEN:-${GH_TOKEN:?Error: No token provided in args or environment}}}"

if [ -z "$user" ] || [ -z "$token" ]; then
  echo "Error: Username and token are required." >&2
  echo "Usage: $0 <username> <token>" >&2
  exit 1
fi

# 1. Create the Basic Auth string (user:token)
auth_string=$(echo -n "${user}:${token}" | base64 | tr -d '\n')

# 2. Construct the Docker Config JSON
JSON_CONFIG="{\"auths\":{\"${registry}\":{\"auth\":\"${auth_string}\"}}}"

# 3. Output the Final Base64 encoded JSON (for SIGNING_DOCKERCFG_BASE64)
echo -n "${JSON_CONFIG}" | base64 | tr -d '\n'
