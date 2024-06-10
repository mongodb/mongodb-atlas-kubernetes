# Image SBOMs

Starting from version 2.2.0 and onward, all Atlas Kubernetes Operator images are attached SBOMs files per image platform released. SBOM stands for Software Bill Of Materials, a recursive list of all dependencies within of a software binary or image image that is useful to evaluate potential security vulnerabilities that might be affected that particular version.

These SBOMs attached after release, as they need the images to be published for the SBOMs to be computed.

This document describes how the project computes those SBOMs as well as how end users can compute them on their own.

## Scripts computing the SBOMs for the CI

The main script to check is [scripts/generate_upload_sbom.sh](../../scripts/generate_upload_sbom.sh):

```shell
$ ./scripts/generate_upload_sbom.sh -h
Generates and uploads an SBOM to an S3 bucket.

Usage:
  generate_upload_sbom.sh [-h]
  generate_upload_sbom.sh -i <image_name>

Options:
  -h                   (optional) Shows this screen.
  -i <image_name>      (required) Image to be processed.
  -b                   (optional) S3 bucket name.
  -p                   (optional) An array of platforms, for example 'linux/arm64,linux/amd64'. The script **doesn't** fail if a particular architecture is not found.
  -o <output_folder>   (optional) Folder to output SBOM to.
```

As you can see one what you use it will be:

```shell
$ ./scripts/generate_upload_sbom.sh -i mongodb/mongodb-atlas-kubernetes-operator:2.3.0
```

When given no platforms it will default to `linux/amd64` & `linux/arm64` and try to download them and produce the SBOMS files.

## DIY SBOMs

To compute the SBOMs manually the only complication, other than having `docker` with the `sbom` plug-in installed, is that to get SBOMs from multi-architecture images require the full SHA nomenclature to successfully produce the SBOM regardless of the host architecture the `docker sbom` where command is run. The gist of it getting the SHA of the desired platform and then getting the SBOM for that particular image SHA:

```shell
export digest=$(docker manifest inspect "${img}" |jq -r '.manifests[] | select(.platform.os == "'"${os}"'" and .platform.architecture=="'"${arch}"'") | .digest')
docker sbom --platform "${os}/${arch}" --format "cyclonedx-json" "${img}@${digest}"
```

For example:
```shell
$ export os=linux
$ export arch=amd64
$ export img=mongodb/mongodb-atlas-kubernetes-operator:2.3.0
$ export digest=$(docker manifest inspect "${img}" |jq -r '.manifests[] | select(.platform.os == "'"${os}"'" and .platform.architecture=="'"${arch}"'") | .digest')
$ docker sbom --platform "${os}/${arch}" --format "cyclonedx-json" "${img}@${digest}"
{
  "bomFormat": "CycloneDX",
  "specVersion": "1.4",
...
```
