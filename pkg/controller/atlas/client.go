package atlas

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/httputil"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
)

// TODO proper version passed on compile time
var userAgent = fmt.Sprintf("%s/%s (%s;%s)", "MongoDBAtlasKubernetesOperator", "version TODO", runtime.GOOS, runtime.GOARCH)

// Client is the central place to create a client for Atlas using specified API keys and a server URL.
// Note, that the default HTTP transport is reused globally by Go so all caching, keep-alive etc will be in action.
func Client(connection Connection, log *zap.SugaredLogger) (*mongodbatlas.Client, error) {
	withDigest := httputil.Digest(connection.PublicKey, connection.PrivateKey)
	withLogging := httputil.LoggingTransport(log)

	httpClient, err := httputil.DecorateClient(basicClient(), withDigest, withLogging)
	if err != nil {
		return nil, err
	}
	client, err := mongodbatlas.New(httpClient, mongodbatlas.SetBaseURL(config.OperatorConfig.AtlasDomain+"/api/atlas/v1.0/"))
	if err != nil {
		return nil, err
	}
	client.UserAgent = userAgent
	return client, nil
}

func basicClient() *http.Client {
	// Do we need any custom configuration of timeout etc?
	return &http.Client{Transport: http.DefaultTransport}
}
