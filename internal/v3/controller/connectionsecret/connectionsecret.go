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
	"hash/fnv"
	"net/url"
	"strings"

	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
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

	ProjectLabelKey      = "atlas.mongodb.com/project-id"
	ClusterLabelKey      = "atlas.mongodb.com/cluster-name"
	TypeLabelKey         = "atlas.mongodb.com/type"
	DatabaseUserLabelKey = "atlas.mongodb.com/database-user-name"
	ConnectionTypelKey   = "atlas.mongodb.com/connection-type"
	CredLabelVal         = "credentials"

	userNameKey     = "username"
	passwordKey     = "password"
	standardKey     = "connectionStringStandard"
	standardKeySrv  = "connectionStringStandardSrv"
	privateKey      = "connectionStringPrivate"
	privateSrvKey   = "connectionStringPrivateSrv"
	privateShardKey = "connectionStringPrivateShard"
)

var (
	ConnectionSecretGoFieldOwner = client.FieldOwner("connectionsecret")
	ErrInternalFormatErr         = errors.New("identifiers could not be loaded from internal format")
	ErrK8SFormatErr              = errors.New("identifiers could not be loaded from k8s format")
	ErrMissingPairing            = errors.New("missing user/endpoint")
	ErrAmbiguousPairing          = errors.New("multiple users/endpoints with the same name found")
	ErrUnresolvedProjectID       = errors.New("could not resolve the project id")
)

