package atlas_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
)

func TestClientUserAgent(t *testing.T) {
	require := require.New(t)

	atlas.ProductVersion = "v1.2.3-testing"
	c, err := atlas.Client("https://cloud.mongodb.com", atlas.Connection{}, nil)
	require.NoError(err)
	require.Regexp(`^MongoDBAtlasKubernetesOperatorRHODA/v1\.2\.3-testing \(\w+;\w+\)$`, c.UserAgent)
}
