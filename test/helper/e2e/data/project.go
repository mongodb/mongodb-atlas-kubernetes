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

package data

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

const ProjectName = "my-project"

func DefaultProject() *akov2.AtlasProject {
	return &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name: ProjectName,
		},
		Spec: akov2.AtlasProjectSpec{
			Name: utils.RandomName("Test Atlas Operator Project"),
			ProjectIPAccessList: []project.IPAccessList{
				{
					CIDRBlock: "0.0.0.0/1",
					Comment:   "Everyone has access. For the test purpose only.",
				},
				{
					CIDRBlock: "128.0.0.0/1",
					Comment:   "Everyone has access. For the test purpose only.",
				},
			},
		},
	}
}
