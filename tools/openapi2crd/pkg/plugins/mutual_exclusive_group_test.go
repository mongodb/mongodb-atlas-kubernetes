package plugins

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

func TestMutualExclusiveGroupName(t *testing.T) {
	p := &MutualExclusiveGroup{}
	assert.Equal(t, "mutual_exclusive_group", p.Name())
}

func TestMutualExclusiveGroupProcess(t *testing.T) {
	tests := map[string]struct {
		request            *MappingProcessorRequest
		expectedValidation apiextensions.ValidationRules
		expectedErr        error
	}{
		"add mutual exclusive validation for groupId and groupRef": {
			request: mappingRequestWithReferences(t, groupBaseCRDWithParameters(t)),
			expectedValidation: apiextensions.ValidationRules{
				{
					Rule:    "(has(self.groupId) && !has(self.groupRef)) || (!has(self.groupId) && has(self.groupRef))",
					Message: "groupId and groupRef are mutually exclusive; only one of them can be set",
				},
			},
		},
		"no validation added when neither groupId nor groupRef exist": {
			request: mappingRequestWithReferences(t, orgBaseCRDWithParameters(t)),
		},
		"error when version property is missing": {
			request:     mappingRequestWithReferences(t, groupBaseCRD(t)),
			expectedErr: errors.New("version v20250312 not found in spec properties"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &MutualExclusiveGroup{}
			err := p.Process(tt.request)
			assert.Equal(t, tt.expectedErr, err)
			validations := tt.request.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[tt.request.MappingConfig.MajorVersion].XValidations
			assert.Equal(t, tt.expectedValidation, validations)
		})
	}
}
