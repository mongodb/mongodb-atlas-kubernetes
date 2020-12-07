package atlas

import (
	"github.com/mongodb-forks/digest"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/httputil"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
)

// AtlasClient is the central place to create a client for Atlas using specified API keys and a server URL.
// Note, that the default HTTP transport is reused globally by Go so all caching, keep-alive etc will be in action.
func AtlasClient(connection Connection, log *zap.SugaredLogger) (*mongodbatlas.Client, error) {
	t := digest.NewTransport(connection.PublicKey, connection.PrivateKey)
	tc, err := t.Client()
	if err != nil {
		return nil, err
	}
	loggingClient := httputil.NewLoggingClient(*tc, log)
	// TODO configuration for base URL (as a global Operator config?)
	return mongodbatlas.New(&loggingClient, mongodbatlas.SetBaseURL("https://cloud-qa.mongodb.com/api/atlas/v1.0/"))
}
