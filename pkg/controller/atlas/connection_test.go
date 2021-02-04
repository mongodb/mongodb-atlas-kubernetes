package atlas

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

func Test_validateConnectionSecret(t *testing.T) {
	err := validateConnectionSecret(kube.ObjectKey("testNs", "testSecret"), map[string]string{})
	assert.EqualError(t, err, "the following fields are missing in the Secret testNs/testSecret: [orgId publicApiKey privateApiKey]")

	err = validateConnectionSecret(kube.ObjectKey("testNs", "testSecret"), map[string]string{"publicApiKey": "foo"})
	assert.EqualError(t, err, "the following fields are missing in the Secret testNs/testSecret: [orgId privateApiKey]")

	err = validateConnectionSecret(kube.ObjectKey("testNs", "testSecret"), map[string]string{"orgId": "some", "publicApiKey": "foo"})
	assert.EqualError(t, err, "the following fields are missing in the Secret testNs/testSecret: [privateApiKey]")

	assert.NoError(t, validateConnectionSecret(kube.ObjectKey("testNs", "testSecret"), map[string]string{"orgId": "some", "publicApiKey": "foo", "privateApiKey": "bla"}))
}
