package provider

const (
	ProviderAWS        ProviderName = "AWS"
	ProviderGCP        ProviderName = "GCP"
	ProviderAzure      ProviderName = "AZURE"
	ProviderTenant     ProviderName = "TENANT"
	ProviderServerless ProviderName = "SERVERLESS"
)

type ProviderName string
type CloudProviders map[ProviderName]struct{}

func (cp *CloudProviders) IsSupported(name ProviderName) bool {
	_, ok := (*cp)[name]

	return ok
}

func SupportedProviders() CloudProviders {
	return CloudProviders{
		ProviderAWS:   {},
		ProviderGCP:   {},
		ProviderAzure: {},
	}
}
