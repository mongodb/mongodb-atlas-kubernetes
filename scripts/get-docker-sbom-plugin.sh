#!/bin/bash

set -euxo pipefail

version=${DOCKER_SBOM_PLUGIN_VERSION}
os=${OS:-linux}
arch=${ARCH:-amd64}
target=tmp/sbom-cli-plugin.tgz

mkdir -p tmp
download_url_base=https://github.com/docker/sbom-cli-plugin/releases/download
url="${download_url_base}/v${version}/sbom-cli-plugin_${version}_${os}_${arch}.tar.gz"

curl -L "${url}" -o "${target}"
tar zxvf "${target}" docker-sbom
chmod +x docker-sbom
./docker-sbom
