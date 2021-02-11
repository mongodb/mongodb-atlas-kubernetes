#!/bin/sh

set -eou pipefail

# simple changelog: get all commits beetween 2 tags and put them into file

echo "Commits: " > changelog.md
#find previos tag
start_tag=$(git describe --abbrev=0 --tags --match 'v*')
if [ -n "${start_tag}" ]; then
    git log "${start_tag}...HEAD" --pretty=format:"- %s" >> changelog.md
else
    # first tag
    git log --pretty=format:"- %s" >> changelog.md
fi

cat changelog.md
