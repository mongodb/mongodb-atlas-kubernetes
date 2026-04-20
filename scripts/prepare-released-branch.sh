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

# Prepare released branch by checking out the commit and replacing CI tooling
# Usage: ./scripts/prepare-released-branch.sh <commit_sha>
# Example: ./scripts/prepare-released-branch.sh abc1234

set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
PROJECT_ROOT=$(cd -- "${SCRIPT_DIR}/.." &> /dev/null && pwd)

COMMIT_SHA="${1:-}"

if [ -z "${COMMIT_SHA}" ]; then
    echo "Error: Commit SHA is required" >&2
    echo "Usage: $0 <commit_sha>" >&2
    exit 1
fi

RELEASED_BRANCH_DIR="${PROJECT_ROOT}/released-branch"
CI_TOOLING_DIR="${PROJECT_ROOT}/.ci-tooling"

# Checkout released commit to released-branch/
echo "Checking out commit ${COMMIT_SHA} to released-branch/..."

# Fetch the commit if needed
git fetch origin "${COMMIT_SHA}" 2>/dev/null || git fetch origin

# Remove existing worktree/directory if it exists
git worktree remove "${RELEASED_BRANCH_DIR}" 2>/dev/null || rm -rf "${RELEASED_BRANCH_DIR}"

# Create worktree for the released commit
git worktree add -f "${RELEASED_BRANCH_DIR}" "${COMMIT_SHA}"

# Backup released-branch's Makefile and scripts if they exist
[ -f "${RELEASED_BRANCH_DIR}/Makefile" ] && mv "${RELEASED_BRANCH_DIR}/Makefile" "${RELEASED_BRANCH_DIR}/Makefile.bak" || true
[ -d "${RELEASED_BRANCH_DIR}/scripts" ] && mv "${RELEASED_BRANCH_DIR}/scripts" "${RELEASED_BRANCH_DIR}/scripts.bak" || true

# Copy CI branch's Makefile and scripts to a stable location
mkdir -p "${CI_TOOLING_DIR}"
cp "${PROJECT_ROOT}/Makefile" "${CI_TOOLING_DIR}/"
cp -r "${PROJECT_ROOT}/scripts" "${CI_TOOLING_DIR}/"

# Replace with symlinks to CI branch versions (use absolute path via .ci-tooling)
ln -sf "${CI_TOOLING_DIR}/Makefile" "${RELEASED_BRANCH_DIR}/Makefile"
ln -sf "${CI_TOOLING_DIR}/scripts" "${RELEASED_BRANCH_DIR}/scripts"

# Copy devbox.json and devbox.lock to released-branch so devbox can set up the environment correctly
cp "${PROJECT_ROOT}/devbox.json" "${RELEASED_BRANCH_DIR}/"
cp "${PROJECT_ROOT}/devbox.lock" "${RELEASED_BRANCH_DIR}/"

echo "✓ Released branch prepared at ${RELEASED_BRANCH_DIR}"

