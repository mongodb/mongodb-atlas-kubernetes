package datafederation

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	atlasmocks "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
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
				SdkClientFunc: func(_ *client.ObjectKey, _ *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return &admin.APIClient{}, "", nil
				},
			},
		},
		{
			name: "failure",
			provider: &atlasmocks.TestProvider{
				SdkClientFunc: func(_ *client.ObjectKey, _ *zap.SugaredLogger) (*admin.APIClient, string, error) {
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
