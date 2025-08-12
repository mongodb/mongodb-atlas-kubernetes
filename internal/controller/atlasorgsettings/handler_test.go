package atlasorgsettings

//
//import (
//	"context"
//	"errors"
//	atlas_controllers "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
//	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
//	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
//	"go.mongodb.org/atlas-sdk/v20250312006/admin"
//	adminmock "go.mongodb.org/atlas-sdk/v20250312006/mockadmin"
//	"go.uber.org/zap"
//	"net/http"
//	"sigs.k8s.io/controller-runtime/pkg/client"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//
//	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
//	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/atlasorgsettings"
//	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
//	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
//)
//
//type fakeAtlasOrgSettingsService struct {
//	mock.Mock
//}
//
//func (f *fakeAtlasOrgSettingsService) Update(ctx context.Context, orgID string, settings *atlasorgsettings.AtlasOrgSettings) (*atlasorgsettings.AtlasOrgSettings, error) {
//	args := f.Called(ctx, orgID, settings)
//	if args.Get(0) == nil {
//		return nil, args.Error(1)
//	}
//	return args.Get(0).(*atlasorgsettings.AtlasOrgSettings), args.Error(1)
//}
//
//func (f *fakeAtlasOrgSettingsService) Get(ctx context.Context, orgID string) (*atlasorgsettings.AtlasOrgSettings, error) {
//	args := f.Called(ctx, orgID)
//	if args.Get(0) == nil {
//		return nil, args.Error(1)
//	}
//	return args.Get(0).(*atlasorgsettings.AtlasOrgSettings), args.Error(1)
//}
//
//func Test_AtlasOrgSettingsHandler_upsert(t *testing.T) {
//	type fields struct {
//		service *fakeAtlasOrgSettingsService
//	}
//	type args struct {
//		currentState state.ResourceState
//		nextState    state.ResourceState
//		aos          *akov2.AtlasOrgSettings
//	}
//	tests := []struct {
//		name           string
//		fields         fields
//		args           args
//		setupMock      func(svc *fakeAtlasOrgSettingsService)
//		wantErr        bool
//		wantResultZero bool
//	}{
//		{
//			name: "returns error when Update fails",
//			fields: fields{
//				service: new(fakeAtlasOrgSettingsService),
//			},
//			args: args{
//				currentState: state.StateInitial,
//				nextState:    state.StateCreated,
//				aos:          &akov2.AtlasOrgSettings{Spec: akov2.AtlasOrgSettingsSpec{OrgID: "org1"}},
//			},
//			setupMock: func(svc *fakeAtlasOrgSettingsService) {
//				svc.On("Update", mock.Anything, "org1", mock.Anything).Return(nil, errors.New("update failed"))
//			},
//			wantErr:        true,
//			wantResultZero: true,
//		},
//		{
//			name: "returns error when Update returns nil",
//			fields: fields{
//				service: new(fakeAtlasOrgSettingsService),
//			},
//			args: args{
//				currentState: state.StateInitial,
//				nextState:    state.StateCreated,
//				aos:          &akov2.AtlasOrgSettings{Spec: akov2.AtlasOrgSettingsSpec{OrgID: "org1"}},
//			},
//			setupMock: func(svc *fakeAtlasOrgSettingsService) {
//				svc.On("Update", mock.Anything, "org1", mock.Anything).Return(nil, nil)
//			},
//			wantErr:        true,
//			wantResultZero: true,
//		},
//		{
//			name: "returns next state on success",
//			fields: fields{
//				service: new(fakeAtlasOrgSettingsService),
//			},
//			args: args{
//				currentState: state.StateInitial,
//				nextState:    state.StateCreated,
//				aos:          &akov2.AtlasOrgSettings{Spec: akov2.AtlasOrgSettingsSpec{OrgID: "org1"}},
//			},
//			setupMock: func(svc *fakeAtlasOrgSettingsService) {
//				svc.On("Update", mock.Anything, "org1", mock.Anything).Return(&atlasorgsettings.AtlasOrgSettings{}, nil)
//			},
//			wantErr:        false,
//			wantResultZero: false,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.setupMock != nil {
//				tt.setupMock(tt.fields.service)
//			}
//			handler := &AtlasOrgSettingsHandler{
//				StateHandler: nil,
//				AtlasReconciler: reconciler.AtlasReconciler{
//					AtlasProvider: &atlas.TestProvider{
//						ClientFunc: nil,
//						//SdkClientSetFunc: func(ctx context.Context, creds *atlas_controllers.Credentials, log *zap.SugaredLogger) (*atlas_controllers.ClientSet, error) {
//						//	return &atlas_controllers.ClientSet{
//						//		SdkClient20250312006: &admin.APIClient{
//						//			OrganizationsApi: func() *adminmock.OrganizationsApi {
//						//				api := adminmock.NewOrganizationsApi(t)
//						//				api.Wi
//						//				return api
//						//			}(),
//						//		},
//						//	}, nil
//						//},
//						IsCloudGovFunc:  nil,
//						IsSupportedFunc: nil,
//					},
//					Client:          nil,
//					Log:             nil,
//					GlobalSecretRef: client.ObjectKey{},
//				},
//				deletionProtection: false,
//			}
//			got, err := handler.upsert(context.Background(), tt.args.currentState, tt.args.nextState, tt.args.aos)
//			if tt.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//			if tt.wantResultZero {
//				assert.Equal(t, ctrlstate.Result{}, got)
//			} else {
//				assert.NotEqual(t, ctrlstate.Result{}, got)
//			}
//		})
//	}
//}
//
//func Test_AtlasOrgSettingsHandler_HandleInitial_newReconcileContextFails(t *testing.T) {
//	handler := &AtlasOrgSettingsHandler{
//		newReconcileContext: func(ctx context.Context, aos *akov2.AtlasOrgSettings) (*reconcileContext, error) {
//			return nil, errors.New("failed to create reconcile context")
//		},
//	}
//	aos := &akov2.AtlasOrgSettings{}
//	res, err := handler.HandleInitial(context.Background(), aos)
//	assert.Error(t, err)
//	assert.Equal(t, ctrlstate.Result{}, res)
//}
//
//func Test_AtlasOrgSettingsHandler_unmanage(t *testing.T) {
//	tests := []struct {
//		name  string
//		orgID string
//	}{
//		{
//			name:  "returns deleted state",
//			orgID: "org1",
//		},
//		{
//			name:  "returns deleted state with empty orgID",
//			orgID: "",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			handler := &AtlasOrgSettingsHandler{}
//			res, err := handler.unmanage(tt.orgID)
//			assert.NoError(t, err)
//			assert.NotEqual(t, ctrlstate.Result{}, res)
//		})
//	}
//}
