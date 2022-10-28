#!/bin/bash

set -eou pipefail

if [ -z "${IMAGE+x}" ]; then
  echo "IMAGE is not set"
  exit 1
fi

if [ -z "${VERSION+x}" ]; then
  echo "VERSION is not set"
  exit 1
fi

# Path to https://github.com/mongodb-forks/certified-operators
if [ -z "${RH_CERTIFIED_OPENSHIFT_REPO_PATH+x}" ]; then
  echo "RH_CERTIFIED_OPENSHIFT_REPO_PATH is not set"
  exit 1
fi

if [ -z "${RH_CERTIFICATION_OSPID+x}" ]; then
  echo "RH_CERTIFICATION_OSPID is not set"
  exit 1
fi

if [ -z "${REGISTRY_TOKEN+x}" ]; then
  echo "REGISTRY_TOKEN is not set"
  exit 1
fi

if [ -z "${RH_CERTIFICATION_PYXIS_API_TOKEN+x}" ]; then
  echo "RH_CERTIFICATION_PYXIS_API_TOKEN is not set"
  exit 1
fi

if [ -z "${CONTAINER_ENGINE+x}" ]; then
  echo "CONTAINER_ENGINE is not set, defaulting to podman"
  CONTAINER_ENGINE=podman
fi

preflight --version
${CONTAINER_ENGINE} --version

${CONTAINER_ENGINE} login -u unused -p "${REGISTRY_TOKEN}" quay.io --authfile ./authfile.json

REPO="${RH_CERTIFIED_OPENSHIFT_REPO_PATH}/operators/mongodb-atlas-kubernetes"

cd "${REPO}"
git checkout main
git fetch origin main
git reset --hard origin/main
mkdir -p "${REPO}/${VERSION}"
cd -

pwd

cp -r bundle.Dockerfile bundle/manifests bundle/metadata bundle/tests "${REPO}/${VERSION}"

IMG_SHA=$("${CONTAINER_ENGINE}" inspect --format='{{ index .RepoDigests 0}}' "${IMAGE}":"${VERSION}")

# Do the preflight check first
preflight check container "${IMG_SHA}" --docker-config=./authfile.json

# Send results to RedHat if preflight finished without errors
preflight check container "${IMG_SHA}" \
  --submit \
  --pyxis-api-token="${RH_CERTIFICATION_PYXIS_API_TOKEN}" \
  --certification-project-id="${RH_CERTIFICATION_OSPID}" \
  --docker-config=./authfile.json

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
git push mongodb "mongodb-atlas-operator-community-${VERSION}"
cd -
