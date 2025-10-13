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

// MajorVersion is a plugin that adds the major version schema to the CRD.
// It requires the base plugin to be run first.
type MajorVersion struct{}

func (s *MajorVersion) Name() string {
	return "major_version"
}

func (s *MajorVersion) Process(req *MappingProcessorRequest) error {
	req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[req.MappingConfig.MajorVersion] = apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: fmt.Sprintf("The spec of the %v resource for version %v.", req.CRD.Spec.Names.Singular, req.MappingConfig.MajorVersion),
		Properties:  map[string]apiextensions.JSONSchemaProps{},
	}

	return nil
}
