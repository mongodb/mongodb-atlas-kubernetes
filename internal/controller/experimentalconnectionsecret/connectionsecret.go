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
)

const (
	InternalSeparator = "$"

	ProjectLabelKey = "atlas.mongodb.com/project-id"
	ClusterLabelKey = "atlas.mongodb.com/cluster-name"
	TypeLabelKey    = "atlas.mongodb.com/type"
	CredLabelVal    = "credentials"

	userNameKey     = "username"
	passwordKey     = "password"
	standardKey     = "connectionStringStandard"
	standardKeySrv  = "connectionStringStandardSrv"
	privateKey      = "connectionStringPrivate"
	privateSrvKey   = "connectionStringPrivateSrv"
	privateShardKey = "connectionStringPrivateShard"
)

var (
	ErrInternalFormatErr     = errors.New("identifiers could not be loaded from internal format")
	ErrK8SFormatErr          = errors.New("identifiers could not be loaded from k8s format")
	ErrMissingPairing        = errors.New("missing user/endpoint")
	ErrAmbiguousPairing      = errors.New("multiple users/endpoints with the same name found")
	ErrUnresolvedProjectID   = errors.New("could not resolve the project id")
	ErrUnresolvedProjectName = errors.New("could not resolve the project name")
)

// ConnnSecretIdentifiers stores all the necessary information that will
// be needed to identiy and get a K8s connection secret
type ConnSecretIdentifiers struct {
	ProjectID        string
	ProjectName      string
	ClusterName      string
	DatabaseUsername string
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

// loadIdentifiers determines whether the request name is internal or K8s format
// and extracts ProjectID, ClusterName, and DatabaseUsername.
func (r *ConnSecretReconciler) loadIdentifiers(ctx context.Context, req types.NamespacedName) (*ConnSecretIdentifiers, error) {
	if strings.Contains(req.Name, InternalSeparator) {
		return r.identifiersFromInternalName(req)
	}

	return r.identifiersFromK8s(ctx, req)
}

// identifiersFromInternalName loads the identifiers for the internal format
// === Internal format: <ProjectID>$<ClusterName>$<DatabaseUserName>
func (r *ConnSecretReconciler) identifiersFromInternalName(req types.NamespacedName) (*ConnSecretIdentifiers, error) {
	parts := strings.Split(req.Name, InternalSeparator)
	if len(parts) != 3 {
		return nil, ErrInternalFormatErr
	}
	if parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, ErrInternalFormatErr
	}
	return &ConnSecretIdentifiers{
		ProjectID:        parts[0],
		ClusterName:      parts[1],
		DatabaseUsername: parts[2],
	}, nil
}

// identifiersFromK8s loads the identifiers for the k8s format
// === K8s format: <ProjectName>-<ClusterName>-<DatabaseUserName>
// K8s secret must exists in the cluster
func (r *ConnSecretReconciler) identifiersFromK8s(ctx context.Context, req types.NamespacedName) (*ConnSecretIdentifiers, error) {
	var secret corev1.Secret
	if err := r.Client.Get(ctx, req, &secret); err != nil {
		return nil, err
	}
	labels := secret.GetLabels()
	projectID, hasProject := labels[ProjectLabelKey]
	clusterName, hasCluster := labels[ClusterLabelKey]
	if !hasProject || !hasCluster {
		return nil, ErrK8SFormatErr
	}
	if projectID == "" || clusterName == "" {
		return nil, ErrK8SFormatErr
	}
	sep := fmt.Sprintf("-%s-", clusterName)
	parts := strings.Split(req.Name, sep)
	if len(parts) != 2 {
		return nil, ErrK8SFormatErr
	}
	if parts[0] == "" || parts[1] == "" {
		return nil, ErrK8SFormatErr
	}
	return &ConnSecretIdentifiers{
		ProjectID:        projectID,
		ProjectName:      parts[0],
		ClusterName:      clusterName,
		DatabaseUsername: parts[1],
	}, nil
}

// loadPair creates the paired resource that contains the parent AtlasDatabaseUser and the Endpoint.
// Endpoint could be AtlasDeployment or AtlasDataFederation
func (r *ConnSecretReconciler) loadPair(ctx context.Context, ids *ConnSecretIdentifiers) (*ConnSecretPair, error) {
	compositeUserKey := ids.ProjectID + "-" + ids.DatabaseUsername

	// Retrieve the AtlasDatabaseUser using the defined indexers
	users := &akov2.AtlasDatabaseUserList{}
	if err := r.Client.List(ctx, users, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.AtlasDatabaseUserBySpecUsernameAndProjectID, compositeUserKey),
	}); err != nil {
		return nil, err
	}
	usersCount := len(users.Items)

	// Retrieve the Endpoints using the defined indexers
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

	// AmbiguousPairing (more than 1 of either resource)
	if usersCount > 1 || totalEndpoints > 1 {
		return nil, ErrAmbiguousPairing
	}

	// Exactly one of each (OK case)
	if usersCount == 1 && totalEndpoints == 1 {
		return &ConnSecretPair{
			ProjectID: ids.ProjectID,
			User:      &users.Items[0],
			Endpoint:  selected,
		}, nil
	}

	// MissingPairing (one or both missing)
	if usersCount == 0 && totalEndpoints == 0 {
		return nil, ErrMissingPairing
	}
	if usersCount == 0 {
		return &ConnSecretPair{
			ProjectID: ids.ProjectID,
			User:      nil,
			Endpoint:  selected,
		}, ErrMissingPairing
	}
	return &ConnSecretPair{
		ProjectID: ids.ProjectID,
		User:      &users.Items[0],
		Endpoint:  nil,
	}, ErrMissingPairing
}

