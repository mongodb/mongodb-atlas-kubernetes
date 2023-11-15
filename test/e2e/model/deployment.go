package model

import (
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/utils"
)

type AtlasDeployment struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec            DeploymentSpec     `json:"spec,omitempty"`
}

type DeploymentSpec v1.AtlasDeploymentSpec

func (spec DeploymentSpec) GetDeploymentName() string {
	if spec.ServerlessSpec != nil {
		return spec.ServerlessSpec.Name
	}
	return spec.DeploymentSpec.Name
}

// LoadUserDeploymentConfig load configuration into object
func LoadUserDeploymentConfig(path string) AtlasDeployment {
	var config AtlasDeployment
	utils.ReadInYAMLFileAndConvert(path, &config)
	return config
}

func (ad *AtlasDeployment) DeploymentFileName(input UserInputs) string {
	// return "data/deployment-" + ac.ObjectMeta.Name + "-" + ac.Spec.Project.Name + ".yaml"
	return filepath.Dir(input.ProjectPath) + "/" + ad.ObjectMeta.Name + "-" + ad.Spec.Project.Name + ".yaml"
}

func (ad *AtlasDeployment) GetDeploymentNameResource() string {
	return "atlasdeployment.atlas.mongodb.com/" + ad.ObjectMeta.Name
}
