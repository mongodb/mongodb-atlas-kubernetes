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

package atlasproject

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312012/admin"
	"go.mongodb.org/atlas-sdk/v20250312012/mockadmin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	atlasmocks "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/encryptionatrest"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/teams"
)

func TestHandleProject(t *testing.T) {
	deletionTime := metav1.Now()

	tests := map[string]struct {
		atlasSDKMocker         func() *admin.APIClient
		projectServiceMocker   func() project.ProjectService
		teamServiceMocker      func() teams.TeamsService
		encryptionAtRestMocker func() encryptionatrest.EncryptionAtRestService
		interceptors           interceptor.Funcs
		project                *akov2.AtlasProject
		result                 reconcile.Result
		conditions             []api.Condition
		finalizers             []string
		wantErr                bool
	}{
		"should fail to get project from atlas": {
			wantErr: true,
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(nil, errors.New("failed to get project"))

				return service
			},
			teamServiceMocker: func() teams.TeamsService {
				return nil
			},
			encryptionAtRestMocker: func() encryptionatrest.EncryptionAtRestService {
				return nil
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
			},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.ProjectNotCreatedInAtlas)).
					WithMessageRegexp("failed to get project"),
			},
		},
		"should create project": {
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(nil, nil)
				service.EXPECT().CreateProject(context.Background(), mock.AnythingOfType("*project.Project")).
					Return(nil)

				return service
			},
			teamServiceMocker: func() teams.TeamsService {
				return nil
			},
			encryptionAtRestMocker: func() encryptionatrest.EncryptionAtRestService {
				return nil
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
			},
			result: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.ProjectBeingConfiguredInAtlas)).
					WithMessageRegexp("configuring project in Atlas"),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
		"should delete project": {
			atlasSDKMocker: func() *admin.APIClient {
				mockPrivateEndpointAPI := mockadmin.NewPrivateEndpointServicesApi(t)
				mockPrivateEndpointAPI.EXPECT().
					ListPrivateEndpointService(mock.Anything, mock.Anything, mock.Anything).
					Return(admin.ListPrivateEndpointServiceApiRequest{ApiService: mockPrivateEndpointAPI})
				mockPrivateEndpointAPI.EXPECT().
					ListPrivateEndpointServiceExecute(admin.ListPrivateEndpointServiceApiRequest{ApiService: mockPrivateEndpointAPI}).
					Return([]admin.EndpointService{}, nil, nil)

				return &admin.APIClient{
					PrivateEndpointServicesApi: mockPrivateEndpointAPI,
				}
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(&project.Project{ID: "projectID", Name: "my-project"}, nil)
				service.EXPECT().DeleteProject(context.Background(), mock.AnythingOfType("*project.Project")).
					Return(nil)

				return service
			},
			teamServiceMocker: func() teams.TeamsService {
				service := translation.NewTeamsServiceMock(t)
				service.EXPECT().ListProjectTeams(context.Background(), mock.Anything).Return([]teams.AssignedTeam{}, nil)
				return service
			},
			encryptionAtRestMocker: func() encryptionatrest.EncryptionAtRestService {
				return nil
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "my-project",
					Namespace:         "default",
					DeletionTimestamp: &deletionTime,
					Finalizers:        []string{customresource.FinalizerLabel},
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
			},
			result: reconcile.Result{},
		},
		"should delete project when it's was already deleted in atlas": {
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(nil, nil)

				return service
			},
			teamServiceMocker: func() teams.TeamsService {
				return nil
			},
			encryptionAtRestMocker: func() encryptionatrest.EncryptionAtRestService {
				return nil
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "my-project",
					Namespace:         "default",
					DeletionTimestamp: &deletionTime,
					Finalizers:        []string{customresource.FinalizerLabel},
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
			},
			result: reconcile.Result{},
		},
		"should fail to remove finalizer from project when it's was already deleted in atlas": {
			wantErr: true,
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(nil, nil)

				return service
			},
			teamServiceMocker: func() teams.TeamsService {
				return nil
			},
			encryptionAtRestMocker: func() encryptionatrest.EncryptionAtRestService {
				return nil
			},
			interceptors: interceptor.Funcs{
				Patch: func(ctx context.Context, client client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
					return errors.New("failed to remove finalizer")
				},
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "my-project",
					Namespace:         "default",
					DeletionTimestamp: &deletionTime,
					Finalizers:        []string{customresource.FinalizerLabel},
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
			},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.AtlasFinalizerNotRemoved)).
					WithMessageRegexp("failed to remove finalizer"),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
		"should fail to configure authentication modes": {
			wantErr: true,
			atlasSDKMocker: func() *admin.APIClient { //nolint:dupl
				ipAccessList := mockadmin.NewProjectIPAccessListApi(t)
				ipAccessList.EXPECT().ListAccessListEntries(context.Background(), "projectID").
					Return(admin.ListAccessListEntriesApiRequest{ApiService: ipAccessList})
				ipAccessList.EXPECT().ListAccessListEntriesExecute(mock.Anything).
					Return(nil, nil, nil)
				privateEndpoints := mockadmin.NewPrivateEndpointServicesApi(t)
				privateEndpoints.EXPECT().ListPrivateEndpointService(context.Background(), "projectID", mock.Anything).
					Return(admin.ListPrivateEndpointServiceApiRequest{ApiService: privateEndpoints})
				privateEndpoints.EXPECT().ListPrivateEndpointServiceExecute(mock.AnythingOfType("admin.ListPrivateEndpointServiceApiRequest")).
					Return(nil, nil, nil)
				networkPeering := mockadmin.NewNetworkPeeringApi(t)
				networkPeering.EXPECT().ListGroupPeersWithParams(context.Background(), mock.AnythingOfType("*admin.ListGroupPeersApiParams")).
					Return(admin.ListGroupPeersApiRequest{ApiService: networkPeering})
				networkPeering.EXPECT().ListGroupPeersExecute(mock.AnythingOfType("admin.ListGroupPeersApiRequest")).
					Return(nil, nil, nil)
				audit := mockadmin.NewAuditingApi(t)
				audit.EXPECT().GetGroupAuditLog(context.Background(), "projectID").
					Return(admin.GetGroupAuditLogApiRequest{ApiService: audit})
				audit.EXPECT().GetGroupAuditLogExecute(mock.AnythingOfType("admin.GetGroupAuditLogApiRequest")).
					Return(nil, nil, nil)
				customRoles := mockadmin.NewCustomDatabaseRolesApi(t)
				customRoles.EXPECT().ListCustomDbRoles(context.Background(), "projectID").
					Return(admin.ListCustomDbRolesApiRequest{ApiService: customRoles})
				customRoles.EXPECT().ListCustomDbRolesExecute(mock.AnythingOfType("admin.ListCustomDbRolesApiRequest")).
					Return(nil, nil, nil)
				projectAPI := mockadmin.NewProjectsApi(t)
				projectAPI.EXPECT().GetGroupSettings(context.Background(), "projectID").
					Return(admin.GetGroupSettingsApiRequest{ApiService: projectAPI})
				projectAPI.EXPECT().GetGroupSettingsExecute(mock.AnythingOfType("admin.GetGroupSettingsApiRequest")).
					Return(admin.NewGroupSettings(), nil, nil)
				backup := mockadmin.NewCloudBackupsApi(t)
				backup.EXPECT().GetCompliancePolicy(context.Background(), "projectID").
					Return(admin.GetCompliancePolicyApiRequest{ApiService: backup})
				backup.EXPECT().GetCompliancePolicyExecute(mock.AnythingOfType("admin.GetCompliancePolicyApiRequest")).
					Return(nil, nil, nil)
				integrationsApi := mockadmin.NewThirdPartyIntegrationsApi(t)
				integrationsApi.EXPECT().ListGroupIntegrations(context.Background(), "projectID").
					Return(admin.ListGroupIntegrationsApiRequest{ApiService: integrationsApi})
				integrationsApi.EXPECT().ListGroupIntegrationsExecute(mock.AnythingOfType("admin.ListGroupIntegrationsApiRequest")).
					Return(&admin.PaginatedIntegration{}, nil, nil)

				return &admin.APIClient{
					ProjectIPAccessListApi:     ipAccessList,
					PrivateEndpointServicesApi: privateEndpoints,
					NetworkPeeringApi:          networkPeering,
					AuditingApi:                audit,
					CustomDatabaseRolesApi:     customRoles,
					ProjectsApi:                projectAPI,
					CloudBackupsApi:            backup,
					ThirdPartyIntegrationsApi:  integrationsApi,
				}
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(&project.Project{ID: "projectID", Name: "my-project"}, nil)

				return service
			},
			teamServiceMocker: func() teams.TeamsService {
				service := translation.NewTeamsServiceMock(t)
				service.EXPECT().ListProjectTeams(context.Background(), mock.Anything).Return([]teams.AssignedTeam{}, nil)
				return service
			},
			encryptionAtRestMocker: func() encryptionatrest.EncryptionAtRestService {
				service := translation.NewEncryptionAtRestServiceMock(t)
				service.EXPECT().Get(context.Background(), mock.Anything).Return(nil, nil)
				return service
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "my-project",
					Namespace:  "default",
					Finalizers: []string{customresource.FinalizerLabel},
				},
				Spec: akov2.AtlasProjectSpec{
					Name:        "my-project",
					X509CertRef: &common.ResourceRefNamespaced{Name: "invalid-ref"},
				},
				Status: status.AtlasProjectStatus{
					ID: "projectID",
				},
			},
			conditions: []api.Condition{
				api.TrueCondition(api.ProjectReadyType),
				api.FalseCondition(api.X509AuthReadyType).
					WithMessageRegexp("secrets \"invalid-ref\" not found"),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
		"should configure project resources": {
			atlasSDKMocker: func() *admin.APIClient { //nolint:dupl
				ipAccessList := mockadmin.NewProjectIPAccessListApi(t)
				ipAccessList.EXPECT().ListAccessListEntries(context.Background(), "projectID").
					Return(admin.ListAccessListEntriesApiRequest{ApiService: ipAccessList})
				ipAccessList.EXPECT().ListAccessListEntriesExecute(mock.AnythingOfType("admin.ListAccessListEntriesApiRequest")).
					Return(nil, nil, nil)
				privateEndpoints := mockadmin.NewPrivateEndpointServicesApi(t)
				privateEndpoints.EXPECT().ListPrivateEndpointService(context.Background(), "projectID", mock.Anything).
					Return(admin.ListPrivateEndpointServiceApiRequest{ApiService: privateEndpoints})
				privateEndpoints.EXPECT().ListPrivateEndpointServiceExecute(mock.AnythingOfType("admin.ListPrivateEndpointServiceApiRequest")).
					Return(nil, nil, nil)
				networkPeering := mockadmin.NewNetworkPeeringApi(t)
				networkPeering.EXPECT().ListGroupPeersWithParams(context.Background(), mock.Anything).
					Return(admin.ListGroupPeersApiRequest{ApiService: networkPeering})
				networkPeering.EXPECT().ListGroupPeersExecute(mock.Anything).
					Return(nil, nil, nil)
				networkPeering.EXPECT().ListGroupContainerAll(context.Background(), "projectID").
					Return(admin.ListGroupContainerAllApiRequest{ApiService: networkPeering})
				networkPeering.EXPECT().ListGroupContainerAllExecute(mock.Anything).
					Return(nil, nil, nil)
				audit := mockadmin.NewAuditingApi(t)
				audit.EXPECT().GetGroupAuditLog(context.Background(), "projectID").
					Return(admin.GetGroupAuditLogApiRequest{ApiService: audit})
				audit.EXPECT().GetGroupAuditLogExecute(mock.AnythingOfType("admin.GetGroupAuditLogApiRequest")).
					Return(nil, nil, nil)
				customRoles := mockadmin.NewCustomDatabaseRolesApi(t)
				customRoles.EXPECT().ListCustomDbRoles(context.Background(), "projectID").
					Return(admin.ListCustomDbRolesApiRequest{ApiService: customRoles})
				customRoles.EXPECT().ListCustomDbRolesExecute(mock.AnythingOfType("admin.ListCustomDbRolesApiRequest")).
					Return(nil, nil, nil)
				projectAPI := mockadmin.NewProjectsApi(t)
				projectAPI.EXPECT().GetGroupSettings(context.Background(), "projectID").
					Return(admin.GetGroupSettingsApiRequest{ApiService: projectAPI})
				projectAPI.EXPECT().GetGroupSettingsExecute(mock.AnythingOfType("admin.GetGroupSettingsApiRequest")).
					Return(admin.NewGroupSettings(), nil, nil)
				backup := mockadmin.NewCloudBackupsApi(t)
				backup.EXPECT().GetCompliancePolicy(context.Background(), "projectID").
					Return(admin.GetCompliancePolicyApiRequest{ApiService: backup})
				backup.EXPECT().GetCompliancePolicyExecute(mock.AnythingOfType("admin.GetCompliancePolicyApiRequest")).
					Return(nil, nil, nil)
				integrationsApi := mockadmin.NewThirdPartyIntegrationsApi(t)
				integrationsApi.EXPECT().ListGroupIntegrations(context.Background(), "projectID").
					Return(admin.ListGroupIntegrationsApiRequest{ApiService: integrationsApi})
				integrationsApi.EXPECT().ListGroupIntegrationsExecute(mock.AnythingOfType("admin.ListGroupIntegrationsApiRequest")).
					Return(&admin.PaginatedIntegration{}, nil, nil)

				return &admin.APIClient{
					ProjectIPAccessListApi:     ipAccessList,
					PrivateEndpointServicesApi: privateEndpoints,
					NetworkPeeringApi:          networkPeering,
					AuditingApi:                audit,
					CustomDatabaseRolesApi:     customRoles,
					ProjectsApi:                projectAPI,
					CloudBackupsApi:            backup,
					ThirdPartyIntegrationsApi:  integrationsApi,
				}
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(&project.Project{ID: "projectID", Name: "my-project"}, nil)

				return service
			},
			teamServiceMocker: func() teams.TeamsService {
				service := translation.NewTeamsServiceMock(t)
				service.EXPECT().ListProjectTeams(context.Background(), mock.Anything).Return([]teams.AssignedTeam{}, nil)
				return service
			},
			encryptionAtRestMocker: func() encryptionatrest.EncryptionAtRestService {
				service := translation.NewEncryptionAtRestServiceMock(t)
				service.EXPECT().Get(context.Background(), mock.Anything).Return(nil, nil)
				return service
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "my-project",
					Namespace:  "default",
					Finalizers: []string{customresource.FinalizerLabel},
					Annotations: map[string]string{
						// no skipped config, but some fake applied condig on network peerings
						customresource.AnnotationLastAppliedConfiguration: func() string {
							d, _ := json.Marshal(&akov2.AtlasProjectSpec{
								NetworkPeers: []akov2.NetworkPeer{{}},
							})
							return string(d)
						}(),
					},
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
				Status: status.AtlasProjectStatus{
					ID: "projectID",
				},
			},
			result: reconcile.Result{},
			conditions: []api.Condition{
				api.TrueCondition(api.ProjectReadyType),
				api.TrueCondition(api.ReadyType),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
		"should fail to configure project resources": {
			wantErr: true,
			atlasSDKMocker: func() *admin.APIClient {
				ipAccessList := mockadmin.NewProjectIPAccessListApi(t)
				ipAccessList.EXPECT().ListAccessListEntries(context.Background(), "projectID").
					Return(admin.ListAccessListEntriesApiRequest{ApiService: ipAccessList})
				ipAccessList.EXPECT().ListAccessListEntriesExecute(mock.AnythingOfType("admin.ListAccessListEntriesApiRequest")).
					Return(nil, nil, errors.New("failed to list IP Access List"))
				privateEndpoints := mockadmin.NewPrivateEndpointServicesApi(t)
				privateEndpoints.EXPECT().ListPrivateEndpointService(context.Background(), "projectID", mock.Anything).
					Return(admin.ListPrivateEndpointServiceApiRequest{ApiService: privateEndpoints})
				privateEndpoints.EXPECT().ListPrivateEndpointServiceExecute(mock.AnythingOfType("admin.ListPrivateEndpointServiceApiRequest")).
					Return(nil, nil, nil)
				networkPeering := mockadmin.NewNetworkPeeringApi(t)
				networkPeering.EXPECT().ListGroupPeersWithParams(context.Background(), mock.Anything).
					Return(admin.ListGroupPeersApiRequest{ApiService: networkPeering})
				networkPeering.EXPECT().ListGroupPeersExecute(mock.Anything).
					Return(nil, nil, nil)
				networkPeering.EXPECT().ListGroupContainerAll(context.Background(), "projectID").
					Return(admin.ListGroupContainerAllApiRequest{ApiService: networkPeering})
				networkPeering.EXPECT().ListGroupContainerAllExecute(mock.Anything).
					Return(nil, nil, nil)
				audit := mockadmin.NewAuditingApi(t)
				audit.EXPECT().GetGroupAuditLog(context.Background(), "projectID").
					Return(admin.GetGroupAuditLogApiRequest{ApiService: audit})
				audit.EXPECT().GetGroupAuditLogExecute(mock.AnythingOfType("admin.GetGroupAuditLogApiRequest")).
					Return(nil, nil, nil)
				customRoles := mockadmin.NewCustomDatabaseRolesApi(t)
				customRoles.EXPECT().ListCustomDbRoles(context.Background(), "projectID").
					Return(admin.ListCustomDbRolesApiRequest{ApiService: customRoles})
				customRoles.EXPECT().ListCustomDbRolesExecute(mock.AnythingOfType("admin.ListCustomDbRolesApiRequest")).
					Return(nil, nil, nil)
				projectAPI := mockadmin.NewProjectsApi(t)
				projectAPI.EXPECT().GetGroupSettings(context.Background(), "projectID").
					Return(admin.GetGroupSettingsApiRequest{ApiService: projectAPI})
				projectAPI.EXPECT().GetGroupSettingsExecute(mock.AnythingOfType("admin.GetGroupSettingsApiRequest")).
					Return(admin.NewGroupSettings(), nil, nil)
				backup := mockadmin.NewCloudBackupsApi(t)
				backup.EXPECT().GetCompliancePolicy(context.Background(), "projectID").
					Return(admin.GetCompliancePolicyApiRequest{ApiService: backup})
				backup.EXPECT().GetCompliancePolicyExecute(mock.AnythingOfType("admin.GetCompliancePolicyApiRequest")).
					Return(nil, nil, nil)
				integrationsApi := mockadmin.NewThirdPartyIntegrationsApi(t)
				integrationsApi.EXPECT().ListGroupIntegrations(context.Background(), "projectID").
					Return(admin.ListGroupIntegrationsApiRequest{ApiService: integrationsApi})
				integrationsApi.EXPECT().ListGroupIntegrationsExecute(mock.AnythingOfType("admin.ListGroupIntegrationsApiRequest")).
					Return(&admin.PaginatedIntegration{}, nil, nil)

				return &admin.APIClient{
					ProjectIPAccessListApi:     ipAccessList,
					PrivateEndpointServicesApi: privateEndpoints,
					NetworkPeeringApi:          networkPeering,
					AuditingApi:                audit,
					CustomDatabaseRolesApi:     customRoles,
					ProjectsApi:                projectAPI,
					CloudBackupsApi:            backup,
					ThirdPartyIntegrationsApi:  integrationsApi,
				}
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(&project.Project{ID: "projectID", Name: "my-project"}, nil)

				return service
			},
			teamServiceMocker: func() teams.TeamsService {
				service := translation.NewTeamsServiceMock(t)
				service.EXPECT().ListProjectTeams(context.Background(), mock.Anything).Return([]teams.AssignedTeam{}, nil)
				return service
			},
			encryptionAtRestMocker: func() encryptionatrest.EncryptionAtRestService {
				service := translation.NewEncryptionAtRestServiceMock(t)
				service.EXPECT().Get(context.Background(), mock.Anything).Return(nil, nil)
				return service
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "my-project",
					Namespace:  "default",
					Finalizers: []string{customresource.FinalizerLabel},
					// no skipped config, but some fake applied condig on network peerings
					Annotations: map[string]string{
						customresource.AnnotationLastAppliedConfiguration: func() string {
							d, _ := json.Marshal(&akov2.AtlasProjectSpec{
								NetworkPeers: []akov2.NetworkPeer{{}},
							})
							return string(d)
						}(),
					},
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
				Status: status.AtlasProjectStatus{
					ID: "projectID",
				},
			},
			conditions: []api.Condition{
				api.TrueCondition(api.ProjectReadyType),
				api.FalseCondition(api.IPAccessListReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to get ip access list from Atlas: failed to list IP Access List"),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
		"should fail to save last applied config": {
			wantErr: true,
			atlasSDKMocker: func() *admin.APIClient { //nolint:dupl
				ipAccessList := mockadmin.NewProjectIPAccessListApi(t)
				ipAccessList.EXPECT().ListAccessListEntries(context.Background(), "projectID").
					Return(admin.ListAccessListEntriesApiRequest{ApiService: ipAccessList})
				ipAccessList.EXPECT().ListAccessListEntriesExecute(mock.AnythingOfType("admin.ListAccessListEntriesApiRequest")).
					Return(nil, nil, nil)
				privateEndpoints := mockadmin.NewPrivateEndpointServicesApi(t)
				privateEndpoints.EXPECT().ListPrivateEndpointService(context.Background(), "projectID", mock.Anything).
					Return(admin.ListPrivateEndpointServiceApiRequest{ApiService: privateEndpoints})
				privateEndpoints.EXPECT().ListPrivateEndpointServiceExecute(mock.AnythingOfType("admin.ListPrivateEndpointServiceApiRequest")).
					Return(nil, nil, nil)
				networkPeering := mockadmin.NewNetworkPeeringApi(t)
				networkPeering.EXPECT().ListGroupPeersWithParams(context.Background(), mock.Anything).
					Return(admin.ListGroupPeersApiRequest{ApiService: networkPeering})
				networkPeering.EXPECT().ListGroupPeersExecute(mock.AnythingOfType("admin.ListGroupPeersApiRequest")).
					Return(nil, nil, nil)
				networkPeering.EXPECT().ListGroupContainerAll(context.Background(), "projectID").
					Return(admin.ListGroupContainerAllApiRequest{ApiService: networkPeering})
				networkPeering.EXPECT().ListGroupContainerAllExecute(mock.Anything).
					Return(nil, nil, nil)
				audit := mockadmin.NewAuditingApi(t)
				audit.EXPECT().GetGroupAuditLog(context.Background(), "projectID").
					Return(admin.GetGroupAuditLogApiRequest{ApiService: audit})
				audit.EXPECT().GetGroupAuditLogExecute(mock.AnythingOfType("admin.GetGroupAuditLogApiRequest")).
					Return(nil, nil, nil)
				customRoles := mockadmin.NewCustomDatabaseRolesApi(t)
				customRoles.EXPECT().ListCustomDbRoles(context.Background(), "projectID").
					Return(admin.ListCustomDbRolesApiRequest{ApiService: customRoles})
				customRoles.EXPECT().ListCustomDbRolesExecute(mock.AnythingOfType("admin.ListCustomDbRolesApiRequest")).
					Return(nil, nil, nil)
				projectAPI := mockadmin.NewProjectsApi(t)
				projectAPI.EXPECT().GetGroupSettings(context.Background(), "projectID").
					Return(admin.GetGroupSettingsApiRequest{ApiService: projectAPI})
				projectAPI.EXPECT().GetGroupSettingsExecute(mock.AnythingOfType("admin.GetGroupSettingsApiRequest")).
					Return(admin.NewGroupSettings(), nil, nil)
				backup := mockadmin.NewCloudBackupsApi(t)
				backup.EXPECT().GetCompliancePolicy(context.Background(), "projectID").
					Return(admin.GetCompliancePolicyApiRequest{ApiService: backup})
				backup.EXPECT().GetCompliancePolicyExecute(mock.AnythingOfType("admin.GetCompliancePolicyApiRequest")).
					Return(nil, nil, nil)
				integrationsApi := mockadmin.NewThirdPartyIntegrationsApi(t)
				integrationsApi.EXPECT().ListGroupIntegrations(context.Background(), "projectID").
					Return(admin.ListGroupIntegrationsApiRequest{ApiService: integrationsApi})
				integrationsApi.EXPECT().ListGroupIntegrationsExecute(mock.AnythingOfType("admin.ListGroupIntegrationsApiRequest")).
					Return(&admin.PaginatedIntegration{}, nil, nil)

				return &admin.APIClient{
					ProjectIPAccessListApi:     ipAccessList,
					PrivateEndpointServicesApi: privateEndpoints,
					NetworkPeeringApi:          networkPeering,
					AuditingApi:                audit,
					CustomDatabaseRolesApi:     customRoles,
					ProjectsApi:                projectAPI,
					CloudBackupsApi:            backup,
					ThirdPartyIntegrationsApi:  integrationsApi,
				}
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(&project.Project{ID: "projectID", Name: "my-project"}, nil)

				return service
			},
			teamServiceMocker: func() teams.TeamsService {
				service := translation.NewTeamsServiceMock(t)
				service.EXPECT().ListProjectTeams(context.Background(), mock.Anything).Return([]teams.AssignedTeam{}, nil)
				return service
			},
			encryptionAtRestMocker: func() encryptionatrest.EncryptionAtRestService {
				service := translation.NewEncryptionAtRestServiceMock(t)
				service.EXPECT().Get(context.Background(), mock.Anything).Return(nil, nil)
				return service
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "my-project",
					Namespace:  "default",
					Finalizers: []string{customresource.FinalizerLabel},
					// no skipped config, but some fake applied condig on network peerings
					Annotations: map[string]string{
						customresource.AnnotationLastAppliedConfiguration: func() string {
							d, _ := json.Marshal(&akov2.AtlasProjectSpec{
								NetworkPeers: []akov2.NetworkPeer{{}},
							})
							return string(d)
						}(),
					},
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
				Status: status.AtlasProjectStatus{
					ID: "projectID",
				},
			},
			interceptors: interceptor.Funcs{
				Patch: func(ctx context.Context, client client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
					return errors.New("failed to save last applied config")
				},
			},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to save last applied config"),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(testScheme))
			require.NoError(t, corev1.AddToScheme(testScheme))
			logger := zaptest.NewLogger(t).Sugar()
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312012: tt.atlasSDKMocker(),
				},
			}
			instancesIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
			customRoleIndexer := indexer.NewAtlasCustomRoleByProjectIndexer(logger.Desugar())
			peIndexer := indexer.NewAtlasPrivateEndpointByProjectIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.project).
				WithStatusSubresource(tt.project).
				WithIndex(instancesIndexer.Object(), instancesIndexer.Name(), instancesIndexer.Keys).
				WithIndex(customRoleIndexer.Object(), customRoleIndexer.Name(), customRoleIndexer.Keys).
				WithIndex(peIndexer.Object(), peIndexer.Name(), peIndexer.Keys).
				WithInterceptorFuncs(tt.interceptors).
				Build()

			reconciler := &AtlasProjectReconciler{
				Client:        k8sClient,
				Log:           logger,
				EventRecorder: record.NewFakeRecorder(30),
			}
			services := &AtlasProjectServices{
				projectService:          tt.projectServiceMocker(),
				teamsService:            tt.teamServiceMocker(),
				encryptionAtRestService: tt.encryptionAtRestMocker(),
			}

			result, err := reconciler.handleProject(ctx, "my-org-id", tt.project, services)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.result, result)
			assert.True(
				t,
				cmp.Equal(
					tt.conditions,
					ctx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
			assert.Equal(t, tt.finalizers, tt.project.Finalizers)
		})
	}
}

func TestCreate(t *testing.T) {
	tests := map[string]struct {
		projectServiceMocker func() project.ProjectService
		interceptors         interceptor.Funcs
		project              *akov2.AtlasProject
		result               reconcile.Result
		conditions           []api.Condition
		finalizers           []string
		wantErr              bool
	}{
		"should fail to create project": {
			wantErr: true,
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().CreateProject(context.Background(), mock.AnythingOfType("*project.Project")).
					Return(errors.New("failed to create project"))

				return service
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
			},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.ProjectNotCreatedInAtlas)).
					WithMessageRegexp("failed to create project"),
			},
		},
		"should fail to add finalizer when creating a project": {
			wantErr: true,
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().CreateProject(context.Background(), mock.AnythingOfType("*project.Project")).
					Return(nil)

				return service
			},
			interceptors: interceptor.Funcs{Patch: func(ctx context.Context, client client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
				if patch.Type() == types.MergePatchType {
					return nil
				}

				return errors.New("failed to patch project with finalizers")
			}},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
			},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.AtlasFinalizerNotSet)).
					WithMessageRegexp("failed to patch project with finalizers"),
			},
		},
		"should fail to add last applied config when creating a project": {
			wantErr: true,
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().CreateProject(context.Background(), mock.AnythingOfType("*project.Project")).
					Return(nil)

				return service
			},
			interceptors: interceptor.Funcs{Patch: func(ctx context.Context, client client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
				if patch.Type() == types.MergePatchType {
					return errors.New("failed to patch project with last applied config")
				}

				return nil
			}},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
			},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to patch project with last applied config"),
			},
		},
		"should create a project": {
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().CreateProject(context.Background(), mock.AnythingOfType("*project.Project")).
					Return(nil)

				return service
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
			},
			result: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.ProjectBeingConfiguredInAtlas)).
					WithMessageRegexp("configuring project in Atlas"),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(testScheme))
			instancesIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.project).
				WithIndex(
					instancesIndexer.Object(),
					instancesIndexer.Name(),
					instancesIndexer.Keys,
				).
				WithInterceptorFuncs(tt.interceptors).
				Build()
			logger := zaptest.NewLogger(t).Sugar()
			reconciler := &AtlasProjectReconciler{
				Client: k8sClient,
				Log:    logger,
			}
			services := &AtlasProjectServices{
				projectService: tt.projectServiceMocker(),
			}
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}

			result, err := reconciler.create(ctx, "my-org-id", tt.project, services.projectService)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.result, result)
			assert.True(
				t,
				cmp.Equal(
					tt.conditions,
					ctx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
			assert.Equal(t, tt.finalizers, tt.project.Finalizers)
		})
	}
}