// resolveProject attempts to find the project name, required for creating connection secrets
// as it is used in metadata.name
func (r *ConnSecretReconciler) resolveProjectName(
	ctx context.Context,
	ids *ConnSecretIdentifiers,
	pair *ConnSecretPair,
) (string, error) {
	if ids != nil && ids.ProjectName != "" {
		return ids.ProjectName, nil
	}

	// project name resolution requires at least on parent to be available
	if pair == nil {
		return "", ErrUnresolvedProjectName
	}

	var err error
	var projectName string
	// Try resolving from the Endpoint if present
	if pair.Endpoint != nil {
		projectName, err = pair.Endpoint.GetProjectName(ctx)
		if projectName != "" {
			return kube.NormalizeIdentifier(projectName), nil
		}
	}

	// Fallback, try resolving from the User if present
	if pair.User != nil {
		if name, uerr := r.getUserProjectName(ctx, pair.User); name != "" {
			return kube.NormalizeIdentifier(name), nil
		} else if err == nil {
			err = uerr
		}
	}

	if err == nil {
		err = ErrUnresolvedProjectName
	}
	return "", err
}

// handleDelete ensures that the connection secret from the paired resource and identifiers will get deleted
func (r *ConnSecretReconciler) handleDelete(
	ctx context.Context,
	req ctrl.Request,
	ids *ConnSecretIdentifiers,
	pair *ConnSecretPair,
) (ctrl.Result, error) {
	log := r.Log.With("ns", req.Namespace, "name", req.Name)

	// project name is required for metadata.name
	projectName, err := r.resolveProjectName(ctx, ids, pair)
	if err != nil {
		log.Errorw("failed to resolve project name", "error", err)
		return workflow.Terminate(workflow.ConnSecretUnresolvedProjectName, err).ReconcileResult()
	}

	name := CreateK8sFormat(projectName, ids.ClusterName, ids.DatabaseUsername)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: req.Namespace,
		},
	}

	// delete secret in k8s
	if err := r.Client.Delete(ctx, secret); err != nil {
		if apiErrors.IsNotFound(err) {
			log.Debugw("no secret to delete; already gone")
			return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
		}
		log.Errorw("unable to delete secret", "reason", workflow.ConnSecretFailedDeletion, "error", err)
		return workflow.Terminate(workflow.ConnSecretFailedDeletion, err).ReconcileResult()
	}

	log.Debugw("connection secret deleted")
	r.EventRecorder.Event(secret, corev1.EventTypeNormal, "Deleted", "ConnectionSecret deleted")
	return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
}

// handleUpsert ensures that the connection secret from the paired resource and identifiers will be upserted
func (r *ConnSecretReconciler) handleUpsert(
	ctx context.Context,
	req ctrl.Request,
	ids *ConnSecretIdentifiers,
	pair *ConnSecretPair,
) (ctrl.Result, error) {
	log := r.Log.With("ns", req.Namespace, "name", req.Name)

	// project name is required for metadata.name
	projectName, err := r.resolveProjectName(ctx, ids, pair)
	if err != nil {
		log.Errorw("failed to resolve project name", "error", err)
		return workflow.Terminate(workflow.ConnSecretUnresolvedProjectName, err).ReconcileResult()
	}
	ids.ProjectName = projectName
	log.Debugw("project name resolved for upsert", "projectName", projectName)

	// create the connection data that will populate secret.stringData
	data, err := pair.Endpoint.BuildConnData(ctx, pair.User)
	if err != nil {
		log.Errorw("failed to build connection data", "error", err)
		return workflow.Terminate(workflow.ConnSecretFailedToBuildData, err).ReconcileResult()
	}
	log.Debugw("connection data built")
	if err := r.ensureSecret(ctx, ids, pair, data); err != nil {
		return workflow.Terminate(workflow.ConnSecretFailedToUpsertSecret, err).ReconcileResult()
	}

	log.Infow("connection secret upserted")
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

	// fills the secret.stringData with the information stored in ConnSecretData
	if err := fillConnSecretData(secret, ids, data); err != nil {
		log.Errorw("failed to fill secret data", "reason", workflow.ConnSecretFailedToFillData, "error", err)
		return err
	}

	// adds the owner to be the AtlasDatabaseUser for garbage collecting
	if err := controllerutil.SetControllerReference(pair.User, secret, r.Scheme); err != nil {
		log.Errorw("failed to set controller owner", "reason", workflow.ConnSecretFailedToSetOwnerReferences, "error", err)
		return err
	}

	// upsert secret in k8s
	if err := r.Client.Create(ctx, secret); err != nil {
		if apiErrors.IsAlreadyExists(err) {
			current := &corev1.Secret{}
			if err := r.Client.Get(ctx, client.ObjectKeyFromObject(secret), current); err != nil {
				log.Errorw("failed to fetch existing secret", "error", err)
				return err
			}
			secret.ResourceVersion = current.ResourceVersion
			if err := r.Client.Update(ctx, secret); err != nil {
				log.Errorw("failed to update secret", "error", err)
				return err
			}
		} else {
			log.Errorw("failed to create secret", "error", err)
			return err
		}
	}
	return nil
}

// fillConnSecretData converts the ConnSecretData into secret.stringData
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

// CreateURL creates the connection urls given a hostname, user, and password
func CreateURL(hostname, username, password string) (string, error) {
	if hostname == "" {
		return "", nil
	}
	u, err := url.Parse(hostname)
	if err != nil {
		return "", err
	}
	u.User = url.UserPassword(username, password)
	return u.String(), nil
}
