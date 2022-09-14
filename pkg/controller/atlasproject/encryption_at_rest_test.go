package atlasproject

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

func TestIsEncryptionAtlasEmpty(t *testing.T) {
	spec := &v1.EncryptionAtRest{}
	isEmpty := isEncryptionSpecEmpty(spec)
	assert.True(t, isEmpty, "Empty spec should be empty")

	spec.AwsKms.Enabled = toptr.MakePtr(true)
	isEmpty = isEncryptionSpecEmpty(spec)
	assert.False(t, isEmpty, "Non-empty spec")

	spec.AwsKms.Enabled = toptr.MakePtr(false)
	isEmpty = isEncryptionSpecEmpty(spec)
	assert.True(t, isEmpty, "Enabled flag set to false is same as empty")
}
