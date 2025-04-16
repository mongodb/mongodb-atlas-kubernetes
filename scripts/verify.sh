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


set -euo pipefail

REPO=${IMG_REPO:-docker.io/mongodb/mongodb-atlas-kubernetes-operator-prerelease}
img_to_verify=${IMG:-$REPO:$VERSION}
SIGNATURE_REPO=${SIGNATURE_REPO:-$REPO}

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

KEY_FILE=${KEY_FILE:-ako.pem}

COSIGN_REPOSITORY="${SIGNATURE_REPO}" "${SCRIPT_DIR}"/retry.sh cosign verify \
  --insecure-ignore-tlog --key="${KEY_FILE}" "${img_to_verify}"
