#!/bin/bash

set -eou pipefail

git config --global --add safe.directory /github/workspace

mkdir operator
cd operator
gh repo clone https://github.com/mongodb/mongodb-atlas-kubernetes.git
cd ..
mkdir helm
cd helm
gh repo clone https://github.com/Sugar-pack/helm-charts.git
cd ..
cp -a operator/mongodb-atlas-kubernetes/config/crd/bases/. helm/helm-charts/charts/atlas-operator-crds/templates


git config --global commit.gpgsign true
gh auth setup-git
cd helm/helm-charts
git add -A
BRANCH="atlas-operator-release-${VERSION}"
git checkout -b "$BRANCH"
MESSAGE="Atlas Operator Release ${VERSION}"
git commit -m "$MESSAGE"
git push --set-upstream origin "$BRANCH"