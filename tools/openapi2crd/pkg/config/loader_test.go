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

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
)

func TestCompositeLoader_Load(t *testing.T) {
	const testSpec = `openapi: 3.0.0
info:
  title: Test
  version: 1.0.0
paths: {}
`

	t.Run("path without flatten calls load", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		require.NoError(t, afero.WriteFile(fs, "spec.yaml", []byte(testSpec), 0644))

		kinLoader := NewKinOpenAPI(fs)
		atlas := NewAtlas(kinLoader)
		loader := NewLoader(kinLoader, atlas)

		def := v1alpha1.OpenAPIDefinition{Name: "v1", Path: "spec.yaml"}
		got, err := loader.Load(context.Background(), def)
		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "3.0.0", got.OpenAPI)
	})

	t.Run("path with flatten calls loadFlattened", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		require.NoError(t, afero.WriteFile(fs, "spec.yaml", []byte(testSpec), 0644))

		kinLoader := NewKinOpenAPI(fs)
		atlas := NewAtlas(kinLoader)
		loader := NewLoader(kinLoader, atlas)

		def := v1alpha1.OpenAPIDefinition{Name: "v1", Path: "spec.yaml", Flatten: true}
		got, err := loader.Load(context.Background(), def)
		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "3.0.0", got.OpenAPI)
	})

	t.Run("package calls atlas loadFromPackage", func(t *testing.T) {
		fs := afero.NewOsFs()
		kinLoader := NewKinOpenAPI(fs)
		atlas := NewAtlas(kinLoader)
		loader := NewLoader(kinLoader, atlas)

		def := v1alpha1.OpenAPIDefinition{
			Name:    "v1",
			Package: "go.mongodb.org/atlas-sdk/v20250312008/admin",
			Path:    "../openapi/atlas-api-transformed.yaml",
		}
		got, err := loader.Load(context.Background(), def)
		require.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("invalid package returns error", func(t *testing.T) {
		fs := afero.NewOsFs()
		kinLoader := NewKinOpenAPI(fs)
		atlas := NewAtlas(kinLoader)
		loader := NewLoader(kinLoader, atlas)

		def := v1alpha1.OpenAPIDefinition{
			Name:    "v1",
			Package: "invalid/package/name",
			Path:    "../openapi/spec.yaml",
		}
		_, err := loader.Load(context.Background(), def)
		assert.ErrorContains(t, err, "failed to load module path")
	})

	t.Run("missing path returns error", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		kinLoader := NewKinOpenAPI(fs)
		atlas := NewAtlas(kinLoader)
		loader := NewLoader(kinLoader, atlas)

		def := v1alpha1.OpenAPIDefinition{Name: "v1", Path: "nonexistent.yaml"}
		_, err := loader.Load(context.Background(), def)
		assert.Error(t, err)
	})

	t.Run("composite loader satisfies Loader interface", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		kinLoader := NewKinOpenAPI(fs)
		atlas := NewAtlas(kinLoader)
		var _ Loader = NewLoader(kinLoader, atlas)
	})
}

func TestCompositeLoader_LoadFlattened(t *testing.T) {
	const specWithOneOf = `openapi: 3.0.0
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
          type: string
`

	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "spec.yaml", []byte(specWithOneOf), 0644))

	kinLoader := NewKinOpenAPI(fs)
	atlas := NewAtlas(kinLoader)
	loader := NewLoader(kinLoader, atlas)

	def := v1alpha1.OpenAPIDefinition{Name: "v1", Path: "spec.yaml", Flatten: true}
	spec, err := loader.Load(context.Background(), def)
	require.NoError(t, err)
	require.NotNil(t, spec)

	petSchema := spec.Components.Schemas["Pet"]
	require.NotNil(t, petSchema)

	assert.Nil(t, petSchema.Value.OneOf, "oneOf should be flattened away")
	assert.Contains(t, petSchema.Value.Properties, "name")
	assert.Contains(t, petSchema.Value.Properties, "indoor")
	assert.Contains(t, petSchema.Value.Properties, "breed")
}

func TestCompositeLoader_LoadCachesResult(t *testing.T) {
	const testSpec = `openapi: 3.0.0
info:
  title: Test
  version: 1.0.0
paths: {}
`
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "spec.yaml", []byte(testSpec), 0644))

	kinLoader := NewKinOpenAPI(fs)
	atlas := NewAtlas(kinLoader)
	loader := NewLoader(kinLoader, atlas)

	def := v1alpha1.OpenAPIDefinition{Name: "v1", Path: "spec.yaml"}

	first, err := loader.Load(context.Background(), def)
	require.NoError(t, err)

	second, err := loader.Load(context.Background(), def)
	require.NoError(t, err)

	// Same pointer — cached result reused.
	assert.Same(t, first, second)
}
