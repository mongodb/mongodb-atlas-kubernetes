#!/bin/bash  

set -euxo pipefail
  
# Configuration defaults
github_token=${GITHUB_TOKEN:?}
repo_owner="${REPO_OWNER:-mongodb}"
repo_name="${REPO_NAME:-mongodb-atlas-kubernetes}"
branch="${BRANCH:?}"
commit_message="${COMMIT_MESSAGE:?}"

# Fetch the latest commit SHA  
LATEST_COMMIT_SHA=$(curl -s -H "Authorization: token $github_token" \
  "https://api.github.com/repos/$repo_owner/$repo_name/git/ref/heads/$branch" | jq -r '.object.sha')  
  
LATEST_TREE_SHA=$(curl -s -H "Authorization: token $github_token" \
  "https://api.github.com/repos/$repo_owner/$repo_name/git/commits/$LATEST_COMMIT_SHA" | jq -r '.tree.sha')  
  
echo "Latest commit SHA: $LATEST_COMMIT_SHA"  
echo "Latest tree SHA: $LATEST_TREE_SHA"  
  
# Collect all modified files  
MODIFIED_FILES=$(git diff --name-only)  
echo "Modified files: $MODIFIED_FILES"  
  
# Create blob and tree  
NEW_TREE_ARRAY="["  
for FILE_PATH in $MODIFIED_FILES; do  
  # Read file content encoded to base64  
  ENCODED_CONTENT=$(base64 < "${FILE_PATH}")  
  
  # Create blob  
  BLOB_SHA=$(curl -s -X POST -H "Authorization: token $github_token" \
    -H "Accept: application/vnd.github.v3+json" \
    -d "{\"content\": \"$ENCODED_CONTENT\", \"encoding\": \"base64\"}" \
    "https://api.github.com/repos/$repo_owner/$repo_name/git/blobs" | jq -r '.sha')  
  
  # Append file info to tree JSON  
  NEW_TREE_ARRAY="${NEW_TREE_ARRAY}{\"path\": \"$FILE_PATH\", \"mode\": \"100644\", \"type\": \"blob\", \"sha\": \"$BLOB_SHA\"},"  
done  
  
# Remove trailing comma and close the array  
NEW_TREE_ARRAY="${NEW_TREE_ARRAY%,}]"  
  
# Create new tree  
NEW_TREE_SHA=$(curl -s -X POST -H "Authorization: token $github_token" \
  -H "Accept: application/vnd.github.v3+json" \
  -d "{\"base_tree\": \"$LATEST_TREE_SHA\", \"tree\": $NEW_TREE_ARRAY}" \
  "https://api.github.com/repos/$repo_owner/$repo_name/git/trees" | jq -r '.sha')  
  
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
