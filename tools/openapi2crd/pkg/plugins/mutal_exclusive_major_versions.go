package plugins

import (
	"fmt"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"strings"
)

// (has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef))
type MutualExclusiveMajorVersions struct {
	NoOp
	crd *apiextensions.CustomResourceDefinition
}

var _ Plugin = &MutualExclusiveMajorVersions{}

func NewMutualExclusiveMajorVersions(crd *apiextensions.CustomResourceDefinition) *MutualExclusiveMajorVersions {
	return &MutualExclusiveMajorVersions{
		crd: crd,
	}
}

func (s *MutualExclusiveMajorVersions) Name() string {
	return "mutual_exclusive_major_versions"
}

func (m *MutualExclusiveMajorVersions) ProcessCRD(g Generator, crdConfig *configv1alpha1.CRDConfig) error {
	if len(crdConfig.Mappings) <= 1 {
		return nil
	}

	majorVersions := make([]string, 0, len(crdConfig.Mappings))
	for _, mapping := range crdConfig.Mappings {
		majorVersions = append(majorVersions, mapping.MajorVersion)
	}

	cel := mutualExclusiveCEL(majorVersions)
	specProps := m.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"]
	specProps.XValidations = apiextensions.ValidationRules{
		{
			Rule:    cel,
			Message: fmt.Sprintf(`Only one of the following entries can be set: %q`, strings.Join(majorVersions, ", ")),
		},
	}
	m.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"] = specProps

	return nil
}

func mutualExclusiveCEL(fields []string) string {
	clauses := make([]string, 0, len(fields))
	for i, _ := range fields {
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
