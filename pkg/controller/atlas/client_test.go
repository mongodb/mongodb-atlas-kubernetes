package atlas_test

import (
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/version"

	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
)

func TestClientUserAgent(t *testing.T) {
	r := require.New(t)

	c, err := atlas.Client("https://cloud.mongodb.com", atlas.Connection{}, nil)
	r.NoError(err)
	r.Contains(c.UserAgent, version.Version)
}
