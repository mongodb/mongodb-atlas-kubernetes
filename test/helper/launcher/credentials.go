package launcher

import (
	"encoding/base64"
	"fmt"
)

// #nosec G101 -- This is just a template
const k8sSecretFmt = `apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: mongodb-atlas-operator-api-key
  labels:
    atlas.mongodb.com/type: credentials
data:
  orgId: %s
  publicApiKey: %s
  privateApiKey: %s`

type AtlasCredentials struct {
	OrgID      string
	PublicKey  string
	PrivateKey string
}

func credentialsFromEnv() AtlasCredentials {
	return AtlasCredentials{
		OrgID:      MustSetEnv("MCLI_ORG_ID"),
		PublicKey:  MustSetEnv("MCLI_PUBLIC_API_KEY"),
		PrivateKey: MustSetEnv("MCLI_PRIVATE_API_KEY"),
	}
}

func (crds *AtlasCredentials) secretYML() string {
	return fmt.Sprintf(k8sSecretFmt, encodeBase64(crds.OrgID),
		encodeBase64(crds.PublicKey),
		encodeBase64(crds.PrivateKey))
}

func encodeBase64(s string) string {
	return base64.StdEncoding.EncodeToString(([]byte)(s))
}
