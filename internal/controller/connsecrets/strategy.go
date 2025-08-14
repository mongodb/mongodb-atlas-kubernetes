package connsecrets

// import (
// 	"context"
// 	"errors"
// 	"fmt"

// 	"k8s.io/apimachinery/pkg/fields"
// 	"sigs.k8s.io/controller-runtime/pkg/client"

// 	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
// 	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
// 	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/stringutil"
// )

// const InternalSeparator = "$"

// var (
// 	ErrNoPairedResourcesFound = errors.New("no endpoint and no AtlasDatabaseUser found")
// 	ErrNoEndpointFound        = errors.New("no endpoint found")
// 	ErrManyEndpoints          = errors.New("multiple endpoints found")
// 	ErrNoUserFound            = errors.New("no AtlasDatabaseUser found")
// 	ErrManyUsers              = errors.New("multiple AtlasDatabaseUsers found")
// )

// type AnyEndpointStrategy interface {
// 	LoadPair(ctx context.Context, c client.Client, ids *ConnSecretIdentifiers) (*ConnSecretPair[any], error)
// 	Ready(p *ConnSecretPair[any]) bool
// 	ValidScopes(p *ConnSecretPair[any]) bool
// 	BuildConnectionData(ctx context.Context, c client.Client, p *ConnSecretPair[any]) (ConnSecretData, error)
// 	ResolveProjectName(ctx context.Context, p *ConnSecretPair[any]) (string, error)
// }

// type anyEndpointStrategy[T any] struct {
// 	EndpointStrategy[T]
// }

// type ConnSecretPair[T any] struct {
// 	ProjectID string
// 	User      *akov2.AtlasDatabaseUser
// 	Endpoint  T
// }

// type ConnSecretData struct {
// 	DBUserName      string
// 	Password        string
// 	ConnURL         string
// 	SrvConnURL      string
// 	PrivateConnURLs []PrivateLinkConnURLs
// }

// type PrivateLinkConnURLs struct {
// 	PvtConnURL      string
// 	PvtSrvConnURL   string
// 	PvtShardConnURL string
// }

// func NewAnyEndpointStrategy[T any](s EndpointStrategy[T]) AnyEndpointStrategy {
// 	return &anyEndpointStrategy[T]{s}
// }

// func (w *anyEndpointStrategy[T]) LoadPair(ctx context.Context, c client.Client, ids *ConnSecretIdentifiers) (*ConnSecretPair[any], error) {
// 	if err := c.List(ctx, w.List, &client.ListOptions{FieldSelector: w.Selector(ids)}); err != nil {
// 		return nil, err
// 	}
// 	eps, err := w.ExtractList(w.List)
// 	if err != nil {
// 		return nil, err
// 	}

// 	users := &akov2.AtlasDatabaseUserList{}
// 	userSel := fields.OneTermEqualSelector(indexer.AtlasDatabaseUserBySpecUsernameAndProjectID, ids.ProjectID+"-"+ids.DatabaseUsername)
// 	if err := c.List(ctx, users, &client.ListOptions{FieldSelector: userSel}); err != nil {
// 		return nil, err
// 	}

// 	switch {
// 	case len(eps) == 0 && len(users.Items) == 0:
// 		return nil, ErrNoPairedResourcesFound
// 	case len(eps) == 0:
// 		return &ConnSecretPair[any]{ProjectID: ids.ProjectID, User: &users.Items[0], Endpoint: nil}, ErrNoEndpointFound
// 	case len(users.Items) == 0:
// 		return &ConnSecretPair[any]{ProjectID: ids.ProjectID, User: nil, Endpoint: eps[0]}, ErrNoUserFound
// 	case len(eps) > 1:
// 		return nil, ErrManyEndpoints
// 	case len(users.Items) > 1:
// 		return nil, ErrManyUsers
// 	}

// 	return &ConnSecretPair[any]{ProjectID: ids.ProjectID, User: &users.Items[0], Endpoint: eps[0]}, nil
// }

// func (w *anyEndpointStrategy[T]) ValidScopes(p *ConnSecretPair[any]) bool {
// 	if p == nil || p.User == nil {
// 		return false
// 	}
// 	scopes := p.User.GetScopes(akov2.DeploymentScopeType)
// 	if len(scopes) == 0 {
// 		return true
// 	}
// 	t, ok := p.Endpoint.(T)
// 	if !ok || p.Endpoint == nil {
// 		return false
// 	}
// 	name := w.GetName(t)
// 	if name == "" {
// 		return false
// 	}
// 	return stringutil.Contains(scopes, name)
// }

// func (w *anyEndpointStrategy[T]) Ready(p *ConnSecretPair[any]) bool {
// 	if p == nil || p.User == nil || !p.User.IsDatabaseUserReady() {
// 		return false
// 	}
// 	t, ok := p.Endpoint.(T)
// 	if !ok || p.Endpoint == nil {
// 		return false
// 	}
// 	return w.IsReady(t)
// }

// func (w *anyEndpointStrategy[T]) BuildConnectionData(ctx context.Context, c client.Client, p *ConnSecretPair[any]) (ConnSecretData, error) {
// 	if p == nil || p.User == nil || p.Endpoint == nil {
// 		return ConnSecretData{}, fmt.Errorf("invalid pair: nil user or endpoint")
// 	}

// 	password, err := p.User.ReadPassword(ctx, c)
// 	if err != nil {
// 		return ConnSecretData{}, fmt.Errorf("failed to read password for user %q: %w", p.User.Spec.Username, err)
// 	}

// 	t, ok := p.Endpoint.(T)
// 	if !ok {
// 		return ConnSecretData{}, fmt.Errorf("unexpected endpoint type")
// 	}

// 	conn := w.GetConnStrings(t)

// 	data := ConnSecretData{
// 		DBUserName: p.User.Spec.Username,
// 		Password:   password,
// 		ConnURL:    conn.Standard,
// 		SrvConnURL: conn.StandardSrv,
// 	}

// 	if conn.Private != "" {
// 		data.PrivateConnURLs = append(data.PrivateConnURLs, PrivateLinkConnURLs{
// 			PvtConnURL:    conn.Private,
// 			PvtSrvConnURL: conn.PrivateSrv,
// 		})
// 	}

// 	for _, pe := range conn.PrivateEndpoint {
// 		data.PrivateConnURLs = append(data.PrivateConnURLs, PrivateLinkConnURLs{
// 			PvtConnURL:      pe.ConnectionString,
// 			PvtSrvConnURL:   pe.SRVConnectionString,
// 			PvtShardConnURL: pe.SRVShardOptimizedConnectionString,
// 		})
// 	}

// 	return data, nil
// }

// func (a *anyEndpointStrategy[T]) ResolveProjectName(ctx context.Context, p *ConnSecretPair[any]) (string, error) {
// 	t, ok := p.Endpoint.(T)
// 	if !ok || p.Endpoint == nil {
// 		return "", fmt.Errorf("unexpected endpoint type")
// 	}
// 	return a.GetProjectName(ctx, t), nil
// }
