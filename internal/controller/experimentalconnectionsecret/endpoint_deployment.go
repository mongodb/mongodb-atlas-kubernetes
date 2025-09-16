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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
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

func (e DeploymentEndpoint) GetConnectionType() string {
	if e.obj == nil {
		return ""
	}
	return "deployment"
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
	// Step 1: Log basic input details
	if user == nil || e.obj == nil {
		e.log.Errorw("BuildConnData called with nil Deployment or user",
			"DeploymentEndpoint", e.obj,
			"AtlasDatabaseUser", user,
		)
		return ConnSecretData{}, ErrMissingPairing
	}

	e.log.Debugw("Starting BuildConnData",
		"Username", user.Spec.Username,
		"DeploymentName", e.obj.GetDeploymentName(),
	)

	// Step 2: Read the user's password and log outcome
	password, err := user.ReadPassword(ctx, e.k8s)
	if err != nil {
		e.log.Errorw("Failed to read password for user",
			"Username", user.Spec.Username,
			"Error", err,
		)
		return ConnSecretData{}, fmt.Errorf("failed to read password for user %q: %w", user.Spec.Username, err)
	}

	e.log.Debugw("Successfully read password for user",
		"Username", user.Spec.Username,
		"PasswordPresent", len(password) > 0,
	)

	// Initialize ConnSecretData with DBUserName and Password
	data := ConnSecretData{
		DBUserName: user.Spec.Username,
		Password:   password,
	}

	// Step 3: Check and handle connection strings
	if e.obj.Status.ConnectionStrings == nil {
		e.log.Warn("ConnectionStrings is nil for Deployment", "DeploymentName", e.obj.GetDeploymentName())
		return data, nil
	}

	e.log.Debugw("ConnectionStrings found for Deployment",
		"DeploymentName", e.obj.GetDeploymentName(),
		"ConnectionStrings", e.obj.Status.ConnectionStrings,
	)

	conn := e.obj.Status.ConnectionStrings

	// Standard and SRV connection strings
	data.ConnURL = conn.Standard
	data.SrvConnURL = conn.StandardSrv

	e.log.Debugw("Standard/SRV connection strings",
		"StandardConnURL", data.ConnURL,
		"SrvConnURL", data.SrvConnURL,
	)

	// Private connection strings
	if conn.Private != "" {
		e.log.Debugw("Private connection string detected",
			"PrivateConnURL", conn.Private,
			"PrivateSrvConnURL", conn.PrivateSrv,
		)
		data.PrivateConnURLs = append(data.PrivateConnURLs, PrivateLinkConnURLs{
			PvtConnURL:    conn.Private,
			PvtSrvConnURL: conn.PrivateSrv,
		})
	}

	// Iterate through PrivateEndpoint connection strings
	for _, pe := range conn.PrivateEndpoint {
		e.log.Debugw("PrivateEndpoint connection string detected",
			"PvtConnURL", pe.ConnectionString,
			"PvtSrvConnURL", pe.SRVConnectionString,
			"PvtShardConnURL", pe.SRVShardOptimizedConnectionString,
		)
		data.PrivateConnURLs = append(data.PrivateConnURLs, PrivateLinkConnURLs{
			PvtConnURL:      pe.ConnectionString,
			PvtSrvConnURL:   pe.SRVConnectionString,
			PvtShardConnURL: pe.SRVShardOptimizedConnectionString,
		})
	}

	// Step 4: Log final data construction (success path)
	e.log.Infow("ConnSecretData built successfully",
		"Username", data.DBUserName,
		"StandardConnURL", data.ConnURL,
		"SrvConnURL", data.SrvConnURL,
		"NumPrivateConnURLs", len(data.PrivateConnURLs),
	)

	return data, nil
}
