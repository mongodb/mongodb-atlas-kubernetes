#!/bin/bash

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

mkdir -p "docs/releases/v${VERSION}"
img="mongodb/mongodb-atlas-kubernetes-operator:${VERSION}"
IMG_SHAS=$(docker manifest inspect "${img}" | \
  jq -rc '.manifests[] | select(.platform.os != "unknown" and .platform.architecture != "unknown") | .digest')
for sha in ${IMG_SHAS};do
  docker pull "${img}@${sha}"
  os=$(docker inspect "${img}@${sha}" |jq -r '.[0].Os')
  arch=$(docker inspect "${img}@${sha}" |jq -r '.[0].Architecture')
  docker sbom --platform "${os}/${arch}" --format "cyclonedx-json" \
    -o "docs/releases/v${VERSION}/${os}-${arch}.sbom.json" "${img}@${sha}"
done

envsubst < docs/releases/sdlc-compliance.template.md \
  > "docs/releases/v${VERSION}/sdlc-compliance.md"

echo "SDLC checklist ready:"
ls -l "docs/releases/v${VERSION}"
