#!/bin/bash

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
