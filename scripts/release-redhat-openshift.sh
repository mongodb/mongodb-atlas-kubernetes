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

if [ -z "${RH_COMMUNITY_OPERATORHUB_REPO_PATH}" ]; then
	echo "RH_COMMUNITY_OPERATORHUB_REPO_PATH is not set"
	exit 1
fi

operatorhub="${RH_COMMUNITY_OPERATORHUB_REPO_PATH}/operators/mongodb-atlas-kubernetes/${version}"
openshift="${RH_COMMUNITY_OPENSHIFT_REPO_PATH}/operators/mongodb-atlas-kubernetes/${version}"

cd "${RH_COMMUNITY_OPENSHIFT_REPO_PATH}"

# Fetch latest from both upstream and fork
git fetch upstream main
git fetch origin main

# CRITICAL: Reset completely to upstream/main to ensure we're identical to upstream
# This ensures we aren't "carrying" any old differences from our fork
# Workflow files will match upstream exactly, so they won't show as changes
git reset --hard upstream/main

# Create branch from upstream/main state
git checkout -b "mongodb-atlas-operator-community-${version}" || git checkout "mongodb-atlas-operator-community-${version}"
git reset --hard upstream/main

# Copy operator from community-operators repo
cp -r "${operatorhub}" "${openshift}"

# CRITICAL: Ensure workflow files match upstream exactly (no diff)
# This ensures workflow files won't be included in our commit diff
git checkout upstream/main -- .github/ || true

# Commit ONLY operator changes (workflow files are already identical to upstream, so no diff)
git add "operators/mongodb-atlas-kubernetes/${version}"
git commit -m "operator mongodb-atlas-kubernetes (${version})" --signoff

# Verify that our commit only includes operator changes, not workflow files
if git diff --name-only upstream/main HEAD | grep -q "^\.github/"; then
	echo "WARNING: Commit includes workflow file changes. This may cause push to fail."
	echo "Workflow files in commit:"
	git diff --name-only upstream/main HEAD | grep "^\.github/"
fi

if [ "${RH_DRYRUN}" == "false" ]; then
  # Push - should only push operator changes since workflow files match upstream exactly
  git push origin "mongodb-atlas-operator-community-${version}" --force
else
  echo "DRYRUN Push (set RH_DRYRUN=true to push for real)"
  git push -fu --dry-run origin "mongodb-atlas-operator-community-${version}"
fi
