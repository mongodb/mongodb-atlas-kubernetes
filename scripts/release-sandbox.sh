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
# Usage: ./scripts/release-sandbox.sh <command> [args...]
# Example: ./scripts/release-sandbox.sh make push-release-images

set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
PROJECT_ROOT=$(cd -- "${SCRIPT_DIR}/.." &> /dev/null && pwd)

# Check if command is provided
if [ $# -eq 0 ]; then
    echo "Error: No command provided" >&2
    echo "Usage: $0 <command> [args...]" >&2
    exit 1
fi

# Check if SANDBOX_REGISTRY is set, if not generate a default ttl.sh URL
if [ -z "${SANDBOX_REGISTRY:-}" ]; then
    # Generate a random suffix for ttl.sh (8 hex characters)
    RANDOM_SUFFIX=$(openssl rand -hex 4 2>/dev/null || date +%s | sha256sum | head -c 8)
    SANDBOX_REGISTRY="ttl.sh/${RANDOM_SUFFIX}"
    echo "Warning: SANDBOX_REGISTRY not set, using default: ${SANDBOX_REGISTRY}" >&2
    echo "Warning: This default registry may not be suitable for all tests. Consider setting SANDBOX_REGISTRY explicitly." >&2
fi

# 1. Flag to enable sandbox logic
export SANDBOX_MODE=true

# 2. Extract registry host and username from SANDBOX_REGISTRY
# Assuming SANDBOX_REGISTRY is like "ghcr.io/{user}/ako-sandbox" or "ttl.sh/{random}"
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

# 8. Generate signing docker config (skip if GH_TOKEN/GITHUB_TOKEN not set)
if [ -z "${GH_TOKEN:-}" ] && [ -z "${GITHUB_TOKEN:-}" ]; then
    echo "Warning: GH_TOKEN or GITHUB_TOKEN not set, skipping Docker config generation for signing" >&2
    echo "Warning: Image signing operations may fail if SIGNING_DOCKERCFG_BASE64 is required" >&2
    export SIGNING_DOCKERCFG_BASE64=""
else
    TOKEN="${GH_TOKEN:-${GITHUB_TOKEN}}"
    SIGNING_DOCKERCFG_BASE64=$("${SCRIPT_DIR}/gen-dockerconf.sh" "${SANDBOX_USERNAME}" "${SANDBOX_REGISTRY_HOST}")
    export SIGNING_DOCKERCFG_BASE64
    
    # Login Docker to the registry if it's ghcr.io (GitHub Container Registry)
    if [ "${SANDBOX_REGISTRY_HOST}" = "ghcr.io" ]; then
        echo "${TOKEN}" | docker login "${SANDBOX_REGISTRY_HOST}" -u "${SANDBOX_USERNAME}" --password-stdin 2>/dev/null || {
            echo "Warning: Failed to login Docker to ${SANDBOX_REGISTRY_HOST}. You may need to login manually:" >&2
            echo "  echo \$(gh auth token) | docker login ${SANDBOX_REGISTRY_HOST} -u ${SANDBOX_USERNAME} --password-stdin" >&2
        }
    fi
fi

# 9. Run the command with all sandbox environment variables set
exec "$@"

