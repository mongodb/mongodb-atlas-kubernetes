package client

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/mongodb-forks/digest"

	"go.mongodb.org/atlas/mongodbatlas"
)

func SetupClient(publicKey, privateKey, managerUrl string) (mongodbatlas.Client, error) {
	withDigest := func(c *http.Client) error {
		t := &digest.Transport{
			Username:  publicKey,
			Password:  privateKey,
			Transport: c.Transport,
		}
		c.Transport = t
		return nil
	}

	httpClient := &http.Client{Transport: http.DefaultTransport}

	err := withDigest(httpClient)
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
