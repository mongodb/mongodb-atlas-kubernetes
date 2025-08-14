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
	"errors"
	"fmt"
	"net/url"
	"strings"

	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/stringutil"
)

const (
	InternalSeparator string = "$"

	ProjectLabelKey string = "atlas.mongodb.com/project-id"
	ClusterLabelKey string = "atlas.mongodb.com/cluster-name"
	TypeLabelKey           = "atlas.mongodb.com/type"
	CredLabelVal           = "credentials"

	userNameKey     string = "username"
	passwordKey     string = "password"
	standardKey     string = "connectionStringStandard"
	standardKeySrv  string = "connectionStringStandardSrv"
	privateKey      string = "connectionStringPrivate"
	privateSrvKey   string = "connectionStringPrivateSrv"
	privateShardKey string = "connectionStringPrivateShard"
)

var (
	ErrNoPairedResourcesFound = errors.New("no paired resources found")
	ErrNoEndpointFound        = errors.New("no endpoint found")
	ErrNoUserFound            = errors.New("no user found")
	ErrManyEndpoints          = errors.New("multiple endpoints found")
	ErrManyUsers              = errors.New("multiple users found")
)

type ConnSecretIdentifiers struct {
	ProjectID        string
	ProjectName      string
	ClusterName      string
	DatabaseUsername string
}

type ConnSecretData struct {
	DBUserName      string
	Password        string
	ConnURL         string
	SrvConnURL      string
	PrivateConnURLs []PrivateLinkConnURLs
}

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

func (r *ConnSecretReconciler) LoadIdentifiers(ctx context.Context, req types.NamespacedName) (*ConnSecretIdentifiers, error) {
	// === Internal format: <ProjectID>$<ClusterName>$<DatabaseUserName>
	if strings.Contains(req.Name, InternalSeparator) {
		parts := strings.Split(req.Name, InternalSeparator)
		if len(parts) != 3 {
			return nil, fmt.Errorf("internal format expected 3 parts separated by %q", InternalSeparator)
		}
		if parts[0] == "" || parts[1] == "" || parts[2] == "" {
			return nil, fmt.Errorf("internal format got empty value in one or more parts")
		}
		return &ConnSecretIdentifiers{
			ProjectID:        parts[0],
			ClusterName:      parts[1],
			DatabaseUsername: parts[2],
		}, nil
	}

	// === K8s format: <ProjectName>-<ClusterName>-<DatabaseUserName>
	var secret corev1.Secret
	if err := r.Client.Get(ctx, req, &secret); err != nil {
		return nil, err
	}
	labels := secret.GetLabels()
	projectID, hasProject := labels[ProjectLabelKey]
	clusterName, hasCluster := labels[ClusterLabelKey]
	if !hasProject || !hasCluster {
		return nil, fmt.Errorf("k8s format got a missing required label(s)")
	}
	if projectID == "" || clusterName == "" {
		return nil, fmt.Errorf("k8s format got label present but empty")
	}

	sep := fmt.Sprintf("-%s-", clusterName)
	parts := strings.SplitN(req.Name, sep, 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("k8s format expected to separate across -%s-", clusterName)
	}
	if parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("k8s format got empty value in one or more parts")
	}

	return &ConnSecretIdentifiers{
		ProjectID:        projectID,
		ProjectName:      parts[0],
		ClusterName:      clusterName,
		DatabaseUsername: parts[1],
	}, nil
}

