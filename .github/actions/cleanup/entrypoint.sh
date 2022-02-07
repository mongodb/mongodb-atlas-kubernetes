#!/bin/bash

# For Deleting empty(!) PROJECTs which live more then (9) hours
# It deletes all if INPUT_CLEAN_ALL is true

MAX_PROJECT_LIFETIME=9 # in hours
mongocli config set skip_update_check true

delete_endpoints_for_project() {
    projectID=$1
    provider=$2

    endpoints=$(mongocli atlas privateEndpoints "$provider" list --projectId "$projectID" -o json | jq -c . )
    [[ "$provider" == "aws" ]] && field=".interfaceEndpoints" || field=".privateEndpoints"

    # shellcheck disable=SC2068
    # multiline
    for endpoint in $(echo "$endpoints" | jq -cr '.[]'); do
        # echo $endpoint
        service_id=$(echo "$endpoint" | jq -r '.id' )
        points=$(echo "$endpoint" | jq -r "$field")

        if [[ $points != "null" ]]; then
            for pe in $(echo "$points" | jq -r '.[]'); do
                echo "----Delete private endpoint: $pe from $projectID will be deleted"
                pe=$(echo "$pe" | sed 's/\//%2F/g')
                mongocli atlas privateEndpoints "$provider" interfaces delete "$pe" --endpointServiceId "$service_id" --projectId "$projectID" --force
                mongocli atlas privateEndpoints "$provider" interfaces describe "$pe" --endpointServiceId "$service_id" --projectId "$projectID" | awk 'NR!=1{print $1}'
                until [[ "$(mongocli atlas privateEndpoints "$provider" interfaces describe "$pe" --endpointServiceId "$service_id" --projectId "$projectID" | awk 'NR!=1{print $1}')" == "" ]]; do
                    echo "wait..."
                    sleep 1
                done
            done
        fi

        echo "--Private endpoint service will be deleted: $service_id in $projectID"
        # we do not wait PE service deletion - long operation
        mongocli atlas privateEndpoints "$provider" delete "$service_id" --projectId "$projectID" --force
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
        echo "====== Cleaning Project: $id"
        delete_endpoints_for_project "$id" "aws"
        delete_endpoints_for_project "$id" "azure"
        isReady=$(mongocli atlas privateEndpoints aws list --projectId "$projectID" | awk 'NR!=1{print $1}')$(mongocli atlas privateEndpoints azure list --projectId "$projectID" | awk 'NR!=1{print $1}')
        [[ $isReady == "" ]] && mongocli iam projects delete "$id" --force
    fi
}

# delete private endpoints, terminate clusters, delete empty project
delete_all() {
    echo "====== Cleaning Project: $id"
    delete_endpoints_for_project "$id" "aws"
    delete_endpoints_for_project "$id" "azure"
    if [[ -z ${count:-} ]] || [[ ${count:-} == "null" ]]; then
        isReady=$(mongocli atlas privateEndpoints aws list --projectId "$projectID" | awk 'NR!=1{print $1}')$(mongocli atlas privateEndpoints azure list --projectId "$projectID" | awk 'NR!=1{print $1}')
        [[ $isReady == "" ]] && mongocli iam projects delete "$id" --force
    else
        echo "delete only cluster (will not wait)"
        delete_clusters "$id"
    fi
}

projects=$(mongocli iam projects list -o json | jq -c .)
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
