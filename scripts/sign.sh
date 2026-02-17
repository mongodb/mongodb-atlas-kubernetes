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
img=${IMG:-$REPO:$VERSION}
SIGNATURE_REPO=${SIGNATURE_REPO:-$REPO}
TMPDIR=${TMPDIR:-/tmp}

# Useful for setups with credential helpers
SIGNING_DOCKERCFG_BASE64=${SIGNING_DOCKERCFG_BASE64:-EMPTY}
if [ "${SIGNING_DOCKERCFG_BASE64}" != "EMPTY" ]; then
  DOCKER_CFG="${TMPDIR}/signing-docker-config.json"
  echo "${SIGNING_DOCKERCFG_BASE64}" | base64 -d > "${DOCKER_CFG}"
fi

DOCKER_CFG=${DOCKER_CFG:-~/.docker/config.json}

SIGNING_ENVFILE="${TMPDIR}/signing-envfile"

{
  echo "GRS_CONFIG_USER1_USERNAME=${GRS_USERNAME}";
  echo "GRS_CONFIG_USER1_PASSWORD=${GRS_PASSWORD}";
  echo "COSIGN_REPOSITORY=${SIGNATURE_REPO}";
}  > "${SIGNING_ENVFILE}"

docker run \
  --platform linux/amd64 \
  --env-file="${SIGNING_ENVFILE}" \
  -v "${DOCKER_CFG}:/root/.docker/config.json" \
  --rm \
  -v "$(pwd):$(pwd)" \
  -w "$(pwd)" \
  artifactory.corp.mongodb.com/release-tools-container-registry-local/garasign-cosign \
  cosign sign --key "${PKCS11_URI}" \
  --tlog-upload=false --use-signing-config=false --new-bundle-format=false "${img}" && \
  echo "✍️  Signed"
