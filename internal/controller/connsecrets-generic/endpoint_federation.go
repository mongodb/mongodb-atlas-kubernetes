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
	"strings"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/datafederation"
)

type FederationEndpoint struct {
	obj *akov2.AtlasDataFederation
	r   *ConnSecretReconciler
}

// ---- instance methods ----
func (e FederationEndpoint) GetName() string {
	if e.obj == nil {
		return ""
	}
	return e.obj.Spec.Name
}

func (e FederationEndpoint) IsReady() bool {
	return e.obj != nil && api.HasReadyCondition(e.obj.Status.Conditions)
}

func (e FederationEndpoint) GetProjectRef(ctx context.Context, r client.Reader) string {
	return e.obj.Spec.Project.Name
}

func (e FederationEndpoint) GetProjectID(ctx context.Context, r client.Reader) (string, error) {
	if e.obj == nil {
		return "", fmt.Errorf("nil federation")
	}
	if e.obj.Spec.Project.Name != "" {
		proj := &akov2.AtlasProject{}
		if err := r.Get(ctx, kube.ObjectKey(e.obj.Namespace, e.obj.Spec.Project.Name), proj); err != nil {
			return "", err
		}
		return proj.ID(), nil
	}

	return "", fmt.Errorf("project ID not available")
}

func (e FederationEndpoint) GetProjectName(ctx context.Context, r client.Reader, provider atlas.Provider, log *zap.SugaredLogger) (string, error) {
	if e.obj == nil {
		return "", fmt.Errorf("nil federation")
	}
	if e.obj.Spec.Project.Name != "" {
		proj := &akov2.AtlasProject{}
		if err := r.Get(ctx, kube.ObjectKey(e.obj.Namespace, e.obj.Spec.Project.Name), proj); err != nil {
			return "", err
		}
		if proj.Spec.Name != "" {
			return kube.NormalizeIdentifier(proj.Spec.Name), nil
		}
	}

	return "", fmt.Errorf("project name not available")
}

// ---- indexer methods ----
func (FederationEndpoint) ListObj() client.ObjectList { return &akov2.AtlasDataFederationList{} }

func (FederationEndpoint) SelectorByProject(projectRef string) fields.Selector {
	return fields.OneTermEqualSelector(indexer.AtlasDataFederationByProject, projectRef)
}

func (FederationEndpoint) SelectorByProjectAndName(ids *ConnSecretIdentifiers) fields.Selector {
	return fields.OneTermEqualSelector(indexer.AtlasDataFederationBySpecNameAndProjectID, ids.ProjectID+"-"+ids.ClusterName)
}

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

func (e FederationEndpoint) BuildConnData(ctx context.Context, c client.Client, provider atlas.Provider, log *zap.SugaredLogger, user *akov2.AtlasDatabaseUser) (ConnSecretData, error) {
	if user == nil || e.obj == nil {
		return ConnSecretData{}, fmt.Errorf("invalid endpoint or user")
	}
	password, err := user.ReadPassword(ctx, c)
	if err != nil {
		return ConnSecretData{}, fmt.Errorf("failed to read password for user %q: %w", user.Spec.Username, err)
	}

	project := &akov2.AtlasProject{}
	if err := c.Get(ctx, e.obj.AtlasProjectObjectKey(), project); err != nil {
		return ConnSecretData{}, err
	}

	connectionConfig, err := reconciler.GetConnectionConfig(ctx, c, project.ConnectionSecretObjectKey(), &e.r.GlobalSecretRef)
	if err != nil {
		return ConnSecretData{}, err
	}

	clientSet, err := e.r.AtlasProvider.SdkClientSet(ctx, connectionConfig.Credentials, log)
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
	urls := make([]string, 0, len(df.Hostnames))
	for _, h := range df.Hostnames {
		urls = append(urls, fmt.Sprintf("mongodb://%s:%s@%s?ssl=true", user.Spec.Username, password, h))
	}

	return ConnSecretData{
		DBUserName: user.Spec.Username,
		Password:   password,
		ConnURL:    strings.Join(urls, ","),
	}, nil
}
