package datafederation

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	adminv20241113001 "go.mongodb.org/atlas-sdk/v20241113001/admin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	atlasmocks "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
)

func TestNewDatafederationService(t *testing.T) {
	for _, tt := range []struct {
		name     string
		provider atlas.Provider
		wantErr  string
	}{
		{
			name: "success",
			provider: &atlasmocks.TestProvider{
				SdkSetClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error) {
					return &atlas.ClientSet{
						SdkClient20231115008: &admin.APIClient{},
						SdkClient20241113001: &adminv20241113001.APIClient{},
					}, "", nil
				},
			},
		},
		{
			name: "failure",
			provider: &atlasmocks.TestProvider{
				SdkSetClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error) {
					return nil, "", errors.New("fake error")
				},
			},
			wantErr: "failed to create versioned client: failed to instantiate Versioned Atlas client: fake error",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAtlasDataFederationService(context.Background(), tt.provider, &types.NamespacedName{}, zap.S())
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			require.Equal(t, tt.wantErr, gotErr)
		})
	}
}
