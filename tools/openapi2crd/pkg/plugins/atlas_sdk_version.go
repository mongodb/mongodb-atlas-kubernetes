// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
