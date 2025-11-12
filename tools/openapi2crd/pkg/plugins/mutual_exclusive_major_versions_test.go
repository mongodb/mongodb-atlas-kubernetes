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
//

package plugins

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"

	configv1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
)

func TestMutualExclusiveMajorVersionsName(t *testing.T) {
	p := &MutualExclusiveMajorVersions{}
	assert.Equal(t, "mutual_exclusive_major_versions", p.Name())
}

func TestMutualExclusiveMajorVersionsProcess(t *testing.T) {
	tests := map[string]struct {
		request            *CRDProcessorRequest
		expectedValidation apiextensions.ValidationRules
		expectedErr        error
	}{
		"add mutual exclusive major versions validation": {
			request:            groupMultipleVersionsCRDConfig(t, groupBaseCRD(t)),
			expectedValidation: mutualExclusiveMajorVersionsValidation(t),
		},
		"add no validation when there's only one version": {
			request: groupCRDRequest(t, groupBaseCRD(t)),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &MutualExclusiveMajorVersions{}
			err := p.Process(tt.request)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedValidation, tt.request.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].XValidations)
		})
	}
}

func groupMultipleVersionsCRDConfig(t *testing.T, crd *apiextensions.CustomResourceDefinition) *CRDProcessorRequest {
	t.Helper()

	config := groupCRDRequest(t, crd)
	config.CRDConfig.Mappings = append(
		config.CRDConfig.Mappings,
		configv1alpha1.CRDMapping{
			MajorVersion: "v20250219",
			OpenAPIRef: configv1alpha1.LocalObjectReference{
				Name: "v20250219",
			},
			ParametersMapping: configv1alpha1.PropertyMapping{
				Path: configv1alpha1.PropertyPath{
					Name: "/api/atlas/v2/groups",
					Verb: "post",
				},
			},
			EntryMapping: configv1alpha1.PropertyMapping{
				Schema: "Group",
				Filters: configv1alpha1.Filters{
					ReadWriteOnly: true,
				},
			},
			StatusMapping: configv1alpha1.PropertyMapping{
				Schema: "Group",
				Filters: configv1alpha1.Filters{
					ReadOnly:       true,
					SkipProperties: []string{"$.links"},
				},
			},
		},
	)

	return config
}

func mutualExclusiveMajorVersionsValidation(t *testing.T) apiextensions.ValidationRules {
	t.Helper()

	return apiextensions.ValidationRules{
		{
			Rule:    "!has(self.v20250312) && has(self.v20250219) || has(self.v20250312) && !has(self.v20250219)",
			Message: `Only one of the following entries can be set: "v20250312, v20250219"`,
		},
	}
}
