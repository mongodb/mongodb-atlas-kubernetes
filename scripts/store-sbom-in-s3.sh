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

###
# This script is responsible for uploading both the augmented SBOM and SBOM lites in S3
#
# AWS account:              mongodb-mms-testing
# S3 bucket:                kubernetes-operators-sboms
# Canonical path in bucket:
# s3://kubernetes-operators-sboms/sboms/{lite|augmented}]/atlas-kubernetes-operator-linux-${arch}/${version}/linux-${arch}.json
#
# Usage:
#  AWS_... KONDUKTO_BRANCH_PREFIX=... store_ ${VERSION} ${TARGET_ARCH}
# Where:
#   AWS_... means the AWS credentials for the mongodb-mms-testing account need to be present for S3 access to work
#   KONDUKTO_BRANCH_PREFIX is the environment variable with the Kondukto branch common prefix
#   VERSION is the version of the SBOM lites to store and augment
#   TARGET_ARCH is the architecture to store the SBOM for
###

# Constants
base_s3_dir="kubernetes-operators-sboms/sboms"

# Arguments
version=$1
[ -z "${version}" ] && echo "Missing version parameter #1" && exit 1
arch=$2
[ -z "${arch}" ] && echo "Missing arch parameter #2" && exit 1

# Environment inputs
kondukto_branch_prefix="${KONDUKTO_BRANCH_PREFIX:?KONDUKTO_BRANCH_PREFIX must be set}"

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

aws s3 cp "${sbom_lite_json}" "s3://${base_s3_dir}/lite/${kondukto_branch_prefix}-linux-${arch}/${version}/linux-${arch}.json"
aws s3 cp "${sbom_json}" "s3://${base_s3_dir}/augmented/${kondukto_branch_prefix}-linux-${arch}/${version}/linux-${arch}.json"

echo "Uploaded to S3:"
aws s3 ls "s3://${base_s3_dir}/lite/${kondukto_branch_prefix}-linux-${arch}/${version}/"
aws s3 ls "s3://${base_s3_dir}/augmented/${kondukto_branch_prefix}-linux-${arch}/${version}/"
