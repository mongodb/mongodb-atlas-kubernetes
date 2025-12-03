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

echo -n "Determining SHA for arm64 ... "
IMG_SHA_ARM64=$(docker \
  manifest inspect quay.io/mongodb/mongodb-atlas-kubernetes-operator:${VERSION}-certified |
  jq --raw-output '.manifests[] | select(.platform.architecture == "arm64") | .digest')
echo ${IMG_SHA_ARM64}

echo -n "Determining SHA for amd64 ... "
IMG_SHA_AMD64=$(docker \
  manifest inspect quay.io/mongodb/mongodb-atlas-kubernetes-operator:${VERSION}-certified |
  jq --raw-output '.manifests[] | select(.platform.architecture == "amd64") | .digest')
echo ${IMG_SHA_AMD64}

REPO="${RH_CERTIFIED_OPENSHIFT_REPO_PATH}/operators/mongodb-atlas-kubernetes"

cd "${REPO}"
git checkout main
git fetch origin main
git reset --hard origin/main
mkdir -p "${REPO}/${VERSION}"
cd -

pwd

cp -r releases/v${VERSION}/bundle.Dockerfile releases/v${VERSION}/bundle/manifests releases/v${VERSION}/bundle/metadata releases/v${VERSION}/bundle/tests "${REPO}/${VERSION}"

# Replace deployment image version with SHA256
value="${IMG_SHA_AMD64}" yq e -i '.spec.install.spec.deployments[0].spec.template.spec.containers[0].image = "quay.io/mongodb/mongodb-atlas-kubernetes-operator@" + env(value)' \
  "${REPO}/${VERSION}"/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

# set related images
yq e -i '.spec = { "relatedImages": [ { "name": "mongodb-atlas-kubernetes-operator-arm64" }, { "name": "mongodb-atlas-kubernetes-operator-amd64" } ] } + .spec' \
  "${REPO}/${VERSION}"/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

value="${IMG_SHA_ARM64}" yq e -i '.spec.relatedImages[0].image = "quay.io/mongodb/mongodb-atlas-kubernetes-operator@" + env(value)' \
  "${REPO}/${VERSION}"/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

value="${IMG_SHA_AMD64}" yq e -i '.spec.relatedImages[1].image = "quay.io/mongodb/mongodb-atlas-kubernetes-operator@" + env(value)' \
  "${REPO}/${VERSION}"/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

# set containerImage annotation
value="${IMG_SHA_AMD64}" yq e -i '.metadata.annotations.containerImage = "quay.io/mongodb/mongodb-atlas-kubernetes-operator@" + env(value)' \
  "${REPO}/${VERSION}"/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

# set openshift versions
yq e -i '.annotations = .annotations + { "com.redhat.openshift.versions": "v4.8-v4.18" }' \
  "${REPO}/${VERSION}"/metadata/annotations.yaml

cd "${REPO}"
git checkout -b "mongodb-atlas-kubernetes-operator-${VERSION}" origin/main
git pull --rebase upstream main
git add "${REPO}/${VERSION}"
git commit -m "operator mongodb-atlas-kubernetes (${VERSION})" --signoff
git push -u origin "mongodb-atlas-kubernetes-operator-${VERSION}"
cd -
