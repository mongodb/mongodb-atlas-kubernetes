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

package serviceaccount

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/secretservice"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/ratelimit"
)

const (
	clientIDKey     = "clientId"
	clientSecretKey = "clientSecret"
	accessTokenKey  = "accessToken"
	expiryKey       = "expiry"

	refreshFraction = 2.0 / 3.0
	minRequeue      = 10 * time.Second
)

// TokenProvider abstracts the OAuth token acquisition so it can be mocked in tests.
type TokenProvider interface {
	FetchToken(ctx context.Context, clientID, clientSecret string) (token string, expiry time.Time, err error)
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch

type ServiceAccountReconciler struct {
	Client        client.Client
	Scheme        *runtime.Scheme
	Log           *zap.SugaredLogger
	TokenProvider TokenProvider

	maxConcurrentReconciles int
}

func NewServiceAccountReconciler(c cluster.Cluster, logger *zap.Logger, atlasDomain string, maxConcurrentReconciles int) *ServiceAccountReconciler {
	return &ServiceAccountReconciler{
		Client:                  c.GetClient(),
		Scheme:                  c.GetScheme(),
		Log:                     logger.Named("serviceaccount").Sugar(),
		TokenProvider:           NewAtlasTokenProvider(atlasDomain),
		maxConcurrentReconciles: maxConcurrentReconciles,
	}
}

func (r *ServiceAccountReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("secret", req.NamespacedName)

	secret := &corev1.Secret{}
	if err := r.Client.Get(ctx, req.NamespacedName, secret); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !isServiceAccountSecret(secret) {
		return ctrl.Result{}, nil
	}

	log.Info("Reconciling service account credential secret")

	clientID := string(secret.Data[clientIDKey])
	clientSecret := string(secret.Data[clientSecretKey])

	tokenSecretName, hasAnnotation := secret.Annotations[reconciler.AccessTokenAnnotation]

	if hasAnnotation && tokenSecretName != "" {
		existingTokenSecret := &corev1.Secret{}
		tokenRef := client.ObjectKey{Namespace: req.Namespace, Name: tokenSecretName}
		if err := r.Client.Get(ctx, tokenRef, existingTokenSecret); err == nil {
			expiryStr := string(existingTokenSecret.Data[expiryKey])
			expiry, err := time.Parse(time.RFC3339, expiryStr)
			if err == nil {
				remaining := time.Until(expiry)
				refreshAt := time.Duration(float64(remaining) * refreshFraction)
				if refreshAt > minRequeue {
					log.Infof("Token still valid, re-enqueueing in %s", refreshAt)
					return ctrl.Result{RequeueAfter: refreshAt}, nil
				}
			}

			return r.refreshToken(ctx, log, existingTokenSecret, clientID, clientSecret)
		}
	}

	return r.createToken(ctx, log, secret, clientID, clientSecret)
}

func (r *ServiceAccountReconciler) refreshToken(
	ctx context.Context,
	log *zap.SugaredLogger,
	tokenSecret *corev1.Secret,
	clientID, clientSecretValue string,
) (ctrl.Result, error) {
	token, expiry, err := r.TokenProvider.FetchToken(ctx, clientID, clientSecretValue)
	if err != nil {
		log.Errorw("Failed to refresh access token", "error", err)
		return ctrl.Result{RequeueAfter: minRequeue}, err
	}

	tokenSecret.Data[accessTokenKey] = []byte(token)
	tokenSecret.Data[expiryKey] = []byte(expiry.Format(time.RFC3339))

	if err := r.Client.Update(ctx, tokenSecret); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update access token secret: %w", err)
	}

	requeueAfter := requeueDuration(expiry)
	log.Infof("Refreshed access token, re-enqueueing in %s", requeueAfter)

	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

func (r *ServiceAccountReconciler) createToken(
	ctx context.Context,
	log *zap.SugaredLogger,
	credentialSecret *corev1.Secret,
	clientID, clientSecretValue string,
) (ctrl.Result, error) {
	token, expiry, err := r.TokenProvider.FetchToken(ctx, clientID, clientSecretValue)
	if err != nil {
		log.Errorw("Failed to fetch access token", "error", err)
		return ctrl.Result{RequeueAfter: minRequeue}, err
	}

	isController := true
	tokenSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: credentialSecret.Name + "-token-",
			Namespace:    credentialSecret.Namespace,
			Labels: map[string]string{
				secretservice.TypeLabelKey: secretservice.CredLabelVal,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Kind:       "Secret",
					Name:       credentialSecret.Name,
					UID:        credentialSecret.UID,
					Controller: &isController,
				},
			},
		},
		Data: map[string][]byte{
			accessTokenKey: []byte(token),
			expiryKey:      []byte(expiry.Format(time.RFC3339)),
		},
	}

	if err := r.Client.Create(ctx, tokenSecret); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create access token secret: %w", err)
	}

	if credentialSecret.Annotations == nil {
		credentialSecret.Annotations = map[string]string{}
	}
	credentialSecret.Annotations[reconciler.AccessTokenAnnotation] = tokenSecret.Name

	if err := r.Client.Update(ctx, credentialSecret); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to annotate credential secret: %w", err)
	}

	requeueAfter := requeueDuration(expiry)
	log.Infof("Created access token secret %s, re-enqueueing in %s", tokenSecret.Name, requeueAfter)

	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

func (r *ServiceAccountReconciler) For() (client.Object, builder.Predicates) {
	return &corev1.Secret{}, builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})
}

func (r *ServiceAccountReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Secret{}, builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})).
		WithOptions(controller.TypedOptions[reconcile.Request]{
			RateLimiter:             ratelimit.NewRateLimiter[reconcile.Request](),
			SkipNameValidation:      pointer.MakePtr(skipNameValidation),
			MaxConcurrentReconciles: r.maxConcurrentReconciles,
		}).
		Named("serviceaccount").
		Complete(r)
}

func isServiceAccountSecret(secret *corev1.Secret) bool {
	_, hasClientID := secret.Data[clientIDKey]
	_, hasClientSecret := secret.Data[clientSecretKey]
	return hasClientID && hasClientSecret
}

func requeueDuration(expiry time.Time) time.Duration {
	remaining := time.Until(expiry)
	d := time.Duration(float64(remaining) * refreshFraction)
	if d < minRequeue {
		return minRequeue
	}
	return d
}
