package resource

import (
	"context"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

// TODO: (DELETE ME) Add here any errors that need to be identified by consumers of this layer
var (
// // ErrNotFound means an resource is missing
// ErrNotFound = fmt.Errorf("not found")
)

// ResourceService is the interface that consumers use and can mock
type ResourceService interface {
	// TODO: (DELETE ME) this are all fake, actual interface might ot be a clear CRUD like this
	GetResource(ctx context.Context, id string) (*Resource, error)
	CreateResource(ctx context.Context, resource *Resource) (*Resource, error)
	DeleteResource(ctx context.Context, id string) error
	UpdateResource(ctx context.Context, resource *Resource) (*Resource, error)
}

// ProductionResources is the production implementation for the above ResourceService
type ProductionResources struct {
	// Wraps an SDK interface (DO NOT EMBED)
	// Eg:
	// resourceAPI admin.AtlasResourceApi
}

func NewAtlasDatabaseUsersService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (*ProductionResources, error) {
	/*client*/ _, err := translation.NewVersionedClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, err
	}
	return NewProductionResources( /*client.AtlasResourceApi*/ ), nil
}

func NewProductionResources( /*api admin.AtlasResourceApi*/ ) *ProductionResources {
	return &ProductionResources{ /*resourceAPI: api*/ }
}

func (r *ProductionResources) GetResource(ctx context.Context, id string) (*Resource, error) {
	// 1. Call API
	// resp, httpResp, err := si.resourceAPI.GetAtlasResource(ctx, id).Execute()

	// 2. Convert Error as needed
	/*
		if err != nil {
			if httpResp.StatusCode == http.StatusNotFound {
				return nil, errors.Join(err, ErrNotFound)
			}
			return nil, err
		}
	*/

	// 3. Convert the reply from Atlas (convert error as needed) & return
	/*stateInAtlas*/
	_, _ = fromAtlas( /* *resp */ )
	/*
		if err != nil {
			return nil, err
		}
		return stateInAtlas, nil
	*/
	panic("unimplemented")
}

func (r *ProductionResources) CreateResource(ctx context.Context, resource *Resource) (*Resource, error) {
	// 1. Convert to Atlas
	/*atlasResource*/
	_ = resource.toAtlas()
	// 2. Call API
	//resp, httpResp, err := si.resourceAPI.CreateResource(ctx, atlasResource).Execute()
	// 3. Convert Error as needed
	/*
		if err != nil {
			...
			return nil, err
		}
	*/
	// 4. Convert the reply from Atlas (convert error as needed) & return
	/*
		akoResource, err := fromAtlas(*resp)
		if err != nil {
			...
			return nil, err
		}
		return akoResource, nil
	*/
	panic("unimplemented")
}

func (r *ProductionResources) DeleteResource(ctx context.Context, id string) error {
	// 1. Call API
	// _, resp, err := si.resourceAPI.DeleteResource(ctx, id).Execute()
	// 2. Convert Error as needed
	/*
		if err != nil {
			...
			return err
		}
		return nil
	*/
	panic("unimplemented")
}

func (r *ProductionResources) UpdateResource(ctx context.Context, resource *Resource) (*Resource, error) {
	// TODO: (DELETE ME) Update tends to be almost just like create
	// 1. Convert to Atlas
	/*atlasResource, err := resource.toAtlas()
	if err != nil {
		return nil, err
	}*/
	// 2. Call API
	// resp, httpResp, err := si.resourceAPI.UpdateResource(ctx, atlasResource).Execute()
	// 3. Convert Error as needed
	/*
		if err != nil {
			...
			return nil, err
		}
	*/
	// 4. Convert the reply from Atlas (convert error as needed) & return
	/*
		akoResource, err := fromAtlas(*resp)
		if err != nil {
			...
			return nil, err
		}
		return akoResource, nil
	*/
	panic("unimplemented")
}
