/*
Copyright 2022 MongoDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package atlasinstance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"

	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	ptr "k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	dbaasv1alpha1 "github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/dbaas"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasinventory"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// MongoDBAtlasInstanceReconciler reconciles a MongoDBAtlasInstance object
type MongoDBAtlasInstanceReconciler struct {
	Client      client.Client
	Clientset   kubernetes.Interface
	AtlasClient *mongodbatlas.Client
	watch.ResourceWatcher
	Log             *zap.SugaredLogger
	Scheme          *runtime.Scheme
	AtlasDomain     string
	GlobalAPISecret client.ObjectKey
	EventRecorder   record.EventRecorder
}

type InstanceData struct {
	ProjectName      string
	ClusterName      string
	ProviderName     string
	RegionName       string
	InstanceSizeName string
}

// Dev note: duplicate the permissions in both sections below to generate both Role and ClusterRoles

// +kubebuilder:rbac:groups=dbaas.redhat.com,resources=mongodbatlasinstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dbaas.redhat.com,resources=mongodbatlasinstances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// +kubebuilder:rbac:groups=dbaas.redhat.com,namespace=default,resources=mongodbatlasinstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dbaas.redhat.com,namespace=default,resources=mongodbatlasinstances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=get;list;watch

func (r *MongoDBAtlasInstanceReconciler) Reconcile(cx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = cx
	log := r.Log.With("MongoDBAtlasInstance", req.NamespacedName)
	log.Info("Reconciling MongoDBAtlasInstance")

	inst := &dbaas.MongoDBAtlasInstance{}
	if err := r.Client.Get(cx, req.NamespacedName, inst); err != nil {
		if apiErrors.IsNotFound(err) {
			// CR deleted since request queued, child objects getting GC'd, no requeue
			log.Info("MongoDBAtlasInstance resource not found, has been deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Error fetching MongoDBAtlasInstance for reconcile")
		return ctrl.Result{}, err
	}

	// This update will make sure the status is always updated in case of any errors or successful result
	defer func(c *dbaas.MongoDBAtlasInstance) {
		err := r.Client.Status().Update(context.Background(), c)
		if err != nil {
			log.Infof("Could not update resource status:%v", err)
		}
	}(inst)

	inventory := &dbaas.MongoDBAtlasInventory{}
	namespace := inst.Spec.InventoryRef.Namespace
	if len(namespace) == 0 {
		// Namespace is not populated in InventoryRef, default to the request's namespace
		namespace = req.Namespace
	}
	if err := r.Client.Get(cx, types.NamespacedName{Namespace: namespace, Name: inst.Spec.InventoryRef.Name}, inventory); err != nil {
		if apiErrors.IsNotFound(err) {
			// The corresponding inventory is not found, no reqeue.
			log.Info("MongoDBAtlasInventory resource not found, has been deleted")
			result := workflow.InProgress(workflow.MongoDBAtlasInstanceInventoryNotFound, "inventory not found")
			dbaas.SetInstanceCondition(inst, dbaasv1alpha1.DBaaSInstanceProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
			return ctrl.Result{}, nil
		}
		log.Error(err, "Error fetching MongoDBAtlasInventory")
		return ctrl.Result{}, err
	}
	instData, err := getInstanceData(log, inst)
	if err != nil {
		log.Error(err, "Invalid parameters")
		return ctrl.Result{}, err
	}

	atlasProject := getAtlasProject(inst, instData, inventory.Spec.CredentialsRef.Name)

	_, err = controllerutil.CreateOrUpdate(cx, r.Client, atlasProject, func() error {
		return nil
	})
	if err != nil {
		log.Error(err, "Failed to create or update atlas project resource")
		return ctrl.Result{}, err
	}

	atlasCluster := getOwnedAtlasCluster(inst, instData)

	_, err = controllerutil.CreateOrUpdate(cx, r.Client, atlasCluster, instanceMutateFn(inst, atlasCluster, instData))
	if err != nil {
		log.Error(err, "Failed to create or update atlas cluster resource")
		return ctrl.Result{}, err
	}

	// Update the status
	if err := r.Client.Get(cx, types.NamespacedName{Namespace: atlasCluster.Namespace, Name: atlasCluster.Name}, atlasCluster); err != nil {
		if apiErrors.IsNotFound(err) {
			// The corresponding Atlas Cluster is not found, no reqeue.
			log.Info("Atlas Cluster resource not found, has been deleted")
			result := workflow.InProgress(workflow.MongoDBAtlasInstanceClusterNotFound, "Atlas Cluster not found")
			dbaas.SetInstanceCondition(inst, dbaasv1alpha1.DBaaSInstanceProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
			return ctrl.Result{}, nil
		}
		log.Error(err, "Error fetching Atlas Cluster")
		return ctrl.Result{}, err
	}

	atlasConn, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, inventory.ConnectionSecretObjectKey())
	if err != nil {
		result := workflow.Terminate(workflow.MongoDBAtlasInventoryInputError, err.Error())
		dbaas.SetInstanceCondition(inst, dbaasv1alpha1.DBaaSInstanceProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
		return result.ReconcileResult(), nil
	}

	atlasClient := r.AtlasClient
	if atlasClient == nil {
		cl, err := atlas.Client(r.AtlasDomain, atlasConn, log)
		if err != nil {
			result := workflow.Terminate(workflow.MongoDBAtlasInventoryInputError, err.Error())
			dbaas.SetInstanceCondition(inst, dbaasv1alpha1.DBaaSInstanceProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
			return result.ReconcileResult(), nil
		}
		atlasClient = &cl
	}

	result := setInstanceStatusWithClusterInfo(atlasClient, inst, atlasCluster, atlasProject.Spec.Name)
	if !result.IsOk() {
		log.Infof("Error setting instance status: %v", result.Message())
		return ctrl.Result{}, errors.New(result.Message())
	}
	return ctrl.Result{}, nil
}

func (r *MongoDBAtlasInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("MongoDBAtlasInstance", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MongoDBAtlasInstance & handle delete separately
	err = c.Watch(&source.Kind{Type: &dbaas.MongoDBAtlasInstance{}},
		&watch.EventHandlerWithDelete{Controller: r},
		watch.CommonPredicates())
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	// Watch for dependent AtlasCluster resource
	err = c.Watch(
		&source.Kind{
			Type: &v1.AtlasCluster{},
		},
		&handler.EnqueueRequestForOwner{
			OwnerType:    &dbaas.MongoDBAtlasInstance{},
			IsController: true,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// Delete implements a handler for the Delete event.
func (r *MongoDBAtlasInstanceReconciler) Delete(e event.DeleteEvent) error {
	inst, ok := e.Object.(*dbaas.MongoDBAtlasInstance)
	log := r.Log.With("MongoDBAtlasInstance", kube.ObjectKeyFromObject(inst))
	if !ok {
		log.Errorf("Ignoring malformed Delete() call (expected type %T, got %T)", &dbaas.MongoDBAtlasInstance{}, e.Object)
		return nil
	}
	return nil
}

// getAtlasProject returns an AtlasProject for provisioning
// No ownerref is set as the same project can be used to provisiong multiple clusters
func getAtlasProject(instance *dbaas.MongoDBAtlasInstance, data *InstanceData, secretName string) *v1.AtlasProject {
	return &v1.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.ProjectName,
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"created-by": "atlas-operator",
			},
		},
		Spec: v1.AtlasProjectSpec{
			Name:                data.ProjectName,
			ConnectionSecret:    &v1.ResourceRef{Name: secretName},
			ProjectIPAccessList: []project.IPAccessList{},
		},
	}
}

// getAtlasClusterSpec returns the spec for the desired cluster
func getAtlasClusterSpec(instance *dbaas.MongoDBAtlasInstance, data *InstanceData) *v1.AtlasClusterSpec {
	var providerSettingsSpec *v1.ProviderSettingsSpec
	if data.InstanceSizeName == "M0" || data.InstanceSizeName == "M2" || data.InstanceSizeName == "M5" {
		// See Atlas documentation https://docs.atlas.mongodb.com/reference/api/clusters-create-one/
		providerSettingsSpec = &v1.ProviderSettingsSpec{
			InstanceSizeName:    data.InstanceSizeName,
			BackingProviderName: data.ProviderName,
			ProviderName:        provider.ProviderName("TENANT"),
			RegionName:          data.RegionName,
		}
	} else {
		providerSettingsSpec = &v1.ProviderSettingsSpec{
			InstanceSizeName: data.InstanceSizeName,
			ProviderName:     provider.ProviderName(data.ProviderName),
			RegionName:       data.RegionName,
		}
	}
	return &v1.AtlasClusterSpec{
		Project:          v1.ResourceRefNamespaced{Name: data.ProjectName, Namespace: instance.Namespace},
		Name:             data.ClusterName,
		ProviderSettings: providerSettingsSpec,
	}
}

// getOwnedAtlasCluster returns an AtlasCluster object owned by the instance
func getOwnedAtlasCluster(instance *dbaas.MongoDBAtlasInstance, data *InstanceData) *v1.AtlasCluster {
	return &v1.AtlasCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.ClusterName,
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"managed-by":      "atlas-operator",
				"owner":           instance.Name,
				"owner.kind":      instance.Kind,
				"owner.namespace": instance.Namespace,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					UID:                instance.GetUID(),
					APIVersion:         dbaas.GroupVersion.Identifier(),
					BlockOwnerDeletion: ptr.BoolPtr(false),
					Controller:         ptr.BoolPtr(true),
					Kind:               "MongoDBAtlasInstance",
					Name:               instance.Name,
				},
			},
		},
	}
}

func getInstanceData(log *zap.SugaredLogger, inst *dbaas.MongoDBAtlasInstance) (*InstanceData, error) {
	if len(inst.Spec.Name) == 0 {
		log.Errorf("Missing %v", dbaas.ClusterNameKey)
		return nil, fmt.Errorf("missing %v", dbaas.ClusterNameKey)
	}
	projectName, ok := inst.Spec.OtherInstanceParams[dbaas.ProjectNameKey]
	if !ok || len(projectName) == 0 {
		log.Errorf("Missing %v", dbaas.ProjectNameKey)
		return nil, fmt.Errorf("missing %v", dbaas.ProjectNameKey)
	}
	provider := strings.ToUpper(inst.Spec.CloudProvider)
	if len(provider) == 0 {
		provider = "AWS"
		log.Infof("%v is missing, default value of AWS is used", dbaas.CloudProviderKey)
	}
	region := inst.Spec.CloudRegion
	if len(region) == 0 {
		switch provider {
		case "AWS":
			region = "US_EAST_1"
		case "GCE":
			region = "CENTRAL_US"
		case "AZURE":
			region = "US_WEST"
		}
		log.Infof("%v is missing, default value of %s is used", dbaas.CloudProviderKey, region)
	}
	instanceSizeName, ok := inst.Spec.OtherInstanceParams[dbaas.InstanceSizeNameKey]
	if !ok || len(instanceSizeName) == 0 {
		log.Infof("%v is missing, default value of M0 is used", dbaas.InstanceSizeNameKey)
		instanceSizeName = "M0"
	}

	return &InstanceData{
		ProjectName:      projectName,
		ClusterName:      inst.Spec.Name,
		ProviderName:     provider,
		RegionName:       region,
		InstanceSizeName: instanceSizeName,
	}, nil
}

func instanceMutateFn(inst *dbaas.MongoDBAtlasInstance, atlasCluster *v1.AtlasCluster, data *InstanceData) controllerutil.MutateFn {
	return func() error {
		atlasCluster.Spec = *getAtlasClusterSpec(inst, data)
		return nil
	}
}

func setInstanceStatusWithClusterInfo(atlasClient *mongodbatlas.Client, inst *dbaas.MongoDBAtlasInstance, atlasCluster *v1.AtlasCluster, project string) workflow.Result {
	instInfo, result := atlasinventory.GetClusterInfo(atlasClient, project, inst.Spec.Name)
	if result.IsOk() {
		// Set the phase based on AtlasCluster state
		// pending - provisioning not yet started
		// creating - provisioning in progress
		// updating - cluster updating in progress
		// deleting - cluster deletion in progress
		// deleted - cluster has been deleted
		// ready - cluster provisioning complete
		statusFound := false
		inst.Status.InstanceID = instInfo.InstanceID
		inst.Status.InstanceInfo = instInfo.InstanceInfo
		inst.Status.Phase = instInfo.InstanceInfo[dbaas.ProvisionPhaseKey]
		for _, cond := range atlasCluster.Status.Conditions {
			if cond.Type == status.ClusterReadyType {
				statusFound = true
				if cond.Status == "True" {
					dbaas.SetInstanceCondition(inst, dbaasv1alpha1.DBaaSInstanceProviderSyncType, metav1.ConditionStatus(cond.Status), "Ready", cond.Message)
				} else {
					dbaas.SetInstanceCondition(inst, dbaasv1alpha1.DBaaSInstanceProviderSyncType, metav1.ConditionStatus(cond.Status), cond.Reason, cond.Message)
				}
			}
		}
		if !statusFound {
			dbaas.SetInstanceCondition(inst, dbaasv1alpha1.DBaaSInstanceProviderSyncType, metav1.ConditionFalse, "Pending", "Waiting for cluster creation to start")
		}
	} else {
		dbaas.SetInstanceCondition(inst, dbaasv1alpha1.DBaaSInstanceProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
	}
	return result
}
