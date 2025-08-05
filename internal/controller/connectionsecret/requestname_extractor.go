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
	"errors"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/stringutil"
)

const InternalSeparator = "$"

var (
	// Parsing & format
	ErrInternalFormatPartsInvalid = errors.New("internal format expected 3 parts separated by $")
	ErrInternalFormatPartEmpty    = errors.New("internal format got empty value in one or more parts")
	ErrK8sLabelsMissing           = errors.New("k8s format got a missing required label(s)")
	ErrK8sLabelEmpty              = errors.New("k8s format got label present but empty")
	ErrK8sNameSplitInvalid        = errors.New("k8s format expected to separate across -<clusterName>-")
	ErrK8sNameSplitEmpty          = errors.New("k8s format got empty value in one or more parts")

	// Index lookups
	ErrNoPairedResourcesFound = errors.New("no AtlasDeployment and no AtlasDatabaseUser found")
	ErrNoDeploymentFound      = errors.New("no AtlasDeployment found")
	ErrManyDeployments        = errors.New("multiple AtlasDeployments found")
	ErrNoUserFound            = errors.New("no AtlasDatabaseUser found")
	ErrManyUsers              = errors.New("multiple AtlasDatabaseUsers found")
)

// ConnSecretIdentifiers holds the values extracted from a reconcile request name.
type ConnSecretIdentifiers struct {
	ProjectID        string
	ProjectName      string
	ClusterName      string
	DatabaseUsername string
}

// ConnSecretPair represents the pairing of an AtlasDeployment and an AtlasDatabaseUser
// required to construct a ConnectionSecret. It holds resolved identifiers and the corresponding resources.
// NOTE: this struct intentionally stores only ProjectID (not all identifiers) to keep only necessary information.
type ConnSecretPair struct {
	ProjectID  string
	Deployment *akov2.AtlasDeployment
	User       *akov2.AtlasDatabaseUser
}

// ConnectionData contains all connection information required to populate
// the Kubernetes Secret, including standard and SRV URLs and optional Private Link URLs.
type ConnSecretData struct {
	DBUserName      string
	Password        string
	ConnURL         string
	SrvConnURL      string
	PrivateConnURLs []PrivateLinkConnURLs
}

// PrivateLinkConnURLs holds all Private Link connection strings for a single endpoint set.
// Multiple entries allow for multiple private link configurations per deployment.
type PrivateLinkConnURLs struct {
	PvtConnURL      string
	PvtSrvConnURL   string
	PvtShardConnURL string
}

// CreateK8sFormat returns the Secret name in the Kubernetes naming format: <projectName>-<clusterName>-<username>
func CreateK8sFormat(projectName string, clusterName string, databaseUsername string) string {
	return strings.Join([]string{
		kube.NormalizeIdentifier(projectName),
		kube.NormalizeIdentifier(clusterName),
		kube.NormalizeIdentifier(databaseUsername),
	}, "-")
}

// CreateInternalFormat returns the Secret name in the internal format used by watchers: <projectID>$<clusterName>$<username>
func CreateInternalFormat(projectID string, clusterName string, databaseUsername string) string {
	return strings.Join([]string{
		projectID,
		kube.NormalizeIdentifier(clusterName),
		kube.NormalizeIdentifier(databaseUsername),
	}, InternalSeparator)
}

// LoadRequestIdentifiers determines whether the request name is internal or K8s format
// and extracts ProjectID, ClusterName, and DatabaseUsername.
func LoadRequestIdentifiers(ctx context.Context, c client.Client, req types.NamespacedName) (ConnSecretIdentifiers, error) {
	var ids ConnSecretIdentifiers

	// === Internal format: <ProjectID>$<ClusterName>$<DatabaseUserName>
	if strings.Contains(req.Name, InternalSeparator) {
		parts := strings.SplitN(req.Name, InternalSeparator, 3)
		if len(parts) != 3 {
			return ids, ErrInternalFormatPartsInvalid
		}
		if parts[0] == "" || parts[1] == "" || parts[2] == "" {
			return ids, ErrInternalFormatPartEmpty
		}
		return ConnSecretIdentifiers{
			ProjectID:        parts[0],
			ClusterName:      parts[1],
			DatabaseUsername: parts[2],
		}, nil
	}

	// === K8s format: <ProjectName>-<ClusterName>-<DatabaseUserName>
	var secret corev1.Secret
	if err := c.Get(ctx, req, &secret); err != nil {
		return ids, err
	}

	labels := secret.GetLabels()
	projectID, hasProject := labels[ProjectLabelKey]
	clusterName, hasCluster := labels[ClusterLabelKey]

	// Missing labels or values
	if !hasProject || !hasCluster {
		return ids, ErrK8sLabelsMissing
	}
	if projectID == "" || clusterName == "" {
		return ids, ErrK8sLabelEmpty
	}

	sep := fmt.Sprintf("-%s-", clusterName)
	parts := strings.SplitN(req.Name, sep, 2)
	if len(parts) != 2 {
		return ids, ErrK8sNameSplitInvalid
	}
	if parts[0] == "" || parts[1] == "" {
		return ids, ErrK8sNameSplitEmpty
	}

	return ConnSecretIdentifiers{
		ProjectID:        projectID,
		ProjectName:      parts[0],
		ClusterName:      clusterName,
		DatabaseUsername: parts[1],
	}, nil
}

