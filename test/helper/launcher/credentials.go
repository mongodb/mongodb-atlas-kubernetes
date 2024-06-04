package launcher

type AtlasCredentials struct {
	OrgID      string
	PublicKey  string
	PrivateKey string
}

func credentialsFromEnv() AtlasCredentials {
	return AtlasCredentials{
		OrgID:      MustLookupEnv("MCLI_ORG_ID"),
		PublicKey:  MustLookupEnv("MCLI_PUBLIC_API_KEY"),
		PrivateKey: MustLookupEnv("MCLI_PRIVATE_API_KEY"),
	}
}
