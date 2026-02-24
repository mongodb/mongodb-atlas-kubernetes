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

if [ -z "${RH_COMMUNITY_OPENSHIFT_REPO_PATH}" ]; then
    echo "RH_COMMUNITY_OPENSHIFT_REPO_PATH is not set"
    exit 1
fi

operatorhub="${RH_COMMUNITY_OPERATORHUB_REPO_PATH}/operators/mongodb-atlas-kubernetes/${version}"
openshift="${RH_COMMUNITY_OPENSHIFT_REPO_PATH}/operators/mongodb-atlas-kubernetes/${version}"

# Change to OpenShift repo root
cleanup() {
  echo "Returning to original directory: ${PROJECT_ROOT}"
  popd > /dev/null 2>&1 || cd "${PROJECT_ROOT}"
}
trap cleanup EXIT
pushd "${RH_COMMUNITY_OPENSHIFT_REPO_PATH}"

git fetch upstream main
git checkout -B "mongodb-atlas-operator-community-${version}" upstream/main

# Copy operator from community-operators repo
# Remove destination if it exists to avoid nested directory structure
rm -rf "${openshift}"
cp -r "${operatorhub}" "${openshift}"

git add "operators/mongodb-atlas-kubernetes/${version}"
git commit -m "operator mongodb-atlas-kubernetes (${version})" --signoff

if [ "${RH_DRYRUN}" == "false" ]; then
  git push -fu origin "mongodb-atlas-operator-community-${version}"
else
  echo "DRYRUN Push (set RH_DRYRUN=true to push for real)"
  git push -fu --dry-run origin "mongodb-atlas-operator-community-${version}"
fi
