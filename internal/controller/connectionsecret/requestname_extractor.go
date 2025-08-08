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
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

const InternalSeparator = "$"

// RequestNameParts holds the values extracted from a reconcile request name.
type RequestNameParts struct {
	ProjectID        string
	ProjectName      string
	ClusterName      string
	DatabaseUsername string
}

// ConnectionPair represents the pairing of an AtlasDeployment and an AtlasDatabaseUser
// required to construct a ConnectionSecret. It holds resolved identifiers and the corresponding resources.
type ConnectionPair struct {
	RequestNameInfo RequestNameParts
	Deployment      *akov2.AtlasDeployment
	User            *akov2.AtlasDatabaseUser
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

// ParseRequestNameParts determines whether the request name is internal or K8s format
// and extracts ProjectID, ClusterName, and DatabaseUsername.
// === Format Internal: <ProjectID>$<ClusterName>$<DatabaseUserName>
// === Format K8s: <ProjectName>-<ClusterName>-<DatabaseUserName>
func LoadRequestNameParts(ctx context.Context, c client.Client, req types.NamespacedName) (RequestNameParts, error) {
	var ids RequestNameParts

	// Format Internal
	if strings.Contains(req.Name, InternalSeparator) {
		parts := strings.SplitN(req.Name, InternalSeparator, 3)
		if len(parts) != 3 {
			return ids, fmt.Errorf("invalid internal name format for %q: expected 3 parts separated by '%s'", req.Name, InternalSeparator)
		}

		// Error out on incorrect format
		if parts[0] == "" || parts[1] == "" || parts[2] == "" {
			return ids, fmt.Errorf(
				"invalid internal name format for %q: empty value in one or more parts (projectID=%q, clusterName=%q, databaseUsername=%q)",
				req.Name, parts[0], parts[1], parts[2],
			)
		}
		return RequestNameParts{
			ProjectID:        parts[0],
			ClusterName:      parts[1],
			DatabaseUsername: parts[2],
		}, nil
	}

	var secret corev1.Secret
	if err := c.Get(ctx, req, &secret); err != nil {
		return ids, fmt.Errorf(
			"unable to retrieve Secret %q in namespace %q: %w",
			req.Name, req.Namespace, err,
		)
	}

	labels := secret.GetLabels()
	projectID, hasProject := labels[ProjectLabelKey]
	clusterName, hasCluster := labels[ClusterLabelKey]

	// Error out on missing labels or missing values
	var missing []string
	if !hasProject {
		missing = append(missing, ProjectLabelKey)
	}
	if !hasCluster {
		missing = append(missing, ClusterLabelKey)
	}
	if len(missing) > 0 {
		return ids, fmt.Errorf(
			"secret %q is missing required label(s): %v. Current labels: %v",
			req.Name, missing, labels,
		)
	}
	if projectID == "" {
		return ids, fmt.Errorf("secret %q has empty value for label %q", req.Name, ProjectLabelKey)
	}
	if clusterName == "" {
		return ids, fmt.Errorf("secret %q has empty value for label %q", req.Name, ClusterLabelKey)
	}

	sep := fmt.Sprintf("-%s-", clusterName)
	parts := strings.SplitN(req.Name, sep, 2)

	// Error out on incorrect format
	if len(parts) != 2 {
		return ids, fmt.Errorf(
			"invalid K8s name format for %q: expected separator '-%s-'",
			req.Name, clusterName,
		)
	}
	if parts[0] == "" || parts[1] == "" {
		return ids, fmt.Errorf(
			"invalid K8s name format for %q: empty value in one or more parts (projectName=%q, clusterName=%q, databaseUsername=%q)",
			req.Name, parts[0], clusterName, parts[1],
		)
	}

	return RequestNameParts{
		ProjectID:        projectID,
		ProjectName:      parts[0],
		ClusterName:      clusterName,
		DatabaseUsername: parts[1],
	}, nil
}

// LoadPairedResources retrieves the paired AtlasDeployment and AtlasDatabaseUser forming the ConnectionSecret
func LoadPairedResources(ctx context.Context, c client.Client, ids RequestNameParts, namespace string) (*ConnectionPair, error) {
	// Use the indexer composite key to extract AtlasDeployment
	compositeDeploymentKey := ids.ProjectID + "-" + ids.ClusterName
	deployments := &akov2.AtlasDeploymentList{}
	if err := c.List(ctx, deployments, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.AtlasDeploymentBySpecNameAndProjectID, compositeDeploymentKey),
		Namespace:     namespace,
	}); err != nil {
		return nil, fmt.Errorf("failed to list AtlasDeployments: %w", err)
	}
	if len(deployments.Items) != 1 {
		return nil, fmt.Errorf("expected 1 AtlasDeployment for %q, found %d", compositeDeploymentKey, len(deployments.Items))
	}

	// Use the indexer composite key to extract AtlasDatabaseUser
	compositeUserKey := ids.ProjectID + "-" + ids.DatabaseUsername
	users := &akov2.AtlasDatabaseUserList{}
	if err := c.List(ctx, users, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.AtlasDatabaseUserBySpecUsernameAndProjectID, compositeUserKey),
		Namespace:     namespace,
	}); err != nil {
		return nil, fmt.Errorf("failed to list AtlasDatabaseUsers: %w", err)
	}
	if len(users.Items) != 1 {
		return nil, fmt.Errorf("expected 1 AtlasDatabaseUser for %q, found %d", compositeUserKey, len(users.Items))
	}

	return &ConnectionPair{
		RequestNameInfo: ids,
		Deployment:      &deployments.Items[0],
		User:            &users.Items[0],
	}, nil
}

