#!/bin/bash

# For Deleting empty(!) PROJECTs which live more then 1 days

set -e

BASE_URL="https://cloud-qa.mongodb.com/api/atlas/v1.0"

get_projects() {
    curl -s -u "${INPUT_ATLAS_PUBLIC_KEY}:${INPUT_ATLAS_PRIVATE_KEY}" --digest "${BASE_URL}/groups"
}
delete_project() {
    projectID=$1
    curl -s -X DELETE --digest -u "${INPUT_ATLAS_PUBLIC_KEY}:${INPUT_ATLAS_PRIVATE_KEY}" "${BASE_URL}/groups/${projectID}"
}

projects=$(get_projects)
echo "${projects}"
now=$(date '+%s')
for elkey in $(echo "$projects" | jq '.results | keys | .[]'); do
    element=$(echo "$projects" | jq ".results[$elkey]")
    count=$(echo "$element" | jq -r '.clusterCount')
    id=$(echo "$element" | jq -r '.id')
    created=$(echo "$element" | jq -r '.created')
    existance_days=$(( ("$now" - $(date --date="$created" '+%s')) / 86400 ))
    if [[ "$count" = 0 ]] && [[ "$existance_days" -gt 1 ]]; then
        echo "deleting-$id"
        delete_project "$id"
    fi
done
