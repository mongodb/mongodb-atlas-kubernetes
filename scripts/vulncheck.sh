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
