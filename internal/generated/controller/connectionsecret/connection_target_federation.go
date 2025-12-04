// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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
	"net/url"
	"strings"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/datafederation"
)

type DataFederationConnectionTarget struct {
	obj             *akov2.AtlasDataFederation
	client          client.Client
	provider        atlas.Provider
	globalSecretRef client.ObjectKey
	log             *zap.SugaredLogger
}

// GetConnectionTargetType returns the connectionTarget type
func (e DataFederationConnectionTarget) GetConnectionTargetType() string {
	return "data-federation"
}

// GetName resolves the connectionTargets name from the spec
func (e DataFederationConnectionTarget) GetName() string {
	if e.obj == nil {
		return ""
	}
	return e.obj.Spec.Name
}

// IsReady returns true if the connectionTarget is ready
func (e DataFederationConnectionTarget) IsReady() bool {
	return e.obj != nil && api.HasReadyCondition(e.obj.Status.Conditions)
}

// GetScopeType returns the scope type of the connectionTarget to match with the ones from AtlasDatabaseUser
func (e DataFederationConnectionTarget) GetScopeType() akov2.ScopeType {
	return akov2.DataLakeScopeType
}

// GetProjectID resolves parent project's id (only ProjectRef)
func (e DataFederationConnectionTarget) GetProjectID(ctx context.Context) (string, error) {
	if e.obj == nil {
		return "", fmt.Errorf("nil federation")
	}
	if e.obj.Spec.Project.Name != "" {
		return resolveProjectIDByKey(ctx, e.client, e.obj.AtlasProjectObjectKey())
	}
	return "", ErrUnresolvedProjectID
}

// Defines the selector to use for indexer when trying to retrieve all connectionTargets by project
func (DataFederationConnectionTarget) SelectorByProjectID(projectID string) fields.Selector {
	return fields.OneTermEqualSelector(indexer.AtlasDataFederationByProjectID, projectID)
}

// BuildConnectionData builds the ConnectionSecretData for this connectionTarget type
func (e DataFederationConnectionTarget) BuildConnectionData(ctx context.Context, user *akov2.AtlasDatabaseUser) (ConnectionSecretData, error) {
	if user == nil || e.obj == nil {
		return ConnectionSecretData{}, ErrMissingPairing
	}
	e.log.Debugw("Starting BuildConnectionData", "Username", user.Spec.Username, "DataFederationConnectionTarget", e.obj.Spec.Name)

	password, err := user.ReadPassword(ctx, e.client)
	if err != nil {
		return ConnectionSecretData{}, fmt.Errorf("failed to read password for user %q: %w", user.Spec.Username, err)
	}

	project := &akov2.AtlasProject{}
	if err := e.client.Get(ctx, e.obj.AtlasProjectObjectKey(), project); err != nil {
		e.log.Errorw("Failed to fetch project for DataFederationConnectionTarget", "Error", err)
		return ConnectionSecretData{}, err
	}

	connectionConfig, err := reconciler.GetConnectionConfig(ctx, e.client, project.ConnectionSecretObjectKey(), &e.globalSecretRef)
	if err != nil {
		return ConnectionSecretData{}, err
	}

	clientSet, err := e.provider.SdkClientSet(ctx, connectionConfig.Credentials, e.log)
	if err != nil {
		return ConnectionSecretData{}, err
	}

	dataFederationService := datafederation.NewAtlasDataFederation(clientSet.SdkClient20250312009.DataFederationApi)
	df, err := dataFederationService.Get(ctx, project.ID(), e.obj.Spec.Name)
	if err != nil {
		return ConnectionSecretData{}, fmt.Errorf("atlas DF get: %w", err)
	}

	if len(df.Hostnames) == 0 {
		return ConnectionSecretData{}, fmt.Errorf("no DF hostnames")
	}

	hostlist := strings.Join(df.Hostnames, ",")
	e.log.Debugw("Building connection URL for DataFederationConnectionTarget", "Hostlist", hostlist)

	u := &url.URL{
		Scheme:   "mongodb",
		Host:     hostlist,
		Path:     "/",
		RawQuery: "ssl=true",
	}

	connData := ConnectionSecretData{
		DBUserName:    user.Spec.Username,
		Password:      password,
		ConnectionURL: u.String(),
	}

	e.log.Debugw("ConnectionSecret data built successfully",
		"DBUserName", connData.DBUserName,
		"ConnectionURL", connData.ConnectionURL,
		"PasswordIsSet", len(connData.Password) > 0,
	)
	return connData, nil
}
