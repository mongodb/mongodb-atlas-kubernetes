#!/bin/bash

 set -eou pipefail

 git config --global --add safe.directory /github/workspace

 commit_single_file() {
   # Commit to the branch
   file="$1"
   sha=$(git rev-parse "$DESTINATION_BRANCH:$file") || true
   message="Pushing $file using GitHub API"

   echo "$DESTINATION_BRANCH:$file:$sha"
   if [ "$sha" = "$DESTINATION_BRANCH:$file" ]; then
       echo "File does not exist"
       echo gh api --method PUT "/repos/:owner/:repo/contents/$file" \
           --field message="$message" \
           --field content=@<( base64 -i "$file" ) \
           --field branch="$DESTINATION_BRANCH"
       echo "$file pushed to $DESTINATION_BRANCH"
   else
       echo "File exists"
       echo gh api --method PUT "/repos/:owner/:repo/contents/$file" \
           --field message="$message" \
           --field content=@<( base64 -i "$file" ) \
           --field branch="$DESTINATION_BRANCH" \
           --field sha="$sha"
       echo "$file pushed to $DESTINATION_BRANCH"
   fi
 }

 # simple 'for loop' does not work correctly, see https://github.com/koalaman/shellcheck/wiki/SC2044#correct-code
 while IFS= read -r -d '' file
 do
   commit_single_file "$file"
 done <   <(find "${PATH_TO_COMMIT}" -type f -print0)