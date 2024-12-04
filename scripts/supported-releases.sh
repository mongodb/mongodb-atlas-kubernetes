#!/bin/bash

set -euo pipefail

# Prefers GNU semantics for sed, grep, sort, etc

supported_releases() {
    local all_versions
    local last_three_minor_version
    all_versions=$(git tag -l | grep '^v[0-9]*\.[0-9]*\.[0-9]*$' | sort -ru)
    last_three_minor_version=$(git tag -l | grep '^v[0-9]*\.[0-9]*\.[0-9]*$' | awk -F. '{print $1 "." $2}' | sort -ru | head -3)
    echo "${all_versions}" | grep "${last_three_minor_version}" | sed 's/^v//'
}

supported_releases_json() {
    local releases_json
    releases_json=""
    for release in $(supported_releases); do
        releases_json="${releases_json}\"${release}\","
    done
    # shellcheck disable=SC2001
    releases_json=$(echo "${releases_json}" | sed 's/,$//')
    echo "[${releases_json}]"
}

supported_releases_json
