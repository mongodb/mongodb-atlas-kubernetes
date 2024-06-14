#!/usr/bin/env bash

set -eou pipefail

###
# This script is responsible for clearing the Silk asset metadata
#
# See: https://docs.devprod.prod.corp.mongodb.com/mms/python/src/sbom/silkbomb/docs/SILK
###

asset_group="$1"
[ -z "${asset_group}" ] && echo "Missing asset_group def" && exit 1


if [ -f "$HOME/.silk.env" ]; then
    # shellcheck source=/dev/null
    source "$HOME/.silk.env"
fi

if [ -z "${SILK_CLIENT_SECRET}" ]; then
    echo "Need SILK_CLIENT_SECRET env var" >&2
    exit 1
fi

if [ -z "${SILK_CLIENT_ID}" ]; then
    echo "Need SILK_CLIENT_ID env var" >&2
    exit 1
fi

TOKEN=$(curl --silent -X 'POST' \
  'https://silkapi.us1.app.silk.security/api/v1/authenticate' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d "{
  \"client_id\": \"${SILK_CLIENT_ID}\",
  \"client_secret\": \"${SILK_CLIENT_SECRET}\"
}" | jq -r '.token')

JSON_PAYLOAD=$(
  cat << EOF
{
  "metadata": {}
}
EOF
)
curl -X PATCH \
  "https://silkapi.us1.app.silk.security/api/v1/raw/asset_group/${asset_group}" \
  -H "accept: application/json" -H "Authorization: ${TOKEN}" \
  -H 'Content-Type: application/json' -d "${JSON_PAYLOAD}"

