#!/bin/bash
# Copyright 2025 MongoDB Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -eo pipefail
focus_key=$1
build=$2
if [[ -z "${build:-}" ]]; then
  build="true"
fi

public_key=$(grep "ATLAS_PUBLIC_KEY" .actrc | cut -d "=" -f 2)
private_key=$(grep "ATLAS_PRIVATE_KEY" .actrc | cut -d "=" -f 2)
org_id=$(grep "ATLAS_ORG_ID" .actrc | cut -d "=" -f 2)

# this is the format how it's pushed by act -j build-push
branch="$(git rev-parse --abbrev-ref HEAD | awk '{sub(/\//, "-"); print}')"
commit=$(git rev-parse --short HEAD)
image=$(grep "DOCKER_REGISTRY" .env | cut -d "=" -f 2)/$(grep "DOCKER_REPO" .env | cut -d "=" -f 2):${branch}-${commit}
echo "Using docker image: ${image}"
bundle_image=$(grep "DOCKER_REGISTRY" .env | cut -d "=" -f 2)/$(grep "DOCKER_BUNDLES_REPO" .env | cut -d "=" -f 2):${branch}-${commit} #Registry is nessary
export BUNDLE_IMAGE="${bundle_image}"
export INPUT_IMAGE_URL="${image}"
export INPUT_ENV=dev

if [[ "${build}" == "true" ]]; then
    ./.github/actions/gen-install-scripts/entrypoint.sh
    awk '{gsub(/cloud.mongodb.com/, "cloud-qa.mongodb.com", $0); print}' bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml > yaml.tmp && mv yaml.tmp bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml

    make all-platforms
    docker build -f fast.Dockerfile -t "${image}" .
    docker push "${image}"

    #bundles
    docker build -f bundle.Dockerfile -t "${bundle_image}" .
    docker push "${bundle_image}"
fi

kubectl apply -f deploy/crds

export MCLI_OPS_MANAGER_URL="${MCLI_OPS_MANAGER_URL:-https://cloud-qa.mongodb.com/}"
export MCLI_PUBLIC_API_KEY="${MCLI_PUBLIC_API_KEY:-$public_key}"
export MCLI_PRIVATE_API_KEY="${MCLI_PRIVATE_API_KEY:-$private_key}"
export MCLI_ORG_ID="${MCLI_ORG_ID:-$org_id}"
export IMAGE_URL="${image}" #for helm chart
AKO_E2E_TEST=1 ginkgo --race --label-filter="${focus_key}" --timeout 120m -vv test/e2e/
