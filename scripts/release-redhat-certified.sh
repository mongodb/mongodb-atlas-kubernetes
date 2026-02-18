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


set -eou pipefail

version=${1:-$VERSION}

PROJECT_ROOT=$(pwd)

if [ -z "${version}" ]; then
	echo "version is not set as arguiment or VERSION env var"
	exit 1
fi

if [ -z "${RH_CERTIFIED_OPENSHIFT_REPO_PATH}" ]; then
	echo "RH_CERTIFIED_OPENSHIFT_REPO_PATH is not set"
	exit 1
fi

echo -n "Determining SHA for arm64 ... "
IMG_SHA_ARM64=$(docker \
  manifest inspect "quay.io/mongodb/mongodb-atlas-kubernetes-operator:${version}-certified" |
  jq --raw-output '.manifests[] | select(.platform.architecture == "arm64") | .digest')
echo "${IMG_SHA_ARM64}"

echo -n "Determining SHA for amd64 ... "
IMG_SHA_AMD64=$(docker \
  manifest inspect "quay.io/mongodb/mongodb-atlas-kubernetes-operator:${version}-certified" |
  jq --raw-output '.manifests[] | select(.platform.architecture == "amd64") | .digest')
echo "${IMG_SHA_AMD64}"

REPO="${RH_CERTIFIED_OPENSHIFT_REPO_PATH}/operators/mongodb-atlas-kubernetes"

# Change to repo root for git operations
cleanup() {
  echo "Returning to original directory: ${PROJECT_ROOT}"
  popd > /dev/null 2>&1 || cd "${PROJECT_ROOT}"
}
trap cleanup EXIT
pushd "${RH_CERTIFIED_OPENSHIFT_REPO_PATH}"

git checkout main
git fetch origin main
git fetch upstream main

# CRITICAL: Reset completely to upstream/main to ensure we're identical to upstream
# This ensures we aren't "carrying" any old differences from our fork
# Workflow files will match upstream exactly, so they won't show as changes
git reset --hard upstream/main

# Create branch from upstream/main state
git checkout -B "mongodb-atlas-kubernetes-operator-${version}"

mkdir -p "${REPO}/${version}"

cp -r "${PROJECT_ROOT}/releases/v${version}/bundle.Dockerfile" \
      "${PROJECT_ROOT}/releases/v${version}/bundle/manifests" \
      "${PROJECT_ROOT}/releases/v${version}/bundle/metadata" \
      "${PROJECT_ROOT}/releases/v${version}/bundle/tests" "${REPO}/${version}"

# Replace deployment image version with SHA256
value="${IMG_SHA_AMD64}" yq e -i '.spec.install.spec.deployments[0].spec.template.spec.containers[0].image = "quay.io/mongodb/mongodb-atlas-kubernetes-operator@" + env(value)' \
  "${REPO}/${version}"/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

# set related images
yq e -i '.spec = { "relatedImages": [ { "name": "mongodb-atlas-kubernetes-operator-arm64" }, { "name": "mongodb-atlas-kubernetes-operator-amd64" } ] } + .spec' \
  "${REPO}/${version}"/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

value="${IMG_SHA_ARM64}" yq e -i '.spec.relatedImages[0].image = "quay.io/mongodb/mongodb-atlas-kubernetes-operator@" + env(value)' \
  "${REPO}/${version}"/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

value="${IMG_SHA_AMD64}" yq e -i '.spec.relatedImages[1].image = "quay.io/mongodb/mongodb-atlas-kubernetes-operator@" + env(value)' \
  "${REPO}/${version}"/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

# set containerImage annotation
value="${IMG_SHA_AMD64}" yq e -i '.metadata.annotations.containerImage = "quay.io/mongodb/mongodb-atlas-kubernetes-operator@" + env(value)' \
  "${REPO}/${version}"/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

# set openshift versions
yq e -i '.annotations = .annotations + { "com.redhat.openshift.versions": "v4.8-v4.18" }' \
  "${REPO}/${version}"/metadata/annotations.yaml

# CRITICAL: Ensure workflow files match upstream exactly (no diff)
# This ensures workflow files won't be included in our commit diff
cd "${RH_CERTIFIED_OPENSHIFT_REPO_PATH}"
git checkout upstream/main -- .github/ || true

# Commit ONLY operator changes (workflow files are already identical to upstream, so no diff)

# Add dummy file to force a diff for permission testing
repo="${RH_CERTIFIED_OPENSHIFT_REPO_PATH}/operators/mongodb-atlas-kubernetes"
echo "dummy change $(date)" > "${repo}/${version}/dummy.txt"

git add "operators/mongodb-atlas-kubernetes/${version}"
git commit -m "operator mongodb-atlas-kubernetes (${version})" --signoff || true

# Verify that our commit only includes operator changes, not workflow files
if git diff --name-only upstream/main HEAD | grep -q "^\.github/"; then
	echo "WARNING: Commit includes workflow file changes. This may cause push to fail."
	echo "Workflow files in commit:"
	git diff --name-only upstream/main HEAD | grep "^\.github/"
fi

if [ "${RH_DRYRUN}" == "false" ]; then
  # Push - should only push operator changes since workflow files match upstream exactly
  git push -u origin "mongodb-atlas-kubernetes-operator-${version}" --force
else
  echo "DRYRUN Push (set RH_DRYRUN=true to push for real)"
  git push -fu --dry-run -u origin "mongodb-atlas-kubernetes-operator-${version}"
fi
