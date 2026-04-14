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

func TestAtlas_LoadCachesPath(t *testing.T) {
	// Calling LoadFromPackage twice with the same package should resolve the path once
	// and call the file loader twice (once per call), but with the same resolved path.
	openapiLoader := NewLoaderMock(t)
	openapiLoader.EXPECT().Load(context.Background(), mock.AnythingOfType("string")).Return(&openapi3.T{}, nil).Times(2)

	a := NewAtlas(openapiLoader)
	relPath := "../openapi/atlas-api-transformed.yaml"

	_, err := a.LoadFromPackage(context.Background(), "go.mongodb.org/atlas-sdk/v20250312008/admin", relPath)
	assert.NoError(t, err)

	// The path should now be cached.
	assert.Len(t, a.pathCache, 1)

	_, err = a.LoadFromPackage(context.Background(), "go.mongodb.org/atlas-sdk/v20250312008/admin", relPath)
	assert.NoError(t, err)

	// Still only one entry — path was reused, not re-resolved.
	assert.Len(t, a.pathCache, 1)
}

func TestAtlas_LoadFromPackage(t *testing.T) {
	tests := map[string]struct {
		pkg            string
		relPath        string
		expectedSchema *openapi3.T
		expectedErrMsg string
	}{
		"valid package with relative path": {
			pkg:            "go.mongodb.org/atlas-sdk/v20250312008/admin",
			relPath:        "../openapi/atlas-api-transformed.yaml",
			expectedSchema: &openapi3.T{},
		},
		"invalid package": {
			pkg:            "invalid/package/name",
			relPath:        "../openapi/atlas-api-transformed.yaml",
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
			schema, err := a.LoadFromPackage(context.Background(), tt.pkg, tt.relPath)
			if err != nil {
				assert.ErrorContains(t, err, tt.expectedErrMsg)
			}
			assert.Equal(t, tt.expectedSchema, schema)
		})
	}
}
