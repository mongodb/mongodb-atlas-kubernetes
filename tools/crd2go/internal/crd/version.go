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

package crd

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type VersionedCRD struct {
	// Spec    *apiextensionsv1.CustomResourceDefinitionSpec
	Kind    string
	Version *apiextensionsv1.CustomResourceDefinitionVersion
}

func NewVersionedCRD(spec *apiextensionsv1.CustomResourceDefinitionSpec,
	version *apiextensionsv1.CustomResourceDefinitionVersion) *VersionedCRD {
	return &VersionedCRD{
		Kind:    spec.Names.Kind,
		Version: version,
	}
}

func (versionedCRD *VersionedCRD) SpecTypename() string {
	return fmt.Sprintf("%sSpec", versionedCRD.Kind)
}

func (versionedCRD *VersionedCRD) StatusTypename() string {
	return fmt.Sprintf("%sStatus", versionedCRD.Kind)
}

// SelectVersion returns the version from the CRD spec that matches the given version string
func SelectVersion(spec *apiextensionsv1.CustomResourceDefinitionSpec, version string) *VersionedCRD {
	if len(spec.Versions) == 0 {
		return nil
	}
	if version == "" {
		return NewVersionedCRD(spec, &spec.Versions[0])
	}
	for _, v := range spec.Versions {
		if v.Name == version {
			return NewVersionedCRD(spec, &v)
		}
	}
	return nil
}
