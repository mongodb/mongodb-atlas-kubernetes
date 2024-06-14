#!/usr/bin/env bash
set -Eeou pipefail

###
# This script is responsible for creating Silk Assets if needed.
#
# See: https://docs.devprod.prod.corp.mongodb.com/mms/python/src/sbom/silkbomb/docs/SILK
###

ASSET_GROUP_ID=""
ASSET_GROUP_DESCRIPTION=""

function usage() {
  echo "Creates a new Asset Group in Silk.
Usage:
  create_asset_group.sh [-h]
  create_asset_group.sh -a mongodb-enterprise-kubectl-cli -d \"MongoDB Atlas Operator\"
Example:

Options:
  -h                     (optional) Shows this screen.
  -a <asset group>       (required) The name of the Asset Group.
  -d <asset description> (required) Asset Description.
"
}

function validate() {
  if [ -z "$ASSET_GROUP_ID" ]; then
    echo "Missing Asset Group parameter"
    usage
    exit 1
  fi
  if [ -z "$ASSET_GROUP_DESCRIPTION" ]; then
    echo "Missing Asset Group description parameter"
    usage
    exit 1
  fi
  if [ -z "$SILK_CLIENT_SECRET" ]; then
    echo "Missing SILK_CLIENT_SECRET env variable"
    usage
    exit 1
  fi
  if [ -z "$SILK_CLIENT_ID" ]; then
    echo "Missing SILK_CLIENT_ID env variable"
    usage
    exit 1
  fi
}

function create_asset_group() {
  # Almost copy-paste of https://docs.devprod.prod.corp.mongodb.com/mms/python/src/sbom/silkbomb/docs/SILK

  SILK_JWT_TOKEN=$(curl -s -X POST "https://silkapi.us1.app.silk.security/api/v1/authenticate" \
    -H "accept: application/json" -H "Content-Type: application/json" \
    -d '{ "client_id": "'"${SILK_CLIENT_ID}"'", "client_secret": "'"${SILK_CLIENT_SECRET}"'" }' \
    | jq -r '.token')

  asset_group_response_code=$(curl -X 'GET' \
    -s -o /dev/null -w "%{http_code}" \
    "https://silkapi.us1.app.silk.security/api/v1/raw/asset_group/$ASSET_GROUP_ID" \
    -H "accept: application/json" -H "Authorization: ${SILK_JWT_TOKEN}")

  if [[ $asset_group_response_code == 404 ]]; then
    echo "Creating new Asset"
    curl -X 'POST' \
      'https://silkapi.us1.app.silk.security/api/v1/raw/asset_group' \
      -H "accept: application/json" -H "Authorization: ${SILK_JWT_TOKEN}" \
      -H 'Content-Type: application/json' \
      -d "{
      \"active\": true,
      \"name\": \"$ASSET_GROUP_DESCRIPTION\",
      \"code_repo_url\": \"https://github.com/10gen/ops-manager-kubernetes\",
      \"branch\": \"master\",
      \"project_name\": \"10gen/ops-manager-kubernetes\",
      \"asset_id\": \"$ASSET_GROUP_ID\"
    }"
  else
    echo "Asset already exists, skipping..."
  fi
}

while getopts ':d:a:h' opt; do
  case $opt in
    d) ASSET_GROUP_DESCRIPTION=$OPTARG ;;
    a) ASSET_GROUP_ID=$OPTARG ;;
    h) usage && exit 0;;
    *) usage && exit 0;;
  esac
done
shift "$((OPTIND - 1))"

validate
create_asset_group
