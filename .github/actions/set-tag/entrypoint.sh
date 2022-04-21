#!/bin/sh

#set -eou pipefail

git config --global --add safe.directory /github/workspace

# Setup tag name
commit_id=$(git rev-parse --short HEAD)
branch_name=${GITHUB_HEAD_REF-}
if [ -z "${branch_name}" ]; then
    branch_name=$(echo "$GITHUB_REF" | awk -F'/' '{print $3}')
fi
branch_name=$(echo "${branch_name}" | awk '{print substr($0, 1, 15)}' | sed 's/\//-/g')
tag="${branch_name}-${commit_id}"
echo "::set-output name=tag::$tag"
