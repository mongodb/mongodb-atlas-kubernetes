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

#!/bin/bash
# retry.sh: Executes a command with retries and exponential backoff.
# Usage: ./retry.sh <command> [args...]

set +o pipefail

max_retries=${MAX_RETRIES:-3}
backoff=${BACKOFF:-1}

retries=0
exit_status=0

while true; do
    "${@}" 2>&1 | cat
    
    exit_status="${PIPESTATUS[0]}"

    if [ "$exit_status" -eq 0 ]; then
        exit 0
    fi
    
    if (( retries == max_retries )); then
        echo "❌ Command failed permanently with code $exit_status after $((retries)) retries." >&2
        exit "$exit_status"
    fi

    wait_time=$(( (retries++) * backoff ))
    echo "⚠️ Command failed (exit code $exit_status). Retrying in ${wait_time} seconds..." >&2
    sleep "${wait_time}"
done
