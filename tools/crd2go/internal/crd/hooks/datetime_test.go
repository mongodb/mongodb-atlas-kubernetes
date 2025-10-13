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

package hooks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/crd"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/gotype"
)

func TestDatetimeHookFn(t *testing.T) {
	tests := map[string]struct {
		crdType      *crd.CRDType
		expectedType *gotype.GoType
		expectedErr  error
	}{
		"not a datetime string type": {
			crdType: &crd.CRDType{
				Name: "NotADatetime",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type: crd.OpenAPIString,
				},
			},
			expectedErr: fmt.Errorf("string is not a date time (format is ): %w", crd.ErrNotProcessed),
		},
		"string type with format date-time": {
			crdType: &crd.CRDType{
				Name: "DatetimeString",
				Schema: &apiextensionsv1.JSONSchemaProps{
					Type:   crd.OpenAPIString,
					Format: "date-time",
				},
			},
			expectedType: gotype.TimeType,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := DatetimeHookFn(nil, nil, tt.crdType)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedType, got)
		})
	}
}
