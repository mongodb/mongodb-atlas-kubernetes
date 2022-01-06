#!/bin/bash

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
