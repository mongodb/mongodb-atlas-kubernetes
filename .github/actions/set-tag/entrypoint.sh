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


#set -eou pipefail

git config --global --add safe.directory /github/workspace

# Setup tag name
commit_id=$(git rev-parse --short HEAD)
branch_name=${GITHUB_HEAD_REF-}
if [ -z "${branch_name}" ]; then
    branch_name=$(echo "$GITHUB_REF" | awk -F'/' '{print $3}')
fi
branch_name=$(echo "${branch_name}" | awk '{print substr($0, 1, 15)}' | sed 's/\//-/g; s/\./-/g')
tag="${branch_name}-${commit_id}"
echo "tag=$tag" >> "$GITHUB_OUTPUT"
