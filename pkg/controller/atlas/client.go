package atlas

import (
	"os"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

func NewClient(domain, publicKey, privateKey string) (*admin.APIClient, error) {
	debug := os.Getenv("AKO_ATLAS_SDK_DEBUG") == "1"
	return admin.NewClient(
		admin.UseDebug(debug),
		admin.UseBaseURL(domain),
		admin.UseDigestAuth(publicKey, privateKey),
		admin.UseUserAgent(operatorUserAgent()),
	)
}
