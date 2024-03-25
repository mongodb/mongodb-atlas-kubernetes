#!/bin/bash

set -eou pipefail

version=${1:?"pass the version as the parameter, e.g \"0.5.0\""}
repo="${RH_COMMUNITY_OPERATORHUB_REPO_PATH}/operators/mongodb-atlas-kubernetes"
mkdir "${repo}/${version}"
cp -r bundle.Dockerfile bundle/manifests bundle/metadata bundle/tests "${repo}/${version}"

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
