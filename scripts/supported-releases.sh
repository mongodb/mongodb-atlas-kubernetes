#!/bin/bash

set -euo pipefail

# Prefers GNU semantics for sed, grep, sort, etc

supported_versions() {
    local one_year_ago
    one_year_ago=$(date +%s -d "1 year ago" 2> /dev/null || date -j -v-1y +%s)
    # Filter out version X.Y.0 happening before a year ago
    while IFS='' read -r line; do
      local release
      local released_date
      released_date=$(echo "${line}" | awk '{print $2}')
      if (( released_date > one_year_ago)); then
        release=$(echo "${line}" | awk '{print $1}' |awk -F/ '{print $3}' | sed 's/\.0$//')
        echo "${release}"
      fi
    done < <(git for-each-ref --sort=creatordate --format '%(refname) %(creatordate:raw)' refs/tags|grep 'v[0-9]*\.[0-9]*\.0 ')
}

supported_releases() {
    local releases
    releases=$(git tag -l |grep '^v[0-9]*\.[0-9]*\.[0-9]*$' |sort -ru)
    versions=$(supported_versions)
    echo "${releases}" |grep "${versions}" |sed 's/^v//'
}

supported_releases_json() {
    local releases
    local releases_json
    releases=$(supported_releases)
    releases_json=""
    for release in ${releases}; do
        releases_json="${releases_json}\"${release}\","
    done
    # shellcheck disable=SC2001
    releases_json=$(echo "${releases_json}" | sed 's/,$//')
    echo "[${releases_json}]"
}

supported_releases_json