func (r *ConnSecretReconciler) LoadPair(ctx context.Context, ids *ConnSecretIdentifiers) (*ConnSecretPair, error) {
	compositeUserKey := ids.ProjectID + "-" + ids.DatabaseUsername
	users := &akov2.AtlasDatabaseUserList{}
	if err := r.Client.List(ctx, users, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.AtlasDatabaseUserBySpecUsernameAndProjectID, compositeUserKey),
	}); err != nil {
		return nil, err
	}

	totalEndpoints := 0
	var selected Endpoint
	for _, kind := range r.EndpointKinds {
		list := kind.ListObj()
		if err := r.Client.List(ctx, list, &client.ListOptions{FieldSelector: kind.SelectorByProjectAndName(ids)}); err != nil {
			return nil, err
		}
		eps, err := kind.ExtractList(list)
		if err != nil {
			return nil, err
		}
		if len(eps) == 1 {
			selected = eps[0]
		}
		totalEndpoints += len(eps)
	}

	switch {
	case totalEndpoints > 1:
		return nil, ErrManyEndpoints
	case len(users.Items) > 1:
		return nil, ErrManyUsers
	case totalEndpoints == 0 && len(users.Items) == 0:
		return nil, ErrNoPairedResourcesFound
	case totalEndpoints == 0:
		return &ConnSecretPair{
			ProjectID: ids.ProjectID,
			User:      &users.Items[0],
			Endpoint:  nil,
		}, ErrNoEndpointFound
	case len(users.Items) == 0:
		return &ConnSecretPair{
			ProjectID: ids.ProjectID,
			User:      nil,
			Endpoint:  selected,
		}, ErrNoUserFound
	}

	return &ConnSecretPair{
		ProjectID: ids.ProjectID,
		User:      &users.Items[0],
		Endpoint:  selected,
	}, nil
}

func (r *ConnSecretReconciler) handleDelete(
	ctx context.Context,
	req ctrl.Request,
	ids *ConnSecretIdentifiers,
	pair *ConnSecretPair,
) (ctrl.Result, error) {
	log := r.Log.With("ns", req.Namespace, "name", req.Name)

	projectName, err := pair.Endpoint.GetProjectName(ctx, r.Client, r.AtlasProvider, r.Log)
	if projectName == "" {
		err = fmt.Errorf("project name is empty")
	}
	if err != nil {
		log.Errorw("failed to resolve project name", "reason", workflow.ConnSecretUnresolvedProjectName, "error", err)
		return workflow.Terminate(workflow.ConnSecretUnresolvedProjectName, err).ReconcileResult()
	}

	log.Debugw("project name resolved for delete")

	name := CreateK8sFormat(projectName, ids.ClusterName, ids.DatabaseUsername)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: req.Namespace,
		},
	}

	if err := r.Client.Delete(ctx, secret); err != nil {
		if apiErrors.IsNotFound(err) {
			log.Debugw("no secret to delete; already gone")
			return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
		}
		log.Errorw("unable to delete secret", "reason", workflow.ConnSecretFailedDeletion, "error", err)
		return workflow.Terminate(workflow.ConnSecretFailedDeletion, err).ReconcileResult()
	}

	log.Infow("secret deleted", "reason", workflow.ConnSecretDeleted)
	r.EventRecorder.Event(secret, corev1.EventTypeNormal, "Deleted", "ConnectionSecret deleted")
	return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
}

func (r *ConnSecretReconciler) handleUpsert(
	ctx context.Context,
	req ctrl.Request,
	ids *ConnSecretIdentifiers,
	pair *ConnSecretPair,
) (ctrl.Result, error) {
	log := r.Log.With("ns", req.Namespace, "name", req.Name)

	projectName, err := pair.Endpoint.GetProjectName(ctx, r.Client, r.AtlasProvider, r.Log)
	if projectName == "" {
		err = fmt.Errorf("project name is empty")
	}
	if err != nil {
		log.Errorw("failed to resolve project name", "reason", workflow.ConnSecretFailedToResolveProjectName, "error", err)
		return workflow.Terminate(workflow.ConnSecretFailedToResolveProjectName, err).ReconcileResult()
	}
	ids.ProjectName = projectName
	log.Debugw("project name resolved for upsert")

	data, err := r.buildConnectionData(ctx, pair)
	if err != nil {
		log.Errorw("failed to build connection data", "reason", workflow.ConnSecretFailedToBuildData, "error", err)
		return workflow.Terminate(workflow.ConnSecretFailedToBuildData, err).ReconcileResult()
	}
	log.Debugw("connection data built")

	if err := r.ensureSecret(ctx, ids, pair, data); err != nil {
		return workflow.Terminate(workflow.ConnSecretFailedToCreateSecret, err).ReconcileResult()
	}

	log.Infow("secret upserted", "reason", workflow.ConnSecretUpsert)
	return workflow.OK().ReconcileResult()
}

