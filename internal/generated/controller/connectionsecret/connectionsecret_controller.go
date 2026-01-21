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
	"reflect"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlcluster "sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/cluster"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/flexcluster"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/target"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	controllerstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	mckpredicate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/predicate"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/ratelimit"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

const (
	FieldOwner            = "mongodb-atlas-kubernetes-connection-secret-handler"
	ConnectionSecretReady = "ConnectionSecretReady"

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
)

// ConnnSecretIdentifiers stores all the necessary information that will
// be needed to identiy and get a K8s connection secret
type ConnectionSecretIdentifiers struct {
	ProjectID        string
	ProjectName      string
	TargetName       string
	DatabaseUsername string
	ConnectionType   string
}

type ConnectionSecretReconciler struct {
	AtlasProvider         atlas.Provider
	Client                client.Client
	Scheme                *runtime.Scheme
	GlobalPredicates      []predicate.Predicate
	ConnectionTargetKinds []target.ConnectionTarget
	GlobalSecretRef       client.ObjectKey
	Logger                *zap.Logger
}

func NewConnectionSecretReconciler(c ctrlcluster.Cluster, predicates []predicate.Predicate, atlasProvider atlas.Provider, logger *zap.Logger, globalSecretRef client.ObjectKey) *ConnectionSecretReconciler {
	r := &ConnectionSecretReconciler{
		Client:           c.GetClient(),
		AtlasProvider:    atlasProvider,
		Scheme:           c.GetScheme(),
		GlobalPredicates: predicates,
		Logger:           logger,
		GlobalSecretRef:  globalSecretRef,
	}

	// Register all the connectionTarget types
	r.ConnectionTargetKinds = []target.ConnectionTarget{
		flexcluster.NewFlexClusterTarget(r.Client),
		cluster.NewClusterTarget(r.Client),
	}

	return r
}

func (r *ConnectionSecretReconciler) For() (client.Object, builder.Predicates) {
	return &generatedv1.DatabaseUser{}, builder.WithPredicates(
		predicate.GenerationChangedPredicate{},
		mckpredicate.IgnoreDeletedPredicate[client.Object](),
	)
}

func (r *ConnectionSecretReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("ConnectionSecret").
		For(r.For()).
		Owns(&corev1.Secret{}, builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})).
		Watches(
			&generatedv1.FlexCluster{},
			handler.EnqueueRequestsFromMapFunc(r.newConnectionTargetMapFunc),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&generatedv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(r.newConnectionTargetMapFunc),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{
			RateLimiter:        ratelimit.NewRateLimiter[reconcile.Request](),
			SkipNameValidation: pointer.MakePtr(skipNameValidation),
		}).
		Complete(r)
}

func (r *ConnectionSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	result, err := r.reconcile(ctx, req)

	if err != nil {
		// Re-fetch the user object to get the latest ResourceVersion before patching.
		// This is necessary because handleBatchUpsert may have already patched the status,
		// and we need the latest ResourceVersion for the patcher to work correctly.
		user := &generatedv1.DatabaseUser{}
		if fetchErr := r.Client.Get(ctx, req.NamespacedName, user); fetchErr != nil {
			// If we can't fetch the user (e.g., it was deleted), return the original error
			return result, err
		}

		errorCondition := metav1.Condition{
			Type:               ConnectionSecretReady,
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             "Error",
			Message:            err.Error(),
		}

		patcher := controllerstate.NewPatcher(user).
			WithFieldOwner(FieldOwner).
			UpdateConditions([]metav1.Condition{errorCondition})

		if err := patcher.Patch(ctx, r.Client); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to update status condition on error: %w", err)
		}
	}

	return result, err
}

