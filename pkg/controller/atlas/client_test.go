package atlas_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/version"
)

func TestClientUserAgent(t *testing.T) {
	r := require.New(t)

	c, err := atlas.Client("https://cloud.mongodb.com", atlas.Connection{}, nil)
	r.NoError(err)
	r.Contains(c.UserAgent, version.Version)
}

// RoundTrip implements http.RoundTripper registering if it got used
type testTransport struct {
	used bool
}

func (tt *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	tt.used = true
	return nil, nil
}

func TestCustomTransport(t *testing.T) {
	tt := &testTransport{used: false}

	client, err := atlas.Client("https://cloud.mongodb.com", atlas.Connection{}, nil, httputil.CustomTransport(tt))
	require.NoError(t, err)
	client.Projects.GetAllProjects(context.Background(), &mongodbatlas.ListOptions{})
	assert.True(t, tt.used)
}
