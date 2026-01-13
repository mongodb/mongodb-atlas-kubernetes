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


set -euo pipefail

reposDir=${1:-$REPOS_DIR}

repos=(
    "k8s-operatorhub/community-operators"
    "redhat-openshift-ecosystem/community-operators-prod"
    "redhat-openshift-ecosystem/certified-operators"
)

mkdir -p "${reposDir}"
pushd "${reposDir}"
for repo in "${repos[@]}"; do
	mirror_repo=$(basename "${repo}")
	echo "Cloning ${mirror_repo} from ${repo}"
	git clone --depth 1 "https://github.com/mongodb-forks/${mirror_repo}.git"
	pushd "${mirror_repo}"
	git remote add upstream "https://github.com/${repo}"
	popd
done
