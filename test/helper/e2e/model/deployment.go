package model

import (
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

type AtlasDeployment struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec            DeploymentSpec     `json:"spec,omitempty"`
}

type DeploymentSpec akov2.AtlasDeploymentSpec

func (spec DeploymentSpec) GetDeploymentName() string {
	if spec.FlexSpec != nil {
		return spec.FlexSpec.Name
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
	return filepath.Dir(input.ProjectPath) + "/" + ad.ObjectMeta.Name + "-" + ad.Spec.ProjectRef.Name + ".yaml"
}

func (ad *AtlasDeployment) GetDeploymentNameResource() string {
	return "atlasdeployment.atlas.mongodb.com/" + ad.ObjectMeta.Name
}
