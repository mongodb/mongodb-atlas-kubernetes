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

package connectionsecret

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
)

type DeploymentConnectionTarget struct {
	obj             *akov2.AtlasDeployment
	client          client.Client
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

// GetConnectionTargetType returns the connectionTarget type
func (e DeploymentConnectionTarget) GetConnectionTargetType() string {
	return "deployment"
}

// GetName resolves the connectionTargets name from the spec
func (e DeploymentConnectionTarget) GetName() string {
	if e.obj == nil {
		return ""
	}
	return e.obj.GetDeploymentName()
}

// IsReady returns true if the connectionTarget is ready
func (e DeploymentConnectionTarget) IsReady() bool {
	return e.obj != nil && api.HasReadyCondition(e.obj.Status.Conditions)
}

func (e DeploymentConnectionTarget) GetOwnerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: e.obj.APIVersion,
			Kind:       e.obj.Kind,
			Name:       e.obj.Name,
			UID:        e.obj.UID,
		},
	}
}

// GetScopeType returns the scope type of the connectionTarget to match with the ones from AtlasDatabaseUser
func (e DeploymentConnectionTarget) GetScopeType() akov2.ScopeType {
	return akov2.DeploymentScopeType
}

// GetProjectID resolves parent project's id (ProjectRef or ExternalRef)
func (e DeploymentConnectionTarget) GetProjectID(ctx context.Context) (string, error) {
	if e.obj == nil {
		return "", fmt.Errorf("nil deployment")
	}
	if e.obj.Spec.ExternalProjectRef != nil && e.obj.Spec.ExternalProjectRef.ID != "" {
		return e.obj.Spec.ExternalProjectRef.ID, nil
	}
	if e.obj.Spec.ProjectRef != nil && e.obj.Spec.ProjectRef.Name != "" {
		return resolveProjectIDByKey(ctx, e.client, e.obj.AtlasProjectObjectKey())
	}
	return "", ErrUnresolvedProjectID
}

// Defines the selector to use for indexer when trying to retrieve all connectionTargets by project
func (DeploymentConnectionTarget) SelectorByProjectID(projectID string) fields.Selector {
	return fields.OneTermEqualSelector(indexer.AtlasDeploymentByProject, projectID)
}

// Defines the selector to use for indexer when trying to retrieve all connectionTargets by project and spec name
func (DeploymentConnectionTarget) SelectorByTargetIdentifierFields(ids *ConnectionSecretIdentifiers) fields.Selector {
	return fields.OneTermEqualSelector(indexer.AtlasDeploymentBySpecNameAndProjectID, ids.ProjectID+"-"+ids.TargetName)
}

// BuildConnectionData defines the specific function/way for building the ConnectionSecretData given this type of connectionTarget
// AtlasDeployment stores connection strings in the status field
func (e DeploymentConnectionTarget) BuildConnectionData(ctx context.Context, user *akov2.AtlasDatabaseUser) (ConnectionSecretData, error) {
	// Step 1: Log basic input details
	if user == nil || e.obj == nil {
		return ConnectionSecretData{}, ErrMissingPairing
	}

	e.log.Debugw("Starting BuildConnectionData",
		"Username", user.Spec.Username,
		"DeploymentName", e.obj.GetDeploymentName(),
	)

	// Step 2: Read the user's password and log outcome
	password, err := user.ReadPassword(ctx, e.client)
	if err != nil {
		return ConnectionSecretData{}, fmt.Errorf("failed to read password for user %q: %w", user.Spec.Username, err)
	}

	e.log.Debugw("Successfully read password for user",
		"Username", user.Spec.Username,
		"PasswordPresent", len(password) > 0,
	)

	// Initialize ConnectionSecretData with DBUserName and Password
	data := ConnectionSecretData{
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
	data.ConnectionURL = conn.Standard
	data.SrvConnectionURL = conn.StandardSrv

	e.log.Debugw("Standard/SRV connection strings",
		"StandardConnectionURL", data.ConnectionURL,
		"SrvConnectionURL", data.SrvConnectionURL,
	)

	// Private connection strings
	if conn.Private != "" {
		e.log.Debugw("Private connection string detected",
			"PrivateConnURL", conn.Private,
			"PrivateSrvConnURL", conn.PrivateSrv,
		)
		data.PrivateConnectionURLs = append(data.PrivateConnectionURLs, PrivateLinkConnectionURLs{
			ConnectionURL:    conn.Private,
			SrvConnectionURL: conn.PrivateSrv,
		})
	}

	// Iterate through PrivateEndpoint connection strings
	for _, pe := range conn.PrivateEndpoint {
		e.log.Debugw("PrivateEndpoint connection string detected",
			"ConnectionURL", pe.ConnectionString,
			"SrvConnectionURL", pe.SRVConnectionString,
			"ShardConnectionURL", pe.SRVShardOptimizedConnectionString,
		)
		data.PrivateConnectionURLs = append(data.PrivateConnectionURLs, PrivateLinkConnectionURLs{
			ConnectionURL:      pe.ConnectionString,
			SrvConnectionURL:   pe.SRVConnectionString,
			ShardConnectionURL: pe.SRVShardOptimizedConnectionString,
		})
	}

	// Step 4: Log final data construction (success path)
	e.log.Debugw("ConnectionSecretData built successfully",
		"Username", data.DBUserName,
		"StandardConnectionURL", data.ConnectionURL,
		"SrvConnectionURL", data.SrvConnectionURL,
		"NumPrivateConnectionURLs", len(data.PrivateConnectionURLs),
	)

	return data, nil
}
