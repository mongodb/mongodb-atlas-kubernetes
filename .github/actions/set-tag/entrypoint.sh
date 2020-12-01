#!/bin/sh
# Setup tag name
commit_id=$(git rev-parse --short HEAD)
branch_name=$GITHUB_HEAD_REF
if [ -z "${branch_name}" ]; then
    branch_name=$(echo "$GITHUB_REF" | awk -F'/' '{print $3}')
fi
tag="${branch_name}-${commit_id}"
echo "::set-output name=tag::$tag"