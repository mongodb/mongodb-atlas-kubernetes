#!/usr/bin/env bash
# Copyright 2026 MongoDB Inc
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

# Create a pull request, or update the open pull request for the given head and base.

set -euo pipefail

usage() {
  echo "Usage: $0 <base-branch> <head-branch> <title> <body>" >&2
}

if [[ $# -ne 4 ]]; then
  usage
  echo "Error: expected four non-empty arguments." >&2
  exit 2
fi

base="$1"
head="$2"
title="$3"
body="$4"

if [[ -z "$base" || -z "$head" || -z "$title" || -z "$body" ]]; then
  usage
  echo "Error: all four arguments must be non-empty." >&2
  exit 2
fi

if ! command -v gh >/dev/null 2>&1; then
  echo "Error: GitHub CLI (gh) is not installed." >&2
  exit 1
fi

pr_number=$(gh pr list --head "$head" --base "$base" --state open --json number --jq '.[0].number // empty')

if [[ -n "$pr_number" ]]; then
  gh pr edit "$pr_number" --title "$title" --body "$body"
else
  gh pr create --base "$base" --head "$head" --title "$title" --body "$body"
fi
