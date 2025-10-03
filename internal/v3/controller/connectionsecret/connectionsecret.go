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

	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
)

const (
	ProjectLabelKey      = "atlas.mongodb.com/project-id"
	TargetLabelKey       = "atlas.mongodb.com/target-name"
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
	ErrMissingPairing            = errors.New("missing user/connectionTarget")
	ErrAmbiguousPairing          = errors.New("multiple users/connectionTargets with the same name found")
	ErrUnresolvedProjectID       = errors.New("could not resolve the project id")
)

// ConnnSecretIdentifiers stores all the necessary information that will
// be needed to identiy and get a K8s connection secret
type ConnectionSecretIdentifiers struct {
	ProjectID        string
	TargetName       string
	DatabaseUsername string
	ConnectionType   string
}

// ConnectionData contains all connection information required to populate
// the Kubernetes Secret, including standard and SRV URLs and optional Private Link URLs.
type ConnectionSecretData struct {
	DBUserName            string
	Password              string
	ConnectionURL         string
	SrvConnectionURL      string
	PrivateConnectionURLs []PrivateLinkConnectionURLs
}

type PrivateLinkConnectionURLs struct {
	ConnectionURL      string
	SrvConnectionURL   string
	ShardConnectionURL string
}

func (r *ConnectionSecretReconciler) handleUpsert(
	ctx context.Context,
	req ctrl.Request,
	ids *ConnectionSecretIdentifiers,
	user *akov2.AtlasDatabaseUser,
	connectionTarget ConnectionTarget,
) (reconcile.Result, error) {
	log := r.Log.With("ns", req.Namespace, "name", req.Name)
	log.Debugw("Starting handleUpsert", "ConnectionSecretIdentifiers", ids, "AtlasDatabaseUser", user)
	// create the connection data that will populate secret.stringData
	data, err := connectionTarget.BuildConnectionData(ctx, user)
	if err != nil {
		log.Errorw("failed to build connection data", "error", err)
		return workflow.Terminate(workflow.ConnectionSecretFailedToBuildData, err).ReconcileResult()
	}
	log.Debugw("connection data built")
	if err := r.ensureSecret(ctx, ids, user, connectionTarget, data); err != nil {
		return workflow.Terminate(workflow.ConnectionSecretFailedToUpsertSecret, err).ReconcileResult()
	}

	log.Debugw("connection secret upserted")
	return workflow.OK().ReconcileResult()
}

func (r *ConnectionSecretReconciler) handleBatchUpsert(
	ctx context.Context,
	req ctrl.Request,
	user *akov2.AtlasDatabaseUser,
	projectID string,
	connectionTargets []ConnectionTarget,
) (ctrl.Result, error) {
	log := r.Log.With("namespace", req.Namespace, "name", req.Name)

	for _, connectionTarget := range connectionTargets {
		// Construct connection secret identifier for the current connection target.
		connectionSecretIdentifier := ConnectionSecretIdentifiers{
			ProjectID:        projectID,
			TargetName:       connectionTarget.GetName(),
			DatabaseUsername: user.Spec.Username,
			ConnectionType:   connectionTarget.GetConnectionTargetType(),
		}

		// Check if the user is expired.
		expired, err := timeutil.IsExpired(user.Spec.DeleteAfterDate)
		if err != nil {
			log.Errorw("failed to check expiration date on user", "error", err)
			return workflow.Terminate(workflow.ConnectionSecretUserExpired, err).ReconcileResult()
		}

		//Delete secret if the user is expired.
		if expired {
			log.Debugw("User expired; deleting connection secrets")
			result, err := r.handleDelete(ctx, req, &connectionSecretIdentifier)
			if err != nil {
				log.Errorw("Failed to delete connection secrets for expired user", "error", err)
				return result, err
			}
			continue
		}

		// Check that scopes are still valid.
		if !allowsByScopes(user, connectionTarget.GetName(), connectionTarget.GetScopeType()) {
			log.Infow("invalid scope; scheduling deletion of connection secrets")
			result, err := r.handleDelete(ctx, req, &connectionSecretIdentifier)
			if err != nil {
				log.Errorw("failed to delete connection secrets for invalid scope", "error", err)
				return result, err
			}
			continue
		}

		// Ensure connectionTarget readiness
		if !(connectionTarget.IsReady()) {
			continue
		}

		// Handle the upsert of the connection secret.
		result, err := r.handleUpsert(ctx, req, &connectionSecretIdentifier, user, connectionTarget)
		if err != nil {
			log.Errorw("failed to upsert connection secret", "error", err)
			return result, err
		}
	}

	log.Debugw("batch processing completed successfully")
	return workflow.OK().ReconcileResult()
}

