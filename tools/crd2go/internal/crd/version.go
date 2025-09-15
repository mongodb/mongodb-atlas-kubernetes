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
