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

# Certify OpenShift images using Red Hat preflight (supports Linux and Darwin)
#
# Validates container images against Red Hat's OpenShift certification requirements.
# When SUBMIT=true, writes results to Red Hat Connect (Pyxis API) - DESTRUCTIVE.
# When SUBMIT=false, only runs checks locally (read-only) - safe for testing.
#
# Note: Registry login must be performed before running this script.

set -euo pipefail

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
: "${REGISTRY:?Missing REGISTRY}"
: "${REPOSITORY:?Missing REPOSITORY}"
: "${VERSION:?Missing VERSION}"
: "${RHCC_TOKEN:?Missing RHCC_TOKEN}"
: "${RHCC_PROJECT:?Missing RHCC_PROJECT}"

# Optional environment variables
SUBMIT="${SUBMIT:-false}"
PREFLIGHT_BIN="${PREFLIGHT_BIN:-preflight}"

# Function to check if preflight is available
check_preflight() {
    if command -v "${PREFLIGHT_BIN}" >/dev/null 2>&1; then
        preflight_path=$(command -v "${PREFLIGHT_BIN}")
        log_info "Preflight found: ${preflight_path}"
        "${PREFLIGHT_BIN}" -v || true
        return 0
    fi
    return 1
}

# Function to install preflight
install_preflight() {
    local os
    local arch
    
    # Detect OS
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "${os}" in
        linux)
            os="linux"
            ;;
        darwin)
            os="darwin"
            ;;
        *)
            log_error "Unsupported OS: ${os}"
            return 1
            ;;
    esac
    
    # Detect architecture
    arch=$(uname -m)
    case "${arch}" in
        x86_64|amd64)
            arch="amd64"
            ;;
        aarch64|arm64)
            arch="arm64"
            ;;
        *)
            log_error "Unsupported architecture: ${arch}"
            return 1
            ;;
    esac

    local preflight_url="https://github.com/redhat-openshift-ecosystem/openshift-preflight/releases/latest/download/preflight-${os}-${arch}"
    local install_dir="${HOME}/.local/bin"
    local preflight_path="${install_dir}/preflight"

    log_info "Installing preflight for ${os}/${arch} to ${preflight_path}..."
    mkdir -p "${install_dir}"
    curl -L -o "${preflight_path}" "${preflight_url}"
    chmod +x "${preflight_path}"
    
    if [ -d "${install_dir}" ] && [[ ":$PATH:" != *":${install_dir}:"* ]]; then
        export PATH="${install_dir}:${PATH}"
        log_info "Added ${install_dir} to PATH for this session"
    fi
    
    PREFLIGHT_BIN="${preflight_path}"
    log_info "Preflight installed successfully"
}

# Main certification function
certify_image() {
    local image_ref="${REGISTRY}/${REPOSITORY}:${VERSION}"
    
    log_info "Certifying image: ${image_ref}"
    log_info "Registry: ${REGISTRY}"
    log_info "Repository: ${REPOSITORY}"
    log_info "Version: ${VERSION}"
    log_info "Submit: ${SUBMIT}"

    # Note: Registry login is assumed to have already been performed
    log_info "Using existing registry authentication for ${REGISTRY}"

    # Build submit flag
    local submit_flag=""
    if [ "${SUBMIT}" == "true" ]; then
        submit_flag="--submit"
        log_warn "SUBMIT mode enabled - results will be sent to Red Hat Connect"
    else
        log_info "Submit mode disabled - results will NOT be sent to Red Hat Connect"
    fi

    # Run preflight check
    log_info "Running preflight check..."
    "${PREFLIGHT_BIN}" check container "${image_ref}" \
        --pyxis-api-token="${RHCC_TOKEN}" \
        --certification-component-id="${RHCC_PROJECT}" \
        --docker-config="${HOME}/.docker/config.json" \
        ${submit_flag}

    log_info "✓ Certification check completed"
}

# Main execution
main() {
    log_info "=== OpenShift Image Certification ==="
    
    # Check for preflight
    if ! check_preflight; then
        log_warn "Preflight not found in PATH"
        log_info "Attempting to install preflight..."
        install_preflight
        if ! check_preflight; then
            log_error "Failed to install or verify preflight"
            exit 1
        fi
    fi

    # Run certification
    certify_image
}

main "$@"

