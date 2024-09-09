package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cel"
)

type projectReferenceTestCase map[string]struct {
	object         AtlasCustomResource
	expectedErrors []string
}

func assertProjectReference(t *testing.T, crdPath string, tests projectReferenceTestCase) {
	t.Helper()

	validator, err := cel.VersionValidatorFromFile(t, crdPath, "v1")
	require.NoError(t, err)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tt.object)
			require.NoError(t, err)

			errs := validator(unstructuredObject, nil)

			require.Equal(t, len(tt.expectedErrors), len(errs))

			for i, err := range errs {
				assert.Equal(t, tt.expectedErrors[i], err.Error())
			}
		})
	}
}

func assertExternalProjectReferenceConnectionSecret(t *testing.T, crdPath string, tests projectReferenceTestCase) {
	t.Helper()

	validator, err := cel.VersionValidatorFromFile(t, crdPath, "v1")
	require.NoError(t, err)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tt.object)
			require.NoError(t, err)

			errs := validator(unstructuredObject, nil)

			require.Equal(t, len(tt.expectedErrors), len(errs))

			for i, err := range errs {
				assert.Equal(t, tt.expectedErrors[i], err.Error())
			}
		})
	}
}
