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

package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

type ProjectSpec akov2.AtlasProjectSpec

type AProject struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec            ProjectSpec        `json:"spec"`
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
	p.Spec.ConnectionSecret = &common.ResourceRefNamespaced{Name: name, Namespace: p.ObjectMeta.Namespace}
	return p
}

func (p *AProject) WithSecretRefNamespaced(name, namespace string) *AProject {
	p.Spec.ConnectionSecret = &common.ResourceRefNamespaced{Name: name, Namespace: namespace}
	return p
}

func (p *AProject) WithIpAccess(cidrBlock, comment string) *AProject {
	access := project.NewIPAccessList().
		WithCIDR(cidrBlock).
		WithComment(comment)
	p.Spec.ProjectIPAccessList = append(p.Spec.ProjectIPAccessList, access)
	return p
}

func (p *AProject) WithPrivateLink(provider provider.ProviderName, region string) *AProject {
	link := akov2.PrivateEndpoint{
		Provider: provider,
		Region:   region,
	}
	p.Spec.PrivateEndpoints = append(p.Spec.PrivateEndpoints, link)
	return p
}

func (p *AProject) WithNetworkPeer(peer akov2.NetworkPeer) *AProject {
	p.Spec.NetworkPeers = append(p.Spec.NetworkPeers, peer)
	return p
}

func (p *AProject) WithEncryptionAtRest(spec *akov2.EncryptionAtRest) *AProject {
	p.Spec.EncryptionAtRest = spec
	return p
}

func (p *AProject) WithCloudProviderIntegration(role akov2.CloudProviderIntegration) *AProject {
	p.Spec.CloudProviderIntegrations = append(p.Spec.CloudProviderIntegrations, role)
	return p
}

func (p *AProject) WithIntegration(spec ProjectIntegration) *AProject {
	p.Spec.Integrations = append(p.Spec.Integrations, project.Integration(spec))
	return p
}

func (p *AProject) WithX509(certRef *common.ResourceRefNamespaced) *AProject {
	p.Spec.X509CertRef = certRef
	return p
}

func (p *AProject) WithAuditing(auditing *akov2.Auditing) *AProject {
	p.Spec.Auditing = auditing
	return p
}

func (p *AProject) UpdatePrivateLinkByOrder(i int, id string) *AProject {
	p.Spec.PrivateEndpoints[i].ID = id
	return p
}

func (p *AProject) UpdatePrivateLinkID(test akov2.PrivateEndpoint) *AProject {
	for i, peItem := range p.Spec.PrivateEndpoints {
		if (peItem.Provider == test.Provider) && (peItem.Region == test.Region) {
			p.Spec.PrivateEndpoints[i] = test
		}
	}
	return p
}

func (p *AProject) GetPrivateIDByProviderRegion(statusItem status.ProjectPrivateEndpoint) string {
	if statusItem.Provider == provider.ProviderAWS {
		return statusItem.InterfaceEndpointID
	}
	return statusItem.ID
}

func (p *AProject) DeletePrivateLink(id string) *AProject {
	var peList []akov2.PrivateEndpoint
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
