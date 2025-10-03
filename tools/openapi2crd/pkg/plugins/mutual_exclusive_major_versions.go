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
	"strings"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

// MutualExclusiveMajorVersions is a plugin that adds a CEL validation to the CRD to ensure that only one of the major
// versions is set in the spec. It requires base plugin to be run first.
type MutualExclusiveMajorVersions struct{}

func (p *MutualExclusiveMajorVersions) Name() string {
	return "mutual_exclusive_major_versions"
}

func (p *MutualExclusiveMajorVersions) Process(req *CRDProcessorRequest) error {
	if len(req.CRDConfig.Mappings) <= 1 {
		return nil
	}

	versions := make([]string, 0, len(req.CRDConfig.Mappings))
	for _, mapping := range req.CRDConfig.Mappings {
		versions = append(versions, mapping.MajorVersion)
	}

	cel := mutualExclusiveCEL(versions)
	specProps := req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"]
	specProps.XValidations = apiextensions.ValidationRules{
		{
			Rule:    cel,
			Message: fmt.Sprintf(`Only one of the following entries can be set: %q`, strings.Join(versions, ", ")),
		},
	}
	req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"] = specProps

	return nil
}

func mutualExclusiveCEL(fields []string) string {
	clauses := make([]string, 0, len(fields))
	for i := range fields {
		parts := make([]string, len(fields))
		for j, name := range fields {
			if i == j {
				parts[j] = fmt.Sprintf("!has(self.%s)", name)
			} else {
				parts[j] = fmt.Sprintf("has(self.%s)", name)
			}
		}
		clauses = append(clauses, strings.Join(parts, " && "))
	}
	return strings.Join(clauses, " || ")
}
