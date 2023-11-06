#!/bin/bash
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

