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


set -eu

release_date=${DATE:-$(date -u '+%Y-%m-%d')}
release_type=${RELEASE_TYPE:-Minor}

export DATE="${release_date}"
export VERSION="${VERSION}"
export AUTHORS="${AUTHORS}"
export RELEASE_TYPE="${release_type}"

ignored_list=""
ignored_vulns=$(grep '^# ' vuln-ignore |grep '\S' | sed 's/^# /    - /')
if [ "${ignored_vulns}" != "" ];then
  printf -v ignored_list "\n  - List of explicitly ignored vulnerabilities:\n%s" "${ignored_vulns}"
else
  printf -v ignored_list "\n  - No vulnerabilities were ignored for this release."
fi
export IGNORED_VULNERABILITIES="${ignored_list}"

mkdir -p "docs/releases/v${VERSION}/"
envsubst < docs/releases/sdlc-compliance.template.md \
  > "docs/releases/v${VERSION}/sdlc-compliance.md"

echo "SDLC checklist ready. Files at docs/releases/v${VERSION}:"
ls -l "docs/releases/v${VERSION}"
