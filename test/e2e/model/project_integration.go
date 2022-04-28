package model

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
)

type ProjectIntegration project.Integration

func NewPIntegration(iType string) *ProjectIntegration {
	return &ProjectIntegration{
		Type: iType,
	}
}

func (i *ProjectIntegration) WithLicenseKeyRef(name, ns string) *ProjectIntegration {
	i.LicenseKeyRef.Name = name
	i.LicenseKeyRef.Namespace = ns
	return i
}

func (i *ProjectIntegration) WithAccountID(id string) *ProjectIntegration {
	i.AccountID = id
	return i
}

func (i *ProjectIntegration) WithWriteTokenRef(name, ns string) *ProjectIntegration {
	i.WriteTokenRef.Name = name
	i.WriteTokenRef.Namespace = ns
	return i
}

func (i *ProjectIntegration) WithReadTokenRef(name, ns string) *ProjectIntegration {
	i.ReadTokenRef.Name = name
	i.ReadTokenRef.Namespace = ns
	return i
}

func (i *ProjectIntegration) WithAPIKeyRef(name, ns string) *ProjectIntegration {
	i.APIKeyRef.Name = name
	i.APIKeyRef.Namespace = ns
	return i
}

func (i *ProjectIntegration) WithRegion(region string) *ProjectIntegration {
	i.Region = region
	return i
}

func (i *ProjectIntegration) WithServiceKeyRef(name, ns string) *ProjectIntegration {
	i.ServiceKeyRef.Name = name
	i.ServiceKeyRef.Namespace = ns
	return i
}

func (i *ProjectIntegration) WithAPITokenRef(name, ns string) *ProjectIntegration {
	i.APITokenRef.Name = name
	i.APITokenRef.Namespace = ns
	return i
}

func (i *ProjectIntegration) WithTeamName(t string) *ProjectIntegration {
	i.TeamName = t
	return i
}

func (i *ProjectIntegration) WithChannelName(c string) *ProjectIntegration {
	i.ChannelName = c
	return i
}

func (i *ProjectIntegration) WithRoutingKeyRef(name, ns string) *ProjectIntegration {
	i.RoutingKeyRef.Name = name
	i.RoutingKeyRef.Namespace = ns
	return i
}

func (i *ProjectIntegration) WithFlowName(f string) *ProjectIntegration {
	i.FlowName = f
	return i
}

func (i *ProjectIntegration) WithOrgName(o string) *ProjectIntegration {
	i.OrgName = o
	return i
}

func (i *ProjectIntegration) WithURL(url string) *ProjectIntegration {
	i.URL = url
	return i
}

func (i *ProjectIntegration) WithSecretRef(name, ns string) *ProjectIntegration {
	i.SecretRef.Name = name
	i.SecretRef.Namespace = ns
	return i
}
