#!/usr/bin/env bash
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

set -Eeou pipefail

platforms=("linux/arm64" "linux/amd64")
image_pull_spec=""
image_name=""
tag_name=""
repo_name=""
bucket_name=""
registry_name=""
s3_path=""
output_folder="$PWD"
docker_sbom_binary="docker-sbom"

function usage() {
  echo "Generates and uploads an SBOM to an S3 bucket.

Usage:
  generate_upload_sbom.sh [-h]
  generate_upload_sbom.sh -i <image_name>

Options:
  -h                   (optional) Shows this screen.
  -i <image_name>      (required) Image to be processed.
  -b                   (optional) S3 bucket name.
  -p                   (optional) An array of platforms, for example 'linux/arm64,linux/amd64'. The script **doesn't** fail if a particular architecture is not found.
  -o <output_folder>   (optional) Folder to output SBOM to.
"
}

function validate() {
  if ! command -v aws &> /dev/null
  then
    echo "AWS CLI not found. Please install the AWS CLI before running this script."
    exit 1
  fi
  if [ -z "$image_pull_spec" ]; then
    echo "Missing image"
    usage
    exit 1

  fi
}

function generate_sbom() {
  local image_pull_spec=$1
  local platform=$2
  local digest=$3
  local file_name=$4
  set +Ee
  "${docker_sbom_binary}" sbom --platform "$platform" -o "$file_name" --format "cyclonedx-json" "$image_pull_spec@$digest"
  docker_sbom_return_code=$?
  set -Ee
  if ((docker_sbom_return_code != 0)); then
      echo "Image $image_pull_spec with platform $platform doesn't exist. Ignoring."
      return 1
  fi
  return 0
}

while getopts ':p:i:b:o:h' opt; do
  case $opt in
    i) image_pull_spec=$OPTARG ;;
    b) bucket_name=$OPTARG ;;
    p) IFS=',' read -ra platforms <<< "$OPTARG" ;;
    o) output_folder=$OPTARG ;;
    h) usage && exit 0;;
    *) usage && exit 0;;
  esac
done
shift "$((OPTIND - 1))"

validate

# To follow the logic here, please check https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html
image_name="${image_pull_spec##*/}"
image_name="${image_name%%:*}"

registry_name="${image_pull_spec%%/*}"

repo_name=${image_pull_spec%:*}
repo_name="${repo_name%/*}"
repo_name="${repo_name##*/}"

tag_name=${image_pull_spec##*:}

s3_path="s3://${bucket_name}/sboms/$registry_name/$repo_name/$image_name/$tag_name"

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "Generating and uploading SBOM for $image_pull_spec"
echo "Image Pull Spec: $image_pull_spec"
echo "Registry: $registry_name"
echo "Repo: $repo_name"
echo "Image: $image_name"
echo "Tag: $tag_name"
echo "Platforms:" "${platforms[@]}"
echo "S3 Path: $s3_path"

for platform in "${platforms[@]}"; do
  os=${platform%/*}
  arch=${platform#*/}

  s3_path_platform_dependent="$s3_path/${os}_${arch}"
  file_name="${os}_${arch}.sbom.json"
  
  digest=$(docker manifest inspect "$image_pull_spec" | jq '.manifests[] | select(.platform.architecture == "'"$arch"'" and .platform.os == "'"$os"'")' | jq -r .digest)

  # Verify signature if SKIP_SIGNATURE_VERIFY is not set (allow verification to fail for unsigned images)
  if [ "${SKIP_SIGNATURE_VERIFY:-false}" != "true" ]; then
    echo "Verifying image signature before generating SBOM for $image_pull_spec ($platform)"
    if ! IMG=$image_pull_spec@$digest SIGNATURE_REPO=$repo_name/$image_name "${SCRIPT_DIR}"/verify.sh 2>/dev/null; then
      echo "Warning: Signature verification failed or signature not found. Continuing with SBOM generation..."
    fi
  else
    echo "Skipping signature verification (SKIP_SIGNATURE_VERIFY=true)"
  fi

  echo "Generating SBOM for $image_pull_spec ($platform) and uploading to $s3_path_platform_dependent"

  if generate_sbom "$image_pull_spec" "$platform" "$digest" "$output_folder/$file_name"; then
    echo "Done generating SBOM for $image_pull_spec ($platform)"
    if [ -z "$bucket_name" ]; then
      echo "Skipping S3 Upload (no bucket specified)"
    else 
    echo "Enabling S3 Bucket ($bucket_name) versioning"
    aws s3api put-bucket-versioning --bucket "${bucket_name}" --versioning-configuration Status=Enabled

    echo "Copying SBOM file $file_name to $s3_path_platform_dependent"
    aws s3 cp "$file_name" "$s3_path_platform_dependent"

    echo "Done uploading SBOM for $image_pull_spec ($platform)"
    fi
  fi
done
