package connsecretsgeneric

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

type DeploymentEndpoint struct {
	obj *akov2.AtlasDeployment
	r   *ConnSecretReconciler
}

// ---- instance methods ----
func (e DeploymentEndpoint) GetName() string {
	if e.obj == nil {
		return ""
	}
	return e.obj.GetDeploymentName()
}

func (e DeploymentEndpoint) IsReady() bool {
	return e.obj != nil && api.HasReadyCondition(e.obj.Status.Conditions)
}

func (e DeploymentEndpoint) GetConnStrings() *status.ConnectionStrings {
	if e.obj == nil {
		return nil
	}
	return e.obj.Status.ConnectionStrings
}

func (e DeploymentEndpoint) GetProjectID(ctx context.Context, r client.Reader) (string, error) {
	if e.obj == nil {
		return "", fmt.Errorf("nil deployment")
	}
	if e.obj.Spec.ExternalProjectRef != nil && e.obj.Spec.ExternalProjectRef.ID != "" {
		return e.obj.Spec.ExternalProjectRef.ID, nil
	}
	if e.obj.Spec.ProjectRef != nil && e.obj.Spec.ProjectRef.Name != "" {
		proj := &akov2.AtlasProject{}
		if err := r.Get(ctx, kube.ObjectKey(e.obj.Namespace, e.obj.Spec.ProjectRef.Name), proj); err != nil {
			return "", err
		}
		return proj.ID(), nil
	}

	return "", fmt.Errorf("project ID not available")
}

func (e DeploymentEndpoint) GetProjectName(ctx context.Context, r client.Reader, provider atlas.Provider, log *zap.SugaredLogger) (string, error) {
	if e.obj == nil {
		return "", fmt.Errorf("nil deployment")
	}
	if e.obj.Spec.ProjectRef != nil && e.obj.Spec.ProjectRef.Name != "" {
		proj := &akov2.AtlasProject{}
		if err := r.Get(ctx, kube.ObjectKey(e.obj.Namespace, e.obj.Spec.ProjectRef.Name), proj); err != nil {
			return "", err
		}
		if proj.Spec.Name != "" {
			return kube.NormalizeIdentifier(proj.Spec.Name), nil
		}
	}
	// SDK fallback (optional)
	if e.r != nil {
		cfg, err := e.r.ResolveConnectionConfig(ctx, e.obj)
		if err != nil {
			return "", err
		}
		sdk, err := e.r.AtlasProvider.SdkClientSet(ctx, cfg.Credentials, log)
		if err != nil {
			return "", err
		}
		ap, err := e.r.ResolveProject(ctx, sdk.SdkClient20250312002, e.obj)
		if err != nil {
			return "", err
		}
		return kube.NormalizeIdentifier(ap.Name), nil
	}
	return "", fmt.Errorf("project name not available")
}

// ---- indexer methods (ignore e.obj) ----
func (DeploymentEndpoint) ListObj() client.ObjectList { return &akov2.AtlasDeploymentList{} }

func (DeploymentEndpoint) Selector(ids *ConnSecretIdentifiers) fields.Selector {
	return fields.OneTermEqualSelector(indexer.AtlasDeploymentBySpecNameAndProjectID, ids.ProjectID+"-"+ids.ClusterName)
}

func (e DeploymentEndpoint) ExtractList(ol client.ObjectList) ([]Endpoint, error) {
	l, ok := ol.(*akov2.AtlasDeploymentList)
	if !ok {
		return nil, fmt.Errorf("unexpected list type %T", ol)
	}
	out := make([]Endpoint, 0, len(l.Items))
	for i := range l.Items {
		// wrap each item as an Endpoint object
		out = append(out, DeploymentEndpoint{obj: &l.Items[i], r: e.r})
	}
	return out, nil
}
