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

# GitHub will sign commits if the API request is authenticated and lacks
# author and committer arguments. See:
# https://github.com/peter-evans/create-pull-request/issues/1241#issuecomment-1232477512

set -euo pipefail

# Configuration defaults
github_token=${GITHUB_TOKEN:?}
repo_owner="${REPO_OWNER:-mongodb}"
repo_name="${REPO_NAME:-mongodb-atlas-kubernetes}"
branch="${BRANCH:?}"
commit_message="${COMMIT_MESSAGE:?}"
remote=${REMOTE:-origin}

# Ensure remote branch exists
echo "Check remote branch ${branch} exists..."

function create_branch() {
  echo "Creating branch ${remote}/${branch}"
  git stash
  git push "${remote}" "${branch}"
  git stash pop
  git add .
}

git ls-remote --exit-code --heads origin "${branch}" >/dev/null || create_branch

# Fetch the latest commit SHA
LATEST_COMMIT_SHA=$(curl -s -H "Authorization: token $github_token" \
  "https://api.github.com/repos/$repo_owner/$repo_name/git/ref/heads/$branch" | jq -r '.object.sha')

LATEST_TREE_SHA=$(curl -s -H "Authorization: token $github_token" \
  "https://api.github.com/repos/$repo_owner/$repo_name/git/commits/$LATEST_COMMIT_SHA" | jq -r '.tree.sha')

echo "Creating a signed commit in GitHub."
echo "Latest commit SHA: $LATEST_COMMIT_SHA"  
echo "Latest tree SHA: $LATEST_TREE_SHA"

# Collect all modified files  
MODIFIED_FILES=$(git diff --name-only --cached)  
echo "Modified files: $MODIFIED_FILES"  

# Create blob and tree  
NEW_TREE_ARRAY="["  
for FILE_PATH in $MODIFIED_FILES; do  
  # Read file content encoded to base64  
  ENCODED_CONTENT=$(base64 -w0 < "${FILE_PATH}")

  # Create blob
  echo "{\"content\": \"$ENCODED_CONTENT\", \"encoding\": \"base64\"}" > content.blob
  BLOB_JSON=$(curl -s -X POST -H "Authorization: token $github_token" \
    -H "Accept: application/vnd.github.v3+json" \
    -d @content.blob \
    "https://api.github.com/repos/$repo_owner/$repo_name/git/blobs")
  BLOB_SHA=$(echo "${BLOB_JSON}" | jq -r '.sha')
  rm content.blob

  # Append file info to tree JSON  
  NEW_TREE_ARRAY="${NEW_TREE_ARRAY}{\"path\": \"$FILE_PATH\", \"mode\": \"100644\", \"type\": \"blob\", \"sha\": \"$BLOB_SHA\"},"  
done

# Remove trailing comma and close the array  
NEW_TREE_ARRAY="${NEW_TREE_ARRAY%,}]"

# Create new tree
echo "{\"base_tree\": \"$LATEST_TREE_SHA\", \"tree\": $NEW_TREE_ARRAY}" > new_tree_spec.json
NEW_TREE_SHA=$(curl -s -X POST -H "Authorization: token $github_token" \
  -H "Accept: application/vnd.github.v3+json" \
  -d @new_tree_spec.json \
  "https://api.github.com/repos/$repo_owner/$repo_name/git/trees" | jq -r '.sha')  
rm new_tree_spec.json

echo "New tree SHA: $NEW_TREE_SHA"  

# Create a new commit
NEW_COMMIT_SHA=$(curl -s -X POST -H "Authorization: token $github_token" \
  -H "Accept: application/vnd.github.v3+json" \
  -d "{\"message\": \"$commit_message\", \"tree\": \"$NEW_TREE_SHA\", \"parents\": [\"$LATEST_COMMIT_SHA\"]}" \
  "https://api.github.com/repos/$repo_owner/$repo_name/git/commits" | jq -r '.sha')
echo "New commit SHA: $NEW_COMMIT_SHA"  

# Update the reference of the branch to point to the new commit
curl -s -X PATCH -H "Authorization: token $github_token" \
  -H "Accept: application/vnd.github.v3+json" -d "{\"sha\": \"$NEW_COMMIT_SHA\"}" \
  "https://api.github.com/repos/$repo_owner/$repo_name/git/refs/heads/$branch"

echo "Branch ${branch} updated to new commit ${NEW_COMMIT_SHA}."  
git restore --staged .
git restore .
git clean -f
git fetch "${remote}"
git checkout -b "${branch}" "${remote}/${branch}" || git checkout "${branch}"

echo "Local branch set to ${remote}/${branch}"
