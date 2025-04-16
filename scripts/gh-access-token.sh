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

repo_installation_id() {
	REPO=$1
	JWT=$2
	curl -s -X GET -H "Accept: application/vnd.github+json" \
		-H "Authorization: Bearer ${JWT}" \
		-H "X-GitHub-Api-Version: 2022-11-28" \
		"https://api.github.com/repos/${REPO}/installation" | jq .id
}

repo_access_token() {
	REPO=$1
	JWT=$2
	INSTALL_ID=$(repo_installation_id "${REPO}" "${JWT}")
	curl -s -X POST -H "Accept: application/vnd.github+json" \
		-H "Authorization: Bearer ${JWT}" \
		-H "X-GitHub-Api-Version: 2022-11-28" \
		"https://api.github.com/app/installations/${INSTALL_ID}/access_tokens" | jq -rc .token
}

JWT=$(tools/makejwt/makejwt -appId="${APP_ID}" -key="${RSA_PEM_KEY_BASE64}")

repo_access_token "${REPO}" "${JWT}"
