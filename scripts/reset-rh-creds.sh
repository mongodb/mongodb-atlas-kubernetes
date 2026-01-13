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

if [ -z "${GH_TOKEN}" ]; then
	echo "GH_TOKEN is not set"
	exit 1
fi

repos=(
    "community-operators"
    "community-operators-prod"
    "certified-operators"
)

mkdir -p "${reposDir}"
pushd "${reposDir}"
for repo in "${repos[@]}"; do
	pushd "${repo}"
	git config --local --unset-all http.https://github.com/.extraheader || true
	set +x
	git remote set-url origin "https://x-access-token:${GH_TOKEN}@github.com/${repo}.git"
	# set -x # if needed
	popd
done
popd

