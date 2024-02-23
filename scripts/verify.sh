#!/bin/bash

set -euo pipefail

REPO=${IMG_REPO:-docker.io/mongodb/mongodb-atlas-kubernetes-operator-prerelease}
img_to_verify=${IMG:-$REPO:$VERSION}
SIGNATURE_REPO=${SIGNATURE_REPO:-$REPO}

KEY_FILE=${KEY_FILE:-ako.pem}

COSIGN_REPOSITORY="${SIGNATURE_REPO}" cosign verify \
  --key="${KEY_FILE}" "${img_to_verify}"
