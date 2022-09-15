#!/bin/bash
# shellcheck shell=bash disable=SC1091

# For Deleting empty(!) PROJECTs which live more then (MAX_PROJECT_LIFETIME_INPUT) hours
# It deletes all if INPUT_CLEAN_ALL is true

mongocli config set skip_update_check true
set -e

# ------------------------------------------------------------------------------
# delete global API key by API request (mongocli does not support it yet)
BASE_URL="https://cloud-qa.mongodb.com/api/atlas/v1.0"

get_api_keys() {
    curl -s -u "${MCLI_PUBLIC_API_KEY}:${MCLI_PRIVATE_API_KEY}" --digest "${BASE_URL}/orgs/${MCLI_ORG_ID}/apiKeys"
}

delete_test_apikeys() {
    API_KEY_ID=$1
    curl -s -u "${MCLI_PUBLIC_API_KEY}:${MCLI_PRIVATE_API_KEY}" --digest --request DELETE "${BASE_URL}/orgs/${MCLI_ORG_ID}/apiKeys/${API_KEY_ID}"
}

get_all_serverless() {
    curl -s -u "${MCLI_PUBLIC_API_KEY}:${MCLI_PRIVATE_API_KEY}" --digest "${BASE_URL}/groups/${projectID}/serverless"
}

delete_serverless_request() {
    instance=$1
    curl -s -u "${MCLI_PUBLIC_API_KEY}:${MCLI_PRIVATE_API_KEY}" --digest --request DELETE "${BASE_URL}/groups/${projectID}/serverless/${instance}"
}

get_serverless_instance_request() {
    instance=$1
    curl -s -u "${MCLI_PUBLIC_API_KEY}:${MCLI_PRIVATE_API_KEY}" --digest "${BASE_URL}/groups/${projectID}/serverless/${instance}"
}

# ------------------------------------------------------------------------------

command_result_waiter() {
    command=$1
    match=$2
    interval=$3
    maxCount=$4

    try=1
    echo "wait result...."
    until [[ "$( eval " $command" | grep "$match" && echo "alive" || echo "deleted")" == "deleted" ]] || [[ "$try" -gt $maxCount ]]; do
        echo "wait...try $try..."
        ((try++))
        sleep "$interval"
    done
}

delete_endpoints_for_project() {
    projectID=$1
    provider=$2

    endpoints=$(mongocli atlas privateEndpoints "$provider" list --projectId "$projectID" -o json | jq -c . )
    [[ "$provider" == "aws" ]] && field=".interfaceEndpoints"
    [[ "$provider" == "azure" ]] && field=".privateEndpoints"
    [[ "$provider" == "gcp" ]] && field=".endpointGroupNames"
    [[ "$field" == "" ]] && echo "Please check provider" && exit 1

    # shellcheck disable=SC2068
    # multiline
    for endpoint in $(echo "$endpoints" | jq -cr '.[]'); do
        # echo $endpoint
        service_id=$(echo "$endpoint" | jq -r '.id' )
        points=$(echo "$endpoint" | jq -r "$field")

        if [[ $points != "null" ]]; then
            for pe in $(echo "$points" | jq -r '.[]'); do
                echo "----Delete private endpoint: $pe from $projectID will be deleted"
                mongocli atlas privateEndpoints "$provider" interfaces delete "$pe" --endpointServiceId "$service_id" --projectId "$projectID" --force

                command="mongocli atlas privateEndpoints $provider list --projectId $projectID -o json"
                command_result_waiter "$command" "$pe" 2 150
            done
        fi

        echo "--Private endpoint service will be deleted: $service_id in $projectID"
        # we do not wait PE service deletion - long operation
        mongocli atlas privateEndpoints "$provider" delete "$service_id" --projectId "$projectID" --force
        command="mongocli atlas privateEndpoints $provider list --projectId $projectID -o json"
        command_result_waiter "$command" "$service_id" 10 60
    done
}

delete_clusters() {
    projectID=$1
    echo "====== Cleaning Clusters in Project: $projectID"
    clusters=$(mongocli atlas cluster list --projectId "$projectID" | awk 'NR!=1{print $2}')
    # shellcheck disable=SC2068
    # multiline
    for cluster in ${clusters[@]}; do
        echo "====== Cleaning Clusters: $cluster"
        state=$(mongocli atlas cluster describe "$cluster" --projectId "$projectID" -o json | jq -r '.stateName')
        echo "Current cluster: $cluster. State: $state"
        if [[ -n $state ]] && [[ $state != "DELETING" ]]; then
            echo "deleting cluster $cluster in $projectID"
            mongocli atlas cluster delete "$cluster" --projectId "$projectID" --force
            command="mongocli atlas cluster list --projectId $projectID"
            echo "$command"
            command_result_waiter "$command" "$cluster" 10 90
        fi
    done
}

