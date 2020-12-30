#!/bin/sh

# simple changelog: get all commits beetween 2 tags and put them into file

echo "Commits: " > changelog.md
#find previos tag
start_tag=$(git for-each-ref refs/tags/ --count=2 --sort=-version:refname --format='%(refname:short)' | awk 'NR==2')
if [ -n "${start_tag}" ]; then
    git log "${start_tag}...HEAD" --pretty=format:"- %s" >> changelog.md
else
    # first tag
    git log --pretty=format:"- %s" >> changelog.md
fi

cat changelog.md
