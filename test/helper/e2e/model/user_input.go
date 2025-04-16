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
	"fmt"
	"path/filepath"

	. "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

type UserInputs struct {
	TestID             string
	AtlasKeyAccessType AtlasKeyType
	ProjectID          string
	KeyName            string
	Namespace          string
	ProjectPath        string
	Deployments        []AtlasDeployment
	Users              []DBUser
	Project            *AProject
}

// NewUserInputs prepare users inputs
func NewUserInputs(keyTestPrefix string, project AProject, users []DBUser, r *AtlasKeyType) UserInputs {
	testID := utils.GenID()
	projectName := fmt.Sprintf("%s-%s", keyTestPrefix, testID)
	input := UserInputs{
		TestID:             testID,
		AtlasKeyAccessType: *r,
		ProjectID:          "",
		KeyName:            keyTestPrefix,
		Namespace:          "ns-" + projectName,
		ProjectPath:        filepath.Join(DataGenFolder, projectName, "resources", projectName+".yaml"),
	}

	input.Project = &project
	input.Project = NewProject("k-" + projectName).ProjectName(projectName)
	if len(input.Project.Spec.ProjectIPAccessList) == 0 {
		input.Project = input.Project.WithIpAccess("0.0.0.0/0", "everyone")
	}

	if !r.GlobalLevelKey {
		input.Project = input.Project.WithSecretRef(keyTestPrefix)
	}

	for _, user := range users {
		input.Users = append(input.Users, *user.WithProjectRef(input.Project.GetK8sMetaName()))
	}
	return input
}

// NewSimpleUserInputs prepare users inputs
func NewSimpleUserInputs(keyTestPrefix string, r *AtlasKeyType) UserInputs {
	testID := utils.GenID()
	projectName := fmt.Sprintf("%s-%s", keyTestPrefix, testID)
	input := UserInputs{
		TestID:             testID,
		AtlasKeyAccessType: *r,
		ProjectID:          "",
		KeyName:            keyTestPrefix,
		Namespace:          "ns-" + projectName,
		ProjectPath:        filepath.Join(DataGenFolder, projectName, "resources", projectName+".yaml"),
	}

	input.Project = NewProject("k-" + projectName).ProjectName(projectName)
	if len(input.Project.Spec.ProjectIPAccessList) == 0 {
		input.Project = input.Project.WithIpAccess("0.0.0.0/0", "everyone")
	}

	if !r.GlobalLevelKey {
		input.Project = input.Project.WithSecretRef(keyTestPrefix)
	}
	return input
}

func (u *UserInputs) GetAppFolder() string {
	return filepath.Join(DataGenFolder, u.Project.Spec.Name, "app")
}

func (u *UserInputs) GetOperatorFolder() string {
	return filepath.Join(DataGenFolder, u.Project.Spec.Name, "operator")
}

func (u *UserInputs) GetResourceFolder() string {
	return filepath.Dir(u.ProjectPath)
}

func (u *UserInputs) GetUsersFolder() string {
	return filepath.Join(u.GetResourceFolder(), "user")
}

func (u *UserInputs) GetServiceCatalogSourceFolder() string {
	return filepath.Join(DataGenFolder, u.Project.Spec.Name, "catalog")
}

func (u *UserInputs) GetAtlasProjectFullKubeName() string {
	return fmt.Sprintf("atlasproject/%s", u.Project.ObjectMeta.Name)
}
