#!/bin/bash

set -euxo pipefail

version=${DOCKER_SBOM_PLUGIN_VERSION:-latest}
os=${OS:-linux}
arch=${ARCH:-amd64}
target=$TMPDIR/sbom-cli-plugin.tgz
docker_path=$(which docker)
docker_dir=$(dirname "${docker_path}")

download_url_base=https://github.com/docker/sbom-cli-plugin/releases/download
url="${download_url_base}/v${version}/sbom-cli-plugin_${version}_${os}_${arch}.tar.gz"

curl -L "${url}" -o "${target}"
pushd "${TMPDIR}"
tar zxvf "${target}" docker-sbom
chmod +x docker-sbom
popd
sudo cp "${TMPDIR}/docker-sbom" "${docker_dir}"
