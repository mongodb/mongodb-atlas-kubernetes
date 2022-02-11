#!/bin/bash

set -eo pipefail
focus_key=$1

public_key=$(grep "ATLAS_PUBLIC_KEY" .actrc | cut -d "=" -f 2)
private_key=$(grep "ATLAS_PRIVATE_KEY" .actrc | cut -d "=" -f 2)
org_id=$(grep "ATLAS_ORG_ID" .actrc | cut -d "=" -f 2)
# this is the format how it's pushed by act -j build-push
branch="$(git rev-parse --abbrev-ref HEAD | awk '{sub(/\//, "_"); print}')"
commit=$(git rev-parse --short HEAD)
image=$(grep "DOCKER_REGISTRY" .env | cut -d "=" -f 2)/$(grep "DOCKER_REPO" .env | cut -d "=" -f 2):${branch}-${commit}
echo "Using docker image: ${image}"

export INPUT_IMAGE_URL_DOCKER="${image}"
export INPUT_ENV=dev
./.github/actions/gen-install-scripts/entrypoint.sh
# TODO temporary change line

awk '{gsub(/cloud.mongodb.com/, "cloud-qa.mongodb.com", $0); print}' bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml > tmp && mv tmp bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

docker build -t "${image}" .
docker push "${image}"

#bundles
bundle_image=$(grep "DOCKER_REGISTRY" .env | cut -d "=" -f 2)/$(grep "DOCKER_BUNDLES_REPO" .env | cut -d "=" -f 2):${branch}-${commit} #Registry is nessary
export BUNDLE_IMAGE="${bundle_image}"
docker build -f bundle.Dockerfile -t "${bundle_image}" .
docker push "${bundle_image}"

export MCLI_OPS_MANAGER_URL="https://cloud-qa.mongodb.com/"
export MCLI_PUBLIC_API_KEY="${public_key}"
export MCLI_PRIVATE_API_KEY="${private_key}"
export MCLI_ORG_ID="${org_id}"
export IMAGE_URL="${image}" #for helm chart
ginkgo -focus "${focus_key}" -nodes=3 -stream -v test/e2e/
