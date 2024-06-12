#!/bin/bash

set -euxo pipefail

# Constants
base_s3_dir="kubernetes-operators-sboms/sboms"

# Arguments
version=$1
arch=$2

# Environment inputs
asset_group_prefix="${SILK_ASSET_GROUP}"

# Computed values
sbom_lite_json="docs/releases/v${version}/linux_${arch}.sbom.json"
sbom_json="tmp/linux-${arch}.sbom.json"
lite_name=$(jq -r < "${sbom_lite_json}" '.metadata.component.name')
name=$(jq -r < "${sbom_json}" '.metadata.component.name')

if [[ "${lite_name}" != "${name}" ]]; then
    echo "SBOM name expected to be ${lite_name} but got ${name}"
    exit 1
fi

if [[ "${name}" != mongodb/mongodb-atlas-kubernetes-operator:${version}@sha256:* ]]; then
    echo "Expected to have version tag ${version} the container name is ${name}"
    exit 1
fi

aws s3 cp "${sbom_lite_json}" "s3://${base_s3_dir}/lite/${asset_group_prefix}-linux-${arch}/${version}/linux-${arch}.json"
aws s3 cp "${sbom_json}" "s3://${base_s3_dir}/augmented/${asset_group_prefix}-linux-${arch}/${version}/linux-${arch}.json"

echo "Uploaded to S3:"
aws s3 ls "s3://${base_s3_dir}/lite/${asset_group_prefix}-linux-${arch}/${version}/"
aws s3 ls "s3://${base_s3_dir}/augmented/${asset_group_prefix}-linux-${arch}/${version}/"
