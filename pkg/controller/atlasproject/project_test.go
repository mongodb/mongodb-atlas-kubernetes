package atlasproject

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.mongodb.org/atlas/mongodbatlas"
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

	atlasmocks "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/indexer"
)

func TestHandleProject(t *testing.T) {
	deletionTime := metav1.Now()

	tests := map[string]struct {
		atlasClientMocker    func() *mongodbatlas.Client
		atlasSDKMocker       func() *admin.APIClient
		projectServiceMocker func() project.ProjectService
		interceptors         interceptor.Funcs
		project              *akov2.AtlasProject
		result               reconcile.Result
		conditions           []api.Condition
		finalizers           []string
	}{
		"should fail to get project from atlas": {
			atlasClientMocker: func() *mongodbatlas.Client {
				return nil
			},
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(nil, errors.New("failed to get project"))

				return service
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
					WithReason(string(workflow.ProjectNotCreatedInAtlas)).
					WithMessageRegexp("failed to get project"),
			},
		},
		"should create project": {
			atlasClientMocker: func() *mongodbatlas.Client {
				return nil
			},
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
			atlasClientMocker: func() *mongodbatlas.Client {
				projectsMock := &atlasmocks.ProjectsClientMock{
					DeleteFunc: func(projectID string) (*mongodbatlas.Response, error) {
						return nil, nil
					},
				}

				return &mongodbatlas.Client{
					Projects: projectsMock,
				}
			},
			atlasSDKMocker: func() *admin.APIClient {
				mockPrivateEndpointAPI := mockadmin.NewPrivateEndpointServicesApi(t)
				mockPrivateEndpointAPI.EXPECT().
					ListPrivateEndpointServices(mock.Anything, mock.Anything, mock.Anything).
					Return(admin.ListPrivateEndpointServicesApiRequest{ApiService: mockPrivateEndpointAPI})
				mockPrivateEndpointAPI.EXPECT().
					ListPrivateEndpointServicesExecute(admin.ListPrivateEndpointServicesApiRequest{ApiService: mockPrivateEndpointAPI}).
					Return([]admin.EndpointService{}, nil, nil)

				mockPeeringEndpointAPI := mockadmin.NewNetworkPeeringApi(t)
				mockPeeringEndpointAPI.EXPECT().ListPeeringConnectionsWithParams(mock.Anything, mock.Anything).
					Return(admin.ListPeeringConnectionsApiRequest{ApiService: mockPeeringEndpointAPI})
				mockPeeringEndpointAPI.EXPECT().
					ListPeeringConnectionsExecute(admin.ListPeeringConnectionsApiRequest{ApiService: mockPeeringEndpointAPI}).
					Return(&admin.PaginatedContainerPeer{}, nil, nil)
				mockTeamAPI := mockadmin.NewTeamsApi(t)
				mockTeamAPI.EXPECT().ListProjectTeams(context.Background(), mock.Anything).
					Return(admin.ListProjectTeamsApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().ListProjectTeamsExecute(mock.Anything).
					Return(&admin.PaginatedTeamRole{}, &http.Response{}, nil)

				return &admin.APIClient{
					PrivateEndpointServicesApi: mockPrivateEndpointAPI,
					NetworkPeeringApi:          mockPeeringEndpointAPI,
					TeamsApi:                   mockTeamAPI,
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
			atlasClientMocker: func() *mongodbatlas.Client {
				return nil
			},
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(nil, nil)

				return service
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
			atlasClientMocker: func() *mongodbatlas.Client {
				return nil
			},
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(nil, nil)

				return service
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
			result: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.AtlasFinalizerNotRemoved)).
					WithMessageRegexp("failed to remove finalizer"),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
		"should fail to configure authentication modes": {
			atlasClientMocker: func() *mongodbatlas.Client {
				return nil
			},
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(&project.Project{ID: "projectID", Name: "my-project"}, nil)

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
			result: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("secrets \"invalid-ref\" not found"),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
		"should configure project resources": {
			atlasClientMocker: func() *mongodbatlas.Client {
				integrations := &atlasmocks.ThirdPartyIntegrationsClientMock{
					ListFunc: func(projectID string) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
						return &mongodbatlas.ThirdPartyIntegrations{}, nil, nil
					},
				}
				encryptionAtRest := &atlasmocks.EncryptionAtRestClientMock{
					GetFunc: func(projectID string) (*mongodbatlas.EncryptionAtRest, *mongodbatlas.Response, error) {
						return &mongodbatlas.EncryptionAtRest{}, nil, nil
					},
				}
				projectAPI := &atlasmocks.ProjectsClientMock{}

				return &mongodbatlas.Client{
					Integrations:      integrations,
					EncryptionsAtRest: encryptionAtRest,
					Projects:          projectAPI,
				}
			},
			atlasSDKMocker: func() *admin.APIClient {
				ipAccessList := mockadmin.NewProjectIPAccessListApi(t)
				ipAccessList.EXPECT().ListProjectIpAccessLists(context.Background(), "projectID").
					Return(admin.ListProjectIpAccessListsApiRequest{ApiService: ipAccessList})
				ipAccessList.EXPECT().ListProjectIpAccessListsExecute(mock.AnythingOfType("admin.ListProjectIpAccessListsApiRequest")).
					Return(nil, nil, nil)
				privateEndpoints := mockadmin.NewPrivateEndpointServicesApi(t)
				privateEndpoints.EXPECT().ListPrivateEndpointServices(context.Background(), "projectID", mock.Anything).
					Return(admin.ListPrivateEndpointServicesApiRequest{ApiService: privateEndpoints})
				privateEndpoints.EXPECT().ListPrivateEndpointServicesExecute(mock.AnythingOfType("admin.ListPrivateEndpointServicesApiRequest")).
					Return(nil, nil, nil)
				networkPeering := mockadmin.NewNetworkPeeringApi(t)
				networkPeering.EXPECT().ListPeeringConnectionsWithParams(context.Background(), mock.AnythingOfType("*admin.ListPeeringConnectionsApiParams")).
					Return(admin.ListPeeringConnectionsApiRequest{ApiService: networkPeering})
				networkPeering.EXPECT().ListPeeringConnectionsExecute(mock.AnythingOfType("admin.ListPeeringConnectionsApiRequest")).
					Return(nil, nil, nil)
				networkPeering.EXPECT().ListPeeringContainers(context.Background(), "projectID").
					Return(admin.ListPeeringContainersApiRequest{ApiService: networkPeering})
				networkPeering.EXPECT().ListPeeringContainersExecute(mock.AnythingOfType("admin.ListPeeringContainersApiRequest")).
					Return(nil, nil, nil)
				audit := mockadmin.NewAuditingApi(t)
				audit.EXPECT().GetAuditingConfiguration(context.Background(), "projectID").
					Return(admin.GetAuditingConfigurationApiRequest{ApiService: audit})
				audit.EXPECT().GetAuditingConfigurationExecute(mock.AnythingOfType("admin.GetAuditingConfigurationApiRequest")).
					Return(nil, nil, nil)
				customRoles := mockadmin.NewCustomDatabaseRolesApi(t)
				customRoles.EXPECT().ListCustomDatabaseRoles(context.Background(), "projectID").
					Return(admin.ListCustomDatabaseRolesApiRequest{ApiService: customRoles})
				customRoles.EXPECT().ListCustomDatabaseRolesExecute(mock.AnythingOfType("admin.ListCustomDatabaseRolesApiRequest")).
					Return(nil, nil, nil)
				projectAPI := mockadmin.NewProjectsApi(t)
				projectAPI.EXPECT().GetProjectSettings(context.Background(), "projectID").
					Return(admin.GetProjectSettingsApiRequest{ApiService: projectAPI})
				projectAPI.EXPECT().GetProjectSettingsExecute(mock.AnythingOfType("admin.GetProjectSettingsApiRequest")).
					Return(admin.NewGroupSettings(), nil, nil)
				backup := mockadmin.NewCloudBackupsApi(t)
				backup.EXPECT().GetDataProtectionSettings(context.Background(), "projectID").
					Return(admin.GetDataProtectionSettingsApiRequest{ApiService: backup})
				backup.EXPECT().GetDataProtectionSettingsExecute(mock.AnythingOfType("admin.GetDataProtectionSettingsApiRequest")).
					Return(nil, nil, nil)
				mockTeamAPI := mockadmin.NewTeamsApi(t)
				mockTeamAPI.EXPECT().ListProjectTeams(context.Background(), mock.Anything).
					Return(admin.ListProjectTeamsApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().ListProjectTeamsExecute(mock.Anything).
					Return(nil, &http.Response{}, nil)

				return &admin.APIClient{
					ProjectIPAccessListApi:     ipAccessList,
					PrivateEndpointServicesApi: privateEndpoints,
					NetworkPeeringApi:          networkPeering,
					AuditingApi:                audit,
					CustomDatabaseRolesApi:     customRoles,
					ProjectsApi:                projectAPI,
					CloudBackupsApi:            backup,
					TeamsApi:                   mockTeamAPI,
				}
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(&project.Project{ID: "projectID", Name: "my-project"}, nil)

				return service
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "my-project",
					Namespace:  "default",
					Finalizers: []string{customresource.FinalizerLabel},
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
			atlasClientMocker: func() *mongodbatlas.Client {
				integrations := &atlasmocks.ThirdPartyIntegrationsClientMock{
					ListFunc: func(projectID string) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
						return &mongodbatlas.ThirdPartyIntegrations{}, nil, nil
					},
				}
				encryptionAtRest := &atlasmocks.EncryptionAtRestClientMock{
					GetFunc: func(projectID string) (*mongodbatlas.EncryptionAtRest, *mongodbatlas.Response, error) {
						return &mongodbatlas.EncryptionAtRest{}, nil, nil
					},
				}

				return &mongodbatlas.Client{
					Integrations:      integrations,
					EncryptionsAtRest: encryptionAtRest,
				}
			},
			atlasSDKMocker: func() *admin.APIClient {
				ipAccessList := mockadmin.NewProjectIPAccessListApi(t)
				ipAccessList.EXPECT().ListProjectIpAccessLists(context.Background(), "projectID").
					Return(admin.ListProjectIpAccessListsApiRequest{ApiService: ipAccessList})
				ipAccessList.EXPECT().ListProjectIpAccessListsExecute(mock.AnythingOfType("admin.ListProjectIpAccessListsApiRequest")).
					Return(nil, nil, errors.New("failed to list IP Access List"))
				privateEndpoints := mockadmin.NewPrivateEndpointServicesApi(t)
				privateEndpoints.EXPECT().ListPrivateEndpointServices(context.Background(), "projectID", mock.Anything).
					Return(admin.ListPrivateEndpointServicesApiRequest{ApiService: privateEndpoints})
				privateEndpoints.EXPECT().ListPrivateEndpointServicesExecute(mock.AnythingOfType("admin.ListPrivateEndpointServicesApiRequest")).
					Return(nil, nil, nil)
				networkPeering := mockadmin.NewNetworkPeeringApi(t)
				networkPeering.EXPECT().ListPeeringConnectionsWithParams(context.Background(), mock.AnythingOfType("*admin.ListPeeringConnectionsApiParams")).
					Return(admin.ListPeeringConnectionsApiRequest{ApiService: networkPeering})
				networkPeering.EXPECT().ListPeeringConnectionsExecute(mock.AnythingOfType("admin.ListPeeringConnectionsApiRequest")).
					Return(nil, nil, nil)
				networkPeering.EXPECT().ListPeeringContainers(context.Background(), "projectID").
					Return(admin.ListPeeringContainersApiRequest{ApiService: networkPeering})
				networkPeering.EXPECT().ListPeeringContainersExecute(mock.AnythingOfType("admin.ListPeeringContainersApiRequest")).
					Return(nil, nil, nil)
				audit := mockadmin.NewAuditingApi(t)
				audit.EXPECT().GetAuditingConfiguration(context.Background(), "projectID").
					Return(admin.GetAuditingConfigurationApiRequest{ApiService: audit})
				audit.EXPECT().GetAuditingConfigurationExecute(mock.AnythingOfType("admin.GetAuditingConfigurationApiRequest")).
					Return(nil, nil, nil)
				customRoles := mockadmin.NewCustomDatabaseRolesApi(t)
				customRoles.EXPECT().ListCustomDatabaseRoles(context.Background(), "projectID").
					Return(admin.ListCustomDatabaseRolesApiRequest{ApiService: customRoles})
				customRoles.EXPECT().ListCustomDatabaseRolesExecute(mock.AnythingOfType("admin.ListCustomDatabaseRolesApiRequest")).
					Return(nil, nil, nil)
				projectAPI := mockadmin.NewProjectsApi(t)
				projectAPI.EXPECT().GetProjectSettings(context.Background(), "projectID").
					Return(admin.GetProjectSettingsApiRequest{ApiService: projectAPI})
				projectAPI.EXPECT().GetProjectSettingsExecute(mock.AnythingOfType("admin.GetProjectSettingsApiRequest")).
					Return(admin.NewGroupSettings(), nil, nil)
				backup := mockadmin.NewCloudBackupsApi(t)
				backup.EXPECT().GetDataProtectionSettings(context.Background(), "projectID").
					Return(admin.GetDataProtectionSettingsApiRequest{ApiService: backup})
				backup.EXPECT().GetDataProtectionSettingsExecute(mock.AnythingOfType("admin.GetDataProtectionSettingsApiRequest")).
					Return(nil, nil, nil)
				mockTeamAPI := mockadmin.NewTeamsApi(t)
				mockTeamAPI.EXPECT().ListProjectTeams(context.Background(), mock.Anything).
					Return(admin.ListProjectTeamsApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().ListProjectTeamsExecute(mock.Anything).
					Return(nil, &http.Response{}, nil)

				return &admin.APIClient{
					ProjectIPAccessListApi:     ipAccessList,
					PrivateEndpointServicesApi: privateEndpoints,
					NetworkPeeringApi:          networkPeering,
					AuditingApi:                audit,
					CustomDatabaseRolesApi:     customRoles,
					ProjectsApi:                projectAPI,
					CloudBackupsApi:            backup,
					TeamsApi:                   mockTeamAPI,
				}
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(&project.Project{ID: "projectID", Name: "my-project"}, nil)

				return service
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "my-project",
					Namespace:  "default",
					Finalizers: []string{customresource.FinalizerLabel},
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
				Status: status.AtlasProjectStatus{
					ID: "projectID",
				},
			},
			result: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			conditions: []api.Condition{
				api.TrueCondition(api.ProjectReadyType),
				api.FalseCondition(api.IPAccessListReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to get ip access list from Atlas: failed to list IP Access List"),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
		"should fail to save last applied config": {
			atlasClientMocker: func() *mongodbatlas.Client {
				integrations := &atlasmocks.ThirdPartyIntegrationsClientMock{
					ListFunc: func(projectID string) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
						return &mongodbatlas.ThirdPartyIntegrations{}, nil, nil
					},
				}
				encryptionAtRest := &atlasmocks.EncryptionAtRestClientMock{
					GetFunc: func(projectID string) (*mongodbatlas.EncryptionAtRest, *mongodbatlas.Response, error) {
						return &mongodbatlas.EncryptionAtRest{}, nil, nil
					},
				}

				return &mongodbatlas.Client{
					Integrations:      integrations,
					EncryptionsAtRest: encryptionAtRest,
				}
			},
			atlasSDKMocker: func() *admin.APIClient {
				ipAccessList := mockadmin.NewProjectIPAccessListApi(t)
				ipAccessList.EXPECT().ListProjectIpAccessLists(context.Background(), "projectID").
					Return(admin.ListProjectIpAccessListsApiRequest{ApiService: ipAccessList})
				ipAccessList.EXPECT().ListProjectIpAccessListsExecute(mock.AnythingOfType("admin.ListProjectIpAccessListsApiRequest")).
					Return(nil, nil, nil)
				privateEndpoints := mockadmin.NewPrivateEndpointServicesApi(t)
				privateEndpoints.EXPECT().ListPrivateEndpointServices(context.Background(), "projectID", mock.Anything).
					Return(admin.ListPrivateEndpointServicesApiRequest{ApiService: privateEndpoints})
				privateEndpoints.EXPECT().ListPrivateEndpointServicesExecute(mock.AnythingOfType("admin.ListPrivateEndpointServicesApiRequest")).
					Return(nil, nil, nil)
				networkPeering := mockadmin.NewNetworkPeeringApi(t)
				networkPeering.EXPECT().ListPeeringConnectionsWithParams(context.Background(), mock.AnythingOfType("*admin.ListPeeringConnectionsApiParams")).
					Return(admin.ListPeeringConnectionsApiRequest{ApiService: networkPeering})
				networkPeering.EXPECT().ListPeeringConnectionsExecute(mock.AnythingOfType("admin.ListPeeringConnectionsApiRequest")).
					Return(nil, nil, nil)
				networkPeering.EXPECT().ListPeeringContainers(context.Background(), "projectID").
					Return(admin.ListPeeringContainersApiRequest{ApiService: networkPeering})
				networkPeering.EXPECT().ListPeeringContainersExecute(mock.AnythingOfType("admin.ListPeeringContainersApiRequest")).
					Return(nil, nil, nil)
				audit := mockadmin.NewAuditingApi(t)
				audit.EXPECT().GetAuditingConfiguration(context.Background(), "projectID").
					Return(admin.GetAuditingConfigurationApiRequest{ApiService: audit})
				audit.EXPECT().GetAuditingConfigurationExecute(mock.AnythingOfType("admin.GetAuditingConfigurationApiRequest")).
					Return(nil, nil, nil)
				customRoles := mockadmin.NewCustomDatabaseRolesApi(t)
				customRoles.EXPECT().ListCustomDatabaseRoles(context.Background(), "projectID").
					Return(admin.ListCustomDatabaseRolesApiRequest{ApiService: customRoles})
				customRoles.EXPECT().ListCustomDatabaseRolesExecute(mock.AnythingOfType("admin.ListCustomDatabaseRolesApiRequest")).
					Return(nil, nil, nil)
				projectAPI := mockadmin.NewProjectsApi(t)
				projectAPI.EXPECT().GetProjectSettings(context.Background(), "projectID").
					Return(admin.GetProjectSettingsApiRequest{ApiService: projectAPI})
				projectAPI.EXPECT().GetProjectSettingsExecute(mock.AnythingOfType("admin.GetProjectSettingsApiRequest")).
					Return(admin.NewGroupSettings(), nil, nil)
				backup := mockadmin.NewCloudBackupsApi(t)
				backup.EXPECT().GetDataProtectionSettings(context.Background(), "projectID").
					Return(admin.GetDataProtectionSettingsApiRequest{ApiService: backup})
				backup.EXPECT().GetDataProtectionSettingsExecute(mock.AnythingOfType("admin.GetDataProtectionSettingsApiRequest")).
					Return(nil, nil, nil)
				mockTeamAPI := mockadmin.NewTeamsApi(t)
				mockTeamAPI.EXPECT().ListProjectTeams(context.Background(), mock.Anything).
					Return(admin.ListProjectTeamsApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().ListProjectTeamsExecute(mock.Anything).
					Return(nil, &http.Response{}, nil)

				return &admin.APIClient{
					ProjectIPAccessListApi:     ipAccessList,
					PrivateEndpointServicesApi: privateEndpoints,
					NetworkPeeringApi:          networkPeering,
					AuditingApi:                audit,
					CustomDatabaseRolesApi:     customRoles,
					ProjectsApi:                projectAPI,
					CloudBackupsApi:            backup,
					TeamsApi:                   mockTeamAPI,
				}
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(&project.Project{ID: "projectID", Name: "my-project"}, nil)

				return service
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "my-project",
					Namespace:  "default",
					Finalizers: []string{customresource.FinalizerLabel},
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
			result: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
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
			instancesIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.project).
				WithStatusSubresource(tt.project).
				WithIndex(
					instancesIndexer.Object(),
					instancesIndexer.Name(),
					instancesIndexer.Keys,
				).
				WithInterceptorFuncs(tt.interceptors).
				Build()
			logger := zaptest.NewLogger(t).Sugar()
			reconciler := &AtlasProjectReconciler{
				Client:         k8sClient,
				Log:            logger,
				projectService: tt.projectServiceMocker(),
				EventRecorder:  record.NewFakeRecorder(30),
			}
			ctx := &workflow.Context{
				Context:   context.Background(),
				Log:       logger,
				Client:    tt.atlasClientMocker(),
				SdkClient: tt.atlasSDKMocker(),
			}

			result, err := reconciler.handleProject(ctx, "my-org-id", tt.project)
			require.NoError(t, err)
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
	}{
		"should fail to create project": {
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
			result: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.ProjectNotCreatedInAtlas)).
					WithMessageRegexp("failed to create project"),
			},
		},
		"should fail to add finalizer when creating a project": {
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
			result: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.AtlasFinalizerNotSet)).
					WithMessageRegexp("failed to patch project with finalizers"),
			},
		},
		"should fail to add last applied config when creating a project": {
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
			result: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
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
				Client:         k8sClient,
				Log:            logger,
				projectService: tt.projectServiceMocker(),
			}
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}

			result, err := reconciler.create(ctx, "my-org-id", tt.project)
			require.NoError(t, err)
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
		atlasClientMocker    func() *mongodbatlas.Client
		atlasSDKMocker       func() *admin.APIClient
		projectServiceMocker func() project.ProjectService
		interceptors         interceptor.Funcs
		objects              []client.Object
		result               reconcile.Result
		conditions           []api.Condition
		finalizers           []string
	}{
		"should fail when unable to check project dependencies": {
			atlasClientMocker: func() *mongodbatlas.Client {
				return nil
			},
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
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
			result: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to determine if project has dependencies: failed to list streams instances"),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
		"should fail when project was deleted but it has dependencies": {
			atlasClientMocker: func() *mongodbatlas.Client {
				return nil
			},
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
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
			result: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			conditions: []api.Condition{
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("the project cannot be deleted until dependencies were removed"),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
		"should do soft deletion when deletion protection is enabled": {
			atlasClientMocker: func() *mongodbatlas.Client {
				return nil
			},
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
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
			atlasClientMocker: func() *mongodbatlas.Client {
				return nil
			},
			atlasSDKMocker: func() *admin.APIClient {
				return nil
			},
			projectServiceMocker: func() project.ProjectService {
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
			atlasClientMocker: func() *mongodbatlas.Client {
				projectsMock := &atlasmocks.ProjectsClientMock{
					DeleteFunc: func(projectID string) (*mongodbatlas.Response, error) {
						return nil, nil
					},
				}

				return &mongodbatlas.Client{
					Projects: projectsMock,
				}
			},
			atlasSDKMocker: func() *admin.APIClient {
				mockPrivateEndpointAPI := mockadmin.NewPrivateEndpointServicesApi(t)
				mockPrivateEndpointAPI.EXPECT().
					ListPrivateEndpointServices(mock.Anything, mock.Anything, mock.Anything).
					Return(admin.ListPrivateEndpointServicesApiRequest{ApiService: mockPrivateEndpointAPI})
				mockPrivateEndpointAPI.EXPECT().
					ListPrivateEndpointServicesExecute(admin.ListPrivateEndpointServicesApiRequest{ApiService: mockPrivateEndpointAPI}).
					Return([]admin.EndpointService{}, nil, nil)

				mockPeeringEndpointAPI := mockadmin.NewNetworkPeeringApi(t)
				mockPeeringEndpointAPI.EXPECT().ListPeeringConnectionsWithParams(mock.Anything, mock.Anything).
					Return(admin.ListPeeringConnectionsApiRequest{ApiService: mockPeeringEndpointAPI})
				mockPeeringEndpointAPI.EXPECT().
					ListPeeringConnectionsExecute(admin.ListPeeringConnectionsApiRequest{ApiService: mockPeeringEndpointAPI}).
					Return(&admin.PaginatedContainerPeer{}, nil, nil)
				mockTeamAPI := mockadmin.NewTeamsApi(t)
				mockTeamAPI.EXPECT().ListProjectTeams(context.Background(), mock.Anything).
					Return(admin.ListProjectTeamsApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().ListProjectTeamsExecute(mock.Anything).
					Return(nil, &http.Response{}, nil)
				return &admin.APIClient{
					PrivateEndpointServicesApi: mockPrivateEndpointAPI,
					NetworkPeeringApi:          mockPeeringEndpointAPI,
					TeamsApi:                   mockTeamAPI,
				}
			},
			projectServiceMocker: func() project.ProjectService {
				service := translation.NewProjectServiceMock(t)
				service.EXPECT().DeleteProject(context.Background(), mock.AnythingOfType("*project.Project")).Return(nil)

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
			instancesIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.objects...).
				WithStatusSubresource(tt.objects...).
				WithIndex(
					instancesIndexer.Object(),
					instancesIndexer.Name(),
					instancesIndexer.Keys,
				).
				WithInterceptorFuncs(tt.interceptors).
				Build()
			logger := zaptest.NewLogger(t).Sugar()
			reconciler := &AtlasProjectReconciler{
				Client:                   k8sClient,
				ObjectDeletionProtection: tt.deletionProtection,
				Log:                      logger,
				AtlasProvider: &atlasmocks.TestProvider{
					ClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
						return tt.atlasClientMocker(), "", nil
					},
				},
				projectService: tt.projectServiceMocker(),
				EventRecorder:  record.NewFakeRecorder(1),
			}
			ctx := &workflow.Context{
				Context:   context.Background(),
				Client:    tt.atlasClientMocker(),
				SdkClient: tt.atlasSDKMocker(),
				Log:       logger,
			}

			atlasProject := tt.objects[0].(*akov2.AtlasProject)
			result, err := reconciler.delete(ctx, "my-org-id", atlasProject)
			require.NoError(t, err)
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
		instanceIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(p).
			WithIndex(
				instanceIndexer.Object(),
				instanceIndexer.Name(),
				instanceIndexer.Keys,
			).
			Build()
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		ok, err := reconciler.hasDependencies(ctx, p)
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("should return true when project has dependencies", func(t *testing.T) {
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
		instanceIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(p, streamsInstance).
			WithIndex(
				instanceIndexer.Object(),
				instanceIndexer.Name(),
				instanceIndexer.Keys,
			).
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
}