func (r *ConnectionSecretReconciler) reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Fetch the DatabaseUser resource.
	user := &generatedv1.DatabaseUser{}
	err := r.Client.Get(ctx, req.NamespacedName, user)
	if apierrors.IsNotFound(err) {
		// object is already gone, nothing to do.
		return reconcile.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("unable to get object: %w", err)
	}

	// Retrieve the project ID associated with the user.
	projectID, err := r.getUserGroupId(ctx, user)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("unable to get project ID: %w", err)
	}

	// Load the connection targets for the project.
	connectionTargetInstances, err := r.listConnectionTargetsByProject(ctx, projectID)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("unable to list connection targets: %w", err)
	}

	// Cleanup stale Secrets
	if err := r.cleanupStaleSecrets(ctx, req.Namespace, connectionTargetInstances, projectID); err != nil {
		return ctrl.Result{}, fmt.Errorf("unable to list connection targets: %w", err)
	}

	// Verify if the AtlasDatabaseUser is ready.
	ready := meta.FindStatusCondition(user.GetConditions(), state.ReadyCondition)
	isUserReady := ready != nil && ready.Status == metav1.ConditionTrue
	if !isUserReady {
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// Delegate the batch upsert logic to handleBatchUpsert.
	return r.handleBatchUpsert(ctx, req, user, projectID, connectionTargetInstances)
}

func (r *ConnectionSecretReconciler) generateConnectionSecretRequests(users []generatedv1.DatabaseUser) []reconcile.Request {
	reqs := make([]reconcile.Request, 0, len(users))
	for _, u := range users {
		reqs = append(reqs, reconcile.Request{
			NamespacedName: types.NamespacedName{Namespace: u.Namespace, Name: u.Name},
		})
	}
	return reqs
}

// listConnectionTargetsByProject retrieves all of the connectionTargets that live under an AtlasProject
func (r *ConnectionSecretReconciler) listConnectionTargetsByProject(ctx context.Context, projectID string) ([]target.ConnectionTargetInstance, error) {
	var out []target.ConnectionTargetInstance

	for _, kind := range r.ConnectionTargetKinds {
		targets, err := kind.ListForProject(ctx, projectID)
		if err != nil {
			return nil, err
		}
		out = append(out, targets...)
	}

	return out, nil
}

func (r *ConnectionSecretReconciler) cleanupStaleSecrets(ctx context.Context, namespace string, connectionTargets []target.ConnectionTargetInstance, projectID string) error {
	// Define the label selector to find relevant secrets.
	labelSelector := &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{Key: TypeLabelKey, Operator: metav1.LabelSelectorOpExists},
			{Key: ProjectLabelKey, Operator: metav1.LabelSelectorOpIn, Values: []string{projectID}},
			{Key: TargetLabelKey, Operator: metav1.LabelSelectorOpExists},
			{Key: DatabaseUserLabelKey, Operator: metav1.LabelSelectorOpExists},
		},
	}

	// Convert the label selector into a client-compatible format.
	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return fmt.Errorf("failed to convert label selector: %w", err)
	}

	// Fetch all secrets in the specified namespace that match the label selector.
	secretList := &corev1.SecretList{}
	if err := r.Client.List(ctx, secretList, client.InNamespace(namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return fmt.Errorf("failed to list secrets: %w", err)
	}

	// Iterate through secrets and delete any that are stale.
	for _, secret := range secretList.Items {
		if err := r.checkAndDeleteStaleSecret(ctx, &secret, connectionTargets); err != nil {
			return err
		}
	}

	return nil
}

// checkAndDeleteStaleSecret deletes a secret if its associated user or connected resource no longer exists.
func (r *ConnectionSecretReconciler) checkAndDeleteStaleSecret(ctx context.Context, secret *corev1.Secret, connectionTargets []target.ConnectionTargetInstance) error {
	pendingDeletion := true
	for _, connectionTarget := range connectionTargets {
		if connectionTarget.GetName() == secret.Labels[TargetLabelKey] && connectionTarget.GetConnectionTargetType() == secret.Annotations[ConnectionTypelKey] {
			pendingDeletion = false
		}
	}

	if pendingDeletion {
		if err := r.Client.Delete(ctx, secret); err != nil {
			if apierrors.IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("failed to delete secret: %w", err)
		}
	}

	return nil
}

