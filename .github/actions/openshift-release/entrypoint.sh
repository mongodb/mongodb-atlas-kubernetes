#!/bin/bash

set -eou pipefail

if [ -z "${VERSION+x}" ]; then
  echo "Operator version is not set"
  exit 1
fi

if [ -z "${REPOSITORY+x}" ]; then
  echo "Repository to create PR is not set"
  exit 1
fi

echo $VERSION
echo $REPOSITORY
echo $CERTIFIED

git config --global --add safe.directory /github/workspace

mkdir -p "../${REPOSITORY}"

gh repo fork --clone "${REPOSITORY}" "../${REPOSITORY}"

REPO_PATH="../${REPOSITORY}/operators/mongodb-atlas-kubernetes"

ls -lsa "${REPO_PATH}"

if [ -d "${REPO_PATH}/${VERSION}" ]; then
  echo "version already exist in repository"
  exit 1
fi

mkdir "${REPO_PATH}/${VERSION}"
cp -r bundle.Dockerfile bundle/manifests bundle/metadata bundle/tests "${REPO_PATH}/${VERSION}"

cd "../${REPOSITORY}"
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
