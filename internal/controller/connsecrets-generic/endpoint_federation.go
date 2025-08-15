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

package connsecretsgeneric

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/datafederation"
)

type FederationEndpoint struct {
	obj *akov2.AtlasDataFederation
	r   *ConnSecretReconciler
}

// GetName resolves the endpoints name from the spec
func (e FederationEndpoint) GetName() string {
	if e.obj == nil {
		return ""
	}
	return e.obj.Spec.Name
}

// IsReady returns true if the endpoint is ready
func (e FederationEndpoint) IsReady() bool {
	return e.obj != nil && api.HasReadyCondition(e.obj.Status.Conditions)
}

// GetScopeType returns the scope type of the endpoint to match with the ones from AtlasDatabaseUser
func (e FederationEndpoint) GetScopeType() akov2.ScopeType {
	return akov2.DataLakeScopeType
}

// GetProjectID resolves parent project's id (only ProjectRef)
func (e FederationEndpoint) GetProjectID(ctx context.Context) (string, error) {
	if e.obj == nil {
		return "", fmt.Errorf("nil federation")
	}
	if e.obj.Spec.Project.Name != "" {
		proj := &akov2.AtlasProject{}
		key := e.obj.Spec.Project.GetObject(e.obj.GetNamespace())
		if err := e.r.Client.Get(ctx, *key, proj); err != nil {
			return "", err
		}
		return proj.ID(), nil
	}

	return "", fmt.Errorf("project ID not available")
}

// GetProjectName returns the parent project's name (only by getting K8s AtlasProject)
func (e FederationEndpoint) GetProjectName(ctx context.Context) (string, error) {
	if e.obj == nil {
		return "", fmt.Errorf("nil federation")
	}
	if e.obj.Spec.Project.Name != "" {
		proj := &akov2.AtlasProject{}
		key := e.obj.Spec.Project.GetObject(e.obj.GetNamespace())
		if err := e.r.Client.Get(ctx, *key, proj); err != nil {
			return "", err
		}
		if proj.Spec.Name != "" {
			return kube.NormalizeIdentifier(proj.Spec.Name), nil
		}
	}

	return "", fmt.Errorf("project name not available")
}

// Defines the list type
func (FederationEndpoint) ListObj() client.ObjectList { return &akov2.AtlasDataFederationList{} }

// Defines the selector to use for indexer when trying to retrieve all endpoints by project
func (FederationEndpoint) SelectorByProject(projectRef string) fields.Selector {
	return fields.OneTermEqualSelector(indexer.AtlasDataFederationByProject, projectRef)
}

// Defines the selector to use for indexer when trying to retrieve all endpoints by project and spec name
func (FederationEndpoint) SelectorByProjectAndName(ids *ConnSecretIdentifiers) fields.Selector {
	return fields.OneTermEqualSelector(indexer.AtlasDataFederationBySpecNameAndProjectID, ids.ProjectID+"-"+ids.ClusterName)
}

// ExtractList creates a list of Endpoint types to preserve the abstraction
func (e FederationEndpoint) ExtractList(ol client.ObjectList) ([]Endpoint, error) {
	l, ok := ol.(*akov2.AtlasDataFederationList)
	if !ok {
		return nil, fmt.Errorf("unexpected list type %T", ol)
	}
	out := make([]Endpoint, 0, len(l.Items))
	for i := range l.Items {
		out = append(out, FederationEndpoint{obj: &l.Items[i], r: e.r})
	}
	return out, nil
}

// BuildConnData defines the specific function/way for building the ConnSecretData given this type of endpoint
// AtlasDataFederation uses SDK calls for getting the hostnames
func (e FederationEndpoint) BuildConnData(ctx context.Context, user *akov2.AtlasDatabaseUser) (ConnSecretData, error) {
	if user == nil || e.obj == nil {
		return ConnSecretData{}, fmt.Errorf("invalid endpoint or user")
	}
	password, err := user.ReadPassword(ctx, e.r.Client)
	if err != nil {
		return ConnSecretData{}, fmt.Errorf("failed to read password for user %q: %w", user.Spec.Username, err)
	}

	project := &akov2.AtlasProject{}
	if err := e.r.Client.Get(ctx, e.obj.AtlasProjectObjectKey(), project); err != nil {
		return ConnSecretData{}, err
	}

	connectionConfig, err := reconciler.GetConnectionConfig(ctx, e.r.Client, project.ConnectionSecretObjectKey(), &e.r.GlobalSecretRef)
	if err != nil {
		return ConnSecretData{}, err
	}
	clientSet, err := e.r.AtlasProvider.SdkClientSet(ctx, connectionConfig.Credentials, e.r.Log)
	if err != nil {
		return ConnSecretData{}, err
	}

	dataFederationService := datafederation.NewAtlasDataFederation(clientSet.SdkClient20250312002.DataFederationApi)
	df, err := dataFederationService.Get(ctx, project.ID(), e.obj.Spec.Name)
	if err != nil {
		return ConnSecretData{}, fmt.Errorf("atlas DF get: %w", err)
	}

	if len(df.Hostnames) == 0 {
		return ConnSecretData{}, fmt.Errorf("no DF hostnames")
	}

	hostlist := strings.Join(df.Hostnames, ",")
	u := &url.URL{
		Scheme:   "mongodb",
		Host:     hostlist,
		Path:     "/",
		RawQuery: "ssl=true",
	}

	return ConnSecretData{
		DBUserName: user.Spec.Username,
		Password:   password,
		ConnURL:    u.String(),
	}, nil
}
