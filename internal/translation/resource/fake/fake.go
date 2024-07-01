package fake

import (
	"context"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/resource"
)

type FakeResource struct {
	GetResourceFn    func(ctx context.Context, id string) (*resource.Resource, error)
	CreateResourceFn func(ctx context.Context, resource *resource.Resource) (*resource.Resource, error)
	DeleteResourceFn func(ctx context.Context, id string) error
	UpdateResourceFn func(ctx context.Context, resource *resource.Resource) (*resource.Resource, error)
}

func (frs *FakeResource) GetIndex(ctx context.Context, id string) (*resource.Resource, error) {
	if frs.GetResourceFn == nil {
		panic("unimplemented")
	}
	return frs.GetResourceFn(ctx, id)
}

func (frs *FakeResource) CreateIndex(ctx context.Context, resource *resource.Resource) (*resource.Resource, error) {
	if frs.CreateResourceFn == nil {
		panic("unimplemented")
	}
	return frs.CreateResourceFn(ctx, resource)
}

func (frs *FakeResource) DeleteIndex(ctx context.Context, id string) error {
	if frs.DeleteResourceFn == nil {
		panic("unimplemented")
	}
	return frs.DeleteResourceFn(ctx, id)
}

func (frs *FakeResource) UpdateIndex(ctx context.Context, resource *resource.Resource) (*resource.Resource, error) {
	if frs.UpdateResourceFn == nil {
		panic("unimplemented")
	}
	return frs.UpdateResourceFn(ctx, resource)
}
