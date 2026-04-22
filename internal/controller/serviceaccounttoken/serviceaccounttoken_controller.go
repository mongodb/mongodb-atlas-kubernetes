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

package serviceaccounttoken

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/accesstoken"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/secretservice"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/ratelimit"
)

const (
	// refreshFraction picks the re-enqueue point well before expiry: with a 1h
	// token TTL the controller requeues at ~40 minutes so a freshly refreshed
	// token is always available to downstream reconcilers before the previous
	// one approaches expiry.
	refreshFraction = 2.0 / 3.0
	minRequeue      = 10 * time.Second
)

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=get;list;watch;create;update;patch

type ServiceAccountTokenReconciler struct {
	Client        client.Client
	Scheme        *runtime.Scheme
	Log           *zap.SugaredLogger
	TokenProvider TokenProvider

	maxConcurrentReconciles int
}

func NewServiceAccountTokenReconciler(c cluster.Cluster, logger *zap.Logger, atlasDomain string, maxConcurrentReconciles int) *ServiceAccountTokenReconciler {
	return &ServiceAccountTokenReconciler{
		Client:                  c.GetClient(),
		Scheme:                  c.GetScheme(),
		Log:                     logger.Named("serviceaccounttoken").Sugar(),
		TokenProvider:           NewAtlasTokenProvider(atlasDomain),
		maxConcurrentReconciles: maxConcurrentReconciles,
	}
}

func (r *ServiceAccountTokenReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("secret", req.NamespacedName)

	secret := &corev1.Secret{}
	if err := r.Client.Get(ctx, req.NamespacedName, secret); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	clientID := string(secret.Data[reconciler.ClientIDKey])
	clientSecret := string(secret.Data[reconciler.ClientSecretKey])
	// Skip if not Service Account credentials secret
	if clientID == "" || clientSecret == "" {
		return ctrl.Result{}, nil
	}

	log.Info("Reconciling service account credential secret")

	tokenSecretName, err := accesstoken.DeriveSecretName(secret.Namespace, secret.Name)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to derive access token secret name: %w", err)
	}
	tokenRef := client.ObjectKey{Namespace: secret.Namespace, Name: tokenSecretName}

	currentHash, err := accesstoken.CredentialsHash(clientID, clientSecret)
	if err != nil {
		return ctrl.Result{}, err
	}

	existingTokenSecret := &corev1.Secret{}
	if err := r.Client.Get(ctx, tokenRef, existingTokenSecret); err != nil {
		if apierrors.IsNotFound(err) {
			return r.createToken(ctx, log, secret, tokenSecretName, clientID, clientSecret, currentHash)
		}

		return ctrl.Result{}, fmt.Errorf("failed to read access token secret %s: %w", tokenRef.String(), err)
	}

	// If the credential Secret has been rotated since this token was issued,
	// refresh immediately so reconcilers downstream do not keep using a token
	// minted from revoked credentials.
	if string(existingTokenSecret.Data[accesstoken.CredentialsHashKey]) != currentHash {
		log.Info("Credential secret changed since token was issued; refreshing token")
		return r.refreshToken(ctx, log, existingTokenSecret, clientID, clientSecret, currentHash)
	}

	expiryStr := string(existingTokenSecret.Data[accesstoken.ExpiryKey])
	expiry, parseErr := time.Parse(time.RFC3339, expiryStr)
	if parseErr != nil {
		log.Warnw("Failed to parse access token expiry; falling through to refresh",
			"expiry", expiryStr, "error", parseErr)
	} else {
		remaining := time.Until(expiry)
		refreshAt := time.Duration(float64(remaining) * refreshFraction)
		if refreshAt > minRequeue {
			log.Infof("Token still valid, re-enqueueing in %s", refreshAt)
			return ctrl.Result{RequeueAfter: refreshAt}, nil
		}
	}

	return r.refreshToken(ctx, log, existingTokenSecret, clientID, clientSecret, currentHash)
}

