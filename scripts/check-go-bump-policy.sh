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

# Go toolchain bump policy gate. Delegates latest-version detection, the
# soak window, and editing every go.mod to the gobump tool (declared in
# tools/toolbox/go.mod, run via `go tool -modfile=tools/toolbox/go.mod`).
#
# Behavior:
#   - If an open PR already exists on a auto/bump-go-* branch, skip (dedup).
#   - Otherwise run gobump with -soak ${POLICY_SOAK_DAYS}d (default 90).
#   - If gobump changed the root go.mod version, hand off to
#     scripts/create-go-bump-pr.sh to open a single bump PR.

set -euo pipefail

POLICY_SOAK_DAYS="${POLICY_SOAK_DAYS:-90}"

if ! [[ "${POLICY_SOAK_DAYS}" =~ ^[0-9]+$ ]] || [[ "${POLICY_SOAK_DAYS}" -eq 0 ]]; then
  echo "check-go-bump-policy: error: POLICY_SOAK_DAYS must be a positive integer, got '${POLICY_SOAK_DAYS}'" >&2
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
GO_MOD="${ROOT_DIR}/go.mod"
PR_SCRIPT="${ROOT_DIR}/scripts/create-go-bump-pr.sh"

[[ -f "${GO_MOD}" ]] || {
  echo "check-go-bump-policy: error: missing ${GO_MOD}" >&2
  exit 1
}
[[ -f "${PR_SCRIPT}" ]] || {
  echo "check-go-bump-policy: error: missing ${PR_SCRIPT}" >&2
  exit 1
}

# Prints "N.N.N" from the go directive of the given go.mod.
go_mod_version() {
  local v
  v="$(grep -E '^go[[:space:]]+[0-9]' "$1" | head -1 | awk '{print $2}' | tr -d '\r')"
  [[ -n "${v}" ]] || {
    echo "check-go-bump-policy: error: no go directive in $1" >&2
    return 1
  }
  echo "${v}"
}

find_open_go_bump_pull_request() {
  # Anchor on the branch name created by scripts/create-go-bump-pr.sh
  # (auto/bump-go-<version>) — PR titles can be edited/prefixed by reviewers,
  # branch names set by the automation cannot.
  local raw
  raw=$(gh pr list --state open --limit 100 --json number,title,url,headRefName) || {
    echo "check-go-bump-policy: error: gh pr list" >&2
    return 2
  }
  echo "${raw}" | jq -r '.[] | select(.headRefName | startswith("auto/bump-go-")) | "\(.number)\t\(.title)\t\(.url)"' | head -1
}

# --- main

command -v gh >/dev/null 2>&1 || {
  echo "check-go-bump-policy: error: gh is required" >&2
  exit 1
}
command -v jq >/dev/null 2>&1 || {
  echo "check-go-bump-policy: error: jq is required" >&2
  exit 1
}

pr_line="$(find_open_go_bump_pull_request)" || {
  # A gh failure must abort: bumping without dedup could open a duplicate PR.
  exit 1
}
if [[ -n "${pr_line}" ]]; then
  echo "check-go-bump-policy: open Go bump PR found — skipping (dedup)" >&2
  echo "check-go-bump-policy:   ${pr_line}" >&2
  exit 0
fi

before="$(go_mod_version "${GO_MOD}")"

gobump_args=(-soak "${POLICY_SOAK_DAYS}d")
echo "check-go-bump-policy: running gobump ${gobump_args[*]}" >&2
go tool -modfile=tools/toolbox/go.mod gobump "${gobump_args[@]}"

after="$(go_mod_version "${GO_MOD}")"

if [[ "${before}" == "${after}" ]]; then
  echo "check-go-bump-policy: no version change (${before}) — nothing to do"
  exit 0
fi

echo "check-go-bump-policy: bumped ${before} -> ${after} — opening PR"
"${PR_SCRIPT}" "${after}"
