#!/bin/bash

set -eou pipefail

REPO="${RH_CERTIFIED_OPENSHIFT_REPO_PATH}/operators/mongodb-atlas-kubernetes"

cd "${REPO}"
git checkout main
git fetch origin main
git reset --hard origin/main
mkdir -p "${REPO}/${VERSION}"
cd -

pwd

cp -r bundle.Dockerfile bundle/manifests bundle/metadata bundle/tests "${REPO}/${VERSION}"

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
yq e -i '.annotations = .annotations + { "com.redhat.openshift.versions": "v4.8" }' \
  "${REPO}/${VERSION}"/metadata/annotations.yaml

cd "${REPO}"
git checkout -b origin main
git pull --rebase upstream main
git checkout -b "mongodb-atlas-kubernetes-operator-${VERSION}"
git add "${REPO}/${VERSION}"
git commit -m "operator mongodb-atlas-kubernetes (${VERSION})" --signoff
git push -u origin "mongodb-atlas-kubernetes-operator-${VERSION}"
cd -