func (r *ServiceAccountTokenReconciler) refreshToken(
	ctx context.Context,
	log *zap.SugaredLogger,
	tokenSecret *corev1.Secret,
	clientID, clientSecretValue, credsHash string,
) (ctrl.Result, error) {
	token, expiry, err := r.TokenProvider.FetchToken(ctx, clientID, clientSecretValue)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to refresh access token: %w", err)
	}

	tokenSecret.Data[accesstoken.AccessTokenKey] = []byte(token)
	tokenSecret.Data[accesstoken.ExpiryKey] = []byte(expiry.Format(time.RFC3339))
	tokenSecret.Data[accesstoken.CredentialsHashKey] = []byte(credsHash)

	if err := r.Client.Update(ctx, tokenSecret); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update access token secret: %w", err)
	}

	requeueAfter := requeueDuration(expiry)
	log.Infof("Refreshed access token, re-enqueueing in %s", requeueAfter)

	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

func (r *ServiceAccountTokenReconciler) createToken(
	ctx context.Context,
	log *zap.SugaredLogger,
	credentialSecret *corev1.Secret,
	tokenSecretName string,
	clientID, clientSecretValue, credsHash string,
) (ctrl.Result, error) {
	token, expiry, err := r.TokenProvider.FetchToken(ctx, clientID, clientSecretValue)
	if err != nil {
		log.Errorw("Failed to fetch access token", "error", err)
		return ctrl.Result{}, err
	}

	tokenSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tokenSecretName,
			Namespace: credentialSecret.Namespace,
			Labels: map[string]string{
				secretservice.TypeLabelKey: secretservice.CredLabelVal,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Kind:       "Secret",
					Name:       credentialSecret.Name,
					UID:        credentialSecret.UID,
					Controller: pointer.MakePtr(true),
				},
			},
		},
		Data: map[string][]byte{
			accesstoken.AccessTokenKey:     []byte(token),
			accesstoken.ExpiryKey:          []byte(expiry.Format(time.RFC3339)),
			accesstoken.CredentialsHashKey: []byte(credsHash),
		},
	}

	if err := r.Client.Create(ctx, tokenSecret); err != nil {
		if apierrors.IsAlreadyExists(err) {
			log.Infof("Access token secret %s already exists; will refresh on next reconcile", tokenSecretName)
			return ctrl.Result{RequeueAfter: minRequeue}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to create access token secret: %w", err)
	}

	requeueAfter := requeueDuration(expiry)
	log.Infof("Created access token secret %s, re-enqueueing in %s", tokenSecret.Name, requeueAfter)
	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

func requeueDuration(expiry time.Time) time.Duration {
	remaining := time.Until(expiry)
	d := time.Duration(float64(remaining) * refreshFraction)
	if d < minRequeue {
		return minRequeue
	}
	return d
}

// credentialsLabelPredicate filters Secret events down to those carrying the
// operator's credentials label. The global informer cache in
// internal/operator/builder.go already applies this label selector in
// cluster-wide mode, but not in namespaced mode. Declaring the predicate on
// the controller itself keeps the behavior uniform and cheap.
func credentialsLabelPredicate() predicate.Predicate {
	return predicate.NewPredicateFuncs(func(obj client.Object) bool {
		return obj.GetLabels()[secretservice.TypeLabelKey] == secretservice.CredLabelVal
	})
}

func (r *ServiceAccountTokenReconciler) For() (client.Object, builder.Predicates) {
	return &corev1.Secret{}, builder.WithPredicates(
		credentialsLabelPredicate(),
		predicate.ResourceVersionChangedPredicate{},
	)
}

// MapAccessTokenSecretToOwner returns a Reconcile request for the Connection Secret that owns the given Secret,
// or nil if the Secret has no Secret owner.  This is used by SetupWithManager to re-enqueue the owner when the
// Access Token Secret is deleted or updated, so that an accidentally deleted token Secret is recreated without
// waiting for the next scheduled requeue.
func MapAccessTokenSecretToOwner(_ context.Context, obj client.Object) []reconcile.Request {
	for _, owner := range obj.GetOwnerReferences() {
		if owner.APIVersion == "v1" && owner.Kind == "Secret" {
			return []reconcile.Request{{
				NamespacedName: types.NamespacedName{
					Namespace: obj.GetNamespace(),
					Name:      owner.Name,
				},
			}}
		}
	}
	return nil
}

func (r *ServiceAccountTokenReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(r.For()).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(MapAccessTokenSecretToOwner),
			builder.WithPredicates(credentialsLabelPredicate()),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{
			RateLimiter:             ratelimit.NewRateLimiter[reconcile.Request](),
			SkipNameValidation:      pointer.MakePtr(skipNameValidation),
			MaxConcurrentReconciles: r.maxConcurrentReconciles,
		}).
		Named("serviceaccounttoken").
		Complete(r)
}
