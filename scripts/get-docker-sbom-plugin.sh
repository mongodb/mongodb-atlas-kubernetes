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