// AreResourcesReady checks that both AtlasDeployment and AtlasDatabaseUser are ready
func (p *ConnectionPair) AreResourcesReady() (bool, []string) {
	notReady := []string{}

	if !IsDeploymentReady(p.Deployment) {
		notReady = append(notReady, fmt.Sprintf("AtlasDeployment/%s", p.Deployment.GetName()))
	}
	if !IsDatabaseUserReady(p.User) {
		notReady = append(notReady, fmt.Sprintf("AtlasDatabaseUser/%s", p.User.GetName()))
	}

	return len(notReady) == 0, notReady
}

// NeedsSDKProjectResolution checks if we need to use SDK to retrieve the projectName
func (p *ConnectionPair) NeedsSDKProjectResolution() bool {
	return p.Deployment.Spec.ExternalProjectRef != nil &&
		p.Deployment.Spec.ExternalProjectRef.ID != "" &&
		p.User.Spec.ExternalProjectRef != nil &&
		p.User.Spec.ExternalProjectRef.ID != ""
}

// ResolveProjectNameK8s retrieves the ProjectName by K8s AtlasProject resource
func (p *ConnectionPair) ResolveProjectNameK8s(ctx context.Context, c client.Client, namespace string) error {
	var name string
	if p.Deployment.Spec.ExternalProjectRef == nil {
		name = p.Deployment.Spec.ProjectRef.Name
	} else {
		name = p.User.Spec.ProjectRef.Name
	}

	project := &akov2.AtlasProject{}
	if err := c.Get(ctx, kube.ObjectKey(namespace, name), project); err != nil {
		return fmt.Errorf("failed to retrieve AtlasProject %q: %w", name, err)
	}

	p.RequestNameInfo.ProjectName = kube.NormalizeIdentifier(project.Spec.Name)
	return nil
}

// ResolveProjectNameSDK resolves the ProjectName by calling the Atlas SDK
func (p *ConnectionPair) ResolveProjectNameSDK(ctx context.Context, projectService project.ProjectService) error {
	atlasProject, err := projectService.GetProject(ctx, p.RequestNameInfo.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to fetch project from Atlas API for %q: %w", p.RequestNameInfo.ProjectID, err)
	}

	p.RequestNameInfo.ProjectName = kube.NormalizeIdentifier(atlasProject.Name)
	return nil
}

// BuildConnectionData constructs the secret data that will be passed in the secret
func (p *ConnectionPair) BuildConnectionData(ctx context.Context, c client.Client) (ConnectionData, error) {
	password, err := p.User.ReadPassword(ctx, c)
	if err != nil {
		return ConnectionData{}, fmt.Errorf("failed to read password for user %q: %w", p.User.Spec.Username, err)
	}

	conn := p.Deployment.Status.ConnectionStrings

	data := ConnectionData{
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

// HandleSecret ensures the ConnectionSecret resource is created or updated with the given connection data.
// It wraps the Ensure(...) helper and constructs the secret name from the internal pairing.
func (p *ConnectionPair) HandleSecret(ctx context.Context, c client.Client, data ConnectionData) error {
	_, err := Ensure(ctx, c, p.User.Namespace, p.RequestNameInfo.ProjectName, p.RequestNameInfo.ProjectID, p.RequestNameInfo.ClusterName, data)
	if err != nil {
		return fmt.Errorf("ensure failed for secret (projectName=%q, clusterName=%q, user=%q): %w",
			p.RequestNameInfo.ProjectName, p.RequestNameInfo.ClusterName, p.User.Spec.Username, err)
	}
	return nil
}
