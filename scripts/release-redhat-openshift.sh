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

operatorhub="${RH_COMMUNITY_OPERATORHUB_REPO_PATH}/operators/mongodb-atlas-kubernetes/${version}"
openshift="${RH_COMMUNITY_OPENSHIFT_REPO_PATH}/operators/mongodb-atlas-kubernetes/${version}"

cd "${RH_COMMUNITY_OPENSHIFT_REPO_PATH}"

git fetch upstream main
git reset --hard upstream/main

cp -r "${operatorhub}" "${openshift}"

git checkout -b "mongodb-atlas-operator-community-${version}"
git add "operators/mongodb-atlas-kubernetes/${version}"
git commit -m "MongoDB Atlas Operator ${version}" --signoff
git push origin "mongodb-atlas-operator-community-${version}"
