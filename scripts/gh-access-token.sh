#!/bin/bash

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
