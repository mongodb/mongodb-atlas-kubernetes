package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

type ProjectSpec v1.AtlasProjectSpec

type AProject struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec            ProjectSpec        `json:"spec,omitempty"`
}

type AProjectWithStatus struct {
	AProject
	Status status.AtlasProjectStatus
}

// LoadUserProjectConfig load configuration from file into object
func LoadUserProjectConfig(path string) AProject {
	var config AProject
	utils.ReadInYAMLFileAndConvert(path, &config)
	return config
}

func NewProject(k8sname string) *AProject {
	var t AProject
	t.TypeMeta = metav1.TypeMeta{
		APIVersion: "atlas.mongodb.com/v1",
		Kind:       "AtlasProject",
	}
	t.ObjectMeta = &metav1.ObjectMeta{
		Name: k8sname,
	}
	return &t
}

func (p *AProject) ProjectName(name string) *AProject {
	p.Spec.Name = name
	return p
}

func (p *AProject) WithSecretRef(name string) *AProject {
	p.Spec.ConnectionSecret = &common.ResourceRef{Name: name}
	return p
}

func (p *AProject) WithIpAccess(ipAdress, comment string) *AProject {
	access := project.NewIPAccessList().
		WithIP(ipAdress).
		WithComment(comment)
	p.Spec.ProjectIPAccessList = append(p.Spec.ProjectIPAccessList, access)
	return p
}

func (p *AProject) WithPrivateLink(provider provider.ProviderName, region string) *AProject {
	link := v1.PrivateEndpoint{
		Provider: provider,
		Region:   region,
	}
	p.Spec.PrivateEndpoints = append(p.Spec.PrivateEndpoints, link)
	return p
}

func (p *AProject) WithIntegration(spec ProjectIntegration) *AProject {
	p.Spec.Integrations = append(p.Spec.Integrations, project.Integration(spec))
	return p
}

func (p *AProject) UpdatePrivateLinkByOrder(i int, id string) *AProject {
	p.Spec.PrivateEndpoints[i].ID = id
	return p
}

func (p *AProject) UpdatePrivateLinkID(test v1.PrivateEndpoint) *AProject {
	for i, peItem := range p.Spec.PrivateEndpoints {
		if (peItem.Provider == test.Provider) && (peItem.Region == test.Region) {
			p.Spec.PrivateEndpoints[i] = test
		}
	}
	return p
}

func (p *AProject) GetPrivateIDByProviderRegion(statusItem status.ProjectPrivateEndpoint) string {
	if statusItem.Provider == provider.ProviderAWS {
		for i, peItem := range p.Spec.PrivateEndpoints {
			if (peItem.Provider == statusItem.Provider) && (peItem.Region == statusItem.Region) {
				return p.Spec.PrivateEndpoints[i].ID
			}
		}
	}
	return statusItem.ID
}

func (p *AProject) DeletePrivateLink(id string) *AProject {
	var peList []v1.PrivateEndpoint
	for _, peItem := range p.Spec.PrivateEndpoints {
		if peItem.ID != id {
			peList = append(peList, peItem)
		}
	}
	p.Spec.PrivateEndpoints = peList
	return p
}

func (p *AProject) GetK8sMetaName() string {
	return p.ObjectMeta.Name
}

func (p *AProject) GetProjectName() string {
	return p.Spec.Name
}

func (p *AProject) ConvertByte() []byte {
	yamlConf := utils.JSONToYAMLConvert(p)
	return yamlConf
}
