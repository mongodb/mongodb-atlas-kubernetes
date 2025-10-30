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

import (
	"fmt"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

// ConnectionSecret is a plugin that adds the Atlas Credentials secret to .
// It requires the parameters and references plugins to be run first.
type ConnectionSecret struct{}

func (p *ConnectionSecret) Name() string {
	return "connection_secret"
}

func (p *ConnectionSecret) Process(req *MappingProcessorRequest) error {
	specProps := req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"]

	if _, exists := specProps.Properties["connectionSecretRef"]; !exists {
		specProps.Properties["connectionSecretRef"] = apiextensions.JSONSchemaProps{
			Type:        "object",
			Description: "SENSITIVE FIELD\n\nReference to a secret containing the credentials to setup the connection to Atlas.",
			Properties: map[string]apiextensions.JSONSchemaProps{
				"name": {
					Type:        "string",
					Description: "Name of the secret containing the Atlas credentials.",
				},
			},
		}
	}

	version := req.MappingConfig.MajorVersion
	versionProps, ok := specProps.Properties[version]
	if !ok {
		return fmt.Errorf("version %s not found in spec", version)
	}

	_, groupIDExists := versionProps.Properties["groupId"]
	if groupIDExists {
		if specProps.XValidations == nil {
			specProps.XValidations = apiextensions.ValidationRules{}
		}
		specProps.XValidations = append(specProps.XValidations, apiextensions.ValidationRule{
			Rule:    "(has(self." + version + ".groupId) && has(self.connectionSecretRef)) || (!has(self." + version + ".groupId))",
			Message: fmt.Sprintf("spec.connectionSecretRef must be set if spec.%v.groupId is set.", version),
		})
	}
	req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"] = specProps

	return nil
}
