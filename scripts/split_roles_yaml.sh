#!/bin/bash

set -eou pipefail

subdir=${1:-}

config_dir="config"
if [ "${subdir}" != "" ]; then
  config_dir+="/${subdir}"
fi

# This is the script that allows to avoid the restrictions from the controller-gen tool that puts both Role and ClusterRole
# to the same role.yaml file (and kustomize doesn't provide an easy way to use only a single resource from file as a base)
# So we simply split the '"${config_dir}rbac/roles.yaml' file into two new files
if [[ -f "${config_dir}/rbac/role.yaml" ]]; then
	awk '/---/{f="xx0"int(++i);} {if(NF!=0)print > f};' "${config_dir}/rbac/role.yaml"
	# csplit "${config_dir}/rbac/role.yaml" '/---/' '{*}' &> /dev/null - infinite repetition '{*}' is not working on BSD/OSx
	mv xx01 "${config_dir}/rbac/clusterwide/role.yaml"
	mv xx02 "${config_dir}/rbac/namespaced/role.yaml"
	rm "${config_dir}/rbac/role.yaml"
fi
