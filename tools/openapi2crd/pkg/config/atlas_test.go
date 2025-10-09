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

package config

import (
	"context"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAtlas_Load(t *testing.T) {
	tests := map[string]struct {
		pkg            string
		expectedSchema *openapi3.T
		expectedErrMsg string
	}{
		"valid package": {
			pkg:            "go.mongodb.org/atlas-sdk/v20250312005/admin",
			expectedSchema: &openapi3.T{},
		},
		"invalid package": {
			pkg:            "invalid/package/name",
			expectedErrMsg: "failed to load module path: failed to run 'go list' for module 'invalid/package/name'",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			openapiLoader := NewLoaderMock(t)
			if tt.expectedErrMsg == "" {
				openapiLoader.EXPECT().Load(context.Background(), mock.AnythingOfType("string")).Return(&openapi3.T{}, nil)
			}

			a := NewAtlas(openapiLoader)
			schema, err := a.Load(context.Background(), tt.pkg)
			if err != nil {
				assert.ErrorContains(t, err, tt.expectedErrMsg)
			}
			assert.Equal(t, tt.expectedSchema, schema)
		})
	}
}
