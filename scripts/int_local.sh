#!/bin/bash

set -eou pipefail

if [ -z "${label+x}" ]; then
  label=""
fi


public_key=$(grep "ATLAS_PUBLIC_KEY" .actrc | cut -d "=" -f 2)
private_key=$(grep "ATLAS_PRIVATE_KEY" .actrc | cut -d "=" -f 2)
org_id=$(grep "ATLAS_ORG_ID" .actrc | cut -d "=" -f 2)

export MCLI_OPS_MANAGER_URL="${MCLI_OPS_MANAGER_URL:-https://cloud-qa.mongodb.com/}"
export MCLI_PUBLIC_API_KEY="${MCLI_PUBLIC_API_KEY:-$public_key}"
export MCLI_PRIVATE_API_KEY="${MCLI_PRIVATE_API_KEY:-$private_key}"
export MCLI_ORG_ID="${MCLI_ORG_ID:-$org_id}"

AKO_INT_TEST=1 ginkgo run --race --label-filter="${label}" --timeout 80m -v ./test/int ./test/int/clusterwide -coverprofile cover.out
