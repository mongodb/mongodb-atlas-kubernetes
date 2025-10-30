package plugins

import (
	"fmt"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type MutualExclusiveGroup struct{}

func (p *MutualExclusiveGroup) Name() string {
	return "mutual_exclusive_group"
}

func (p *MutualExclusiveGroup) Process(req *MappingProcessorRequest) error {
	version := req.MappingConfig.MajorVersion
	if _, ok := req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[version]; !ok {
		return fmt.Errorf("version %s not found in spec properties", version)
	}

	versionProps := req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[version]
	_, groupIDExists := versionProps.Properties["groupId"]
	_, groupRefExists := versionProps.Properties["groupRef"]
	if groupIDExists || groupRefExists {
		versionProps.XValidations = append(versionProps.XValidations, apiextensions.ValidationRule{
			Rule:    "(has(self.groupId) && !has(self.groupRef)) || (!has(self.groupId) && has(self.groupRef))",
			Message: "groupId and groupRef are mutually exclusive; only one of them can be set",
		})
	}

	req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[version] = versionProps

	return nil
}
