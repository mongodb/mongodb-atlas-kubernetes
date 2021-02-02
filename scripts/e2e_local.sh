#!/bin/sh
act -j build-push

public_key=$(grep "ATLAS_PUBLIC_KEY" .actrc | cut -d "=" -f 2)
private_key=$(grep "ATLAS_PRIVATE_KEY" .actrc | cut -d "=" -f 2)
org_id=$(grep "ATLAS_ORG_ID" .actrc | cut -d "=" -f 2)
image=$(grep "DOCKER_REPO" .actrc | cut -d "=" -f 2):$(git rev-parse --abbrev-ref HEAD)-$(git rev-parse --short HEAD)

export MCLI_OPS_MANAGER_URL="https://cloud-qa.mongodb.com/"
export MCLI_PUBLIC_API_KEY="${public_key}"
export MCLI_PRIVATE_API_KEY="${private_key}"
export MCLI_ORG_ID=org_id="${org_id}"
export INPUT_IMAGE_URL="${image}"
export INPUT_ENV=dev

./.github/actions/gen-install-scripts/entrypoint.sh

ginkgo -v -x test/e2e