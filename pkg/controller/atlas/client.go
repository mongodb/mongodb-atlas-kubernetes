package atlas

import (
	"fmt"
	"net/http"
	"runtime"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/httputil"
)

const (
	basicVersion    = "1.0"
	advancedVersion = "1.5"
)

// ProductVersion is used for sending the current Operator version in the User-Agent string
var ProductVersion = "unknown"

// AllClients returns all atlas clients required by the operator.
func AllClients(atlasDomain string, connection Connection, log *zap.SugaredLogger) (mongodbatlas.Client, mongodbatlas.Client, error) {
	v1Client, err := Client(atlasDomain, connection, log)
	if err != nil {
		return mongodbatlas.Client{}, mongodbatlas.Client{}, fmt.Errorf("failed creating basic client: %w", err)
	}
	advancedClient, err := AdvancedClient(atlasDomain, connection, log)
	if err != nil {
		return mongodbatlas.Client{}, mongodbatlas.Client{}, fmt.Errorf("failed creating advanced client: %w", err)
	}

	return v1Client, advancedClient, nil
}

// Client is the central place to create a client for Atlas using specified API keys and a server URL.
// Note, that the default HTTP transport is reused globally by Go so all caching, keep-alive etc will be in action.
func Client(atlasDomain string, connection Connection, log *zap.SugaredLogger) (mongodbatlas.Client, error) {
	return versionedClient(atlasDomain, connection, basicVersion, log)
}

func AdvancedClient(atlasDomain string, connection Connection, log *zap.SugaredLogger) (mongodbatlas.Client, error) {
	return versionedClient(atlasDomain, connection, advancedVersion, log)
}

func versionedClient(atlasDomain string, connection Connection, version string, log *zap.SugaredLogger) (mongodbatlas.Client, error) {
	withDigest := httputil.Digest(connection.PublicKey, connection.PrivateKey)
	withLogging := httputil.LoggingTransport(log)

	httpClient, err := httputil.DecorateClient(basicClient(), withDigest, withLogging)
	if err != nil {
		return mongodbatlas.Client{}, err
	}
	client, err := mongodbatlas.New(httpClient, mongodbatlas.SetBaseURL(fmt.Sprintf("%s/api/atlas/v%s/", atlasDomain, version)))
	if err != nil {
		return mongodbatlas.Client{}, err
	}
	client.UserAgent = fmt.Sprintf("%s/%s (%s;%s)", "MongoDBAtlasKubernetesOperator", ProductVersion, runtime.GOOS, runtime.GOARCH)

	return *client, nil
}

func basicClient() *http.Client {
	// Do we need any custom configuration of timeout etc?
	return &http.Client{Transport: http.DefaultTransport}
}