// ensureSecret creates or updates the Secret for the given identifiers and connection data
func (r *ConnSecretReconciler) ensureSecret(
	ctx context.Context,
	ids *ConnSecretIdentifiers,
	pair *ConnSecretPair,
	data ConnSecretData,
) error {
	namespace := pair.User.GetNamespace()
	log := r.Log.With("ns", namespace, "project", ids.ProjectName)

	name := CreateK8sFormat(ids.ProjectName, ids.ClusterName, ids.DatabaseUsername)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	if err := fillConnSecretData(secret, ids, data); err != nil {
		log.Errorw("failed to fill secret data", "reason", workflow.ConnSecretFailedToFillData, "error", err)
		return err
	}

	if err := controllerutil.SetControllerReference(pair.User, secret, r.Scheme); err != nil {
		log.Errorw("failed to set controller owner", "reason", workflow.ConnSecretFailedToSetOwnerReferences, "error", err)
		return err
	}

	if err := r.Client.Create(ctx, secret); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			current := &corev1.Secret{}
			if err := r.Client.Get(ctx, client.ObjectKeyFromObject(secret), current); err != nil {
				log.Errorw("failed to fetch existing secret", "reason", workflow.ConnSecretFailedToGetSecret, "error", err)
				return err
			}
			secret.ResourceVersion = current.ResourceVersion
			if err := r.Client.Update(ctx, secret); err != nil {
				log.Errorw("failed to update secret", "reason", workflow.ConnSecretFailedToUpdateSecret, "error", err)
				return err
			}
		} else {
			log.Errorw("failed to create secret", "reason", workflow.ConnSecretFailedToCreateSecret, "error", err)
			return err
		}
	}
	return nil
}

func fillConnSecretData(secret *corev1.Secret, ids *ConnSecretIdentifiers, data ConnSecretData) error {
	var err error
	username := data.DBUserName
	password := data.Password

	if data.ConnURL, err = CreateURL(data.ConnURL, username, password); err != nil {
		return err
	}
	if data.SrvConnURL, err = CreateURL(data.SrvConnURL, username, password); err != nil {
		return err
	}
	for i, pe := range data.PrivateConnURLs {
		if data.PrivateConnURLs[i].PvtConnURL, err = CreateURL(pe.PvtConnURL, username, password); err != nil {
			return err
		}
		if data.PrivateConnURLs[i].PvtSrvConnURL, err = CreateURL(pe.PvtSrvConnURL, username, password); err != nil {
			return err
		}
		if data.PrivateConnURLs[i].PvtShardConnURL, err = CreateURL(pe.PvtShardConnURL, username, password); err != nil {
			return err
		}
	}

	secret.Labels = map[string]string{
		TypeLabelKey:    CredLabelVal,
		ProjectLabelKey: ids.ProjectID,
		ClusterLabelKey: ids.ClusterName,
	}

	secret.Data = map[string][]byte{
		userNameKey:    []byte(data.DBUserName),
		passwordKey:    []byte(data.Password),
		standardKey:    []byte(data.ConnURL),
		standardKeySrv: []byte(data.SrvConnURL),
		privateKey:     []byte(""),
		privateSrvKey:  []byte(""),
	}

	for i, pe := range data.PrivateConnURLs {
		suffix := ""
		if i != 0 {
			suffix = fmt.Sprint(i)
		}
		secret.Data[privateKey+suffix] = []byte(pe.PvtConnURL)
		secret.Data[privateSrvKey+suffix] = []byte(pe.PvtSrvConnURL)
		secret.Data[privateShardKey+suffix] = []byte(pe.PvtShardConnURL)
	}

	return nil
}

func CreateURL(connURL, username, password string) (string, error) {
	if connURL == "" {
		return "", nil
	}
	u, err := url.Parse(connURL)
	if err != nil {
		return "", err
	}
	u.User = url.UserPassword(username, password)
	return u.String(), nil
}

func (r *ConnSecretReconciler) buildConnectionData(ctx context.Context, p *ConnSecretPair) (ConnSecretData, error) {
	if p == nil || p.User == nil || p.Endpoint == nil {
		return ConnSecretData{}, fmt.Errorf("invalid pair: nil user or endpoint")
	}

	return p.Endpoint.BuildConnData(ctx, r.Client, r.AtlasProvider, r.Log, p.User)
}

func allowsByScopes(u *akov2.AtlasDatabaseUser, epName string) bool {
	scopes := u.GetScopes(akov2.DeploymentScopeType)
	if len(scopes) == 0 || stringutil.Contains(scopes, epName) {
		return true
	}

	return false
}
