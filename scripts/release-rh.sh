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

version=${1:-$VERSION}

if [ -z "${version}" ]; then
	echo "missing version arg or VERSION env var"
	exit 1
fi

vars=(
	RH_COMMUNITY_OPERATORHUB_REPO_PATH
	RH_COMMUNITY_OPENSHIFT_REPO_PATH
	RH_CERTIFIED_OPENSHIFT_REPO_PATH
)
for envar in "${vars[@]}"; do
	if [ -z "${envar}" ]; then
		echo "missing required ${envar} env var"
		exit 1
	fi
done

echo "Releasing to RedHat: ${version}"

./scripts/release-redhat.sh "${version}"
./scripts/release-redhat-openshift.sh "${version}"
./scripts/release-redhat-certified.sh "${version}"

echo "All releases PR posted successfully."

