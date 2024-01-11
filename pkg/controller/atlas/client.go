package atlas

import "go.mongodb.org/atlas-sdk/v20231115002/admin"

func NewClient(domain, publicKey, privateKey string) (*admin.APIClient, error) {
	return admin.NewClient(
		admin.UseBaseURL(domain),
		admin.UseDigestAuth(publicKey, privateKey),
		admin.UseUserAgent(operatorUserAgent()),
	)
}
