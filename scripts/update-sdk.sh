#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CURRENT_SDK_RELEASE=$(cat "${SCRIPT_DIR}/../.atlas-sdk-version")
echo "CURRENT_SDK_RELEASE: $CURRENT_SDK_RELEASE"

LATEST_SDK_TAG=$(curl -sSfL -X GET  https://api.github.com/repos/mongodb/atlas-sdk-go/releases/latest | jq -r '.tag_name')
echo "LATEST_SDK_TAG: $LATEST_SDK_TAG"

LATEST_SDK_RELEASE=$(echo "${LATEST_SDK_TAG}" | cut -d '.' -f 1)
echo "LATEST_SDK_RELEASE: $LATEST_SDK_RELEASE"
echo  "==> Updating SDK to latest major version ${LATEST_SDK_TAG}"

go tool --modfile tools/toolbox/go.mod gomajor get --rewrite "go.mongodb.org/atlas-sdk/${CURRENT_SDK_RELEASE}" "go.mongodb.org/atlas-sdk/${LATEST_SDK_RELEASE}@${LATEST_SDK_TAG}"
go mod tidy

echo "$LATEST_SDK_RELEASE" > ".atlas-sdk-version"
echo "Done"