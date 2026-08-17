#!/usr/bin/env bash
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

# Opens a pull request for a Go version bump after bump-go.sh has updated files.
# Called by check-go-bump-policy.sh; can also be run standalone.
#
# Usage: create-go-bump-pr.sh <version>
#   <version> is the exact go directive (e.g. 1.26.2), no "go" prefix.
#
# Environment:
#   TEST_BUMP_DRY_RUN=1   Print what would happen without touching git or gh.

set -euo pipefail

if [[ $# -lt 1 || -z "${1}" ]]; then
  echo "usage: create-go-bump-pr.sh <version>" >&2
  echo "  example: create-go-bump-pr.sh 1.26.2" >&2
  exit 1
fi

version="${1#go}"
branch="auto/bump-go-${version}"
title="Bump Go version to ${version}"

if [[ "${TEST_BUMP_DRY_RUN:-}" == "1" ]]; then
  printf 'create-go-bump-pr: dry-run: would open PR "%s" from branch %s\n' "${title}" "${branch}"
  exit 0
fi

command -v gh >/dev/null 2>&1 || {
  echo "create-go-bump-pr: error: gh is required" >&2
  exit 1
}

git checkout -b "${branch}"
git add -A

# Create a signed commit via the GitHub API (auto-signs with the token identity).
# The branch is force-updated remotely via the API; a stale remote branch in the
# auto/bump-go-* namespace is safe to overwrite.
BRANCH="${branch}" COMMIT_MESSAGE="${title}" "$(dirname "$0")/create-signed-commit.sh"

gh pr create \
  --title "${title}" \
  --body "$(cat <<'EOF'
## Summary

Automated Go version bump triggered by the go-bump-policy schedule.

The policy (see `scripts/check-go-bump-policy.sh`) bumps when a newer
stable release is available on go.dev **and** its `.0` release is at
least 90 days old (soak window). A 2-minor gap skips the soak and
bumps immediately (past Go's N-1 support window).

## Checklist

- [ ] CI passes
- [ ] Review propagated version in Dockerfiles, `.tool-versions`, and secondary `go.mod` files
EOF
  )" \
  --base main \
  --head "${branch}"
