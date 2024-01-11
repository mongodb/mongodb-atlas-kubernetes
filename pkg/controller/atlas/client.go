package atlas

import (
	"net/http"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/httputil"
)

func NewClient(domain, publicKey, privateKey string, logger *zap.SugaredLogger) (*mongodbatlas.Client, error) {
	clientCfg := []httputil.ClientOpt{
		httputil.Digest(publicKey, privateKey),
	}

	if logger != nil {
		clientCfg = append(clientCfg, httputil.LoggingTransport(logger))
	}

	httpClient, err := httputil.DecorateClient(&http.Client{Transport: http.DefaultTransport}, clientCfg...)
	if err != nil {
		return nil, err
	}

	return mongodbatlas.New(httpClient, mongodbatlas.SetBaseURL(domain), mongodbatlas.SetUserAgent(operatorUserAgent()))
}
