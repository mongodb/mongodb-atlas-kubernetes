package model

import (
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

type AtlasDeployment struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec            DeploymentSpec     `json:"spec,omitempty"`
}

type DeploymentSpec v1.AtlasDeploymentSpec

func (c DeploymentSpec) GetDeploymentName() string {
	if c.AdvancedDeploymentSpec != nil {
		return c.AdvancedDeploymentSpec.Name
	}
	if c.ServerlessSpec != nil {
		return c.ServerlessSpec.Name
	}
	return c.DeploymentSpec.Name
}

// LoadUserDeploymentConfig load configuration into object
func LoadUserDeploymentConfig(path string) AtlasDeployment {
	var config AtlasDeployment
	utils.ReadInYAMLFileAndConvert(path, &config)
	return config
}

func (ac *AtlasDeployment) DeploymentFileName(input UserInputs) string {
	// return "data/deployment-" + ac.ObjectMeta.Name + "-" + ac.Spec.Project.Name + ".yaml"
	return filepath.Dir(input.ProjectPath) + "/" + ac.ObjectMeta.Name + "-" + ac.Spec.Project.Name + ".yaml"
}

func (ac *AtlasDeployment) GetDeploymentNameResource() string {
	return "atlasdeployment.atlas.mongodb.com/" + ac.ObjectMeta.Name
}
