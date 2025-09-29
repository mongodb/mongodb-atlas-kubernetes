package cmd

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunOpenapi2crd(t *testing.T) {
	tests := map[string]struct {
		input       string
		output      string
		overwrite   bool
		expectedErr error
	}{
		"generates CRD successfully": {
			input:     "./testdata/config.yaml",
			output:    "./testdata/output.yaml",
			overwrite: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			_, err := fs.Create(tt.input)
			require.NoError(t, err)
			err = afero.WriteFile(fs, tt.input, []byte(configFile()), 0o644)
			require.NoError(t, err)
			err = afero.WriteFile(fs, "./testdata/openapi.yaml", []byte(openapiFile()), 0o644)
			require.NoError(t, err)

			err = runOpenapi2crd(context.Background(), fs, tt.input, tt.output, tt.overwrite)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func configFile() string {
	return `kind: Config
apiVersion: atlas2crd.mongodb.com/v1alpha1
spec:
  pluginSets:
    - name: example
      default: true
      plugins:
        - base
        - major_version
        - parameter
        - entry
        - status
  openapi:
    - name: v1
      path: ./testdata/openapi.yaml
  crd:
    - gvk:
      version: v1
      kind: Example
      group: example.generated.mongodb.com
      categories:
        - example
      shortNames:
        - ex
      mappings:
        - majorVersion: v1
          openAPIRef:
            name: v1
          entry:
            schema: "ExampleRequest"
          status:
            schema: "ExampleResponse"`
}

func openapiFile() string {
	return `openapi: 3.0.0
info:
  title: Example API
  version: 1.0.0
components:
  schemas:
    ExampleRequest:
      type: object
      properties:
        name:
          type: string
        description:
          type: string
    ExampleResponse:
      type: object
      properties:
        id:
          type: string
        name:
          type: string
        description:
          type: string
paths:
  /example:
    post:
      summary: Create an example
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ExampleRequest'
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema: 
                $ref: '#/components/schemas/ExampleResponse'`
}
