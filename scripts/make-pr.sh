#!/bin/bash

set -euo pipefail

repo=${REPO:-mongodb/mongodb-atlas-kubernetes}
github_token=${GITHUB_TOKEN:-unset}

if [[ $(git diff --stat) = '' ]]; then
  echo 'No PR needed, git is clean'
  exit 0
fi

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

update_timestamp=$(date -u --iso-8601='minutes')
title="AKOBot Automatic Update ${update_timestamp}"

echo "Touching pkg/controller/touch.go to force the full test suite to run"
echo "package controller // updated at: ${update_timestamp}" \
  > "${SCRIPT_DIR}/../pkg/controller/touch.go"

echo "Creating branch akobot-update and pushing it"
git branch -m akobot-update
git add .
git commit -m "${title}"
git push -fu origin akobot-update

echo "Creating PR"
curl -L \
  -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${github_token}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  "https://api.github.com/repos/${repo}/pulls" -d "{\"title\":\"${title}\"}"
