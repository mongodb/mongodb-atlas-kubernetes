package plugins

// AtlasSdkVersionPlugin is a plugin that adds the Atlas SDK version information as an OpenAPI extension in the CRD.
// It requires the entry plugin to be run first.
type AtlasSdkVersionPlugin struct{}

func (p *AtlasSdkVersionPlugin) Name() string {
	return "atlas_sdk_version"
}

func (p *AtlasSdkVersionPlugin) Process(req *ExtensionProcessorRequest) error {
	pkg := req.ApiDefinitions[req.MappingConfig.OpenAPIRef.Name].Package
	if pkg == "" {
		return nil
	}

	extensions := req.ExtensionsSchema.Properties["spec"].Value.Properties[req.MappingConfig.MajorVersion].Value.Extensions
	if extensions == nil {
		extensions = map[string]interface{}{}
	}
	extensions["x-atlas-sdk-version"] = pkg
	req.ExtensionsSchema.Properties["spec"].Value.Properties[req.MappingConfig.MajorVersion].Value.Extensions = extensions

	return nil
}
