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

package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cel"
)

func TestDeploymentCELChecks(t *testing.T) {
	for _, tc := range []struct {
		title          string
		old, obj       *AtlasDeployment
		expectedErrors []string
	}{
		{
			title: "Cannot rename a deployment",
			old: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					DeploymentSpec: &AdvancedDeploymentSpec{
						Name: "name-old",
					},
				},
			},
			obj: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					DeploymentSpec: &AdvancedDeploymentSpec{
						Name: "name-new",
					},
				},
			},
			expectedErrors: []string{"spec.deploymentSpec.name: Invalid value: \"string\": Name cannot be modified after deployment creation"},
		},
		{
			title: "Cannot add a serverless deployment",
			old:   nil,
			obj: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					ServerlessSpec: &ServerlessSpec{
						Name: "my-serverless",
					},
				},
			},
			expectedErrors: []string{"spec.serverlessSpec: Invalid value: \"object\": serverlessSpec cannot be added - serverless instances are deprecated"},
		},
		{
			title: "Can modify to a serverless deployment",
			old: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					ServerlessSpec: &ServerlessSpec{
						Name:                         "my-serverless",
						TerminationProtectionEnabled: false,
					},
				},
			},
			obj: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					ServerlessSpec: &ServerlessSpec{
						Name:                         "my-serverless",
						TerminationProtectionEnabled: true,
					},
				},
			},
		},
		{
			title: "Existing serverless deployment can continue existing when not modified",
			old: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					ServerlessSpec: &ServerlessSpec{
						Name:                         "my-serverless",
						TerminationProtectionEnabled: false,
					},
				},
			},
			obj: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					BackupScheduleRef: common.ResourceRefNamespaced{},
					ServerlessSpec: &ServerlessSpec{
						Name:                         "my-serverless",
						TerminationProtectionEnabled: false,
					},
				},
			},
		},
		{
			title: "can migrate from serverless to flex cluster",
			old: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					ServerlessSpec: &ServerlessSpec{
						Name: "my-serverless",
					},
				},
			},
			obj: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					FlexSpec: &FlexSpec{
						Name: "my-serverless",
					},
				},
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			// inject a project to avoid other CEL validations being hit
			tc.obj.Spec.ProjectRef = &common.ResourceRefNamespaced{Name: "some-project"}
			unstructuredOldObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tc.old)
			require.NoError(t, err)
			unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tc.obj)
			require.NoError(t, err)

			crdPath := "../../config/crd/bases/atlas.mongodb.com_atlasdeployments.yaml"
			validator, err := cel.VersionValidatorFromFile(t, crdPath, "v1")
			assert.NoError(t, err)
			errs := validator(unstructuredObject, unstructuredOldObject)

			require.Equal(t, tc.expectedErrors, cel.ErrorListAsStrings(errs))
		})
	}
}
