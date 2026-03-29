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
	"sync"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testOpenAPISpec = `openapi: 3.0.0
info:
  title: Test
  version: 1.0.0
paths: {}`

func writeSpec(t *testing.T, fs afero.Fs, path string) {
	t.Helper()
	require.NoError(t, afero.WriteFile(fs, path, []byte(testOpenAPISpec), 0644))
}

func TestKinOpenAPILoadCachesResult(t *testing.T) {
	fs := afero.NewMemMapFs()
	writeSpec(t, fs, "spec.yaml")

	loader := NewKinOpenAPI(fs)

	first, err := loader.Load(context.Background(), "spec.yaml")
	require.NoError(t, err)

	second, err := loader.Load(context.Background(), "spec.yaml")
	require.NoError(t, err)

	// Same pointer — the second call returned the cached result.
	assert.Same(t, first, second)
}

func TestKinOpenAPILoadConcurrent(t *testing.T) {
	fs := afero.NewMemMapFs()
	writeSpec(t, fs, "spec.yaml")

	loader := NewKinOpenAPI(fs)

	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)
	results := make([]*openapi3.T, goroutines)

	for i := range goroutines {
		go func() {
			defer wg.Done()
			spec, err := loader.Load(context.Background(), "spec.yaml")
			require.NoError(t, err)
			results[i] = spec
		}()
	}
	wg.Wait()

	// All goroutines got the same pointer.
	for i := 1; i < goroutines; i++ {
		assert.Same(t, results[0], results[i])
	}
}

func TestKinOpenAPILoadFlattened(t *testing.T) {
	input := `openapi: 3.0.0
info:
  title: Test
  version: 1.0.0
paths:
  /pets:
    get:
      operationId: listPets
      responses:
        '200':
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Pet'
components:
  schemas:
    Pet:
      type: object
      oneOf:
        - $ref: '#/components/schemas/Cat'
        - $ref: '#/components/schemas/Dog'
      properties:
        name:
          type: string
    Cat:
      type: object
      properties:
        indoor:
          type: boolean
    Dog:
      type: object
      properties:
        breed:
          type: string`

	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "test.yaml", []byte(input), 0644))

	loader := NewKinOpenAPI(fs)
	spec, err := loader.LoadFlattened(nil, "test.yaml")
	require.NoError(t, err)
	require.NotNil(t, spec)

	petSchema := spec.Components.Schemas["Pet"]
	require.NotNil(t, petSchema)

	assert.Nil(t, petSchema.Value.OneOf, "oneOf should be flattened away")
	assert.Contains(t, petSchema.Value.Properties, "name")
	assert.Contains(t, petSchema.Value.Properties, "indoor")
	assert.Contains(t, petSchema.Value.Properties, "breed")
}

func TestKinOpenAPILoad(t *testing.T) {
	tests := map[string]struct {
		file            string
		filePath        string
		expectedOpenAPI *openapi3.T
		expectError     error
	}{
		"valid openapi file": {
			file: `openapi: 3.0.0
info:
  title: Swagger Petstore
  version: 1.0.0
paths:
  /pets:
    get:
      operationId: listPets
      x-xgen-changelog:
        2025-05-08: Corrects an issue where the endpoint would include Atlas internal entries.`,
			filePath: "testdata/petstore.yaml",
			expectedOpenAPI: &openapi3.T{
				OpenAPI: "3.0.0",
				Info: &openapi3.Info{
					Title:   "Swagger Petstore",
					Version: "1.0.0",
				},
				Paths: openapi3.NewPaths(
					openapi3.WithPath(
						"/pets",
						&openapi3.PathItem{
							Get: &openapi3.Operation{
								OperationID: "listPets",
							},
						},
					),
				),
			},
			expectError: nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			_, err := fs.Create(tt.filePath)
			require.NoError(t, err)
			err = afero.WriteFile(fs, tt.filePath, []byte(tt.file), 0644)
			require.NoError(t, err)

			tt.expectedOpenAPI.Paths.Extensions = map[string]any{}

			loader := NewKinOpenAPI(fs)
			openapi, err := loader.Load(context.Background(), tt.filePath)
			assert.Equal(t, tt.expectError, err)
			assert.Equal(t, tt.expectedOpenAPI, openapi)
		})
	}
}
