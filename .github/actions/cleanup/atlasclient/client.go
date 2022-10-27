package atlasclient

import (
	"fmt"
	"net/http"
	"runtime"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/httputil"
)

func SetupClient(publicKey, privateKey, managerUrl string) (mongodbatlas.Client, error) {
	withDigest := httputil.Digest(publicKey, privateKey)

	httpClient, err := httputil.DecorateClient(&http.Client{Transport: http.DefaultTransport}, withDigest)
	if err != nil {
		return mongodbatlas.Client{}, err
	}
	client, err := mongodbatlas.New(httpClient, mongodbatlas.SetBaseURL(managerUrl))
	if err != nil {
		return mongodbatlas.Client{}, err
	}
	client.UserAgent = fmt.Sprintf("%s/%s (%s;%s)", "MongoDBAtlasKubernetesOperator", "unknown", runtime.GOOS, runtime.GOARCH)
	return *client, nil
}
