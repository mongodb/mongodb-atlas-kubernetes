// Copyright 2024 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package privateendpoint

import (
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

const (
	ProviderAWS   = "AWS"
	ProviderAzure = "AZURE"
	ProviderGCP   = "GCP"

	StatusInitiating        = "INITIATING"
	StatusPending           = "PENDING"
	StatusPendingAcceptance = "PENDING_ACCEPTANCE"
	StatusWaitingForUser    = "WAITING_FOR_USER"
	StatusVerified          = "VERIFIED"
	StatusFailed            = "FAILED"
	StatusRejected          = "REJECTED"
	StatusDeleting          = "DELETING"
)

type EndpointService interface {
	ServiceID() string
	EndpointInterfaces() EndpointInterfaces
	Provider() string
	Region() string
	Status() string
	ErrorMessage() string
}

type CommonEndpointService struct {
	ID            string
	CloudRegion   string
	ServiceStatus string
	Error         string
	Interfaces    EndpointInterfaces
}

func (s *CommonEndpointService) ServiceID() string {
	return s.ID
}

func (s *CommonEndpointService) EndpointInterfaces() EndpointInterfaces {
	return s.Interfaces
}

func (s *CommonEndpointService) Region() string {
	return s.CloudRegion
}

func (s *CommonEndpointService) Status() string {
	return s.ServiceStatus
}

func (s *CommonEndpointService) ErrorMessage() string {
	return s.Error
}

type AWSService struct {
	CommonEndpointService

	ServiceName string
}

func (s *AWSService) Provider() string {
	return ProviderAWS
}

type AzureService struct {
	CommonEndpointService
	ServiceName string
	ResourceID  string
}

func (s *AzureService) Provider() string {
	return ProviderAzure
}

type GCPService struct {
	CommonEndpointService
	AttachmentNames []string
}

func (s *GCPService) Provider() string {
	return ProviderGCP
}

type EndpointInterface interface {
	InterfaceID() string
	Status() string
	ErrorMessage() string
}

type CommonEndpointInterface struct {
	ID              string
	InterfaceStatus string
	Error           string
}

func (i *CommonEndpointInterface) InterfaceID() string {
	return i.ID
}

func (i *CommonEndpointInterface) Status() string {
	return i.InterfaceStatus
}

func (i *CommonEndpointInterface) ErrorMessage() string {
	return i.Error
}

type AWSInterface struct {
	CommonEndpointInterface
}

type AzureInterface struct {
	CommonEndpointInterface
	IP             string
	ConnectionName string
}

type GCPInterface struct {
	CommonEndpointInterface
	Endpoints []GCPInterfaceEndpoint
}

type GCPInterfaceEndpoint struct {
	Name   string
	IP     string
	Status string
}

type EndpointInterfaces []EndpointInterface

func (ei EndpointInterfaces) Get(ID string) EndpointInterface {
	if ei == nil {
		return nil
	}

	for _, i := range ei {
		if i.InterfaceID() == ID {
			return i
		}
	}

	return nil
}

func NewPrivateEndpoint(akoPrivateEndpoint *akov2.AtlasPrivateEndpoint) EndpointService {
	switch akoPrivateEndpoint.Spec.Provider {
	case ProviderAWS:
		return &AWSService{
			CommonEndpointService: CommonEndpointService{
				ID:            akoPrivateEndpoint.Status.ServiceID,
				CloudRegion:   akoPrivateEndpoint.Spec.Region,
				ServiceStatus: akoPrivateEndpoint.Status.ServiceStatus,
				Error:         akoPrivateEndpoint.Status.Error,
				Interfaces:    newPrivateEndpointInterface(akoPrivateEndpoint),
			},
			ServiceName: akoPrivateEndpoint.Status.ServiceName,
		}
	case ProviderAzure:
		return &AzureService{
			CommonEndpointService: CommonEndpointService{
				ID:            akoPrivateEndpoint.Status.ServiceID,
				CloudRegion:   akoPrivateEndpoint.Spec.Region,
				ServiceStatus: akoPrivateEndpoint.Status.ServiceStatus,
				Error:         akoPrivateEndpoint.Status.Error,
				Interfaces:    newPrivateEndpointInterface(akoPrivateEndpoint),
			},
			ServiceName: akoPrivateEndpoint.Status.ServiceName,
			ResourceID:  akoPrivateEndpoint.Status.ResourceID,
		}
	case ProviderGCP:
		return &GCPService{
			CommonEndpointService: CommonEndpointService{
				ID:            akoPrivateEndpoint.Status.ServiceID,
				CloudRegion:   akoPrivateEndpoint.Spec.Region,
				ServiceStatus: akoPrivateEndpoint.Status.ServiceStatus,
				Error:         akoPrivateEndpoint.Status.Error,
				Interfaces:    newPrivateEndpointInterface(akoPrivateEndpoint),
			},
			AttachmentNames: akoPrivateEndpoint.Status.ServiceAttachmentNames,
		}
	}

	return nil
}

func NewPrivateEndpointStatus(peService EndpointService) status.AtlasPrivateEndpointStatusOption {
	return func(s *status.AtlasPrivateEndpointStatus) {
		endpoints := make([]status.EndpointInterfaceStatus, 0, len(peService.EndpointInterfaces()))
		for _, i := range peService.EndpointInterfaces() {
			connName := ""
			if azureInterface, ok := i.(*AzureInterface); ok {
				connName = azureInterface.ConnectionName
			}

			var gcpForwardRules []status.GCPForwardingRule
			if gcpInterface, ok := i.(*GCPInterface); ok {
				for _, fr := range gcpInterface.Endpoints {
					gcpForwardRules = append(
						gcpForwardRules,
						status.GCPForwardingRule{
							Name:   fr.Name,
							Status: fr.Status,
						},
					)
				}
			}

			endpoints = append(
				endpoints,
				status.EndpointInterfaceStatus{
					ID:                 i.InterfaceID(),
					ConnectionName:     connName,
					GCPForwardingRules: gcpForwardRules,
					Status:             i.Status(),
					Error:              i.ErrorMessage(),
				},
			)
		}

		s.ServiceID = peService.ServiceID()
		s.Endpoints = endpoints
		s.ServiceStatus = peService.Status()
		s.Error = peService.ErrorMessage()

		switch pe := peService.(type) {
		case *AWSService:
			s.ServiceName = pe.ServiceName
		case *AzureService:
			s.ServiceName = pe.ServiceName
			s.ResourceID = pe.ResourceID
		case *GCPService:
			s.ServiceAttachmentNames = pe.AttachmentNames
		}
	}
}

func newPrivateEndpointInterface(akoPrivateEndpoint *akov2.AtlasPrivateEndpoint) EndpointInterfaces {
	endpoints := EndpointInterfaces{}

	switch akoPrivateEndpoint.Spec.Provider {
	case ProviderAWS:
		for _, endpoint := range akoPrivateEndpoint.Spec.AWSConfiguration {
			ep := &AWSInterface{
				CommonEndpointInterface: CommonEndpointInterface{
					ID: endpoint.ID,
				},
			}

			for _, epStatus := range akoPrivateEndpoint.Status.Endpoints {
				if epStatus.ID == endpoint.ID {
					ep.InterfaceStatus = epStatus.Status
					ep.Error = epStatus.Error
				}
			}

			endpoints = append(endpoints, ep)
		}
	case ProviderAzure:
		for _, endpoint := range akoPrivateEndpoint.Spec.AzureConfiguration {
			ep := &AzureInterface{
				CommonEndpointInterface: CommonEndpointInterface{
					ID: endpoint.ID,
				},
				IP: endpoint.IP,
			}

			for _, epStatus := range akoPrivateEndpoint.Status.Endpoints {
				if epStatus.ID == endpoint.ID {
					ep.InterfaceStatus = epStatus.Status
					ep.Error = epStatus.Error
					ep.ConnectionName = epStatus.ConnectionName
				}
			}

			endpoints = append(endpoints, ep)
		}
	case ProviderGCP:
		for _, endpoint := range akoPrivateEndpoint.Spec.GCPConfiguration {
			gcpEPs := make([]GCPInterfaceEndpoint, 0, len(endpoint.Endpoints))
			for _, gcpEP := range endpoint.Endpoints {
				gcpEPs = append(
					gcpEPs,
					GCPInterfaceEndpoint{
						Name: gcpEP.Name,
						IP:   gcpEP.IP,
					},
				)
			}

			ep := &GCPInterface{
				CommonEndpointInterface: CommonEndpointInterface{
					ID: endpoint.GroupName,
				},
				Endpoints: gcpEPs,
			}

			for _, epStatus := range akoPrivateEndpoint.Status.Endpoints {
				if epStatus.ID == endpoint.GroupName {
					ep.InterfaceStatus = epStatus.Status
					ep.Error = epStatus.Error

					for i, gcpEP := range ep.Endpoints {
						for _, gcpEPStatus := range epStatus.GCPForwardingRules {
							if gcpEP.Name == gcpEPStatus.Name {
								ep.Endpoints[i].Status = gcpEPStatus.Status
							}
						}
					}
				}
			}

			endpoints = append(endpoints, ep)
		}
	}

	cmp.NormalizeSlice(endpoints, func(a, b EndpointInterface) int {
		return strings.Compare(a.InterfaceID(), b.InterfaceID())
	})

	return endpoints
}

func serviceFromAtlas(peService *admin.EndpointService, endpoints EndpointInterfaces) EndpointService {
	switch peService.GetCloudProvider() {
	case ProviderAWS:
		return &AWSService{
			CommonEndpointService: CommonEndpointService{
				ID:            peService.GetId(),
				CloudRegion:   peService.GetRegionName(),
				ServiceStatus: peService.GetStatus(),
				Error:         peService.GetErrorMessage(),
				Interfaces:    endpoints,
			},
			ServiceName: peService.GetEndpointServiceName(),
		}
	case ProviderAzure:
		return &AzureService{
			CommonEndpointService: CommonEndpointService{
				ID:            peService.GetId(),
				CloudRegion:   peService.GetRegionName(),
				ServiceStatus: peService.GetStatus(),
				Error:         peService.GetErrorMessage(),
				Interfaces:    endpoints,
			},
			ServiceName: peService.GetPrivateLinkServiceName(),
			ResourceID:  peService.GetPrivateLinkServiceResourceId(),
		}
	case ProviderGCP:
		return &GCPService{
			CommonEndpointService: CommonEndpointService{
				ID:            peService.GetId(),
				CloudRegion:   peService.GetRegionName(),
				ServiceStatus: peService.GetStatus(),
				Error:         peService.GetErrorMessage(),
				Interfaces:    endpoints,
			},
			AttachmentNames: peService.GetServiceAttachmentNames(),
		}
	}

	return nil
}

func serviceCreateToAtlas(peService EndpointService) *admin.CloudProviderEndpointServiceRequest {
	return &admin.CloudProviderEndpointServiceRequest{
		ProviderName: peService.Provider(),
		Region:       peService.Region(),
	}
}

func interfaceFromAtlas(peInterface *admin.PrivateLinkEndpoint) EndpointInterface {
	switch peInterface.GetCloudProvider() {
	case ProviderAWS:
		return &AWSInterface{
			CommonEndpointInterface: CommonEndpointInterface{
				ID:              peInterface.GetInterfaceEndpointId(),
				InterfaceStatus: peInterface.GetConnectionStatus(),
				Error:           peInterface.GetErrorMessage(),
			},
		}
	case ProviderAzure:
		return &AzureInterface{
			CommonEndpointInterface: CommonEndpointInterface{
				ID:              peInterface.GetPrivateEndpointResourceId(),
				InterfaceStatus: peInterface.GetStatus(),
				Error:           peInterface.GetErrorMessage(),
			},
			IP:             peInterface.GetPrivateEndpointIPAddress(),
			ConnectionName: peInterface.GetPrivateEndpointConnectionName(),
		}
	case ProviderGCP:
		endpoints := make([]GCPInterfaceEndpoint, 0, len(peInterface.GetEndpoints()))
		for _, ep := range peInterface.GetEndpoints() {
			endpoints = append(
				endpoints,
				GCPInterfaceEndpoint{
					Name:   ep.GetEndpointName(),
					IP:     ep.GetIpAddress(),
					Status: ep.GetStatus(),
				},
			)
		}

		return &GCPInterface{
			CommonEndpointInterface: CommonEndpointInterface{
				ID:              peInterface.GetEndpointGroupName(),
				InterfaceStatus: peInterface.GetStatus(),
				Error:           peInterface.GetErrorMessage(),
			},
			Endpoints: endpoints,
		}
	}

	return nil
}

func interfaceCreateToAtlas(peInterface EndpointInterface, gcpProjectID string) *admin.CreateEndpointRequest {
	switch i := peInterface.(type) {
	case *AWSInterface:
		return &admin.CreateEndpointRequest{
			Id: pointer.MakePtr(i.InterfaceID()),
		}
	case *AzureInterface:
		return &admin.CreateEndpointRequest{
			Id:                       pointer.MakePtr(i.InterfaceID()),
			PrivateEndpointIPAddress: pointer.MakePtr(i.IP),
		}
	case *GCPInterface:
		gcpEPs := make([]admin.CreateGCPForwardingRuleRequest, 0, len(i.Endpoints))
		for _, ep := range i.Endpoints {
			gcpEPs = append(
				gcpEPs,
				admin.CreateGCPForwardingRuleRequest{
					EndpointName: pointer.MakePtr(ep.Name),
					IpAddress:    pointer.MakePtr(ep.IP),
				},
			)
		}

		return &admin.CreateEndpointRequest{
			GcpProjectId:      pointer.MakePtr(gcpProjectID),
			EndpointGroupName: pointer.MakePtr(i.InterfaceID()),
			Endpoints:         &gcpEPs,
		}
	}

	return nil
}

type CompositeEndpointInterface struct {
	AKO   EndpointInterface
	Atlas EndpointInterface
}

func MapPrivateEndpoints(akoInterfaces, atlasInterfaces []EndpointInterface) map[string]CompositeEndpointInterface {
	m := map[string]CompositeEndpointInterface{}

	for _, akoInterface := range akoInterfaces {
		m[akoInterface.InterfaceID()] = CompositeEndpointInterface{
			AKO: akoInterface,
		}
	}

	for _, atlasInterface := range atlasInterfaces {
		i := CompositeEndpointInterface{}
		if existing, ok := m[atlasInterface.InterfaceID()]; ok {
			i = existing
		}

		i.Atlas = atlasInterface
		m[atlasInterface.InterfaceID()] = i
	}

	return m
}