func (r *ConnectionSecretReconciler) handleDelete(
	ctx context.Context,
	req ctrl.Request,
	ids *ConnectionSecretIdentifiers,
) (ctrl.Result, error) {
	log := r.Log.With("ns", req.Namespace, "name", req.Name)

	name := K8sConnectionSecretName(ids.ProjectID, ids.TargetName, ids.DatabaseUsername, ids.ConnectionType)
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
		log.Errorw("unable to delete secret", "reason", workflow.ConnectionSecretFailedDeletion, "error", err)
		return workflow.Terminate(workflow.ConnectionSecretFailedDeletion, err).ReconcileResult()
	}

	log.Debugw("connection secret deleted")
	r.EventRecorder.Event(secret, corev1.EventTypeNormal, "Deleted", "ConnectionSecret deleted")
	return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
}

// ensureSecret creates or updates the Secret for the given identifiers and connection data
func (r *ConnectionSecretReconciler) ensureSecret(
	ctx context.Context,
	ids *ConnectionSecretIdentifiers,
	user *akov2.AtlasDatabaseUser,
	connectionTarget ConnectionTarget,
	data ConnectionSecretData,
) error {
	namespace := user.GetNamespace()
	log := r.Log.With("namespace", namespace, "project", ids.ProjectID)

	name := K8sConnectionSecretName(ids.ProjectID, ids.TargetName, ids.DatabaseUsername, connectionTarget.GetConnectionTargetType())

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

	// Fills the secret.stringData with the information stored in ConnectionSecretData
	if err := fillConnSecretData(secret, ids, data, connectionTarget.GetConnectionTargetType()); err != nil {
		log.Errorw("failed to fill secret data", "reason", workflow.ConnectionSecretFailedToFillData, "error", err)
		return err
	}

	// Add the owner to be the AtlasDatabaseUser for garbage collection
	if err := controllerutil.SetControllerReference(user, secret, r.Scheme); err != nil {
		log.Errorw("failed to set controller owner", "reason", workflow.ConnectionSecretFailedToSetOwnerReferences, "error", err)
		return err
	}

	// Upsert the secret in Kubernetes
	if err := r.Client.Patch(ctx, secret, client.Apply, client.ForceOwnership, ConnectionSecretGoFieldOwner); err != nil {
		log.Errorw("failed to create/update secret via apply", "error", err)
		return err
	}

	return nil
}

// fillConnSecretData converts the ConnectionSecretData into secret.stringData
func fillConnSecretData(secret *corev1.Secret, ids *ConnectionSecretIdentifiers, data ConnectionSecretData, connectionTargetType string) error {
	var err error
	username := data.DBUserName
	password := data.Password

	if data.ConnectionURL, err = CreateURL(data.ConnectionURL, username, password); err != nil {
		return err
	}
	if data.SrvConnectionURL, err = CreateURL(data.SrvConnectionURL, username, password); err != nil {
		return err
	}
	for i, pe := range data.PrivateConnectionURLs {
		if data.PrivateConnectionURLs[i].ConnectionURL, err = CreateURL(pe.ConnectionURL, username, password); err != nil {
			return err
		}
		if data.PrivateConnectionURLs[i].SrvConnectionURL, err = CreateURL(pe.SrvConnectionURL, username, password); err != nil {
			return err
		}
		if data.PrivateConnectionURLs[i].ShardConnectionURL, err = CreateURL(pe.ShardConnectionURL, username, password); err != nil {
			return err
		}
	}

	secret.Labels = map[string]string{
		TypeLabelKey:         CredLabelVal,
		ProjectLabelKey:      ids.ProjectID,
		TargetLabelKey:       ids.TargetName,
		DatabaseUserLabelKey: ids.DatabaseUsername,
	}

	secret.Annotations = map[string]string{
		ConnectionTypelKey: connectionTargetType,
	}

	secret.Data = map[string][]byte{
		userNameKey:    []byte(data.DBUserName),
		passwordKey:    []byte(data.Password),
		standardKey:    []byte(data.ConnectionURL),
		standardKeySrv: []byte(data.SrvConnectionURL),
		privateKey:     []byte(""),
		privateSrvKey:  []byte(""),
	}

	for i, pe := range data.PrivateConnectionURLs {
		suffix := ""
		if i != 0 {
			suffix = fmt.Sprint(i)
		}
		secret.Data[privateKey+suffix] = []byte(pe.ConnectionURL)
		secret.Data[privateSrvKey+suffix] = []byte(pe.SrvConnectionURL)
		secret.Data[privateShardKey+suffix] = []byte(pe.ShardConnectionURL)
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
func ComputeHash(projectID, targetName, userName, connectionTargetType string) string {
	hashInput := fmt.Sprintf("%s-%s-%s-%s", projectID, targetName, userName, connectionTargetType)
	hasher := fnv.New64a()

	hasher.Write([]byte(hashInput))
	rawHash := hasher.Sum64()

	encodedHash := rand.SafeEncodeString(fmt.Sprint(rawHash))
	return encodedHash
}

func K8sConnectionSecretName(projectID, targetName, userName, connectionTargetType string) string {
	hash := ComputeHash(projectID, targetName, userName, connectionTargetType)
	return fmt.Sprintf("connection-%s", hash)
}
