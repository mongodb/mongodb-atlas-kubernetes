package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	project "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

type ISpec interface {
	ProjectName(string) ISpec
	SecretRef(string) ISpec
	WithIpAccess(string, string) ISpec
	CompleteK8sConfig(string) []byte
}

type ProjectSpec v1.AtlasProjectSpec

type ap struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec            ProjectSpec        `json:"spec,omitempty"`
}

// LoadUserProjectConfig load configuration from file into object
func LoadUserProjectConfig(path string) ap {
	var config ap
	utils.ReadInYAMLFileAndConvert(path, &config)
	return config
}

func NewProject() ISpec {
	return &ProjectSpec{}
}

func (s *ProjectSpec) ProjectName(name string) ISpec {
	s.Name = name
	return s
}

func (s *ProjectSpec) SecretRef(name string) ISpec {
	s.ConnectionSecret = &v1.ResourceRef{Name: name}
	return s
}

func (s *ProjectSpec) WithIpAccess(ipAdress, comment string) ISpec {
	access := project.NewIPAccessList().
		WithIP(ipAdress).
		WithComment(comment)
	s.ProjectIPAccessList = append(s.ProjectIPAccessList, access)
	return s
}

func (s ProjectSpec) CompleteK8sConfig(k8sname string) []byte {
	var t ap
	t.TypeMeta = metav1.TypeMeta{
		APIVersion: "atlas.mongodb.com/v1",
		Kind:       "AtlasProject",
	}
	t.ObjectMeta = &metav1.ObjectMeta{
		Name: k8sname,
	}
	t.Spec = s
	yamlConf := utils.JSONToYAMLConvert(t)
	return yamlConf
}
