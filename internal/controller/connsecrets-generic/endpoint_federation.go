package connsecretsgeneric

// import (
// 	"context"
// 	"fmt"

// 	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
// 	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
// 	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
// 	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
// 	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
// 	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
// 	"go.uber.org/zap"
// 	"k8s.io/apimachinery/pkg/fields"
// 	"sigs.k8s.io/controller-runtime/pkg/client"
// )

// type FederationEndpoint struct {
// 	obj *akov2.AtlasDataFederation
// 	r   *ConnSecretReconciler
// }

// // ---- instance methods ----
// func (e FederationEndpoint) GetName() string {
// 	if e.obj == nil {
// 		return ""
// 	}
// 	return e.obj.GetDFName() // adjust if your CR has a different getter
// }

// func (e FederationEndpoint) IsReady() bool {
// 	return e.obj != nil && api.HasReadyCondition(e.obj.Status.Conditions)
// }

// func (e FederationEndpoint) GetConnStrings() *status.ConnectionStrings {
// 	if e.obj == nil {
// 		return nil
// 	}
// 	return e.obj.Status.ConnectionStrings // or nil if federation doesn’t expose this
// }

// func (e FederationEndpoint) GetProjectID(ctx context.Context, r client.Reader) (string, error) {
// 	if e.obj == nil {
// 		return "", fmt.Errorf("nil federation")
// 	}
// 	if e.obj.Spec.ProjectRef != nil && e.obj.Spec.ProjectRef.Name != "" {
// 		proj := &akov2.AtlasProject{}
// 		if err := r.Get(ctx, kube.ObjectKey(e.obj.Namespace, e.obj.Spec.ProjectRef.Name), proj); err != nil {
// 			return "", err
// 		}
// 		return proj.ID(), nil
// 	}
// 	if id := e.obj.Status.ProjectID; id != "" {
// 		return id, nil
// 	}
// 	return "", fmt.Errorf("project ID not available")
// }

// func (e FederationEndpoint) GetProjectName(ctx context.Context, r client.Reader, provider atlas.Provider, log *zap.SugaredLogger) (string, error) {
// 	if e.obj == nil {
// 		return "", fmt.Errorf("nil federation")
// 	}
// 	if e.obj.Spec.ProjectRef != nil && e.obj.Spec.ProjectRef.Name != "" {
// 		proj := &akov2.AtlasProject{}
// 		if err := r.Get(ctx, kube.ObjectKey(e.obj.Namespace, e.obj.Spec.ProjectRef.Name), proj); err != nil {
// 			return "", err
// 		}
// 		if proj.Spec.Name != "" {
// 			return kube.NormalizeIdentifier(proj.Spec.Name), nil
// 		}
// 	}
// 	// SDK fallback (optional)
// 	if e.r != nil {
// 		cfg, err := e.r.ResolveConnectionConfig(ctx, e.obj)
// 		if err != nil {
// 			return "", err
// 		}
// 		sdk, err := e.r.AtlasProvider.SdkClientSet(ctx, cfg.Credentials, log)
// 		if err != nil {
// 			return "", err
// 		}
// 		ap, err := e.r.ResolveProject(ctx, sdk.SdkClient20250312002, e.obj)
// 		if err != nil {
// 			return "", err
// 		}
// 		return kube.NormalizeIdentifier(ap.Name), nil
// 	}
// 	return "", fmt.Errorf("project name not available")
// }

// // ---- indexer methods ----
// func (FederationEndpoint) ListObj() client.ObjectList { return &akov2.AtlasDataFederationList{} }

// func (FederationEndpoint) Selector(ids *ConnSecretIdentifiers) fields.Selector {
// 	return fields.OneTermEqualSelector(indexer.AtlasDataFederationBySpecNameAndProjectID, ids.ProjectID+"-"+ids.ClusterName)
// }

// func (e FederationEndpoint) ExtractList(ol client.ObjectList) ([]Endpoint, error) {
// 	l, ok := ol.(*akov2.AtlasDataFederationList)
// 	if !ok {
// 		return nil, fmt.Errorf("unexpected list type %T", ol)
// 	}
// 	out := make([]Endpoint, 0, len(l.Items))
// 	for i := range l.Items {
// 		out = append(out, FederationEndpoint{obj: &l.Items[i], r: e.r})
// 	}
// 	return out, nil
// }
