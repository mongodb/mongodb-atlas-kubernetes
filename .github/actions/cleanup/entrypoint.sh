#!/bin/bash

# For Deleting empty(!) PROJECTs which live more then 9 hours
# It deletes all if INPUT_CLEAN_ALL is true

set -eou pipefail
MAX_PROJECT_LIFETIME=9 # in hours

delete_endpoints_for_project() {
    projectID=$1
    provider=$2
    endpoints=$(mongocli atlas privateEndpoints "$provider" list --projectId "$projectID" | awk 'NR!=1{print $1}')
    # shellcheck disable=SC2068
    # multiline
    for endpoint in ${endpoints[@]}; do
        echo "deleting endpoint $endpoint in $projectID"
        mongocli atlas privateEndpoints "$provider" delete "$endpoint" --projectId "$projectID" --force
    done
}

delete_clusters() {
    projectID=$1
    clusters=$(mongocli atlas cluster list --projectId "$projectID" | awk 'NR!=1{print $2}')
    # shellcheck disable=SC2068
    # multiline
    for cluster in ${clusters[@]}; do
        echo "deleting cluster $cluster in $projectID"
        mongocli atlas cluster delete "$cluster" --projectId "$projectID" --force
    done
}

# delete only old projects (older than 9 hours)
delete_old_project() {
    if [[ -z "${count:-}" ]] || [[ ${count:-} == "null"  ]] && [[ "$existance_hours" -gt $MAX_PROJECT_LIFETIME ]]; then
        echo "deleting-$id"
        delete_endpoints_for_project "$id" "aws"
        delete_endpoints_for_project "$id" "azure"
        mongocli iam projects delete "$id" --force
    fi
}

# delete private endpoints, terminate clusters, delete empty project
delete_all() {
    delete_endpoints_for_project "$id" "aws"
    delete_endpoints_for_project "$id" "azure"
    if [[ -z ${count:-} ]] || [[ ${count:-} == "null" ]]; then
        echo "deleting-$id"
        mongocli iam projects delete "$id" --force
    else
        echo "delete only cluster (will not wait)"
        delete_clusters "$id"
    fi
}

projects=$(mongocli iam projects list -o json)
if [[ $projects == *"error"* ]]; then
    echo "Error: $projects"
    exit 1
fi

echo "${projects}"
now=$(date '+%s')

for elkey in $(echo "$projects" | jq '.results | keys | .[]'); do
    element=$(echo "$projects" | jq ".results[$elkey]")
    count=$(echo "$element" | jq -r '.clusterCount')
    id=$(echo "$element" | jq -r '.id')
    created=$(echo "$element" | jq -r '.created')
    existance_hours=$(( ("$now" - $(date --date="$created" '+%s')) / 3600 % 24 ))

    # by default delete only old projects
    if [[ "${INPUT_CLEAN_ALL:-}" == "true" ]]; then
        delete_all
    else
        delete_old_project
    fi
done
