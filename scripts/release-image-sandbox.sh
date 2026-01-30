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

# Helper script to set sandbox environment variables and run a command
# Usage: ./scripts/release-image-sandbox.sh <command> [args...]
# Example: ./scripts/release-image-sandbox.sh make push-release-images

set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
PROJECT_ROOT=$(cd -- "${SCRIPT_DIR}/.." &> /dev/null && pwd)

# Check if SANDBOX_REGISTRY is set
if [ -z "${SANDBOX_REGISTRY:-}" ]; then
    echo "Error: SANDBOX_REGISTRY is not set" >&2
    echo "Usage: SANDBOX_REGISTRY=<registry> $0 <command> [args...]" >&2
    exit 1
fi

# Check if command is provided
if [ $# -eq 0 ]; then
    echo "Error: No command provided" >&2
    echo "Usage: $0 <command> [args...]" >&2
    exit 1
fi

# 1. Flag to enable sandbox logic
export SANDBOX_MODE=true

# 2. Extract registry host and username from SANDBOX_REGISTRY
# Assuming SANDBOX_REGISTRY is like "ghcr.io/{user}/ako-sandbox"
SANDBOX_REGISTRY_HOST=$(echo "${SANDBOX_REGISTRY}" | cut -d/ -f1)
SANDBOX_USERNAME=$(echo "${SANDBOX_REGISTRY}" | cut -d/ -f2)

# 3. Get NEXT_VERSION from version.json (defaults to current-test if not available)
VERSION_FILE="${PROJECT_ROOT}/version.json"
if [ -f "${VERSION_FILE}" ] && command -v jq >/dev/null 2>&1; then
    NEXT_VERSION=$(jq -r .next "${VERSION_FILE}")
else
    NEXT_VERSION="unknown-test"
fi

# 4. Set version and promoted tag (can be overridden by env vars)
export VERSION="${VERSION:-${NEXT_VERSION}-test}"
export PROMOTED_TAG="${PROMOTED_TAG:-promoted-latest}"
export SANDBOX_REGISTRY

# 5. Upstream Repos (Defaults)
export DOCKER_PRERELEASE_REPO="${DOCKER_PRERELEASE_REPO:-docker.io/mongodb/mongodb-atlas-kubernetes-operator-prerelease}"
export QUAY_PRERELEASE_REPO="${QUAY_PRERELEASE_REPO:-quay.io/mongodb/mongodb-atlas-kubernetes-operator-prerelease}"

# 6. Sandbox Destinations
#    We simply prefix the existing repo names with our SANDBOX_REGISTRY.
#    This creates paths like: ttl.sh/abc123/docker.io/mongodb/...
export DEST_PRERELEASE_REPO="${SANDBOX_REGISTRY}/${DOCKER_PRERELEASE_REPO}"
export DOCKER_RELEASE_REPO="${SANDBOX_REGISTRY}/docker.io/mongodb/mongodb-atlas-kubernetes-operator"
export QUAY_RELEASE_REPO="${SANDBOX_REGISTRY}/quay.io/mongodb/mongodb-atlas-kubernetes-operator"

# 7. Signature Config
#    Aligning these ensures consistency.
export DOCKER_SIGNATURE_REPO="${SANDBOX_REGISTRY}/docker.io/mongodb/signatures"
export SIGNATURE_REPO="${SANDBOX_REGISTRY}/mongodb/signature"

# 8. Generate signing docker config
SIGNING_DOCKERCFG_BASE64=$("${SCRIPT_DIR}/gen-dockerconf.sh" "${SANDBOX_USERNAME}" "${SANDBOX_REGISTRY_HOST}")
export SIGNING_DOCKERCFG_BASE64

# 9. Run the command with all sandbox environment variables set
exec "$@"

