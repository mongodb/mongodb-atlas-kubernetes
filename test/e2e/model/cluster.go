package model

import (
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

type AC struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec            ClusterSpec        `json:"spec,omitempty"`
}

type ClusterSpec v1.AtlasDeploymentSpec

func (c ClusterSpec) GetClusterName() string {
	if c.AdvancedDeploymentSpec != nil {
		return c.AdvancedDeploymentSpec.Name
	}
	if c.ServerlessSpec != nil {
		return c.ServerlessSpec.Name
	}
	return c.DeploymentSpec.Name
}

// LoadUserClusterConfig load configuration into object
func LoadUserClusterConfig(path string) AC {
	var config AC
	utils.ReadInYAMLFileAndConvert(path, &config)
	return config
}

func (ac *AC) ClusterFileName(input UserInputs) string {
	// return "data/cluster-" + ac.ObjectMeta.Name + "-" + ac.Spec.Project.Name + ".yaml"
	return filepath.Dir(input.ProjectPath) + "/" + ac.ObjectMeta.Name + "-" + ac.Spec.Project.Name + ".yaml"
}

func (ac *AC) GetClusterNameResource() string {
	return "atlascluster.atlas.mongodb.com/" + ac.ObjectMeta.Name
}
