package config

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKinOpeAPILoad(t *testing.T) {
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

			loader := NewKinOpeAPI(fs)
			openapi, err := loader.Load(tt.filePath)
			assert.Equal(t, tt.expectError, err)
			assert.Equal(t, tt.expectedOpenAPI, openapi)
		})
	}
}
