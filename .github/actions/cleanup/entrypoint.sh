#!/bin/bash

# For Deleting empty(!) PROJECTs which live more then (3) hours
# It deletes all if INPUT_CLEAN_ALL is true

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
                echo "mongocli atlas privateEndpoints $provider interfaces delete $pe --endpointServiceId $service_id --projectId $projectID --force"
                mongocli atlas privateEndpoints "$provider" interfaces delete "$pe" --endpointServiceId "$service_id" --projectId "$projectID" --force
                try=1
                until [[ "$(mongocli atlas privateEndpoints "$provider" list --projectId "$projectID" -o json | grep "$pe" && echo "alive" || echo "deleted")" == "deleted" ]] || [[ "$try" -gt 150 ]]; do
                    echo "wait...$try..."
                    ((try++))
                    sleep 2
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
        state=$(mongocli atlas cluster describe "$cluster" --projectId "$projectID" -o json | jq -r '.stateName')
        echo "Current cluster: $cluster. State: $state"
        if [[ -n $state ]] || [[ $state != "DELETING" ]]; then
            echo "deleting cluster $cluster in $projectID"
            mongocli atlas cluster delete "$cluster" --projectId "$projectID" --force
        fi
    done
}

delete_project() {
    peDeleted=$(mongocli atlas privateEndpoints aws list --projectId "$projectID" | awk 'NR!=1{print $1}')$(mongocli atlas privateEndpoints azure list --projectId "$projectID" | awk 'NR!=1{print $1}')
    [[ $peDeleted == "" ]] && mongocli iam projects delete "$id" --force
}

# delete only old projects (older than MAX_PROJECT_LIFETIME_INPUT hours)
delete_old_project() {
    exist=is_project_exist
    if [[ -n ${exist:-} ]] || [[ -z "${count:-}" ]] || [[ ${count:-} == "null"  ]] && [[ "$existance_hours" -gt $MAX_PROJECT_LIFETIME_INPUT ]]; then
        echo "====== Cleaning Project: $id"
        delete_endpoints_for_project "$id" "aws"
        delete_endpoints_for_project "$id" "azure"
        delete_project
    fi
}

# Check if it is still exits. Could be process from another source (test, manual work, CI/CD pipelines)
is_project_exist() {
    mongocli iam projects list | awk '/'"$id"'/{print "true"}'
}

# delete private endpoints, terminate clusters, delete empty project
delete_all() {
    exist=is_project_exist
    if [[ -n ${exist:-} ]]; then
        echo "====== Cleaning Project: $id"
        delete_endpoints_for_project "$id" "aws"
        delete_endpoints_for_project "$id" "azure"
        if [[ -z ${count:-} ]] || [[ ${count:-} == "null" ]]; then
            delete_project
        else
            echo "delete only cluster (will not wait)"
            delete_clusters "$id"
        fi
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
