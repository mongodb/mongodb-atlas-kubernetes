package atlasdatafederation

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"go.mongodb.org/atlas/mongodbatlas"
)

type DataFederationServiceOp service

const (
	dataFederationBasePath = "https://cloud-qa.mongodb.com/api/atlas/v1.0/groups" // todo: fix cloud-qa later !!
)

func NewClient(client mongodbatlas.Client) *DataFederationServiceOp {
	return &DataFederationServiceOp{
		Client: client,
	}
}

func (s *DataFederationServiceOp) Get(ctx context.Context, groupID string, tenantName string) (*mdbv1.DataFederationSpec, *mongodbatlas.Response, error) {
	if groupID == "" {
		return nil, nil, errors.New("groupID must be set")
	}
	if tenantName == "" {
		return nil, nil, errors.New("tenantName must be set")
	}

	path := fmt.Sprintf("%s/%s/dataFederation/%s", dataFederationBasePath, groupID, tenantName)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(mdbv1.DataFederationSpec)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

func (s *DataFederationServiceOp) Create(ctx context.Context, groupID string, spec *mdbv1.DataFederationSpec) (*mdbv1.DataFederationSpec, *mongodbatlas.Response, error) {
	if groupID == "" {
		return nil, nil, errors.New("groupID must be set")
	}

	path := fmt.Sprintf("%s/%s/dataFederation", dataFederationBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, spec)
	if err != nil {
		return nil, nil, err
	}

	root := new(mdbv1.DataFederationSpec)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

func (s *DataFederationServiceOp) Update(ctx context.Context, groupID string, spec *mdbv1.DataFederationSpec) (*mdbv1.DataFederationSpec, *mongodbatlas.Response, error) {
	if groupID == "" {
		return nil, nil, errors.New("groupID must be set")
	}

	path := fmt.Sprintf("%s/%s/dataFederation/%s", dataFederationBasePath, groupID, spec.Name)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, spec)
	if err != nil {
		return nil, nil, err
	}

	root := new(mdbv1.DataFederationSpec)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

func (s *DataFederationServiceOp) Delete(ctx context.Context, groupID string, tenantName string) (*mdbv1.DataFederationSpec, *mongodbatlas.Response, error) {
	if groupID == "" {
		return nil, nil, errors.New("groupID must be set")
	}
	if tenantName == "" {
		return nil, nil, errors.New("tenantName must be set")
	}

	path := fmt.Sprintf("%s/%s/dataFederation/%s", dataFederationBasePath, groupID, tenantName)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(mdbv1.DataFederationSpec)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

func (s *DataFederationServiceOp) GetAllPrivateEndpoints(ctx context.Context, groupID string) ([]mdbv1.DataFederationPE, *mongodbatlas.Response, error) {
	if groupID == "" {
		return nil, nil, errors.New("groupID must be set")
	}

	path := fmt.Sprintf("%s/%s/privateNetworkSettings/endpointIds", dataFederationBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new([]mdbv1.DataFederationPE)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return *root, resp, nil
}

func (s *DataFederationServiceOp) CreateOnePrivateEndpoint(ctx context.Context, groupID string, endpoint mdbv1.DataFederationPE) (*mdbv1.DataFederationPE, *mongodbatlas.Response, error) {
	if groupID == "" {
		return nil, nil, errors.New("groupID must be set")
	}

	path := fmt.Sprintf("%s/%s/privateNetworkSettings/endpointIds", dataFederationBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, endpoint)
	if err != nil {
		return nil, nil, err
	}

	root := new(mdbv1.DataFederationPE)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

func (s *DataFederationServiceOp) DeleteOnePrivateEndpoint(ctx context.Context, groupID string, endpointID string) (*mdbv1.DataFederationPE, *mongodbatlas.Response, error) {
	if groupID == "" {
		return nil, nil, errors.New("groupID must be set")
	}
	if endpointID == "" {
		return nil, nil, errors.New("endpointID must be set")
	}

	path := fmt.Sprintf("%s/%s/privateNetworkSettings/endpointIds/%s", dataFederationBasePath, groupID, endpointID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(mdbv1.DataFederationPE)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

type service struct {
	Client mongodbatlas.Client
}