// LoadPairedResources fetches the paired AtlasDeployment and AtlasDatabaseUser forming the ConnectionSecret
// using the registered indexers
func LoadPairedResources(ctx context.Context, c client.Client, ids ConnSecretIdentifiers, namespace string) (*ConnSecretPair, error) {
	compositeDeploymentKey := ids.ProjectID + "-" + ids.ClusterName
	deployments := &akov2.AtlasDeploymentList{}
	if err := c.List(ctx, deployments, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.AtlasDeploymentBySpecNameAndProjectID, compositeDeploymentKey),
		// Namespace:     namespace, // Do not uncomment; we should be able to create connection secrets cross-namespaced
	}); err != nil {
		return nil, err
	}

	compositeUserKey := ids.ProjectID + "-" + ids.DatabaseUsername
	users := &akov2.AtlasDatabaseUserList{}
	if err := c.List(ctx, users, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.AtlasDatabaseUserBySpecUsernameAndProjectID, compositeUserKey),
		Namespace:     namespace,
	}); err != nil {
		return nil, err
	}

	switch {
	case len(deployments.Items) == 0 && len(users.Items) == 0:
		return nil, ErrNoPairedResourcesFound
	case len(deployments.Items) == 0:
		return &ConnSecretPair{
			ProjectID:  ids.ProjectID,
			Deployment: nil,
			User:       &users.Items[0],
		}, ErrNoDeploymentFound
	case len(users.Items) == 0:
		return &ConnSecretPair{
			ProjectID:  ids.ProjectID,
			Deployment: &deployments.Items[0],
			User:       nil,
		}, ErrNoUserFound
	case len(deployments.Items) > 1:
		return nil, ErrManyDeployments
	case len(users.Items) > 1:
		return nil, ErrManyUsers
	}

	return &ConnSecretPair{
		ProjectID:  ids.ProjectID,
		Deployment: &deployments.Items[0],
		User:       &users.Items[0],
	}, nil
}

// InvalidScopes checks whether the Deployment and User have a common scope
func (p *ConnSecretPair) InvalidScopes() bool {
	scopes := p.User.GetScopes(akov2.DeploymentScopeType)
	if len(scopes) != 0 && !stringutil.Contains(scopes, p.Deployment.GetDeploymentName()) {
		return true
	}

	return false
}

// IsReady checks that both AtlasDeployment and AtlasDatabaseUser are ready
func (p *ConnSecretPair) IsReady() (bool, []string) {
	notReady := []string{}

	if p.Deployment == nil || !IsDeploymentReady(p.Deployment) {
		if p.Deployment != nil {
			notReady = append(notReady, fmt.Sprintf("AtlasDeployment/%s", p.Deployment.GetName()))
		} else {
			notReady = append(notReady, "AtlasDeployment/<nil>")
		}
	}
	if p.User == nil || !IsDatabaseUserReady(p.User) {
		if p.User != nil {
			notReady = append(notReady, fmt.Sprintf("AtlasDatabaseUser/%s", p.User.GetName()))
		} else {
			notReady = append(notReady, "AtlasDatabaseUser/<nil>")
		}
	}

	return len(notReady) == 0, notReady
}

// ResolveProjectNameK8s retrieves the ProjectName by K8s AtlasProject resource
func (p *ConnSecretPair) ResolveProjectNameK8s(ctx context.Context, c client.Client, namespace string) (string, error) {
	var name string
	if p.Deployment != nil && p.Deployment.Spec.ProjectRef != nil {
		name = p.Deployment.Spec.ProjectRef.Name
	} else if p.User != nil && p.User.Spec.ProjectRef != nil {
		name = p.User.Spec.ProjectRef.Name
	} else {
		return "", errors.New("no ProjectRef available on Deployment or User")
	}

	proj := &akov2.AtlasProject{}
	if err := c.Get(ctx, kube.ObjectKey(namespace, name), proj); err != nil {
		return "", fmt.Errorf("failed to retrieve AtlasProject %q: %w", name, err)
	}

	return kube.NormalizeIdentifier(proj.Spec.Name), nil
}

// BuildConnectionData constructs the secret data that will be passed in the secret
func (p *ConnSecretPair) BuildConnectionData(ctx context.Context, c client.Client) (ConnSecretData, error) {
	password, err := p.User.ReadPassword(ctx, c)
	if err != nil {
		return ConnSecretData{}, fmt.Errorf("failed to read password for user %q: %w", p.User.Spec.Username, err)
	}

	conn := p.Deployment.Status.ConnectionStrings

	data := ConnSecretData{
		DBUserName: p.User.Spec.Username,
		Password:   password,
		ConnURL:    conn.Standard,
		SrvConnURL: conn.StandardSrv,
	}

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
