// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package connsecretsgeneric

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
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

func (e DeploymentEndpoint) GetProjectRef(ctx context.Context) string {
	return "PROJECTID"
}

func (e DeploymentEndpoint) GetProjectID(ctx context.Context) (string, error) {
	if e.obj == nil {
		return "", fmt.Errorf("nil deployment")
	}
	if e.obj.Spec.ExternalProjectRef != nil && e.obj.Spec.ExternalProjectRef.ID != "" {
		return e.obj.Spec.ExternalProjectRef.ID, nil
	}
	if e.obj.Spec.ProjectRef != nil && e.obj.Spec.ProjectRef.Name != "" {
		proj := &akov2.AtlasProject{}
		if err := e.r.Client.Get(ctx, kube.ObjectKey(e.obj.Namespace, e.obj.Spec.ProjectRef.Name), proj); err != nil {
			return "", err
		}
		return proj.ID(), nil
	}

	return "", fmt.Errorf("project ID not available")
}

func (e DeploymentEndpoint) GetProjectName(ctx context.Context) (string, error) {
	if e.obj == nil {
		return "", fmt.Errorf("nil deployment")
	}
	if e.obj.Spec.ProjectRef != nil && e.obj.Spec.ProjectRef.Name != "" {
		proj := &akov2.AtlasProject{}
		if err := e.r.Client.Get(ctx, kube.ObjectKey(e.obj.Namespace, e.obj.Spec.ProjectRef.Name), proj); err != nil {
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
		sdk, err := e.r.AtlasProvider.SdkClientSet(ctx, cfg.Credentials, e.r.Log)
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

// ---- indexer methods ----
func (DeploymentEndpoint) ListObj() client.ObjectList { return &akov2.AtlasDeploymentList{} }

func (DeploymentEndpoint) SelectorByProject(projectRef string) fields.Selector {
	return fields.OneTermEqualSelector(indexer.AtlasDeploymentByProject, projectRef)
}

func (DeploymentEndpoint) SelectorByProjectAndName(ids *ConnSecretIdentifiers) fields.Selector {
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

func (e DeploymentEndpoint) BuildConnData(ctx context.Context, user *akov2.AtlasDatabaseUser) (ConnSecretData, error) {
	if user == nil || e.obj == nil {
		return ConnSecretData{}, fmt.Errorf("invalid endpoint or user")
	}
	password, err := user.ReadPassword(ctx, e.r.Client)
	if err != nil {
		return ConnSecretData{}, fmt.Errorf("failed to read password for user %q: %w", user.Spec.Username, err)
	}
	data := ConnSecretData{
		DBUserName: user.Spec.Username,
		Password:   password,
	}

	conn := e.obj.Status.ConnectionStrings
	data.ConnURL = conn.Standard
	data.SrvConnURL = conn.StandardSrv
	if conn.Private != "" {
		data.PrivateConnURLs = append(data.PrivateConnURLs, PrivateLinkConnURLs{
			PvtConnURL:    conn.Private,
			PvtSrvConnURL: conn.PrivateSrv,
		})
	}
	for _, pe := range conn.PrivateEndpoint {
		data.PrivateConnURLs = append(data.PrivateConnURLs, PrivateLinkConnURLs{
			PvtConnURL:      pe.ConnectionString,
			PvtSrvConnURL:   pe.SRVConnectionString,
			PvtShardConnURL: pe.SRVShardOptimizedConnectionString,
		})
	}

	return data, nil
}
