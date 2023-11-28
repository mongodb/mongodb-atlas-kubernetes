#!/bin/bash

set -euo pipefail

ACCESS_TOKEN=$(REPO="mongodb/helm-charts" ./scripts/gh-access-token.sh)

curl -s --fail-with-body -X POST -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer ${ACCESS_TOKEN}"\
        -H "X-GitHub-Api-Version: 2022-11-28" \
        -d '{"ref":"main","inputs":{"version":"'"${VERSION}"'"}}' \
        https://api.github.com/repos/mongodb/helm-charts/actions/workflows/post-atlas-operator-release.yaml/dispatches
