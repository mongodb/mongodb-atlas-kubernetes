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

set -euo pipefail
MESSAGE=$(cat)
WEBHOOK="${1:-${SLACK_WEBHOOK:-}}"

if [ -z "${WEBHOOK}" ]; then
  echo "Error: No webhook URL provided. Pass it as an argument or set SLACK_WEBHOOK."
  exit 1
fi

if [ -z "${MESSAGE}" ]; then
  echo "Error: No message to send"
  exit
fi

# Use jq to properly escape the JSON payload
PAYLOAD=$(echo -n "${MESSAGE}" | jq -Rs '{text: .}')
curl -X POST -H "Content-Type: application/json" -d "${PAYLOAD}" "${WEBHOOK}"