func TestDelete(t *testing.T) {
	deletionTime := metav1.Now()
	teamName := "my-team"
	teamID := "teamID"

	tests := map[string]struct {
		deletionProtection   bool
		atlasSDKMocker       func() *admin.APIClient
		projectServiceMocker func() project.ProjectService
		teamServiceMocker    func() teams.TeamsService
		interceptors         interceptor.Funcs
		objects              []client.Object
		result               reconcile.Result
		conditions           []api.Condition
		finalizers           []string
		wantErr              bool
	}{
		"should fail when unable to check project dependencies": {
			wantErr: true,
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
				return nil
			},
			teamServiceMocker: func() teams.TeamsService {
				return nil
			},
			interceptors: interceptor.Funcs{List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
				return errors.New("failed to list streams instances")
			}},
			objects: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "my-project",
						Namespace:         "default",
						Finalizers:        []string{customresource.FinalizerLabel},
						DeletionTimestamp: &deletionTime,
					},
				},
				&akov2.AtlasStreamInstance{ObjectMeta: metav1.ObjectMeta{Name: "instance0"}},
				&akov2.AtlasTeam{ObjectMeta: metav1.ObjectMeta{Name: teamName}},
			},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to determine if project has dependencies: failed to list streams instances"),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
		"should fail when project was deleted but it has dependencies": {
			wantErr: true,
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
				return nil
			},
			teamServiceMocker: func() teams.TeamsService {
				return nil
			},
			objects: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "my-project",
						Namespace:         "default",
						Finalizers:        []string{customresource.FinalizerLabel},
						DeletionTimestamp: &deletionTime,
					},
				},
				&akov2.AtlasStreamInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "instance-0",
						Namespace: "default",
					},
					Spec: akov2.AtlasStreamInstanceSpec{
						Project: common.ResourceRefNamespaced{
							Name:      "my-project",
							Namespace: "default",
						},
					},
				},
				&akov2.AtlasTeam{ObjectMeta: metav1.ObjectMeta{Name: teamName}},
			},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("the project cannot be deleted until dependencies were removed"),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
		"should do soft deletion when deletion protection is enabled": {
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
				return nil
			},
			teamServiceMocker: func() teams.TeamsService {
				return nil
			},
			deletionProtection: true,
			objects: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "my-project",
						Namespace:         "default",
						DeletionTimestamp: &deletionTime,
						Finalizers:        []string{customresource.FinalizerLabel},
						Annotations: map[string]string{
							customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
						},
					},
				},
				&akov2.AtlasStreamInstance{ObjectMeta: metav1.ObjectMeta{Name: "instance0"}},
				&akov2.AtlasTeam{ObjectMeta: metav1.ObjectMeta{Name: teamName}},
			},
			result: reconcile.Result{},
		},
		"should do soft deletion when resource policy is set to keep": {
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
				return nil
			},
			teamServiceMocker: func() teams.TeamsService {
				return nil
			},
			objects: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "my-project",
						Namespace:         "default",
						DeletionTimestamp: &deletionTime,
						Finalizers:        []string{customresource.FinalizerLabel},
						Annotations: map[string]string{
							customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
						},
					},
				},
				&akov2.AtlasStreamInstance{ObjectMeta: metav1.ObjectMeta{Name: "instance0"}},
				&akov2.AtlasTeam{ObjectMeta: metav1.ObjectMeta{Name: teamName}},
			},
		},
		"should update team status when project is deleted": {
			atlasSDKMocker: func() *admin.APIClient {
				mockPrivateEndpointAPI := mockadmin.NewPrivateEndpointServicesApi(t)
				mockPrivateEndpointAPI.EXPECT().
					ListPrivateEndpointService(mock.Anything, mock.Anything, mock.Anything).
					Return(admin.ListPrivateEndpointServiceApiRequest{ApiService: mockPrivateEndpointAPI})
				mockPrivateEndpointAPI.EXPECT().
					ListPrivateEndpointServiceExecute(admin.ListPrivateEndpointServiceApiRequest{ApiService: mockPrivateEndpointAPI}).
					Return([]admin.EndpointService{}, nil, nil)

				return &admin.APIClient{
					PrivateEndpointServicesApi: mockPrivateEndpointAPI,
				}
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().DeleteProject(context.Background(), mock.AnythingOfType("*project.Project")).Return(nil)

				return service
			},
			teamServiceMocker: func() teams.TeamsService {
				service := translation.NewTeamsServiceMock(t)
				service.EXPECT().ListProjectTeams(context.Background(), mock.Anything).Return([]teams.AssignedTeam{}, nil)
				return service
			},
			objects: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "my-project",
						Namespace:         "default",
						DeletionTimestamp: &deletionTime,
						Finalizers:        []string{customresource.FinalizerLabel},
					},
					Spec: akov2.AtlasProjectSpec{
						Teams: []akov2.Team{
							{
								TeamRef: common.ResourceRefNamespaced{
									Name:      "my-team",
									Namespace: "default",
								},
								Roles: []akov2.TeamRole{
									"PROJECT_OWNER",
								},
							},
						},
					},
					Status: status.AtlasProjectStatus{
						ID: "projectID",
						Teams: []status.ProjectTeamStatus{
							{
								ID: teamID,
								TeamRef: common.ResourceRefNamespaced{
									Name:      teamName,
									Namespace: "default",
								},
							},
						},
					},
				},
				&akov2.AtlasStreamInstance{ObjectMeta: metav1.ObjectMeta{Name: "instance0"}},
				&akov2.AtlasTeam{
					ObjectMeta: metav1.ObjectMeta{
						Name:      teamName,
						Namespace: "default",
					},
					Spec: akov2.TeamSpec{
						Name: "teamName",
					},
					Status: status.TeamStatus{
						ID: teamID,
						Projects: []status.TeamProject{
							{
								ID:   "projectID",
								Name: "project",
							},
						},
					},
				},
			},
			result: reconcile.Result{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(testScheme))
			require.NoError(t, corev1.AddToScheme(testScheme))
			logger := zaptest.NewLogger(t).Sugar()
			ctx := &workflow.Context{
				Context: context.Background(),
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312012: tt.atlasSDKMocker(),
				},
				Log: logger,
			}

			instancesIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
			customRoleIndexer := indexer.NewAtlasCustomRoleByProjectIndexer(logger.Desugar())
			peIndexer := indexer.NewAtlasPrivateEndpointByProjectIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.objects...).
				WithStatusSubresource(tt.objects...).
				WithIndex(instancesIndexer.Object(), instancesIndexer.Name(), instancesIndexer.Keys).
				WithIndex(customRoleIndexer.Object(), customRoleIndexer.Name(), customRoleIndexer.Keys).
				WithIndex(peIndexer.Object(), peIndexer.Name(), peIndexer.Keys).
				WithInterceptorFuncs(tt.interceptors).
				Build()

			reconciler := &AtlasProjectReconciler{
				Client:                   k8sClient,
				ObjectDeletionProtection: tt.deletionProtection,
				Log:                      logger,
				AtlasProvider:            &atlasmocks.TestProvider{},
				EventRecorder:            record.NewFakeRecorder(1),
			}

			atlasProject := tt.objects[0].(*akov2.AtlasProject)
			services := &AtlasProjectServices{
				projectService: tt.projectServiceMocker(),
				teamsService:   tt.teamServiceMocker(),
			}
			result, err := reconciler.delete(ctx, services, "my-org-id", atlasProject)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.result, result)
			assert.True(
				t,
				cmp.Equal(
					tt.conditions,
					ctx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
			assert.Equal(t, tt.finalizers, atlasProject.Finalizers)
		})
	}
}

