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

# List used actions recursively, iunclkuding transitively used actions

set -eao pipefail

actions=()
while IFS='' read -r line; do actions+=("${line}"); done < <(grep -r "uses: " .github |awk -F: '{print $3}' |sort -u |grep -v '.github/' | awk -F@ '{print $1}' |sort -u)

for action in "${actions[@]}"; do
        while IFS='' read -r line; do actions+=("${line}"); done < <(curl -s "https://raw.githubusercontent.com/${action}/main/action.yml" | grep 'uses: ' |awk -F: '{print $3}' | awk -F@ '{print $1}' |sort -u)
done

for action in "${actions[@]}"; do
	echo "${action}"
done

