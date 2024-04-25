#!/bin/bash

set -eu

release_date=${DATE:-$(date -u '+%Y-%m-%d')}
release_type=${RELEASE_TYPE:-Minor}
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

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

mkdir -p "docs/releases/v${VERSION}"
img="mongodb/mongodb-atlas-kubernetes-operator:${VERSION}"
${SCRIPT_DIR}/generate_upload_sbom.sh -i $img

envsubst < docs/releases/sdlc-compliance.template.md \
  > "docs/releases/v${VERSION}/sdlc-compliance.md"

echo "SDLC checklist ready:"
ls -l "docs/releases/v${VERSION}"
