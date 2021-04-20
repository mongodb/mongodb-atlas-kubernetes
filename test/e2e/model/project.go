package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	project "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

type ProjectSpec v1.AtlasProjectSpec

type AP struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec            ProjectSpec        `json:"spec,omitempty"`
}

// LoadUserProjectConfig load configuration from file into object
func LoadUserProjectConfig(path string) AP {
	var config AP
	utils.ReadInYAMLFileAndConvert(path, &config)
	return config
}

func NewProject(k8sname string) *AP {
	var t AP
	t.TypeMeta = metav1.TypeMeta{
		APIVersion: "atlas.mongodb.com/v1",
		Kind:       "AtlasProject",
	}
	t.ObjectMeta = &metav1.ObjectMeta{
		Name: k8sname,
	}
	return &t
}

func (p *AP) ProjectName(name string) *AP {
	p.Spec.Name = name
	return p
}

func (p *AP) SecretRef(name string) *AP {
	p.Spec.ConnectionSecret = &v1.ResourceRef{Name: name}
	return p
}

func (p *AP) WithIpAccess(ipAdress, comment string) *AP {
	access := project.NewIPAccessList().
		WithIP(ipAdress).
		WithComment(comment)
	p.Spec.ProjectIPAccessList = append(p.Spec.ProjectIPAccessList, access)
	return p
}

func (p *AP) GetK8sMetaName() string {
	return p.ObjectMeta.Name
}

func (p *AP) GetProjectName() string {
	return p.Spec.Name
}

func (p *AP) ConvertByte() []byte {
	yamlConf := utils.JSONToYAMLConvert(p)
	return yamlConf
}
