package atlas_test

import (
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/version"

	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
)

func TestClientUserAgent(t *testing.T) {
	require := require.New(t)

	c, err := atlas.Client("https://cloud.mongodb.com", atlas.Connection{}, nil)
	require.NoError(err)
	require.Contains(c.UserAgent, version.Version)
	//require.Regexp(`^MongoDBAtlasKubernetesOperator/v1\.2\.3-testing \(\w+;\w+\)$`, c.UserAgent)
}
