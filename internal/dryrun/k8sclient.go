package dryrun

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client struct {
	origin                       client.Client
	subResourceClientConstructor *subResourceClientConstructor
	statusWriter                 *statusWriter
}

// NewClient returns a controller-runtime client that will passes read-only operations to the given client
// or returns a DryRunError in other cases
func NewClient(origin client.Client) *Client {
	return &Client{
		origin:                       origin,
		subResourceClientConstructor: &subResourceClientConstructor{origin: origin},
		statusWriter:                 &statusWriter{},
	}
}

func (r *Client) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	return r.origin.Get(ctx, key, obj, opts...)
}

func (r *Client) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return r.origin.List(ctx, list, opts...)
}

func (r *Client) Create(_ context.Context, obj client.Object, _ ...client.CreateOption) error {
	return NewDryRunError(obj.GetObjectKind(), obj.(metav1.ObjectMetaAccessor), "Normal", DryRunReason, "Would CREATE %s", obj.GetName())
}

func (r *Client) Delete(_ context.Context, obj client.Object, _ ...client.DeleteOption) error {
	return NewDryRunError(obj.GetObjectKind(), obj.(metav1.ObjectMetaAccessor), "Normal", DryRunReason, "Would DELETE %s", obj.GetName())
}

func (r *Client) Update(_ context.Context, obj client.Object, _ ...client.UpdateOption) error {
	return NewDryRunError(obj.GetObjectKind(), obj.(metav1.ObjectMetaAccessor), "Normal", DryRunReason, "Would UPDATE %s", obj.GetName())
}

func (r *Client) Patch(_ context.Context, obj client.Object, _ client.Patch, _ ...client.PatchOption) error {
	return NewDryRunError(obj.GetObjectKind(), obj.(metav1.ObjectMetaAccessor), "Normal", DryRunReason, "Would PATCH %s", obj.GetName())
}

func (r *Client) DeleteAllOf(_ context.Context, obj client.Object, _ ...client.DeleteAllOfOption) error {
	return NewDryRunError(obj.GetObjectKind(), obj.(metav1.ObjectMetaAccessor), "Normal", DryRunReason, "Would DELETE ALL OF %s", obj.GetName())
}

func (r *Client) Status() client.SubResourceWriter {
	return r.statusWriter
}

func (r *Client) SubResource(subResource string) client.SubResourceClient {
	return r.subResourceClientConstructor.SubResource(subResource)
}

func (r *Client) Scheme() *runtime.Scheme {
	return r.origin.Scheme()
}

func (r *Client) RESTMapper() meta.RESTMapper {
	return r.origin.RESTMapper()
}

func (r *Client) GroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	return r.origin.GroupVersionKindFor(obj)
}

func (r *Client) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	return r.origin.IsObjectNamespaced(obj)
}

type subresourceClient struct {
	origin      client.SubResourceClient
	subResource string
}

type subResourceClientConstructor struct {
	origin client.Client
}

func (s *subResourceClientConstructor) SubResource(subResource string) client.SubResourceClient {
	return newSubresourceClient(s.origin.SubResource(subResource), subResource)
}

func newSubresourceClient(origin client.SubResourceClient, subResource string) *subresourceClient {
	return &subresourceClient{
		origin:      origin,
		subResource: subResource,
	}
}

func (s *subresourceClient) Get(ctx context.Context, obj client.Object, subResource client.Object, opts ...client.SubResourceGetOption) error {
	return s.origin.Get(ctx, obj, subResource, opts...)
}

func (s *subresourceClient) Create(ctx context.Context, obj client.Object, _ client.Object, _ ...client.SubResourceCreateOption) error {
	return NewDryRunError(obj.GetObjectKind(), obj.(metav1.ObjectMetaAccessor), "Normal", DryRunReason, "Would CREATE sub-resource %s", obj.GetName())
}

func (s *subresourceClient) Update(ctx context.Context, obj client.Object, _ ...client.SubResourceUpdateOption) error {
	return NewDryRunError(obj.GetObjectKind(), obj.(metav1.ObjectMetaAccessor), "Normal", DryRunReason, "Would UPDATE sub-resource %s", obj.GetName())
}

func (s *subresourceClient) Patch(ctx context.Context, obj client.Object, _ client.Patch, _ ...client.SubResourcePatchOption) error {
	return NewDryRunError(obj.GetObjectKind(), obj.(metav1.ObjectMetaAccessor), "Normal", DryRunReason, "Would PATCH sub-resource %s", obj.GetName())
}

type statusWriter struct{}

func (s *statusWriter) Create(ctx context.Context, obj client.Object, _ client.Object, _ ...client.SubResourceCreateOption) error {
	return NewDryRunError(obj.GetObjectKind(), obj.(metav1.ObjectMetaAccessor), "Normal", DryRunReason, "Would CREATE status %s", obj.GetName())
}

func (s *statusWriter) Update(ctx context.Context, obj client.Object, _ ...client.SubResourceUpdateOption) error {
	return NewDryRunError(obj.GetObjectKind(), obj.(metav1.ObjectMetaAccessor), "Normal", DryRunReason, "Would UPDATE status %s", obj.GetName())
}

func (s *statusWriter) Patch(ctx context.Context, obj client.Object, _ client.Patch, _ ...client.SubResourcePatchOption) error {
	return NewDryRunError(obj.GetObjectKind(), obj.(metav1.ObjectMetaAccessor), "Normal", DryRunReason, "Would PATCH status %s", obj.GetName())
}
