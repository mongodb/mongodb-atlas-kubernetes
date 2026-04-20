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


set -eou pipefail

# controller-gen puts both Role and ClusterRole into a single role.yaml file.
# Kustomize doesn't provide an easy way to use only a single resource from a
# multi-document file as a base, so we split them into separate files under
# clusterwide/ and namespaced/ subdirectories.
#
# Usage: split_roles_yaml.sh [base_dir]
#   base_dir defaults to config/rbac

BASE_DIR="${1:-config/rbac}"

if [[ -f "${BASE_DIR}/role.yaml" ]]; then
	awk '/---/{f="xx0"int(++i);} {if(NF!=0)print > f};' "${BASE_DIR}/role.yaml"
	mv xx01 "${BASE_DIR}/clusterwide/role.yaml"
	mv xx02 "${BASE_DIR}/namespaced/role.yaml"
	rm "${BASE_DIR}/role.yaml"
fi
