package launcher

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
