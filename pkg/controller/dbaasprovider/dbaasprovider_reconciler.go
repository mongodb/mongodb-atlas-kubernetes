/*
 * Copyright (C) 2010 MongoDB, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dbaasprovider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	dbaasoperator "github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
	"go.uber.org/zap"
	v1 "k8s.io/api/apps/v1"
	rbac "k8s.io/api/rbac/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	label "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	resourceName        = "mongodb-atlas-registration"
	dbaasProviderKind   = "DBaaSProvider"
	defaultProviderFile = "./dbaas_provider.yaml"
)

type DBaaSProviderReconciler struct {
	client.Client
	*runtime.Scheme
	Log                      *zap.SugaredLogger
	Clientset                kubernetes.Interface
	operatorNameVersion      string
	operatorInstallNamespace string
	providerFile             string
	cdrChecker               func(groupVersion, kind string) (bool, error)
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch
// +kubebuilder:rbac:groups=dbaas.redhat.com,resources=dbaasproviders,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dbaas.redhat.com,resources=dbaasproviders/status,verbs=get;update;patch

func (r *DBaaSProviderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("DBaaSProvider Reconciler", req.NamespacedName)

	// due to predicate filtering, we'll only reconcile this operator's own deployment when it's seen the first time
	// meaning we have a reconcile entry-point on operator start-up, so now we can create a cluster-scoped resource
	// owned by the operator's ClusterRole to ensure cleanup on uninstall

	dep := &v1.Deployment{}
	if err := r.Get(ctx, req.NamespacedName, dep); err != nil {
		if apiErrors.IsNotFound(err) {
			// CR deleted since request queued, child objects getting GC'd, no requeue
			log.Info("deployment not found, deleted, no requeue")
			return ctrl.Result{}, nil
		}
		// error fetching deployment, requeue and try again
		log.Error(err, "error fetching Deployment CR")
		return ctrl.Result{}, err
	}

	if r.cdrChecker == nil {
		r.cdrChecker = r.checkCrdInstalled
	}
	isCrdInstalled, err := r.cdrChecker(dbaasoperator.GroupVersion.String(), dbaasProviderKind)
	if err != nil {
		log.Error(err, "error discovering GVK")
		return ctrl.Result{}, err
	}
	if !isCrdInstalled {
		log.Info("CRD not found, requeueing with rate limiter")
		// returning with 'Requeue: true' will invoke our custom rate limiter seen in SetupWithManager below
		return ctrl.Result{Requeue: true}, nil
	}

	instance := &dbaasoperator.DBaaSProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceName,
		},
	}
	if err := r.Get(ctx, client.ObjectKeyFromObject(instance), instance); err != nil {
		if apiErrors.IsNotFound(err) {
			log.Info("resource not found, creating now")
			// Atlas registration custom resource isn't present,so create now with ClusterRole owner for GC
			opts := &client.ListOptions{
				LabelSelector: label.SelectorFromSet(map[string]string{
					"olm.owner":      r.operatorNameVersion,
					"olm.owner.kind": "ClusterServiceVersion",
				}),
			}
			clusterRoleList := &rbac.ClusterRoleList{}
			if err := r.List(context.Background(), clusterRoleList, opts); err != nil {
				log.Error(err, "unable to list ClusterRoles to seek potential operand owners")
				return ctrl.Result{}, err
			}

			if len(clusterRoleList.Items) < 1 {
				err := apiErrors.NewNotFound(
					schema.GroupResource{Group: "rbac.authorization.k8s.io", Resource: "ClusterRole"}, "potentialOwner")
				log.Error(err, "could not find ClusterRole owned by CSV to inherit operand")
				return ctrl.Result{}, err
			}

			instance, err := r.getAtlasProviderCR(clusterRoleList)
			if err != nil {
				log.Error(err, "error while constructing new cluster-scoped Atlas DbaaS provider CR")
				return ctrl.Result{}, err
			}
			if err = r.Create(ctx, instance); err != nil {
				log.Error(err, "error while creating new cluster-scoped resource")
				return ctrl.Result{}, err
			} else {
				log.Info("cluster-scoped resource created")
				return ctrl.Result{}, nil
			}
		}
		// error fetching the resource, requeue and try again
		log.Error(err, "error fetching the resource")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// getAtlasProviderCR CR for MongoDB Atlas DBaaS registration
func (r *DBaaSProviderReconciler) getAtlasProviderCR(clusterRoleList *rbac.ClusterRoleList) (*dbaasoperator.DBaaSProvider, error) {
	d, err := ioutil.ReadFile(r.providerFile)
	if err != nil {
		return nil, err
	}
	jsonData, err := yaml.ToJSON(d)
	if err != nil {
		return nil, err
	}
	instance := &dbaasoperator.DBaaSProvider{}
	err = json.Unmarshal(jsonData, instance)
	if err != nil {
		return nil, err
	}
	instance.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion:         "rbac.authorization.k8s.io/v1",
			Kind:               "ClusterRole",
			UID:                clusterRoleList.Items[0].GetUID(), // doesn't really matter which 'item' we use
			Name:               clusterRoleList.Items[0].Name,
			Controller:         pointer.BoolPtr(true),
			BlockOwnerDeletion: pointer.BoolPtr(false),
		},
	}
	return instance, nil
}

// CheckCrdInstalled checks whether dbaas provider CRD, has been created yet
func (r *DBaaSProviderReconciler) checkCrdInstalled(groupVersion, kind string) (bool, error) {
	resources, err := r.Clientset.Discovery().ServerResourcesForGroupVersion(groupVersion)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check DBaaSProvider CRD:%w", err)
	}
	for _, r := range resources.APIResources {
		if r.Kind == kind {
			return true, nil
		}
	}
	return false, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DBaaSProviderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// envVar set in controller-manager's Deployment YAML
	if operatorInstallNamespace, found := os.LookupEnv("OPERATOR_NAMESPACE"); !found {
		return errors.New("OPERATOR_NAMESPACE must be set")
	} else {
		r.operatorInstallNamespace = operatorInstallNamespace
	}

	// envVar set for all operators
	if operatorNameEnvVar, found := os.LookupEnv("OPERATOR_CONDITION_NAME"); !found {
		return errors.New("OPERATOR_CONDITION_NAME must be set")
	} else {
		r.operatorNameVersion = operatorNameEnvVar
	}

	// envVar set for directory for dbaas_provider.yaml
	if providerFileEnvVar, found := os.LookupEnv("DBAAS_PROVIDER_FILE"); !found {
		r.providerFile = defaultProviderFile
	} else {
		r.providerFile = providerFileEnvVar
	}

	customRateLimiter := workqueue.NewItemExponentialFailureRateLimiter(30*time.Second, 30*time.Minute)

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(controller.Options{RateLimiter: customRateLimiter}).
		For(
			&v1.Deployment{},
			builder.WithPredicates(r.ignoreOtherDeployments()),
			builder.OnlyMetadata,
		).
		Complete(r)
}

// ignoreOtherDeployments  only on a 'create' event is issued for the deployment
func (r *DBaaSProviderReconciler) ignoreOtherDeployments() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return r.evaluatePredicateObject(e.Object)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return false
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return false
		},
	}
}

func (r *DBaaSProviderReconciler) evaluatePredicateObject(obj client.Object) bool {
	lbls := obj.GetLabels()
	if obj.GetNamespace() == r.operatorInstallNamespace {
		if val, keyFound := lbls["olm.owner.kind"]; keyFound {
			if val == "ClusterServiceVersion" {
				if val, keyFound := lbls["olm.owner"]; keyFound {
					return val == r.operatorNameVersion
				}
			}
		}
	}
	return false
}
