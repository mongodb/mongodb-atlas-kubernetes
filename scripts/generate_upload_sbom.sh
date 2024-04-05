#!/usr/bin/env bash
set -Eeou pipefail

image_pull_spec=""
image_name=""
tag_name=""
repo_name=""
bucket_name=""
registry_name=""
s3_path=""

function usage() {
  echo "Generates and uploads an SBOM to an S3 bucket.
Usage:
  generate_upload_sbom.sh [-h]
  generate_upload_sbom.sh -i <image_name> -b <bucket_name>
Options:
  -h                   (optional) Shows this screen.
  -i <image_name>      (required) Image to be processed.
  -b                   (required) S3 bucket name.
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
  if [ -z "$bucket_name" ]; then
      echo "Missing bucket name"
      usage
      exit 1
  fi
}

while getopts ':i:b:h' opt; do
  case $opt in
    i) image_pull_spec=$OPTARG ;;
    b) bucket_name=$OPTARG ;;
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

file_name="${image_name}_${tag_name}.json"
s3_path="s3://${bucket_name}/sboms/$registry_name/$repo_name/$image_name/$tag_name"

echo "Generating and uploading SBOM for $image_pull_spec"
echo "Image Pull Spec: $image_pull_spec"
echo "Registry: $registry_name"
echo "Repo: $repo_name"
echo "Image: $image_name"
echo "Tag: $tag_name"
echo "S3 Path: $s3_path"
echo "File name: $file_name"

echo "Generating SBOM for $image_pull_spec"
docker sbom -o "$file_name" --format "cyclonedx-json" "$image_pull_spec"

echo "Enabling S3 Bucket ($bucket_name) versioning"
aws s3api put-bucket-versioning --bucket "${bucket_name}" --versioning-configuration Status=Enabled

echo "Copying SBOM to $s3_path"
aws s3 cp "$file_name" "$s3_path"

echo "Done generating and uploading SBOM for $image_pull_spec"

