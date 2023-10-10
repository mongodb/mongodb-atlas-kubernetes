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


# Replace image version with SHA256
value="${IMG_SHA}" yq e -i '.spec.install.spec.deployments[0].spec.template.spec.containers[0].image = env(value)' \
  "${REPO}/${VERSION}"/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

# Add skip range
value='">=0.8.0"' yq e -i '.spec.skipRange = env(value)' \
  "${REPO}/${VERSION}"/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

cd "${REPO}"
git checkout -b "mongodb-atlas-kubernetes-operator-${VERSION}"
git add "${REPO}/${VERSION}"
git commit -m "operator mongodb-atlas-kubernetes (${VERSION})" --signoff
git push -u origin "mongodb-atlas-operator-community-${VERSION}"
cd -

