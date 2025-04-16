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

# This is the script that allows to avoid the restrictions from the controller-gen tool that puts both Role and ClusterRole
# to the same role.yaml file (and kustomize doesn't provide an easy way to use only a single resource from file as a base)
# So we simply split the 'config/rbac/roles.yaml' file into two new files
if [[ -f config/rbac/role.yaml ]]; then
	awk '/---/{f="xx0"int(++i);} {if(NF!=0)print > f};' config/rbac/role.yaml
	# csplit config/rbac/role.yaml '/---/' '{*}' &> /dev/null - infinite repetition '{*}' is not working on BSD/OSx
	mv xx01 config/rbac/clusterwide/role.yaml
	mv xx02 config/rbac/namespaced/role.yaml
	rm config/rbac/role.yaml
fi
