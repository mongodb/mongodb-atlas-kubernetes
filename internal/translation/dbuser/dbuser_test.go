package dbuser_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/dbuser"
)

func TestNewAtlasDatabaseUsersService(t *testing.T) {
	ctx := context.Background()
	provider := &atlas.TestProvider{
		SdkClientFunc: func(_ *client.ObjectKey, _ *zap.SugaredLogger) (*admin.APIClient, string, error) {
			return &admin.APIClient{}, "", nil
		},
	}
	secretRef := &types.NamespacedName{}
	log := zap.S()
	users, err := dbuser.NewAtlasDatabaseUsersService(ctx, provider, secretRef, log)
	require.NoError(t, err)
	assert.Equal(t, &dbuser.AtlasUsers{}, users)
}

func TestFailedNewAtlasDatabaseUsersService(t *testing.T) {
	expectedErr := errors.New("fake error")
	ctx := context.Background()
	provider := &atlas.TestProvider{
		SdkClientFunc: func(_ *client.ObjectKey, _ *zap.SugaredLogger) (*admin.APIClient, string, error) {
			return nil, "", expectedErr
		},
	}
	secretRef := &types.NamespacedName{}
	log := zap.S()
	users, err := dbuser.NewAtlasDatabaseUsersService(ctx, provider, secretRef, log)
	require.Nil(t, users)
	require.ErrorIs(t, err, expectedErr)
}
