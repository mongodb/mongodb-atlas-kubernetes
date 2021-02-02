#!/bin/sh
# act -j build-push
export MCLI_OPS_MANAGER_URL="https://cloud-qa.mongodb.com/"
export MCLI_PUBLIC_API_KEY=$(grep "ATLAS_PUBLIC_KEY" .actrc | cut -d "=" -f 2)
export MCLI_PRIVATE_API_KEY=$(grep "ATLAS_PRIVATE_KEY" .actrc | cut -d "=" -f 2)
export MCLI_ORG_ID=$(grep "ATLAS_ORG_ID" .actrc | cut -d "=" -f 2)

export INPUT_IMAGE_URL=$(grep "DOCKER_REPO" .actrc | cut -d "=" -f 2):$(git rev-parse --abbrev-ref HEAD)-$(git rev-parse --short HEAD)
export INPUT_ENV=dev
./.github/actions/gen-install-scripts/entrypoint.sh

ginkgo -v -x test/e2e