#!/bin/bash

ignore_file=${1:-./vuln-ignore}

set -e
govulncheck -version
set +e

ignore_lines=$(grep -v '^#' "${ignore_file}")
check_cmd='govulncheck ./... |grep "Vulnerability #"'
while IFS= read -r line; do
  if [ "${line}" != "" ]; then
    check_cmd+="|grep -v \"${line}\""
  fi
done <<< "${ignore_lines}"

echo "${check_cmd}"
vulns=$(eval "${check_cmd}")
count=$(echo "${vulns}" |grep -c "\S")
echo "govulncheck found $((count)) non ignored vulnerabilities"
if (( count != 0 )); then
  echo "${vulns}"
  echo "FAIL"
  exit 1
fi
echo "PASS"
