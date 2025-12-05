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
repo="${RH_COMMUNITY_OPERATORHUB_REPO_PATH}/operators/mongodb-atlas-kubernetes"
mkdir "${repo}/${version}"
cp -r releases/v${version}/bundle.Dockerfile releases/v${version}/bundle/manifests releases/v${version}/bundle/metadata releases/v${version}/bundle/tests "${repo}/${version}"

cd "${repo}"
git fetch upstream main
git reset --hard upstream/main

# replace the move instructions in the docker file
sed -i.bak 's/COPY bundle\/manifests/COPY manifests/' "${version}/bundle.Dockerfile"
sed -i.bak 's/COPY bundle\/metadata/COPY metadata/' "${version}/bundle.Dockerfile"
sed -i.bak 's/COPY bundle\/tests\/scorecard/COPY tests\/scorecard/' "${version}/bundle.Dockerfile"
rm "${version}/bundle.Dockerfile.bak"

yq e -i '.metadata.annotations.containerImage = "quay.io/" + .metadata.annotations.containerImage' \
  "${repo}/${version}"/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

yq e -i '.spec.install.spec.deployments[0].spec.template.spec.containers[0].image = "quay.io/" + .spec.install.spec.deployments[0].spec.template.spec.containers[0].image' \
  "${repo}/${version}"/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

# commit
git checkout -b "mongodb-atlas-operator-community-${version}"
git add "${version}"
git commit -m "MongoDB Atlas Operator ${version}" --signoff
git push origin "mongodb-atlas-operator-community-${version}"
