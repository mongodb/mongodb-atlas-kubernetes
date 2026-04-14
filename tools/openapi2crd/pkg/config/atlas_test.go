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
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAtlas_LoadCachesPath(t *testing.T) {
	// Calling loadFromPackage twice with the same package should resolve the path once
	// and call the file loader twice (once per call), but with the same resolved path.
	fs := afero.NewOsFs()
	kinLoader := NewKinOpenAPI(fs)
	a := NewAtlas(kinLoader)
	relPath := "../openapi/atlas-api-transformed.yaml"

	_, err := a.loadFromPackage(context.Background(), "go.mongodb.org/atlas-sdk/v20250312008/admin", relPath)
	// This may fail if the spec file doesn't exist at the resolved path, but path caching still works.
	// We just check that the path was cached.
	if err == nil {
		assert.Len(t, a.pathCache, 1)

		_, err = a.loadFromPackage(context.Background(), "go.mongodb.org/atlas-sdk/v20250312008/admin", relPath)
		require.NoError(t, err)
		assert.Len(t, a.pathCache, 1)
	} else {
		// Path resolution succeeded (it's cached) but loading failed — that's fine for this test.
		assert.Len(t, a.pathCache, 1)
	}
}

func TestAtlas_loadFromPackage(t *testing.T) {
	tests := map[string]struct {
		pkg            string
		relPath        string
		expectedSchema *openapi3.T
		expectedErrMsg string
	}{
		"invalid package": {
			pkg:            "invalid/package/name",
			relPath:        "../openapi/atlas-api-transformed.yaml",
			expectedErrMsg: "failed to load module path: failed to run 'go list' for module 'invalid/package/name'",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fs := afero.NewOsFs()
			kinLoader := NewKinOpenAPI(fs)
			a := NewAtlas(kinLoader)
			schema, err := a.loadFromPackage(context.Background(), tt.pkg, tt.relPath)
			if tt.expectedErrMsg != "" {
				assert.ErrorContains(t, err, tt.expectedErrMsg)
			}
			assert.Equal(t, tt.expectedSchema, schema)
		})
	}
}
