#!/bin/sh

set -eou pipefail

if [ -z "${VERSION+x}" ]; then
  echo "Operator version is not set"
  exit 1
fi

if [ -z "${REPOSITORY+x}" ]; then
  echo "Repository to create PR is not set"
  exit 1
fi

git config --global --add safe.directory /github/workspace

gh repo fork --clone "${REPOSITORY}" repository

REPO_PATH="repository/operators/mongodb-atlas-kubernetes"
mkdir "${REPO_PATH}/${VERSION}"
cp -r bundle.Dockerfile bundle/manifests bundle/metadata bundle/tests "${REPO_PATH}/${VERSION}"

cd "${REPO_PATH}"
git fetch upstream main
git reset --hard upstream/main

# replace the move instructions in the docker file
sed -i.bak 's/COPY bundle\/manifests/COPY manifests/' "${VERSION}/bundle.Dockerfile"
sed -i.bak 's/COPY bundle\/metadata/COPY metadata/' "${VERSION}/bundle.Dockerfile"
sed -i.bak 's/COPY bundle\/tests\/scorecard/COPY tests\/scorecard/' "${VERSION}/bundle.Dockerfile"
rm "${VERSION}/bundle.Dockerfile.bak"

# commit
git checkout -b "mongodb-atlas-operator-community-${VERSION}"
git add "operators/mongodb-atlas-kubernetes/${VERSION}"
git commit -m "MongoDB Atlas Operator ${VERSION}" --signoff
# git push origin "mongodb-atlas-operator-community-${VERSION}"
