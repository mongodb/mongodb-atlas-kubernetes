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

version=${1:?"pass the version as the parameter, e.g \"0.5.0\""}

PROJECT_ROOT=$(pwd)

if [ -z "${RH_COMMUNITY_OPERATORHUB_REPO_PATH}" ]; then
	echo "RH_COMMUNITY_OPERATORHUB_REPO_PATH is not set"
	exit 1
fi

# Change to repo root for git operations
cleanup() {
  echo "Returning to original directory: ${PROJECT_ROOT}"
  popd > /dev/null 2>&1 || cd "${PROJECT_ROOT}"
}
trap cleanup EXIT
pushd "${RH_COMMUNITY_OPERATORHUB_REPO_PATH}"

git fetch upstream main
git checkout -B "mongodb-atlas-operator-community-${version}" upstream/main

repo="${RH_COMMUNITY_OPERATORHUB_REPO_PATH}/operators/mongodb-atlas-kubernetes"
mkdir -p "${repo}/${version}"
cp -r "${PROJECT_ROOT}/releases/v${version}/bundle.Dockerfile" \
      "${PROJECT_ROOT}/releases/v${version}/bundle/manifests" \
      "${PROJECT_ROOT}/releases/v${version}/bundle/metadata" \
      "${PROJECT_ROOT}/releases/v${version}/bundle/tests" "${repo}/${version}"

# replace the move instructions in the docker file
sed -i.bak 's/COPY bundle\/manifests/COPY manifests/' "${repo}/${version}/bundle.Dockerfile"
sed -i.bak 's/COPY bundle\/metadata/COPY metadata/' "${repo}/${version}/bundle.Dockerfile"
sed -i.bak 's/COPY bundle\/tests\/scorecard/COPY tests\/scorecard/' "${repo}/${version}/bundle.Dockerfile"
rm "${repo}/${version}/bundle.Dockerfile.bak"

export REG_PX="quay.io/"

yq e -i '(.metadata.annotations.containerImage | select(. == (env(REG_PX) + "*") | not)) |= env(REG_PX) + .' \
  "${repo}/${version}/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml"

yq e -i '(.spec.install.spec.deployments[0].spec.template.spec.containers[0].image | select(. == (env(REG_PX) + "*") | not)) |= env(REG_PX) + .' \
  "${repo}/${version}/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml"

cd "${RH_COMMUNITY_OPERATORHUB_REPO_PATH}"
git add "operators/mongodb-atlas-kubernetes/${version}"
git commit -m "operator mongodb-atlas-kubernetes (${version})" --signoff

if [ "${RH_DRYRUN}" == "false" ]; then
  git push -fu origin "mongodb-atlas-operator-community-${version}"
else
  echo "DRYRUN Push (set RH_DRYRUN=true to push for real)"
  git push -fu --dry-run origin "mongodb-atlas-operator-community-${version}"
fi
