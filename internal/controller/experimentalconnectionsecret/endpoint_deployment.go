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

package experimentalconnectionsecret

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

type DeploymentEndpoint struct {
	obj             *akov2.AtlasDeployment
	k8s             client.Client
	provider        atlas.Provider
	globalSecretRef client.ObjectKey
	log             *zap.SugaredLogger
}

// resolveProjectIDByKey returns the project id from the key
func resolveProjectIDByKey(ctx context.Context, c client.Client, key client.ObjectKey) (string, error) {
	proj := &akov2.AtlasProject{}
	if err := c.Get(ctx, key, proj); err != nil {
		return "", err
	}
	if proj.ID() == "" {
		return "", ErrUnresolvedProjectID
	}
	return proj.ID(), nil
}

// resolveProjectNameByKey returns the project name from the key
func resolveProjectNameByKey(ctx context.Context, c client.Client, key client.ObjectKey) (string, error) {
	proj := &akov2.AtlasProject{}
	if err := c.Get(ctx, key, proj); err != nil {
		return "", err
	}
	if proj.Spec.Name == "" {
		return "", ErrUnresolvedProjectName
	}
	return proj.Spec.Name, nil
}

// resolveProjectNameBySDK returns the project name from SDL
func resolveProjectNameBySDK(
	ctx context.Context,
	c client.Client,
	provider atlas.Provider,
	log *zap.SugaredLogger,
	globalSecretRef client.ObjectKey,
	referrer project.ProjectReferrerObject,
) (string, error) {
	pdr := referrer.ProjectDualRef()

	var secretRef *client.ObjectKey
	if pdr.ConnectionSecret != nil && pdr.ConnectionSecret.Name != "" {
		if obj, ok := any(referrer).(client.Object); ok {
			key := client.ObjectKeyFromObject(obj)
			key.Name = pdr.ConnectionSecret.Name
			secretRef = &key
		} else {
			key := client.ObjectKey{Namespace: referrer.GetNamespace(), Name: pdr.ConnectionSecret.Name}
			secretRef = &key
		}
	}

	cfg, err := reconciler.GetConnectionConfig(ctx, c, secretRef, &globalSecretRef)
	if err != nil {
		return "", err
	}

	if pdr.ExternalProjectRef == nil || pdr.ExternalProjectRef.ID == "" {
		return "", ErrUnresolvedProjectName
	}

	cs, err := provider.SdkClientSet(ctx, cfg.Credentials, log)
	if err != nil {
		return "", err
	}

	svc := project.NewProjectAPIService(cs.SdkClient20250312006.ProjectsApi)
	prj, err := svc.GetProject(ctx, pdr.ExternalProjectRef.ID)
	if err != nil {
		return "", err
	}
	if prj.Name == "" {
		return "", ErrUnresolvedProjectName
	}
	return prj.Name, nil
}

// GetName resolves the endpoints name from the spec
func (e DeploymentEndpoint) GetName() string {
	if e.obj == nil {
		return ""
	}
	return e.obj.GetDeploymentName()
}

// IsReady returns true if the endpoint is ready
func (e DeploymentEndpoint) IsReady() bool {
	return e.obj != nil && api.HasReadyCondition(e.obj.Status.Conditions)
}

// GetScopeType returns the scope type of the endpoint to match with the ones from AtlasDatabaseUser
func (e DeploymentEndpoint) GetScopeType() akov2.ScopeType {
	return akov2.DeploymentScopeType
}

// GetProjectID resolves parent project's id (ProjectRef or ExternalRef)
func (e DeploymentEndpoint) GetProjectID(ctx context.Context) (string, error) {
	if e.obj == nil {
		return "", fmt.Errorf("nil deployment")
	}
	if e.obj.Spec.ExternalProjectRef != nil && e.obj.Spec.ExternalProjectRef.ID != "" {
		return e.obj.Spec.ExternalProjectRef.ID, nil
	}
	if e.obj.Spec.ProjectRef != nil && e.obj.Spec.ProjectRef.Name != "" {
		return resolveProjectIDByKey(ctx, e.k8s, e.obj.AtlasProjectObjectKey())
	}
	return "", ErrUnresolvedProjectID
}

// GetProjectName returns the parent project's name (either by getting K8s AtlasProject or SDK calls)
func (e DeploymentEndpoint) GetProjectName(ctx context.Context) (string, error) {
	if e.obj == nil {
		return "", fmt.Errorf("nil deployment")
	}
	if e.obj.Spec.ProjectRef != nil && e.obj.Spec.ProjectRef.Name != "" {
		return resolveProjectNameByKey(ctx, e.k8s, e.obj.AtlasProjectObjectKey())
	}
	if e.obj.Spec.ConnectionSecret != nil && e.obj.Spec.ConnectionSecret.Name != "" {
		return resolveProjectNameBySDK(ctx, e.k8s, e.provider, e.log, e.globalSecretRef, e.obj)
	}
	return "", ErrUnresolvedProjectName
}

// Defines the list type
func (DeploymentEndpoint) ListObj() client.ObjectList { return &akov2.AtlasDeploymentList{} }

// Defines the selector to use for indexer when trying to retrieve all endpoints by project
func (DeploymentEndpoint) SelectorByProject(projectID string) fields.Selector {
	return fields.OneTermEqualSelector(indexer.AtlasDeploymentByProject, projectID)
}

// Defines the selector to use for indexer when trying to retrieve all endpoints by project and spec name
func (DeploymentEndpoint) SelectorByProjectAndName(ids *ConnSecretIdentifiers) fields.Selector {
	return fields.OneTermEqualSelector(indexer.AtlasDeploymentBySpecNameAndProjectID, ids.ProjectID+"-"+ids.ClusterName)
}

// ExtractList creates a list of Endpoint types to preserve the abstraction
func (e DeploymentEndpoint) ExtractList(ol client.ObjectList) ([]Endpoint, error) {
	l, ok := ol.(*akov2.AtlasDeploymentList)
	if !ok {
		return nil, fmt.Errorf("unexpected list type %T", ol)
	}
	out := make([]Endpoint, 0, len(l.Items))
	for i := range l.Items {
		out = append(out, DeploymentEndpoint{
			obj:             &l.Items[i],
			k8s:             e.k8s,
			provider:        e.provider,
			globalSecretRef: e.globalSecretRef,
			log:             e.log,
		})
	}
	return out, nil
}

// BuildConnData defines the specific function/way for building the ConnSecretData given this type of endpoint
// AtlasDeployment stores connection strings in the status field
func (e DeploymentEndpoint) BuildConnData(ctx context.Context, user *akov2.AtlasDatabaseUser) (ConnSecretData, error) {
	if user == nil || e.obj == nil {
		return ConnSecretData{}, ErrMissingPairing
	}
	password, err := user.ReadPassword(ctx, e.k8s)
	if err != nil {
		return ConnSecretData{}, fmt.Errorf("failed to read password for user %q: %w", user.Spec.Username, err)
	}
	data := ConnSecretData{
		DBUserName: user.Spec.Username,
		Password:   password,
	}

	if e.obj.Status.ConnectionStrings == nil {
		return data, nil
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