delete_serverless() {
    projectID=$1
    echo "====== Cleaning Serverless Instances in Project: $projectID"
    all=$(get_all_serverless)
    echo "all instances: $all"
    serverless=$(echo "$all" | jq -r '.results | .[] | .name')
    echo "list name: ${serverless[*]}"
    total_serverless=$(echo "$all" | jq -r '.totalCount // 0')
    echo "total serverless instances: $total_serverless"
    if [[ $total_serverless -gt 0 ]]; then
        # shellcheck disable=SC2068
        # multiline
        for instance in ${serverless[@]}; do
            echo "== Cleaning Serverless instance: $instance"
            instance_atlas=$(get_serverless_instance_request "$instance")
            echo "inside: $instance_atlas"

            state=$(get_serverless_instance_request "$instance" | jq -r '.stateName')
            echo "Current instance: $instance. State: $state"
            if [[ -n $state ]] && [[ $state != "DELETING" ]]; then
                echo "Deleting instance $instance in $projectID"
                delete_serverless_request "$instance"
                command="get_serverless_instance_request $instance"
                command_result_waiter "$command" "providerSettings" 10 20 # if no providerSetting - no serverless instance
            fi
        done
    fi
}

delete_networkpeerings_for_project() {
  projectID=$1

  connections=$(mongocli atlas networking peering list --projectId "$projectID" -o json | jq -c .)

  for connection in $(echo "$connections" | jq -cr '.[]'); do
    id=$(echo "$connection" | jq -r '.id')
    echo "Removing connection $id"
    mongocli atlas networking peering delete "$id" --force --projectId "$projectID"
  done
}

delete_project() {
    peDeleted=$(mongocli atlas privateEndpoints aws list --projectId "$projectID" | awk 'NR!=1{print $1}')$(mongocli atlas privateEndpoints azure list --projectId "$projectID" | awk 'NR!=1{print $1}')
    [[ $peDeleted == "" ]] && mongocli iam projects delete "$id" --force
}

# delete only old projects (older than MAX_PROJECT_LIFETIME_INPUT hours)
delete_old_project() {
    output=$(
        exist=is_project_exist
        if [[ -n ${exist:-} ]] || [[ -z "${count:-}" ]] || [[ ${count:-} == "null"  ]] && [[ "$existance_hours" -gt $MAX_PROJECT_LIFETIME_INPUT ]]; then
            echo "======================= Cleaning Project: $id"
            delete_endpoints_for_project "$id" "aws"
            delete_endpoints_for_project "$id" "azure"
            delete_endpoints_for_project "$id" "gcp"
            delete_networkpeerings_for_project "$id"
            delete_project
        fi
    )
    echo "${output[@]}"
}

# Check if it is still exits. Could be process from another source (test, manual work, CI/CD pipelines)
is_project_exist() {
    mongocli iam projects list | awk '/'"$id"'/{print "true"}'
}

# delete private endpoints, terminate clusters, delete empty project
delete_all() {
    output=$(
        exist=is_project_exist
        if [[ -n ${exist:-} ]]; then
            echo "======================= Cleaning Project: $id"
            delete_endpoints_for_project "$id" "aws"
            delete_endpoints_for_project "$id" "azure"
            delete_endpoints_for_project "$id" "gcp"
            delete_serverless "$id"
            delete_clusters "$id"
            delete_networkpeerings_for_project "$id"
            delete_project
        fi
    )
    echo "${output[@]}"
}

echo "The process could take a while. Please, wait..."

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
        delete_all &
    else
        delete_old_project &
    fi
done
wait
echo "Finish project deletion"

# ------------------------------------------------------------------------------
if [[ "${INPUT_CLEAN_ALL:-}" == "true" ]]; then
    echo "Delete all global API keys with a particular description"
    echo "Please, remember running tests will fail (run CLEAN_ALL = false, if need soft deletion)"
    test_description="created from the AO test"
    all_keys=$(get_api_keys)
    for key in $(echo "$all_keys" | jq 'select(.results | length > 0) | .results | keys | .[]'); do
        element=$(echo "$all_keys" | jq ".results[$key]")
        desc=$(echo "$element" | jq -r '.desc')
        if [[ "${desc}" == "${test_description}" ]]; then
            key_id=$(echo "$element" | jq -r '.id')
            echo "Key to delete: $key_id"
            delete_test_apikeys "$key_id"
        fi
    done
fi
# ------------------------------------------------------------------------------

echo "Job Done"
