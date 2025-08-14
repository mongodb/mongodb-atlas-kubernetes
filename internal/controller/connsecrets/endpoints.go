package connsecrets

// import (
// 	"context"
// 	"fmt"

// 	"k8s.io/apimachinery/pkg/fields"
// 	"k8s.io/apimachinery/pkg/types"
// 	"sigs.k8s.io/controller-runtime/pkg/client"

// 	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
// 	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
// 	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
// 	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
// )

// type EndpointStrategy[T any] struct {
// 	List        client.ObjectList
// 	Selector    func(ids *ConnSecretIdentifiers) fields.Selector
// 	ExtractList func(client.ObjectList) ([]T, error)

// 	GetName        func(obj T) string
// 	IsReady        func(obj T) bool
// 	GetConnStrings func(obj T) *status.ConnectionStrings
// 	GetProjectID   func(ctx context.Context, obj T) string
// 	GetProjectName func(ctx context.Context, obj T) string
// }

// // NewDeploymentEndpoint returns the EndpointStrategy for AtlasDeployment.
// func (r *ConnSecretReconciler) NewDeploymentEndpoint() EndpointStrategy[*akov2.AtlasDeployment] {
// 	return EndpointStrategy[*akov2.AtlasDeployment]{
// 		List: &akov2.AtlasDeploymentList{},
// 		Selector: func(ids *ConnSecretIdentifiers) fields.Selector {
// 			return fields.OneTermEqualSelector(
// 				indexer.AtlasDeploymentBySpecNameAndProjectID,
// 				ids.ProjectID+"-"+ids.ClusterName,
// 			)
// 		},
// 		ExtractList: func(ol client.ObjectList) ([]*akov2.AtlasDeployment, error) {
// 			l, ok := ol.(*akov2.AtlasDeploymentList)
// 			if !ok {
// 				return nil, fmt.Errorf("unexpected list type %T", ol)
// 			}
// 			out := make([]*akov2.AtlasDeployment, 0, len(l.Items))
// 			for i := range l.Items {
// 				out = append(out, &l.Items[i])
// 			}
// 			return out, nil
// 		},
// 		GetName: func(dpl *akov2.AtlasDeployment) string {
// 			return dpl.GetDeploymentName()
// 		},
// 		IsReady: func(dpl *akov2.AtlasDeployment) bool {
// 			return api.HasReadyCondition(dpl.Status.Conditions)
// 		},
// 		GetConnStrings: func(dpl *akov2.AtlasDeployment) *status.ConnectionStrings {
// 			return dpl.Status.ConnectionStrings
// 		},
// 		GetProjectID: func(ctx context.Context, dpl *akov2.AtlasDeployment) string {
// 			if dpl.Spec.ExternalProjectRef != nil && dpl.Spec.ExternalProjectRef.ID != "" {
// 				return dpl.Spec.ExternalProjectRef.ID
// 			}
// 			if dpl.Spec.ProjectRef != nil && dpl.Spec.ProjectRef.Name != "" {
// 				ns := dpl.Spec.ProjectRef.Namespace
// 				var proj akov2.AtlasProject
// 				if err := r.Client.Get(ctx, types.NamespacedName{Namespace: ns, Name: dpl.Spec.ProjectRef.Name}, &proj); err == nil {
// 					return proj.ID()
// 				}
// 			}
// 			return ""
// 		},
// 		GetProjectName: func(ctx context.Context, dpl *akov2.AtlasDeployment) string {
// 			// Prefer K8s project name when ProjectRef is present
// 			if dpl.Spec.ProjectRef != nil && dpl.Spec.ProjectRef.Name != "" {
// 				ns := dpl.Spec.ProjectRef.Namespace
// 				var proj akov2.AtlasProject
// 				if err := r.Client.Get(ctx, types.NamespacedName{Namespace: ns, Name: dpl.Spec.ProjectRef.Name}, &proj); err == nil && proj.Spec.Name != "" {
// 					return proj.Spec.Name
// 				}
// 			}

// 			// SDK fallback
// 			connCfg, err := r.ResolveConnectionConfig(ctx, dpl)
// 			if err != nil {
// 				return ""
// 			}
// 			sdkClientSet, err := r.AtlasProvider.SdkClientSet(ctx, connCfg.Credentials, r.Log)
// 			if err != nil {
// 				return ""
// 			}
// 			ap, err := r.ResolveProject(ctx, sdkClientSet.SdkClient20250312002, dpl)
// 			if err != nil {
// 				return ""
// 			}
// 			return ap.Name
// 		},
// 	}
// }
