/*
Copyright 2021 MongoDB.

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

package atlasconnection

import (
	"context"
	cryptorand "crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"

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
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/source"

	dbaasv1alpha1 "github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/dbaas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

const (
	DBUserNameKey = "username"
	DBPasswordKey = "password"
	digits        = "0123456789"
	specials      = "~=+%^*/()[]{}/!@#$?|"
	all           = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" + digits + specials
)

// MongoDBAtlasConnectionReconciler reconciles a MongoDBAtlasConnection object
type MongoDBAtlasConnectionReconciler struct {
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

// Dev note: duplicate the permissions in both sections below to generate both Role and ClusterRoles

// +kubebuilder:rbac:groups=dbaas.redhat.com,resources=mongodbatlasconnections,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dbaas.redhat.com,resources=mongodbatlasconnections/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// +kubebuilder:rbac:groups=dbaas.redhat.com,namespace=default,resources=mongodbatlasconnections,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dbaas.redhat.com,namespace=default,resources=mongodbatlasconnections/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=get;list;watch

func (r *MongoDBAtlasConnectionReconciler) Reconcile(cx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = cx
	log := r.Log.With("MongoDBAtlasConnection", req.NamespacedName)

	conn := &dbaas.MongoDBAtlasConnection{}
	if err := r.Client.Get(cx, req.NamespacedName, conn); err != nil {
		if apiErrors.IsNotFound(err) {
			// CR deleted since request queued, child objects getting GC'd, no requeue
			log.Info("MongoDBAtlasConnection resource not found, has been deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Error fetching MongoDBAtlasConnection for reconcile")
		return ctrl.Result{}, err
	}

	if isReadyForBinding(conn) {
		// For this release, the InstanceID is mutable, so no reconciliation is needed.
		return ctrl.Result{}, nil
	}

	// This update will make sure the status is always updated in case of any errors or successful result
	defer func(c *dbaas.MongoDBAtlasConnection) {
		err := r.Client.Status().Update(context.Background(), c)
		if err != nil {
			log.Infof("Could not update resource status:%v", err)
		}
	}(conn)

	inventory := &dbaas.MongoDBAtlasInventory{}
	namespace := conn.Spec.InventoryRef.Namespace
	if len(namespace) == 0 {
		// Namespace is not populated in InventoryRef, default to the request's namespace
		namespace = req.Namespace
	}
	if err := r.Client.Get(cx, types.NamespacedName{Namespace: namespace, Name: conn.Spec.InventoryRef.Name}, inventory); err != nil {
		if apiErrors.IsNotFound(err) {
			// The corresponding inventory is not found, no reqeue.
			log.Info("MongoDBAtlasInventory resource not found, has been deleted")
			result := workflow.InProgress(workflow.MongoDBAtlasConnectionInventoryNotFound, "inventory not found")
			dbaas.SetConnectionCondition(conn, dbaasv1alpha1.DBaaSConnectionProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
			return ctrl.Result{}, nil
		}
		log.Error(err, "Error fetching MongoDBAtlasInventory")
		return ctrl.Result{}, err
	}

	if !isInventoryReady(inventory) {
		// The corresponding inventory is not ready yet, requeue
		result := workflow.InProgress(workflow.MongoDBAtlasConnectionInventoryNotReady, "inventory not ready")
		dbaas.SetConnectionCondition(conn, dbaasv1alpha1.DBaaSConnectionProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
		// Requeue
		return result.ReconcileResult(), errors.New("inventory not ready")
	}

	// Retrieve the instance from inventory based on instanceID
	instance := getInstance(inventory, conn.Spec.InstanceID)
	if instance == nil {
		result := workflow.Terminate(workflow.MongoDBAtlasConnectionInstanceIDNotFound, "Atlas database instance not found")
		dbaas.SetConnectionCondition(conn, dbaasv1alpha1.DBaaSConnectionProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
		// No further reconciliation needed
		return result.ReconcileResult(), nil
	}

	projectID := instance.InstanceInfo[dbaas.ProjectIDKey]

	// Now create a configmap for non-sensitive information needed for connecting to the DB instance
	cm := getOwnedConfigMap(conn, instance.InstanceInfo[dbaas.ConnectionStringsStandardSrvKey])
	cmCreated, err := r.Clientset.CoreV1().ConfigMaps(req.Namespace).Create(context.Background(), cm, metav1.CreateOptions{})
	if err != nil {
		result := workflow.Terminate(workflow.MongoDBAtlasConnectionBackendError, err.Error())
		dbaas.SetConnectionCondition(conn, dbaasv1alpha1.DBaaSConnectionProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
		return ctrl.Result{}, fmt.Errorf("failed to create configmap:%w", err)
	}

	// Generate a db username and password
	dbUserName := fmt.Sprintf("atlas-db-user-%v", time.Now().UnixNano())
	dbPassword := generatePassword()
	// Create the db user in Atlas
	res, err := r.createDBUserInAtlas(conn, projectID, dbUserName, dbPassword, inventory, log)
	if err != nil {
		return res, err
	}

	// Now create a secret to store the password locally
	secret := getOwnedSecret(conn, dbUserName, dbPassword)
	secretCreated, err := r.Clientset.CoreV1().Secrets(req.Namespace).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		// Clean up the db user in atlas that was just created
		_ = r.deleteDBUserFromAtlas(instance.InstanceInfo[dbaas.ProjectIDKey], dbUserName, inventory, log)
		result := workflow.Terminate(workflow.MongoDBAtlasConnectionBackendError, err.Error())
		dbaas.SetConnectionCondition(conn, dbaasv1alpha1.DBaaSConnectionProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
		return ctrl.Result{}, fmt.Errorf("failed to create secret:%w", err)
	}

	// Update the status
	dbaas.SetConnectionCondition(conn, dbaasv1alpha1.DBaaSConnectionProviderSyncType, metav1.ConditionTrue, "Ready", "")
	conn.Status.CredentialsRef = &corev1.LocalObjectReference{Name: secretCreated.Name}
	conn.Status.ConnectionInfoRef = &corev1.LocalObjectReference{Name: cmCreated.Name}
	return ctrl.Result{}, nil
}

func (r *MongoDBAtlasConnectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("MongoDBAtlasConnection", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MongoDBAtlasConnection & handle delete separately
	err = c.Watch(&source.Kind{Type: &dbaas.MongoDBAtlasConnection{}}, &watch.EventHandlerWithDelete{Controller: r}, watch.CommonPredicates())
	if err != nil {
		return err
	}

	// Watch for Connection Secrets
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, watch.NewSecretHandler(r.WatchedResources))
	if err != nil {
		return err
	}
	return nil
}

// Delete implements a handler for the Delete event.
func (r *MongoDBAtlasConnectionReconciler) Delete(e event.DeleteEvent) error {
	conn, ok := e.Object.(*dbaas.MongoDBAtlasConnection)
	if !ok {
		r.Log.Errorf("Ignoring malformed Delete() call (expected type %T, got %T)", &dbaas.MongoDBAtlasConnection{}, e.Object)
		return nil
	}
	log := r.Log.With("MongoDBAtlasConnection", kube.ObjectKeyFromObject(conn))
	log.Infow("-> Starting MongoDBAtlasConnection deletion", "spec", conn.Spec)
	_ = r.deleteDBUser(conn, log)
	return nil
}

func (r *MongoDBAtlasConnectionReconciler) deleteDBUser(conn *dbaas.MongoDBAtlasConnection, log *zap.SugaredLogger) error {
	if conn.Status.CredentialsRef == nil {
		log.Infow("No credentialsRef provided. Nothing to delete.")
		return nil
	}
	secret, err := r.getSecret(conn.Namespace, conn.Status.CredentialsRef.Name)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			log.Infow("No secret found for db user credentials. Deletion done.")
			return nil
		}
		return err // For retry
	}

	// Retrieve the db username from the secret
	dbUserName := ""
	dbUser, ok := secret.Data[DBUserNameKey]
	if !ok {
		log.Infow("No db usernmae found for deletion. Deletion done.")
		return nil
	}
	dbUserName = string(dbUser)

	// The corresponding inventory and instanceID in the Spec are immutable, so we can determine the
	// project ID from the instances in the inventory based on the instanceID.
	// First find the corresponding inventory
	inventory := &dbaas.MongoDBAtlasInventory{}
	if err := r.Client.Get(context.Background(), types.NamespacedName{Namespace: conn.Namespace, Name: conn.Spec.InventoryRef.Name}, inventory); err != nil {
		if apiErrors.IsNotFound(err) {
			// CR deleted since request queued, child objects getting GC'd, no requeue
			log.Info("MongoDBAtlasInventory resource not found, has been deleted")
			return nil
		}
		log.Error(err, "Error fetching MongoDBAtlasInventory")
		return err
	}

	if !isInventoryReady(inventory) {
		// The corresponding inventory is not ready yet
		// Nothing to clean up
		return nil
	}

	// Retrieve the instance from inventory based on instanceID
	instance := getInstance(inventory, conn.Spec.InstanceID)
	if instance == nil {
		log.Infow("No instance found in the inventory. Deletion done.")
		return nil
	}

	// Get the projectID from the status
	projectID, ok := instance.InstanceInfo[dbaas.ProjectIDKey]
	if !ok {
		log.Infow("No projectID found. Deletion done.")
		return nil
	}

	// Now delete the db user from Atlas
	if err := r.deleteDBUserFromAtlas(projectID, dbUserName, inventory, log); err != nil {
		log.Errorf("Failed to remove db user from Atlas: %s", err)
		return err
	}

	// db user secret and connectioninfo configmap are automatically deleted based on owner references
	return nil
}

// getSecret gets a secret object
func (r *MongoDBAtlasConnectionReconciler) getSecret(namespace, name string) (*corev1.Secret, error) {
	return r.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

// createDBUserInAtlas create the database user in Atlas
func (r *MongoDBAtlasConnectionReconciler) createDBUserInAtlas(conn *dbaas.MongoDBAtlasConnection, projectID, dbUserName, dbPassword string, inventory *dbaas.MongoDBAtlasInventory, log *zap.SugaredLogger) (ctrl.Result, error) {
	dbUser := &mongodbatlas.DatabaseUser{
		DatabaseName: dbaas.DefaultDatabase,
		GroupID:      projectID,
		Roles: []mongodbatlas.Role{
			{
				DatabaseName: dbaas.DefaultDatabase,
				RoleName:     "readWriteAnyDatabase",
			},
		},
		Username: dbUserName,
		Password: dbPassword,
	}

	atlasConnection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, inventory.ConnectionSecretObjectKey())
	if err != nil {
		result := workflow.Terminate(workflow.MongoDBAtlasConnectionAuthenticationError, err.Error())
		dbaas.SetConnectionCondition(conn, dbaasv1alpha1.DBaaSConnectionProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
		return result.ReconcileResult(), err
	}

	atlasClient := r.AtlasClient
	if atlasClient == nil {
		cl, err := atlas.Client(r.AtlasDomain, atlasConnection, log)
		if err != nil {
			result := workflow.Terminate(workflow.MongoDBAtlasConnectionBackendError, err.Error())
			dbaas.SetInventoryCondition(inventory, dbaasv1alpha1.DBaaSConnectionProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
			return result.ReconcileResult(), nil
		}
		atlasClient = &cl
	}

	// Try to create the db user
	if _, _, err := atlasClient.DatabaseUsers.Create(context.Background(), projectID, dbUser); err != nil {
		result := workflow.Terminate(workflow.DatabaseUserNotCreatedInAtlas, err.Error())
		dbaas.SetConnectionCondition(conn, dbaasv1alpha1.DBaaSConnectionProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// deleteDBUserFromAtlas delete database user from Atlas
func (r *MongoDBAtlasConnectionReconciler) deleteDBUserFromAtlas(projectID, dbUserName string, inventory *dbaas.MongoDBAtlasInventory, log *zap.SugaredLogger) error {
	atlasConnection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, inventory.ConnectionSecretObjectKey())
	if err != nil {
		return err
	}

	atlasClient := r.AtlasClient
	if atlasClient == nil {
		cl, err := atlas.Client(r.AtlasDomain, atlasConnection, log)
		if err != nil {
			return fmt.Errorf("cannot build Atlas client: %w", err)
		}
		atlasClient = &cl
	}

	_, err = atlasClient.DatabaseUsers.Delete(context.Background(), "admin", projectID, dbUserName)
	if err != nil {
		log.Errorw("Cannot delete Atlas database user", "error", err)
		return fmt.Errorf("failed to delete Atlas database user %s: %w", dbUserName, err)
	}
	return nil
}

// getHost retrieves host from the standard srv connection string
func getHost(connectionStringStandardSrv string) string {
	tokens := strings.Split(connectionStringStandardSrv, "//")
	var host string
	if len(tokens) < 2 {
		host = connectionStringStandardSrv
	} else {
		host = tokens[1]
	}
	return host
}

// getOwnedConfigMap returns a configmap object with ownership set
func getOwnedConfigMap(connection *dbaas.MongoDBAtlasConnection, connectionStringStandardSrv string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "atlas-connection-cm-",
			Namespace:    connection.Namespace,
			Labels: map[string]string{
				"managed-by":      "atlas-operator",
				"owner":           connection.Name,
				"owner.kind":      connection.Kind,
				"owner.namespace": connection.Namespace,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					UID:                connection.GetUID(),
					APIVersion:         "dbaas.redhat.com/v1alpha1",
					BlockOwnerDeletion: ptr.BoolPtr(false),
					Controller:         ptr.BoolPtr(true),
					Kind:               "MongoDBAtlasConnection",
					Name:               connection.Name,
				},
			},
		},
		Data: map[string]string{
			dbaas.ProviderKey:           dbaas.Provider,
			dbaas.ServiceBindingTypeKey: dbaas.ServiceBindingType,
			dbaas.HostKey:               getHost(connectionStringStandardSrv),
			dbaas.SrvKey:                "true",
		},
	}
}

// getOwnedSecret returns a secret object for database credentials with ownership set
func getOwnedSecret(connection *dbaas.MongoDBAtlasConnection, username, password string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind: "Opaque",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "atlas-db-user-",
			Namespace:    connection.Namespace,
			Labels: map[string]string{
				"managed-by":      "atlas-operator",
				"owner":           connection.Name,
				"owner.kind":      connection.Kind,
				"owner.namespace": connection.Namespace,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					UID:                connection.GetUID(),
					APIVersion:         "dbaas.redhat.com/v1alpha1",
					BlockOwnerDeletion: ptr.BoolPtr(false),
					Controller:         ptr.BoolPtr(true),
					Kind:               "MongoDBAtlasConnection",
					Name:               connection.Name,
				},
			},
		},
		Data: map[string][]byte{
			DBUserNameKey: []byte(username),
			DBPasswordKey: []byte(password),
		},
	}
}

// isReadyForBinding is the MongoDBAtlasConnection ready for binding already?
func isReadyForBinding(conn *dbaas.MongoDBAtlasConnection) bool {
	cond := dbaas.GetConnectionCondition(conn, dbaasv1alpha1.DBaaSConnectionProviderSyncType)
	return cond != nil && cond.Status == metav1.ConditionTrue
}

// isInventoryReady is the MongoDBAtlasInvenotry ready?
func isInventoryReady(inventory *dbaas.MongoDBAtlasInventory) bool {
	cond := dbaas.GetInventoryCondition(inventory, dbaasv1alpha1.DBaaSInventoryProviderSyncType)
	return cond != nil && cond.Status == metav1.ConditionTrue
}

// getInstance returns an instance from the inventory based on instanceID
func getInstance(inventory *dbaas.MongoDBAtlasInventory, instanceID string) *dbaasv1alpha1.Instance {
	for _, instance := range inventory.Status.Instances {
		if instance.InstanceID == instanceID {
			// Found the instance based on its ID
			return &instance
		}
	}
	return nil
}

// generatePassword generates a random password with at least one digit and one special character.
func generatePassword() string {
	length := 8
	buf := make([]byte, length)
	buf[0] = digits[getRandInt(len(digits))]
	buf[1] = specials[getRandInt(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[getRandInt(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf) // E.g. "3i[g0|)z"
}

func getRandInt(s int) int64 {
	result, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(s)))
	return result.Int64()
}
