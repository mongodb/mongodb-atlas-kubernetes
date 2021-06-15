package atlas

import (
	"context"
	"net/url"
	"os"
	"strings"

	"github.com/mongodb-forks/digest"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
)

type Atlas struct {
	OrgID   string
	Public  string
	Private string
	Client  *mongodbatlas.Client
}

func AClient() (Atlas, error) {
	var A Atlas
	A.OrgID = os.Getenv("MCLI_ORG_ID")
	A.Public = os.Getenv("MCLI_PUBLIC_API_KEY")
	A.Private = os.Getenv("MCLI_PRIVATE_API_KEY")
	t := digest.NewTransport(A.Public, A.Private)
	tc, err := t.Client()
	if err != nil {
		return A, err
	}
	A.Client = mongodbatlas.NewClient(tc)
	u, _ := url.Parse(config.AtlasAPIURL)
	A.Client.BaseURL = u
	return A, nil
}

func (a *Atlas) AddKeyWithAccessList(projectID string, roles, access []string) (public string, private string, err error) {
	createKeyRequest := &mongodbatlas.APIKeyInput{
		Desc:  "created from the AO test",
		Roles: roles,
	}
	newKey, _, err := a.Client.ProjectAPIKeys.Create(context.Background(), projectID, createKeyRequest)
	if err != nil {
		return "", "", err
	}
	createAccessRequest := formAccessRequest(access)
	_, _, err = a.Client.WhitelistAPIKeys.Create(context.Background(), a.OrgID, newKey.ID, createAccessRequest)
	if err != nil {
		return "", "", err
	}
	return newKey.PublicKey, newKey.PrivateKey, nil
}

func formAccessRequest(access []string) []*mongodbatlas.WhitelistAPIKeysReq {
	createRequest := make([]*mongodbatlas.WhitelistAPIKeysReq, 0)
	var req *mongodbatlas.WhitelistAPIKeysReq
	for _, item := range access {
		if strings.Contains(item, "/") {
			req = &mongodbatlas.WhitelistAPIKeysReq{CidrBlock: item}
		} else {
			req = &mongodbatlas.WhitelistAPIKeysReq{IPAddress: item}
		}
		createRequest = append(createRequest, req)
	}
	return createRequest
}
