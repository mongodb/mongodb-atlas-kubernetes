package atlas

import (
	"fmt"
	"runtime"

	"github.com/mongodb-forks/digest"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/httputil"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
)

// TODO proper version passed on compile time
var userAgent = fmt.Sprintf("%s/%s (%s;%s)", "MongoDBAtlasKubernetesOperator", "version TODO", runtime.GOOS, runtime.GOARCH)

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
	client, err := mongodbatlas.New(&loggingClient, mongodbatlas.SetBaseURL("https://cloud-qa.mongodb.com/api/atlas/v1.0/"))
	if err != nil {
		return nil, err
	}
	client.UserAgent = userAgent
	return client, err
}