// newConnectionTargetMapFunc maps a ConnectionTarget to requests by fetching all AtlasDatabaseUsers and creating a request for each
func (r *ConnectionSecretReconciler) newConnectionTargetMapFunc(ctx context.Context, obj client.Object) []reconcile.Request {
	var ep target.ConnectionTargetInstance

	// Find the matching connection target kind for this object
	for _, kind := range r.ConnectionTargetKinds {
		if wrapped := kind.GetConnectionTargetInstance(obj); wrapped != nil {
			ep = wrapped
			break
		}
	}

	if ep == nil {
		return nil
	}

	projectID := ep.GetProjectID(ctx)
	if projectID == "" {
		return nil
	}

	users := &generatedv1.DatabaseUserList{}
	if err := r.Client.List(ctx, users, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.DatabaseUserByGroupId, projectID),
	}); err != nil {
		return nil
	}

	return r.generateConnectionSecretRequests(users.Items)
}

func (r *ConnectionSecretReconciler) getSDKClientSet(ctx context.Context, databaseuser *generatedv1.DatabaseUser) (*atlas.ClientSet, error) {
	var connectionSecretRef *client.ObjectKey
	if databaseuser.Spec.ConnectionSecretRef != nil {
		connectionSecretRef = &client.ObjectKey{
			Name:      databaseuser.Spec.ConnectionSecretRef.Name,
			Namespace: databaseuser.Namespace,
		}
	}

	connectionConfig, err := reconciler.GetConnectionConfig(ctx, r.Client, connectionSecretRef, &r.GlobalSecretRef)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve Atlas credentials: %w", err)
	}

	clientSet, err := r.AtlasProvider.SdkClientSet(ctx, connectionConfig.Credentials, r.Logger.Sugar())
	if err != nil {
		return nil, fmt.Errorf("failed to setup Atlas SDK client: %w", err)
	}

	return clientSet, nil
}

// ensureSecret creates or updates the Secret for the given identifiers and connection data
func (r *ConnectionSecretReconciler) ensureSecret(
	ctx context.Context,
	ids *ConnectionSecretIdentifiers,
	user *generatedv1.DatabaseUser,
	connectionTarget target.ConnectionTargetInstance,
	connData *data.ConnectionSecret,
) error {
	namespace := user.GetNamespace()
	name := K8sConnectionSecretName(ids.ProjectName, ids.TargetName, ids.DatabaseUsername, connectionTarget.GetConnectionTargetType())

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
	if err := fillConnSecretData(secret, ids, connData, connectionTarget.GetConnectionTargetType()); err != nil {
		return err
	}

	// Add the owner to be the AtlasDatabaseUser for garbage collection
	if err := controllerutil.SetControllerReference(user, secret, r.Scheme); err != nil {
		return err
	}

	currentSecret := &corev1.Secret{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, currentSecret)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	if reflect.DeepEqual(secret.Data, currentSecret.Data) {
		return nil // nothing to do
	}

	// Apply the secret using the new Apply() API (replaces deprecated Patch with client.Apply)
	secretUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(secret)
	if err != nil {
		return err
	}
	secretUnstructuredObj := &unstructured.Unstructured{Object: secretUnstructured}
	applyConfig := client.ApplyConfigurationFromUnstructured(secretUnstructuredObj)
	if err := r.Client.Apply(ctx, applyConfig, client.FieldOwner(ConnectionSecretGoFieldOwner), client.ForceOwnership); err != nil {
		return err
	}

	return nil
}

//nolint:unparam
func (r *ConnectionSecretReconciler) handleDelete(ctx context.Context, req ctrl.Request, ids *ConnectionSecretIdentifiers) (ctrl.Result, error) {
	name := K8sConnectionSecretName(ids.ProjectName, ids.TargetName, ids.DatabaseUsername, ids.ConnectionType)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: req.Namespace,
		},
	}

	// delete secret in k8s
	err := r.Client.Delete(ctx, secret)
	if err != nil && apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to delete connection secret: %w", err)
	}

	return ctrl.Result{}, nil
}

func K8sConnectionSecretName(projectName, targetName, userName, connectionTargetType string) string {
	name := fmt.Sprintf("%s-%s-%s",
		kube.NormalizeIdentifier(projectName),
		kube.NormalizeIdentifier(targetName),
		kube.NormalizeIdentifier(userName))
	return kube.NormalizeIdentifier(name)
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

// fillConnSecretData converts the ConnectionSecretData into secret.stringData
func fillConnSecretData(secret *corev1.Secret, ids *ConnectionSecretIdentifiers, data *data.ConnectionSecret, connectionTargetType string) error {
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