func TestHasDependencies(t *testing.T) {
	t.Run("should return error when unable to list stream instances", func(t *testing.T) {
		p := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
		}
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(p).
			WithInterceptorFuncs(interceptor.Funcs{List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
				return errors.New("failed to list instances")
			}}).
			Build()
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		ok, err := reconciler.hasDependencies(ctx, p)
		require.ErrorContains(t, err, "failed to list instances")
		assert.False(t, ok)
	})

	t.Run("should return false when project has no dependencies", func(t *testing.T) {
		p := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}
		logger := zaptest.NewLogger(t)
		instanceIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
		customRoleIndexer := indexer.NewAtlasCustomRoleByProjectIndexer(zap.L())
		peIndexer := indexer.NewAtlasPrivateEndpointByProjectIndexer(logger)
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(p).
			WithIndex(instanceIndexer.Object(), instanceIndexer.Name(), instanceIndexer.Keys).
			WithIndex(customRoleIndexer.Object(), customRoleIndexer.Name(), customRoleIndexer.Keys).
			WithIndex(peIndexer.Object(), peIndexer.Name(), peIndexer.Keys).
			Build()
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}

		ok, err := reconciler.hasDependencies(ctx, p)
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("should return true when project has streams as dependencies", func(t *testing.T) {
		p := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
		}
		streamsInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "instance-0",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
			},
		}
		logger := zaptest.NewLogger(t)
		instanceIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(logger)
		peIndexer := indexer.NewAtlasPrivateEndpointByProjectIndexer(logger)
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(p, streamsInstance).
			WithIndex(instanceIndexer.Object(), instanceIndexer.Name(), instanceIndexer.Keys).
			WithIndex(peIndexer.Object(), peIndexer.Name(), peIndexer.Keys).
			Build()
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		ok, err := reconciler.hasDependencies(ctx, p)
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("should return true when project has private endpoints as dependencies", func(t *testing.T) {
		p := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
		}
		pe := &akov2.AtlasPrivateEndpoint{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pe-0",
				Namespace: "default",
			},
			Spec: akov2.AtlasPrivateEndpointSpec{
				ProjectDualReference: akov2.ProjectDualReference{
					ProjectRef: &common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
					},
				},
			},
		}
		logger := zaptest.NewLogger(t)
		ctx := &workflow.Context{
			Context: context.Background(),
		}
		instanceIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(logger)
		customRolesIndexer := indexer.NewAtlasCustomRoleByProjectIndexer(logger)
		peIndexer := indexer.NewAtlasPrivateEndpointByProjectIndexer(logger)
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(p, pe).
			WithIndex(instanceIndexer.Object(), instanceIndexer.Name(), instanceIndexer.Keys).
			WithIndex(customRolesIndexer.Object(), customRolesIndexer.Name(), customRolesIndexer.Keys).
			WithIndex(peIndexer.Object(), peIndexer.Name(), peIndexer.Keys).
			Build()
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}

		ok, err := reconciler.hasDependencies(ctx, p)
		require.NoError(t, err)
		assert.True(t, ok)
	})
}
