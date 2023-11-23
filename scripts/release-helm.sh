#!/bin/bash

set -euo pipefail

helm_charts_installation_id() {
	JWT=$1
	curl -s -X GET -H "Accept: application/vnd.github+json" \
		-H "Authorization: Bearer ${JWT}" \
		-H "X-GitHub-Api-Version: 2022-11-28" \
		"https://api.github.com/repos/mongodb/helm-charts/installation" | jq .id
}

helm_charts_token() {
	JWT=$1
	INSTALL_ID=$(helm_charts_installation_id "${JWT}")
	curl -s -X POST -H "Accept: application/vnd.github+json" \
		-H "Authorization: Bearer ${JWT}" \
		-H "X-GitHub-Api-Version: 2022-11-28" \
		"https://api.github.com/app/installations/${INSTALL_ID}/access_tokens" | jq -rc .token
}

JWT=$(tools/makejwt/makejwt -appId="${APP_ID}" -key="${RSA_PEM_KEY_BASE64}")

ACCESS_TOKEN=$(helm_charts_token "${JWT}")

curl -s --fail-with-body -X POST -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer ${ACCESS_TOKEN}"\
        -H "X-GitHub-Api-Version: 2022-11-28" \
        -d '{"ref":"main","inputs":{"version":"'"${VERSION}"'"}}' \
        https://api.github.com/repos/mongodb/helm-charts/actions/workflows/post-atlas-operator-release.yaml/dispatches
