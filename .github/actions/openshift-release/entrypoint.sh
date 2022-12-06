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

OPERATOR_PATH=$(pwd)

mkdir -p "../${REPOSITORY}"

git config --global --add safe.directory /github/workspace
gh repo fork --clone "${REPOSITORY}" "../${REPOSITORY}"

REPO_PATH=$(realpath "../${REPOSITORY}/operators/mongodb-atlas-kubernetes")
cd "${REPO_PATH}"
git fetch upstream main
git reset --hard upstream/main

if [ -d "${VERSION}" ]; then
  echo "version already exist in repository"
  exit 1
fi

mkdir "${VERSION}"
cd "${OPERATOR_PATH}"
cp -r bundle.Dockerfile bundle/manifests bundle/metadata bundle/tests "${REPO_PATH}/${VERSION}"

# replace the move instructions in the docker file
cd "${REPO_PATH}"
sed -i.bak 's/COPY bundle\/manifests/COPY manifests/' "${VERSION}/bundle.Dockerfile"
sed -i.bak 's/COPY bundle\/metadata/COPY metadata/' "${VERSION}/bundle.Dockerfile"
sed -i.bak 's/COPY bundle\/tests\/scorecard/COPY tests\/scorecard/' "${VERSION}/bundle.Dockerfile"
rm "${VERSION}/bundle.Dockerfile.bak"

# configure git user
git config --global user.email "41898282+github-actions[bot]@users.noreply.github.com"
git config --global user.name "github-actions[bot]"

# commit, push and open PR
git checkout -b "mongodb-atlas-operator-community-${VERSION}"
git add "${VERSION}"
git status
git commit -m "MongoDB Atlas Operator ${VERSION}" --signoff
# git push origin "mongodb-atlas-operator-community-${VERSION}"
# gh pr create \
#    --title "operator mongodb-atlas-kubernetes (${VERSION})" \
#    --assignee "${ASSIGNEES}"
