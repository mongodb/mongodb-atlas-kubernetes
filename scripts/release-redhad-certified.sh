#!/bin/bash

set -eou pipefail

# Make sure Version is specified
# Make sure Repo exists
# Reset repo and fetch from main
# Copy manifests
# Replace image tag with image SHA

if [ -z "${VERSION+x}" ]; then
  echo "VERSION is not set"
  exit 1
fi

if [ -z "${TAG_IMG}" ]; then
  echo "TAG_IMAGE is not set"
  exit 1
fi

if [ -z "${SHA_IMG}" ]; then
  echo "SHA_IMG is not set"
  exit 1
fi

if [ -z "${RH_API_TOKEN}" ]; then
  echo "SHA_IMG is not set"
  exit 1
fi

if [ -z "${RH_PROJECT_ID}" ]; then
  echo "SHA_IMG is not set"
  exit 1
fi

# Path to https://github.com/mongodb-forks/community-operators
if [ -z "${RH_CERTIFIED_OPERATORS_FORK}" ]; then
  echo "RH_CERTIFIED_OPERATORS_FORK is not set"
  exit 1
fi

preflight --version
podman --version


REPO="${RH_CERTIFIED_OPERATORS_FORK}/operators/mongodb-atlas-kubernetes"
mkdir "${REPO}/${VERSION}"

cp -r bundle.Dockerfile bundle/manifests bundle/metadata bundle/tests "${REPO}/${VERSION}"

cd "${REPO}"
git fetch origin main
git reset --hard origin/main

podman login

# Do the preflight check
preflight check container "${SHA_IMG}" --submit --pyxis-api-token="${RH_API_TOKEN}" --certification-project-id="${RH_PROJECT_ID}"


# Replace image version with SHA256
value="${SHA_IMG}" yq e -i '.spec.install.spec.deployments[0].spec.template.spec.containers[0].image = env(value)' \
  bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

# Add skip range
value='">=0.8.0"' yq e -i '.spec.skipRange = env' \
  bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

git checkout -b "mongodb-atlas-operator-community-${VERSION}"
git add "${REPO}/${VERSION}"
git commit -m "operator mongodb-atlas-kubernetes (${VERSION})" --signoff
git push origin "mongodb-atlas-operator-community-${VERSION}"