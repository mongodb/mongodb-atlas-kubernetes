package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cel"
)

type celTestCase map[string]struct {
	object         AtlasCustomResource
	oldObject      AtlasCustomResource
	expectedErrors []string
}

func assertCELValidation(t *testing.T, crdPath string, tests celTestCase) {
	t.Helper()

	validator, err := cel.VersionValidatorFromFile(t, crdPath, "v1")
	require.NoError(t, err)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tt.object)
			require.NoError(t, err)

			unstructuredOldObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tt.oldObject)
			require.NoError(t, err)

			errs := validator(unstructuredObject, unstructuredOldObject)

			require.Equal(t, len(tt.expectedErrors), len(errs))

			for i, err := range errs {
				assert.Equal(t, tt.expectedErrors[i], err.Error())
			}
		})
	}
}
