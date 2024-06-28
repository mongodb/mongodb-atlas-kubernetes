package ipaccesslist

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

type IPAccessListService interface {
	List(ctx context.Context, projectID string) (IPAccessEntries, error)
	Add(ctx context.Context, projectID string, entries IPAccessEntries) error
	Delete(ctx context.Context, projectID string, entry *IPAccessEntry) error
	Status(ctx context.Context, projectID string, entry *IPAccessEntry) (string, error)
}

type IPAccessList struct {
	ipAccessListAPI admin.ProjectIPAccessListApi
}

func (i *IPAccessList) List(ctx context.Context, projectID string) (IPAccessEntries, error) {
	netPermResult, _, err := i.ipAccessListAPI.ListProjectIpAccessLists(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get ip access list from Atlas: %w", err)
	}

	return fromAtlas(netPermResult.GetResults()), nil
}

func (i *IPAccessList) Add(ctx context.Context, projectID string, entries IPAccessEntries) error {
	_, _, err := i.ipAccessListAPI.CreateProjectIpAccessList(ctx, projectID, toAtlas(entries)).Execute()
	if err != nil {
		return fmt.Errorf("failed to create ip access list from Atlas: %w", err)
	}

	return nil
}

func (i *IPAccessList) Delete(ctx context.Context, projectID string, entry *IPAccessEntry) error {
	_, _, err := i.ipAccessListAPI.DeleteProjectIpAccessList(ctx, projectID, entry.ID()).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete ip access list from Atlas: %w", err)
	}

	return nil
}

func (i *IPAccessList) Status(ctx context.Context, projectID string, entry *IPAccessEntry) (string, error) {
	result, _, err := i.ipAccessListAPI.GetProjectIpAccessListStatus(ctx, projectID, entry.ID()).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to get status of ip access list from Atlas: %w", err)
	}

	return result.GetSTATUS(), nil
}

func NewIPAccessList(api admin.ProjectIPAccessListApi) *IPAccessList {
	return &IPAccessList{
		ipAccessListAPI: api,
	}
}
