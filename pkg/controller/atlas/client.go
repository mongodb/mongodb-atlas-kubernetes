package atlas

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/version"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/httputil"
)

// Client is the central place to create a client for Atlas using specified API keys and a server URL.
// Note, that the default HTTP transport is reused globally by Go so all caching, keep-alive etc will be in action.
func Client(atlasDomain string, connection Connection, log *zap.SugaredLogger) (mongodbatlas.Client, error) {
	withLogging := httputil.LoggingTransport(log)
	return CustomClient(atlasDomain, connection, withLogging)
}

func CustomClient(atlasDomain string, connection Connection, opts ...httputil.ClientOpt) (mongodbatlas.Client, error) {
	withDigest := httputil.Digest(connection.PublicKey, connection.PrivateKey)
	allOptions := []httputil.ClientOpt{withDigest}
	allOptions = append(allOptions, opts...)

	httpClient, err := httputil.DecorateClient(basicClient(), allOptions...)
	if err != nil {
		return mongodbatlas.Client{}, err
	}
	client, err := mongodbatlas.New(httpClient, mongodbatlas.SetBaseURL(atlasDomain))
	if err != nil {
		return mongodbatlas.Client{}, err
	}
	client.UserAgent = fmt.Sprintf("%s/%s (%s;%s)", "MongoDBAtlasKubernetesOperator", version.Version, runtime.GOOS, runtime.GOARCH)

	return *client, nil
}

func basicClient() *http.Client {
	// Do we need any custom configuration of timeout etc?
	return &http.Client{Transport: http.DefaultTransport}
}
