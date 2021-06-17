package provider

type ProviderName string

const (
	ProviderAWS    ProviderName = "AWS"
	ProviderGCP    ProviderName = "GCP"
	ProviderAzure  ProviderName = "AZURE"
	ProviderTenant ProviderName = "TENANT"
)
