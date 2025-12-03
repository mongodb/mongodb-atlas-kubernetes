#!/bin/bash
# Copyright 2025 MongoDB Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Check if an argument is provided
if [ -z "$1" ]; then
    echo "Error: No target specified."
    echo "Usage: $0 [all | <target_string>]"
    echo "  all             : Resets community, openshift, and certified operators."
    echo "  <target_string> : Resets specific operators if the string contains:"
    echo "                    'community', 'openshift', or 'certified'."
    echo "Example: $0 'community certified' (resets both)"
    exit 1
fi

TARGET="$1"

function reset_community() {
    echo "Remove prev version branch locally and remotely"
    pushd "${RH_COMMUNITY_OPERATORHUB_REPO_PATH}"
    git checkout main
    git branch -D "mongodb-atlas-operator-community-${VERSION}"
    git push origin ":mongodb-atlas-operator-community-${VERSION}"
    popd
}

function reset_openshift() {
    echo "Remove prev version branch locally and remotely"
    pushd "${RH_COMMUNITY_OPENSHIFT_REPO_PATH}"
    git checkout main
    git branch -D "mongodb-atlas-operator-community-${VERSION}"
    git push origin ":mongodb-atlas-operator-community-${VERSION}"
    popd
}

function reset_certified() {
    echo "Remove prev version branch locally and remotely"
    pushd "${RH_CERTIFIED_OPENSHIFT_REPO_PATH}"
    git checkout main
    git branch -D "mongodb-atlas-kubernetes-operator-${VERSION}"
    git push origin ":mongodb-atlas-kubernetes-operator-${VERSION}"
    popd
}

TARGET_LOWER=$(echo "$TARGET" | tr '[:upper:]' '[:lower:]')

RUN_COMMUNITY=false
RUN_OPENSHIFT=false
RUN_CERTIFIED=false

if [[ "$TARGET_LOWER" == "all" ]]; then
    RUN_COMMUNITY=true
    RUN_OPENSHIFT=true
    RUN_CERTIFIED=true
else
    if [[ "$TARGET_LOWER" == *"community"* ]]; then
        RUN_COMMUNITY=true
    fi
    if [[ "$TARGET_LOWER" == *"openshift"* ]]; then
        RUN_OPENSHIFT=true
    fi
    if [[ "$TARGET_LOWER" == *"certified"* ]]; then
        RUN_CERTIFIED=true
    fi

    if [ "$RUN_COMMUNITY" = false ] && [ "$RUN_OPENSHIFT" = false ] && [ "$RUN_CERTIFIED" = false ]; then
        echo "Error: Invalid argument '$TARGET'."
        echo "Argument must be 'all' or contain 'community', 'openshift', or 'certified'."
        exit 1
    fi
fi

# Execute the functions based on flags
if [ "$RUN_COMMUNITY" = true ]; then
    reset_community
fi

if [ "$RUN_OPENSHIFT" = true ]; then
    reset_openshift
fi

if [ "$RUN_CERTIFIED" = true ]; then
    reset_certified
fi