// ConnnSecretIdentifiers stores all the necessary information that will
// be needed to identiy and get a K8s connection secret
type ConnSecretIdentifiers struct {
	ProjectID        string
	ClusterName      string
	DatabaseUsername string
	ConnectionType   string
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

// CreateInternalFormat returns the Secret name in the internal format used by watchers: <projectID>$<clusterName>$<username>
func CreateInternalFormat(projectID string, clusterName string, databaseUsername string, connectionType string) string {
	return strings.Join([]string{
		projectID,
		kube.NormalizeIdentifier(clusterName),
		kube.NormalizeIdentifier(databaseUsername),
		kube.NormalizeIdentifier(connectionType),
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

// identifiersFromInternalName parses identifiers from the internal format.
// === Internal format: <ProjectID>$<ClusterName>$<DatabaseUserName>$<ConnectionType>
func (r *ConnSecretReconciler) identifiersFromInternalName(req types.NamespacedName) (*ConnSecretIdentifiers, error) {
	parts := strings.Split(req.Name, InternalSeparator)
	if len(parts) != 4 {
		return nil, ErrInternalFormatErr
	}
	if parts[0] == "" || parts[1] == "" || parts[2] == "" || parts[3] == "" {
		return nil, ErrInternalFormatErr
	}
	return &ConnSecretIdentifiers{
		ProjectID:        parts[0],
		ClusterName:      parts[1],
		DatabaseUsername: parts[2],
		ConnectionType:   parts[3],
	}, nil
}

// identifiersFromK8s retrieves identifiers from labels and annotations instead of parsing the secret name in Kubernetes format.
// === K8s format: Use labels/annotations to extract metadata.
func (r *ConnSecretReconciler) identifiersFromK8s(ctx context.Context, req types.NamespacedName) (*ConnSecretIdentifiers, error) {
	var secret corev1.Secret
	if err := r.Client.Get(ctx, req, &secret); err != nil {
		return nil, err
	}
	labels := secret.GetLabels()
	annotations := secret.GetAnnotations()

	projectID, hasProject := labels[ProjectLabelKey]
	clusterName, hasCluster := labels[ClusterLabelKey]
	databaseUsername, hasUser := labels[DatabaseUserLabelKey]
	connectionType, hasConnectionType := annotations[ConnectionTypelKey]

	// Validate required fields
	if !hasProject || !hasCluster || !hasUser || !hasConnectionType || projectID == "" || clusterName == "" || databaseUsername == "" || connectionType == "" {
		err := ErrK8SFormatErr
		return nil, err
	}
	return &ConnSecretIdentifiers{
		ProjectID:        projectID,
		ClusterName:      clusterName,
		DatabaseUsername: databaseUsername,
		ConnectionType:   connectionType,
	}, nil
}

// loadPair creates the paired resource that contains the parent AtlasDatabaseUser and the Endpoint.
// Endpoint could be AtlasDeployment or AtlasDataFederation
func (r *ConnSecretReconciler) loadPair(ctx context.Context, ids *ConnSecretIdentifiers) (*akov2.AtlasDatabaseUser, Endpoint, error) {
	compositeUserKey := ids.ProjectID + "-" + ids.DatabaseUsername

	// Retrieve the AtlasDatabaseUser using the defined indexers
	users := &akov2.AtlasDatabaseUserList{}
	if err := r.Client.List(ctx, users, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasDatabaseUserBySpecUsernameAndProjectID,
			compositeUserKey,
		),
	}); err != nil {
		return nil, nil, err
	}
	usersCount := len(users.Items)

	// Retrieve Endpoints using the defined indexers
	totalEndpoints := 0
	var selected Endpoint
	for _, kind := range r.EndpointKinds {
		switch kind.(type) {
		case FederationEndpoint:
			list := &akov2.AtlasDataFederationList{}
			if err := r.Client.List(ctx, list, &client.ListOptions{
				FieldSelector: kind.SelectorByProjectAndName(ids),
			}); err != nil {
				return nil, nil, err
			}

			if len(list.Items) == 1 {
				selected = FederationEndpoint{
					obj:             &list.Items[0],
					k8s:             r.Client,
					provider:        r.AtlasProvider,
					globalSecretRef: r.GlobalSecretRef,
					log:             r.Log,
				}
			}
			totalEndpoints += len(list.Items)

		case DeploymentEndpoint:
			// Handle DeploymentEndpoint
			list := &akov2.AtlasDeploymentList{}
			if err := r.Client.List(ctx, list, &client.ListOptions{
				FieldSelector: kind.SelectorByProjectAndName(ids),
			}); err != nil {
				return nil, nil, err
			}

			if len(list.Items) == 1 {
				selected = DeploymentEndpoint{
					obj:             &list.Items[0],
					k8s:             r.Client,
					provider:        r.AtlasProvider,
					globalSecretRef: r.GlobalSecretRef,
					log:             r.Log,
				}
			}
			totalEndpoints += len(list.Items)
		}
	}

	// AmbiguousPairing (more than 1 of either resource)
	if usersCount > 1 || totalEndpoints > 1 {
		return nil, nil, ErrAmbiguousPairing
	}

	// Exactly one of each (OK case)
	if usersCount == 1 && totalEndpoints == 1 {
		return &users.Items[0], selected, nil
	}

	// Handle missing pairing cases
	if usersCount == 0 && totalEndpoints == 0 {
		return nil, nil, ErrMissingPairing
	}
	if usersCount == 0 {
		return nil, selected, ErrMissingPairing
	}
	return &users.Items[0], nil, ErrMissingPairing
}

// handleDelete ensures that the connection secret from the paired resource and identifiers will get deleted
func (r *ConnSecretReconciler) handleDelete(
	ctx context.Context,
	req ctrl.Request,
	ids *ConnSecretIdentifiers,
) (ctrl.Result, error) {
	log := r.Log.With("ns", req.Namespace, "name", req.Name)

	name := K8sConnectionSecretName(ids.ProjectID, ids.ClusterName, ids.DatabaseUsername, ids.ConnectionType)
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
	user *akov2.AtlasDatabaseUser,
	endpoint Endpoint,
) (ctrl.Result, error) {
	log := r.Log.With("ns", req.Namespace, "name", req.Name)
	log.Debugw("Starting handleUpsert", "ConnSecretIdentifiers", ids, "AtlasDatabaseUser", user)
	// create the connection data that will populate secret.stringData
	data, err := endpoint.BuildConnData(ctx, user)
	if err != nil {
		log.Errorw("failed to build connection data", "error", err)
		return workflow.Terminate(workflow.ConnSecretFailedToBuildData, err).ReconcileResult()
	}
	log.Debugw("connection data built")
	if err := r.ensureSecret(ctx, ids, user, endpoint, data); err != nil {
		return workflow.Terminate(workflow.ConnSecretFailedToUpsertSecret, err).ReconcileResult()
	}

	log.Infow("connection secret upserted")
	return workflow.OK().ReconcileResult()
}

// ensureSecret creates or updates the Secret for the given identifiers and connection data
func (r *ConnSecretReconciler) ensureSecret(
	ctx context.Context,
	ids *ConnSecretIdentifiers,
	user *akov2.AtlasDatabaseUser,
	endpoint Endpoint,
	data ConnSecretData,
) error {
	namespace := user.GetNamespace()
	log := r.Log.With("ns", namespace, "project", ids.ProjectID)

	var connectionType string
	switch endpoint.(type) {
	case FederationEndpoint:
		connectionType = "data-federation"
	case DeploymentEndpoint:
		connectionType = "deployment"
	}

	name := K8sConnectionSecretName(ids.ProjectID, ids.ClusterName, ids.DatabaseUsername, connectionType)

	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	// Fills the secret.stringData with the information stored in ConnSecretData
	if err := fillConnSecretData(secret, ids, data, connectionType); err != nil {
		log.Errorw("failed to fill secret data", "reason", workflow.ConnSecretFailedToFillData, "error", err)
		return err
	}

	// Add the owner to be the AtlasDatabaseUser for garbage collection
	if err := controllerutil.SetControllerReference(user, secret, r.Scheme); err != nil {
		log.Errorw("failed to set controller owner", "reason", workflow.ConnSecretFailedToSetOwnerReferences, "error", err)
		return err
	}

	// Upsert the secret in Kubernetes
	if err := r.Client.Patch(ctx, secret, client.Apply, client.ForceOwnership, ConnectionSecretGoFieldOwner); err != nil {
		log.Errorw("failed to create/update secret via apply", "error", err)
		return err
	}

	return nil
}

// fillConnSecretData converts the ConnSecretData into secret.stringData
func fillConnSecretData(secret *corev1.Secret, ids *ConnSecretIdentifiers, data ConnSecretData, endpointType string) error {
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
		TypeLabelKey:         CredLabelVal,
		ProjectLabelKey:      ids.ProjectID,
		ClusterLabelKey:      ids.ClusterName,
		DatabaseUserLabelKey: ids.DatabaseUsername,
	}

	secret.Annotations = map[string]string{
		ConnectionTypelKey: endpointType,
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

// ComputeHash generates a hash based on key connection metadata for immutable secret naming
func ComputeHash(projectID, clusterName, userName, endpointType string) string {
	hashInput := fmt.Sprintf("%s-%s-%s-%s", projectID, clusterName, userName, endpointType)
	hasher := fnv.New64a()

	hasher.Write([]byte(hashInput))
	rawHash := hasher.Sum64()

	encodedHash := rand.SafeEncodeString(fmt.Sprint(rawHash))
	return encodedHash
}

func K8sConnectionSecretName(projectID, clusterName, userName, endpointType string) string {
	hash := ComputeHash(projectID, clusterName, userName, endpointType)
	return fmt.Sprintf("connection-%s", hash)
}
