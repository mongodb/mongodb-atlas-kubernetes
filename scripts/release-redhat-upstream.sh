#!/bin/bash

set -eou pipefail

version=${1:?"pass the version as the parameter, e.g \"0.5.0\""}

cd "${RH_COMMUNITY_REPO_PATH}"
cp -r community-operators/mongodb-atlas-kubernetes/"${version}" upstream-community-operators/mongodb-atlas-kubernetes
git checkout -b "mongodb-atlas-operator-upstream-${version}"
git add upstream-community-operators/mongodb-atlas-kubernetes/"${version}"
git commit -m "[community] MongoDB Atlas Operator ${version}" --signoff upstream-community-operators/mongodb-atlas-kubernetes/"${version}"
git push origin mongodb-atlas-operator-upstream-"${version}"