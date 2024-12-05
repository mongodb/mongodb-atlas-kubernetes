package encryptionatrest

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

type EncryptionAtRestService interface {
	Get(context.Context, string) (*EncryptionAtRest, error)
	Update(context.Context, string, EncryptionAtRest) error
}

type EncryptionAtRestAPI struct {
	encryptionAtRestAPI admin.EncryptionAtRestUsingCustomerKeyManagementApi
}

func NewEncryptionAtRestAPI(api admin.EncryptionAtRestUsingCustomerKeyManagementApi) *EncryptionAtRestAPI {
	return &EncryptionAtRestAPI{encryptionAtRestAPI: api}
}

func (e *EncryptionAtRestAPI) Get(ctx context.Context, projectID string) (*EncryptionAtRest, error) {
	result, _, err := e.encryptionAtRestAPI.GetEncryptionAtRest(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get encryption at rest from Atlas: %w", err)
	}

	return fromAtlas(result), nil
}

func (e *EncryptionAtRestAPI) Update(ctx context.Context, projectID string, ear EncryptionAtRest) error {
	a := toAtlas(&ear)
	_, _, err := e.encryptionAtRestAPI.UpdateEncryptionAtRest(ctx, projectID, a).Execute()
	if err != nil {
		return fmt.Errorf("failed to update encryption at rest in Atlas: %w", err)
	}

	return nil
}
