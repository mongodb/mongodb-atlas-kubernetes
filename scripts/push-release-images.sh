#!/bin/bash
# Copyright 2025 MongoDB Inc
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

# Push, sign, and verify release images per target
# This is Phase 6 of the release process - the point of no return
# Each target is processed atomically: push → sign → verify before moving to next

set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# Color output helpers
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Required environment variables
: "${VERSION:?Missing VERSION}"
: "${PROMOTED_TAG:?Missing PROMOTED_TAG}"
: "${DOCKER_PRERELEASE_REPO:?Missing DOCKER_PRERELEASE_REPO}"
: "${QUAY_PRERELEASE_REPO:?Missing QUAY_PRERELEASE_REPO}"
: "${DOCKER_RELEASE_REPO:?Missing DOCKER_RELEASE_REPO}"
: "${QUAY_RELEASE_REPO:?Missing QUAY_RELEASE_REPO}"
: "${DEST_PRERELEASE_REPO:?Missing DEST_PRERELEASE_REPO}"
: "${DOCKER_SIGNATURE_REPO:?Missing DOCKER_SIGNATURE_REPO}"

# Credentials are required for signing
: "${PKCS11_URI:?Missing PKCS11_URI}"
: "${GRS_USERNAME:?Missing GRS_USERNAME}"
: "${GRS_PASSWORD:?Missing GRS_PASSWORD}"

# Optional environment variables
CERTIFIED_TAG="${CERTIFIED_TAG:-${VERSION}-certified}"

# Function to sign an image reference if signature doesn't already exist
sign_image_if_needed() {
    local img_ref="$1"
    local signature_repo="$2"
    
    log_info "  Signing ${img_ref} to ${signature_repo}..."
    if IMG="${img_ref}" SIGNATURE_REPO="${signature_repo}" \
        MAX_RETRIES=0 "${SCRIPT_DIR}/sign-multiarch.sh" verify >/dev/null 2>&1; then
        unset MAX_RETRIES
        log_info "  Signature already exists, skipping"
    else
        unset MAX_RETRIES
        IMG="${img_ref}" SIGNATURE_REPO="${signature_repo}" \
            "${SCRIPT_DIR}/sign-multiarch.sh"
    fi
}

# Function to verify signature
verify_signature() {
    local img_url="$1"
    local signature_repo="$2"
    
    log_info "  Verifying ${img_url} against ${signature_repo}..."
    if IMG="${img_url}" SIGNATURE_REPO="${signature_repo}" \
        "${SCRIPT_DIR}/sign-multiarch.sh" verify >/dev/null 2>&1; then
        log_info "  ✓ Signature verified"
        return 0
    else
        log_error "  ✗ Signature verification failed"
        return 1
    fi
}

# Function to process a target: push → sign → verify
push_sign_verify() {
    local target_name="$1"
    local src_repo="$2"
    local dest_repo="$3"
    local src_tag="$4"
    local dest_tag="$5"
    local signature_repo="$6"
    
    log_info "Processing ${target_name}..."
    
    # Push image
    log_info "  Pushing ${src_repo}:${src_tag} → ${dest_repo}:${dest_tag}..."
    IMAGE_SRC_REPO="${src_repo}" \
        IMAGE_DEST_REPO="${dest_repo}" \
        IMAGE_SRC_TAG="${src_tag}" \
        IMAGE_DEST_TAG="${dest_tag}" \
        "${SCRIPT_DIR}/move-image.sh"
    
    # Sign image (use promoted image as source since signatures are SHA-based)
    # Sign to both the target repo and the canonical signature repo
    sign_image_if_needed "${PROMOTED_DOCKER_IMAGE}" "${signature_repo}"
    sign_image_if_needed "${PROMOTED_DOCKER_IMAGE}" "${DOCKER_SIGNATURE_REPO}"
    
    # Verify signature
    local dest_image="${dest_repo}:${dest_tag}"
    if ! verify_signature "${dest_image}" "${signature_repo}"; then
        log_error "Failed to verify ${target_name} - aborting"
        exit 1
    fi
    
    log_info "✓ ${target_name} complete (pushed, signed, verified)"
}

log_info "=== PUSH, SIGN, AND VERIFY RELEASE IMAGES ==="
log_warn "This is the point of no return - images will be publicly available"
log_info "Each target is processed atomically: push → sign → verify"

PROMOTED_DOCKER_IMAGE="${DOCKER_PRERELEASE_REPO}:${PROMOTED_TAG}"

# Process each target atomically
push_sign_verify \
    "Docker release image" \
    "${DOCKER_PRERELEASE_REPO}" \
    "${DOCKER_RELEASE_REPO}" \
    "${PROMOTED_TAG}" \
    "${VERSION}" \
    "${DOCKER_RELEASE_REPO}"

push_sign_verify \
    "Quay release image" \
    "${QUAY_PRERELEASE_REPO}" \
    "${QUAY_RELEASE_REPO}" \
    "${PROMOTED_TAG}" \
    "${VERSION}" \
    "${QUAY_RELEASE_REPO}"

push_sign_verify \
    "Quay certified image" \
    "${QUAY_PRERELEASE_REPO}" \
    "${QUAY_RELEASE_REPO}" \
    "${PROMOTED_TAG}" \
    "${CERTIFIED_TAG}" \
    "${QUAY_RELEASE_REPO}"

# Copy image back to prerelease for all-in-one test (no signing needed)
log_info "Copying image back to prerelease for all-in-one test..."
IMAGE_SRC_REPO="${DOCKER_RELEASE_REPO}" \
    IMAGE_DEST_REPO="${DEST_PRERELEASE_REPO}" \
    IMAGE_SRC_TAG="${VERSION}" \
    IMAGE_DEST_TAG="${VERSION}" \
    "${SCRIPT_DIR}/move-image.sh"

log_info "✓ All targets complete (pushed, signed, verified)"

