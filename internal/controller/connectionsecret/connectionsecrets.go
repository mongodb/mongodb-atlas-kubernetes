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
	"net/url"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
)

const (
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

// resolveProjectName finds the respective project name for the given projectID in the identifiers
func (r *ConnectionSecretReconciler) resolveProjectName(ctx context.Context, ids ConnSecretIdentifiers, pair *ConnSecretPair) (string, error) {
	if ids.ProjectName != "" {
		return ids.ProjectName, nil
	}

	if pair.Deployment != nil && pair.Deployment.Spec.ProjectRef != nil {
		projectName, err := pair.ResolveProjectNameK8s(ctx, r.Client, pair.Deployment.Namespace)
		if err != nil {
			return "", err
		}
		if projectName != "" {
			return projectName, nil
		}
	}

	if pair.User != nil && pair.User.Spec.ProjectRef != nil {
		projectName, err := pair.ResolveProjectNameK8s(ctx, r.Client, pair.User.Namespace)
		if err != nil {
			return "", err
		}
		if projectName != "" {
			return projectName, nil
		}
	}

	if pair.Deployment != nil {
		connCfg, err := r.ResolveConnectionConfig(ctx, pair.Deployment)
		if err != nil {
			return "", err
		}
		sdkClientSet, err := r.AtlasProvider.SdkClientSet(ctx, connCfg.Credentials, r.Log)
		if err != nil {
			return "", err
		}
		atlasProject, err := r.ResolveProject(ctx, sdkClientSet.SdkClient20250312002, pair.Deployment)
		if err != nil {
			return "", err
		}
		if atlasProject.Name != "" {
			return atlasProject.Name, nil
		}
	}

	if pair.User != nil {
		connCfg, err := r.ResolveConnectionConfig(ctx, pair.User)
		if err != nil {
			return "", err
		}
		sdkClientSet, err := r.AtlasProvider.SdkClientSet(ctx, connCfg.Credentials, r.Log)
		if err != nil {
			return "", err
		}
		atlasProject, err := r.ResolveProject(ctx, sdkClientSet.SdkClient20250312002, pair.User)
		if err != nil {
			return "", err
		}
		if atlasProject.Name != "" {
			return atlasProject.Name, nil
		}
	}

	return "", fmt.Errorf("unable to resolve ProjectName")
}

// handleDelete manages the case where we will delete the connection secret
func (r *ConnectionSecretReconciler) handleDelete(
	ctx context.Context, req ctrl.Request, ids ConnSecretIdentifiers, pair *ConnSecretPair) (ctrl.Result, error) {
	log := r.Log.With("ns", req.Namespace, "name", req.Name)

	// ProjectName is required for ConnectionSecret metadata.name to delete
	projectName, err := r.resolveProjectName(ctx, ids, pair)
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

	// Delete the secret
	if err := r.Client.Delete(ctx, secret); err != nil {
		if apierrors.IsNotFound(err) {
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

// handleUpdate manages the case where we will create or update the connection secret
func (r *ConnectionSecretReconciler) handleUpdate(
	ctx context.Context, req ctrl.Request, ids ConnSecretIdentifiers, pair *ConnSecretPair) (ctrl.Result, error) {
	log := r.Log.With("ns", req.Namespace, "name", req.Name)

	// ProjectName is required for ConnectionSecret metadata.name to create or update
	projectName, err := r.resolveProjectName(ctx, ids, pair)
	if projectName == "" {
		err = fmt.Errorf("project name is empty")
	}
	if err != nil {
		log.Errorw("failed to resolve project name", "reason", workflow.ConnSecretFailedToResolveProjectName, "error", err)
		return workflow.Terminate(workflow.ConnSecretFailedToResolveProjectName, err).ReconcileResult()
	}
	ids.ProjectName = projectName
	log.Debugw("project name resolved for upsert")

	// Build connection data
	data, err := pair.BuildConnectionData(ctx, r.Client)
	if err != nil {
		log.Errorw("failed to build connection data", "reason", workflow.ConnSecretFailedToBuildData, "error", err)
		return workflow.Terminate(workflow.ConnSecretFailedToBuildData, err).ReconcileResult()
	}

	log.Debugw("connection data built")

	name := CreateK8sFormat(projectName, ids.ClusterName, ids.DatabaseUsername)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: req.Namespace,
		},
	}

	// Populate Secret data/labels
	if err := fillConnSecretData(secret, ids, data); err != nil {
		log.Errorw("failed to fill secret data", "reason", workflow.ConnSecretFailedToFillData, "error", err)
		return workflow.Terminate(workflow.ConnSecretFailedToFillData, err).ReconcileResult()
	}

	// Add owners
	if err := controllerutil.SetOwnerReference(pair.User, secret, r.Scheme); err != nil {
		log.Errorw("failed to set controller owner (DatabaseUser)", "reason", workflow.ConnSecretFailedToSetOwnerReferences, "error", err)
		return workflow.Terminate(workflow.ConnSecretFailedToSetOwnerReferences, err).ReconcileResult()
	}

	// Create or Update the Secret
	if err := r.Client.Create(ctx, secret); err != nil {
		if apierrors.IsAlreadyExists(err) {
			// Fetch existing to get ResourceVersion, then update
			current := &corev1.Secret{}
			if getErr := r.Client.Get(ctx, client.ObjectKeyFromObject(secret), current); getErr != nil {
				log.Errorw("failed to fetch existing secret", "reason", workflow.ConnSecretFailedToGetSecret, "error", getErr)
				return workflow.Terminate(workflow.ConnSecretFailedToGetSecret, getErr).ReconcileResult()
			}
			secret.ResourceVersion = current.ResourceVersion
			if updErr := r.Client.Update(ctx, secret); updErr != nil {
				log.Errorw("failed to update secret", "reason", workflow.ConnSecretFailedToUpdateSecret, "error", updErr)
				return workflow.Terminate(workflow.ConnSecretFailedToUpdateSecret, updErr).ReconcileResult()
			}
		} else {
			log.Errorw("failed to create secret", "reason", workflow.ConnSecretFailedToCreateSecret, "error", err)
			return workflow.Terminate(workflow.ConnSecretFailedToCreateSecret, err).ReconcileResult()
		}
	}

	log.Infow("secret created/updated", "reason", workflow.ConnSecretUpsert)
	r.EventRecorder.Event(secret, corev1.EventTypeNormal, "Updated", "ConnectionSecret updated")
	return workflow.OK().ReconcileResult()
}

func fillConnSecretData(secret *corev1.Secret, ids ConnSecretIdentifiers, data ConnSecretData) error {
	var err error
	username := data.DBUserName
	password := data.Password

	if data.ConnURL, err = CreateURL(data.ConnURL, username, password); err != nil {
		return err
	}
	if data.SrvConnURL, err = CreateURL(data.SrvConnURL, username, password); err != nil {
		return err
	}
	for idx, privateConn := range data.PrivateConnURLs {
		if data.PrivateConnURLs[idx].PvtConnURL, err = CreateURL(privateConn.PvtConnURL, username, password); err != nil {
			return err
		}
		if data.PrivateConnURLs[idx].PvtSrvConnURL, err = CreateURL(privateConn.PvtSrvConnURL, username, password); err != nil {
			return err
		}
		if data.PrivateConnURLs[idx].PvtShardConnURL, err = CreateURL(privateConn.PvtShardConnURL, username, password); err != nil {
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

	for idx, privateConn := range data.PrivateConnURLs {
		var suffix string
		if idx != 0 {
			suffix = fmt.Sprint(idx)
		}
		secret.Data[privateKey+suffix] = []byte(privateConn.PvtConnURL)
		secret.Data[privateSrvKey+suffix] = []byte(privateConn.PvtSrvConnURL)
		secret.Data[privateShardKey+suffix] = []byte(privateConn.PvtShardConnURL)
	}

	return nil
}

func CreateURL(connURL, username, password string) (string, error) {
	cs, err := url.Parse(connURL)
	if err != nil {
		return "", err
	}

	cs.User = url.UserPassword(username, password)
	return cs.String(), nil
}
