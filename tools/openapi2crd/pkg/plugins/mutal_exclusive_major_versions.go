package plugins

import (
	"fmt"
	"strings"

	"github.com/mongodb/atlas2crd/pkg/processor"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

// (has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef))
type MutualExclusiveMajorVersions struct{}

func (m *MutualExclusiveMajorVersions) Name() string {
	return "mutual_exclusive_major_versions"
}

func (m *MutualExclusiveMajorVersions) Process(input processor.Input) error {
	i, ok := input.(*processor.CRDInput)

	if !ok {
		return nil // no operation performed
	}

	crd := i.CRD
	crdConfig := i.CRDConfig

	if len(crdConfig.Mappings) <= 1 {
		return nil
	}

	versions := make([]string, 0, len(crdConfig.Mappings))
	for _, mapping := range crdConfig.Mappings {
		versions = append(versions, mapping.MajorVersion)
	}

	cel := mutualExclusiveCEL(versions)
	specProps := crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"]
	specProps.XValidations = apiextensions.ValidationRules{
		{
			Rule:    cel,
			Message: fmt.Sprintf(`Only one of the following entries can be set: %q`, strings.Join(versions, ", ")),
		},
	}
	crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"] = specProps

	return nil
}

func NewMutualExclusiveMajorVersions() *MutualExclusiveMajorVersions {
	return &MutualExclusiveMajorVersions{}
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
