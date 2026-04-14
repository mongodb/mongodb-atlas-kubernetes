// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	Spec            DeploymentSpec     `json:"spec"`
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
