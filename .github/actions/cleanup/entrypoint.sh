#!/bin/bash

# For Deleting empty(!) PROJECTs which live more then 1 days

set -eou pipefail

BASE_URL="https://cloud-qa.mongodb.com/api/atlas/v1.0"

get_projects() {
    curl -s -u "${INPUT_ATLAS_PUBLIC_KEY}:${INPUT_ATLAS_PRIVATE_KEY}" --digest "${BASE_URL}/groups"
}
delete_project() {
    projectID=$1
    curl -s -X DELETE --digest -u "${INPUT_ATLAS_PUBLIC_KEY}:${INPUT_ATLAS_PRIVATE_KEY}" "${BASE_URL}/groups/${projectID}"
}

delete_pe_service() {
    projectID=$1
    provider=$2
    pe_service=$3
    curl -s -X DELETE --digest -u "${INPUT_ATLAS_PUBLIC_KEY}:${INPUT_ATLAS_PRIVATE_KEY}" "${BASE_URL}/groups/${projectID}/privateEndpoint/${provider}/endpointService/${pe_service}"
}

delete_pe() {
    projectID=$1
    provider=$2
    pe_service=$3
    pe=$4
    curl -s -X DELETE --digest -u "${INPUT_ATLAS_PUBLIC_KEY}:${INPUT_ATLAS_PRIVATE_KEY}" "${BASE_URL}/groups/${projectID}/privateEndpoint/${provider}/endpointService/${pe_service}/endpoint/${pe}"
}

delete_all_pe_for_project() {
    projectID=$1
    provider=$2
    [[ "$provider" == "AWS" ]] && endpoint_field=".interfaceEndpoints" || endpoint_field=".privateEndpoints"

    list=$(get_pe_service_in_project "$projectID" "$provider")
    for elkey in $(echo "$list" | jq -c '.[]'); do
        echo $elkey
        pe_service_id=$(echo "$elkey" | jq -r '.id' )
        for pe in $(echo "$elkey" | jq -r $endpoint_field' | .[]'); do
            echo "private endpoint: $pe from $projectID will be deleted"

            pe=$(echo "$pe" | sed 's/\//%2F/g')
            delete_pe "$projectID" "$provider" "$pe_service_id" "$pe"
        done
        sleep 5 # fast operation we can wait
        echo "pe service will be deleted: $pe_service_id in $projectID"
        delete_pe_service "$projectID" "$provider" "$pe_service_id"
        # we do not wait PE service deletion - long operation
    done
}

get_pe_service_in_project() {
    projectID=$1
    provider=$2
    curl -s -u "${INPUT_ATLAS_PUBLIC_KEY}:${INPUT_ATLAS_PRIVATE_KEY}" --digest "${BASE_URL}/groups/${projectID}/privateEndpoint/${provider}/endpointService"
}

get_clusters() {
    projectID=$1
    curl -s -u "${INPUT_ATLAS_PUBLIC_KEY}:${INPUT_ATLAS_PRIVATE_KEY}" --digest "${BASE_URL}/groups/${projectID}/clusters"
}
delete_cluster() {
    projectID=$1
    name=$2
    echo "${BASE_URL}/groups/${projectID}/clusters/${name}"
    curl -s -u "${INPUT_ATLAS_PUBLIC_KEY}:${INPUT_ATLAS_PRIVATE_KEY}" --digest -X DELETE "${BASE_URL}/groups/${projectID}/clusters/${name}"
}

# delete only old projects
delete_old_project() {
    if [[ "$count" = 0 ]] && [[ "$existance_days" -gt 1 ]]; then
        echo "deleting-$id"
        delete_all_pe_for_project "$id" "AWS"
        delete_all_pe_for_project "$id" "AZURE"
        delete_project "$id"
    fi
}

# terminate cluster and delete empty project
delete_all() {
    delete_all_pe_for_project "$id" "AWS"
    delete_all_pe_for_project "$id" "AZURE"
    if [[ $count != 0 ]]; then
        #delete cluster
        clusters=$(get_clusters "$id")
        for ckey in $(echo "$clusters" | jq '.results | keys | .[]'); do
            cluster=$(echo "$clusters" | jq -r ".results[$ckey]")
            csize=$(echo "$cluster" | jq -r '.providerSettings.instanceSizeName')
            cname=$(echo "$cluster" | jq -r '.name')
            if [[ $csize != "M0" ]]; then
                echo "delete cluster: $id $cname $csize"
                delete_cluster "$id" "$cname"
                #not going to wait for deleting projects
            else
                echo "$cname $csize is M0"
            fi
        done
    else
        echo "deleting-$id"
        delete_project "$id"
    fi
}

projects=$(get_projects)
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
    existance_days=$(( ("$now" - $(date --date="$created" '+%s')) / 86400 ))
    # by default delete only old projects
    if [[ "${INPUT_CLEAN_ALL}" == "true" ]]; then
        delete_all
    else
        delete_old_project
    fi
done
