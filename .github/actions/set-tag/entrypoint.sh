#!/bin/sh
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

set -eou pipefail

git config --global --add safe.directory /github/workspace

# Get the full commit hash and shorten to 6 characters
full_commit_sha="${INPUT_COMMIT_SHA:-}"
if [ -z "$full_commit_sha" ]; then
  full_commit_sha=$(git rev-parse HEAD)
fi
commit_id=$(echo "$full_commit_sha" | cut -c1-6)

# Get the full branch name
branch_name="${INPUT_BRANCH_NAME:-}"
if [ -z "$branch_name" ]; then
  if [ -n "$GITHUB_HEAD_REF" ]; then
    branch_name="$GITHUB_HEAD_REF"
  else
    branch_name="${GITHUB_REF#refs/heads/}"
  fi
fi

# Replace / and . with -
# Then truncate to 15 characters
branch_name=$(echo "$branch_name" | sed 's/[\/\.]/-/g' | awk '{print substr($0, 1, 15)}')

# Create tag as {branch_name}-{6-digit-commit} 
tag="${branch_name}-${commit_id}"
echo "tag=${tag}" >> "$GITHUB_OUTPUT"
