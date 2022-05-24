// Copyright 2021 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atlasconnection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/fgrosse/zaptest"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	krand "k8s.io/apimachinery/pkg/util/rand"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	k8sfakecorev1 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	ktesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	dbaasv1alpha1 "github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/dbaas"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints.
	baseURLPath = "/api-v1"
)

var simulateAtlasDBUserDeleteFailure = false

// setupMockAltasServer sets up a test HTTP server along with a mongodbatlas.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setupMockAltasServer() (client *mongodbatlas.Client, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	router := mux.NewRouter()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, router))

	router.HandleFunc("/api/atlas/v1.0/groups/{group-id}/databaseUsers/{db}/{user}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		groupID, ok := vars["group-id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "group-id is missing in parameters")
			return
		}
		_, ok = vars["db"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "db is missing in parameters")
			return
		}
		_, ok = vars["user"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "user is missing in parameters")
			return
		}
		if simulateAtlasDBUserDeleteFailure {
			// Simulates a db deletion failure
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{\"detail\":\"An invalid group ID %s was specified.\",\"error\":404,\"errorCode\":\"INVALID_GROUP_ID\",\"parameters\":[\"%s\"],\"reason\":\"Not Found\"}", groupID, groupID)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}).Methods(http.MethodDelete)

	router.HandleFunc("/api/atlas/v1.0/groups/{group-id}/databaseUsers", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		groupID, ok := vars["group-id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "group-id is missing in parameters")
			return
		}
		data, err := ioutil.ReadFile(fmt.Sprintf("../../../test/e2e/data/atlasdatabaseuserresp_%s.json", groupID))
		if err == nil {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, string(data))
		} else {
			// Simulates a db creation failure
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{\"detail\":\"An invalid group ID %s was specified.\",\"error\":404,\"errorCode\":\"INVALID_GROUP_ID\",\"parameters\":[\"%s\"],\"reason\":\"Not Found\"}", groupID, groupID)
		}
	}).Methods(http.MethodPost)

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the Atlas client being tested and is
	// configured to use test server.
	client = mongodbatlas.NewClient(nil)
	u, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = u

	return client, server.Close
}

func TestAtlasConnectionReconcile(t *testing.T) {
	s := scheme.Scheme
	utilruntime.Must(scheme.AddToScheme(s))
	utilruntime.Must(mdbv1.AddToScheme(s))
	utilruntime.Must(dbaas.AddToScheme(s))

	atlasClient, teardown := setupMockAltasServer()
	defer teardown()

	logger := zaptest.Logger(t)

	testCase := map[string]struct {
		createConnection    bool
		configMapCreateFail bool
		secretCreateFail    bool
		instanceID          string
		expectedRequeue     bool
		expectedErrString   string
		expectedStatus      string
		expectedReason      string
		inventoryReason     string
		inventoryStatus     string
		instancesPath       string
	}{
		"Nominal": {
			createConnection:    true,
			configMapCreateFail: false,
			secretCreateFail:    false,
			instanceID:          "70b7a72f4877d05880c487ef",
			inventoryReason:     "SyncOK",
			inventoryStatus:     "True",
			instancesPath:       "../../../test/e2e/data/atlasinventoryexpected.json",
			expectedErrString:   "",
			expectedRequeue:     false,
			expectedStatus:      "True",
			expectedReason:      "Ready",
		},
		"ConnectionCRNotFound": {
			createConnection:    false,
			configMapCreateFail: false,
			secretCreateFail:    false,
			instanceID:          "70b7a72f4877d05880c487",
			inventoryReason:     "SyncOK",
			inventoryStatus:     "True",
			instancesPath:       "../../../test/e2e/data/atlasinventoryexpected.json",
			expectedRequeue:     false,
		},
		"InstanceIDNotFound": {
			createConnection:    true,
			configMapCreateFail: false,
			secretCreateFail:    false,
			instanceID:          "70b7a72f4877d05880c487efmissing",
			inventoryReason:     "SyncOK",
			inventoryStatus:     "True",
			instancesPath:       "../../../test/e2e/data/atlasinventoryexpected.json",
			expectedErrString:   "",
			expectedRequeue:     false,
			expectedStatus:      "False",
			expectedReason:      "InstanceIDNotFound",
		},
		"InventoryNotFound": {
			createConnection:    true,
			configMapCreateFail: false,
			secretCreateFail:    false,
			instanceID:          "70b7a72f4877d05880c487ef",
			inventoryReason:     "SyncOK",
			inventoryStatus:     "True",
			instancesPath:       "../../../test/e2e/data/atlasinventoryexpected.json",
			expectedErrString:   "",
			expectedRequeue:     false,
			expectedStatus:      "False",
			expectedReason:      "InventoryNotFound",
		},
		"ConfigMapCreateFail": {
			createConnection:    true,
			configMapCreateFail: true,
			secretCreateFail:    false,
			instanceID:          "70b7a72f4877d05880c487ef",
			inventoryReason:     "SyncOK",
			inventoryStatus:     "True",
			instancesPath:       "../../../test/e2e/data/atlasinventoryexpected.json",
			expectedErrString:   "failed to create configmap",
			expectedRequeue:     false,
			expectedStatus:      "False",
			expectedReason:      "BackendError",
		},
		"SecretCreateFail": {
			createConnection:    true,
			configMapCreateFail: false,
			secretCreateFail:    true,
			instanceID:          "70b7a72f4877d05880c487ef",
			inventoryReason:     "SyncOK",
			inventoryStatus:     "True",
			instancesPath:       "../../../test/e2e/data/atlasinventoryexpected.json",
			expectedErrString:   "failed to create secret",
			expectedRequeue:     false,
			expectedStatus:      "False",
			expectedReason:      "BackendError",
		},
		"AtlasDBUserCreateFail": {
			createConnection:    true,
			configMapCreateFail: false,
			secretCreateFail:    false,
			instanceID:          "60b7a72f4877d05880c487d2",
			inventoryReason:     "SyncOK",
			inventoryStatus:     "True",
			instancesPath:       "../../../test/e2e/data/atlasinventoryexpected.json",
			expectedErrString:   "",
			expectedRequeue:     false,
			expectedStatus:      "False",
			expectedReason:      "DatabaseUserNotCreatedInAtlas",
		},
	}

	for tcName, tc := range testCase {
		t.Run(tcName, func(t *testing.T) {
			instances := []dbaasv1alpha1.Instance{}
			if len(tc.instancesPath) > 0 {
				data, err := ioutil.ReadFile("../../../test/e2e/data/atlasinventoryexpected.json")
				assert.NoError(t, err)
				err = json.Unmarshal(data, &instances)
				assert.NoError(t, err)
			}
			secret := &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind: "Opaque",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("secret-%s", tcName),
					Namespace: "dbaas-operator",
				},
				Data: map[string][]byte{
					"orgId":         []byte("testorgid"),
					"privateApiKey": []byte("testprivatekey"),
					"publicApiKey":  []byte("testpublickey"),
				},
			}

			inventory := &dbaas.MongoDBAtlasInventory{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "dbaas.redhat.com/v1alpha1",
					Kind:       "MongoDBAtlasInventory",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("inventory-%s", tcName),
					Namespace: "dbaas-operator",
				},
				Spec: dbaasv1alpha1.DBaaSInventorySpec{
					CredentialsRef: &dbaasv1alpha1.NamespacedName{
						Name:      fmt.Sprintf("secret-%s", tcName),
						Namespace: "dbaas-operator",
					},
				},
				Status: dbaasv1alpha1.DBaaSInventoryStatus{
					Conditions: []metav1.Condition{
						{
							LastTransitionTime: metav1.Now(),
							Status:             metav1.ConditionStatus(tc.inventoryStatus),
							Reason:             tc.inventoryReason,
							Type:               dbaasv1alpha1.DBaaSInventoryProviderSyncType,
						},
					},
					Instances: instances,
				},
			}

			connection := &dbaas.MongoDBAtlasConnection{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "dbaas.redhat.com/v1alpha1",
					Kind:       "MongoDBAtlasConnection",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("connection-%s", tcName),
					Namespace: "myproject",
				},
				Spec: dbaasv1alpha1.DBaaSConnectionSpec{
					InventoryRef: dbaasv1alpha1.NamespacedName{
						Name:      fmt.Sprintf("inventory-%s", tcName),
						Namespace: "dbaas-operator",
					},
					InstanceID: tc.instanceID,
				},
			}
			objs := []runtime.Object{secret}
			if tcName != "InventoryNotFound" {
				objs = append(objs, inventory)
			}
			if tc.createConnection {
				objs = append(objs, connection)
			}

			// Create a fake client with the objects
			client := fake.NewClientBuilder().WithRuntimeObjects(objs...).WithScheme(s).Build()

			// Mock request to simulate Reconcile() being called on an event for a
			// watched resource .
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      connection.Name,
					Namespace: connection.Namespace,
				},
			}

			// Create a fake clientset
			clientSet := k8sfake.NewSimpleClientset()
			// Fake clientset does not generate resource names based on GenerateName,
			// so we add reactor to generate such names when a secret or configmap is created
			clientSet.PrependReactor("create", "secrets", GenerateNameReactor)
			clientSet.PrependReactor("create", "configmaps", GenerateNameReactor)
			// Simulate configmap or secret creationg failures
			if tc.configMapCreateFail {
				clientSet.CoreV1().(*k8sfakecorev1.FakeCoreV1).PrependReactor("create", "configmaps", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &corev1.ConfigMap{}, errors.New("Error creating configmap")
				})
			} else if tc.secretCreateFail {
				clientSet.CoreV1().(*k8sfakecorev1.FakeCoreV1).PrependReactor("create", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &corev1.Secret{}, errors.New("Error creating secret")
				})
			}
			r := &MongoDBAtlasConnectionReconciler{
				Client:          client,
				Clientset:       clientSet,
				AtlasClient:     atlasClient,
				Scheme:          s,
				Log:             logger.Sugar(),
				ResourceWatcher: watch.NewResourceWatcher(),
			}

			res, err := r.Reconcile(context.Background(), req)
			if err != nil {
				assert.Contains(t, err.Error(), tc.expectedErrString)
			}
			assert.Equal(t, tc.expectedRequeue, res.Requeue)
			connectionUpdated := &dbaas.MongoDBAtlasConnection{}
			err = client.Get(context.Background(),
				types.NamespacedName{
					Name:      connection.Name,
					Namespace: connection.Namespace,
				}, connectionUpdated)
			if tcName == "ConnectionCRNotFound" {
				assert.True(t, apiErrors.IsNotFound(err))
				// Special case: the CR does not exist
				// No further checking is needed
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, connectionUpdated.Status)
			assert.Equal(t, tc.expectedStatus, string(connectionUpdated.Status.Conditions[0].Status))
			assert.Equal(t, tc.expectedReason, connectionUpdated.Status.Conditions[0].Reason)
			if isReadyForBinding(connectionUpdated) {
				assert.NotNil(t, connectionUpdated.Status.ConnectionInfoRef)
				assert.NotNil(t, connectionUpdated.Status.CredentialsRef)
				assert.NotEmpty(t, connectionUpdated.Status.ConnectionInfoRef.Name)
				assert.NotEmpty(t, connectionUpdated.Status.CredentialsRef.Name)
			} else {
				assert.Nil(t, connectionUpdated.Status.ConnectionInfoRef)
				assert.Nil(t, connectionUpdated.Status.CredentialsRef)
			}
		})
	}
}

func TestDBUserDelete(t *testing.T) {
	s := scheme.Scheme
	utilruntime.Must(scheme.AddToScheme(s))
	utilruntime.Must(mdbv1.AddToScheme(s))
	utilruntime.Must(dbaas.AddToScheme(s))

	atlasClient, teardown := setupMockAltasServer()
	defer teardown()

	logger := zaptest.Logger(t)

	testCase := map[string]struct {
		instanceID        string
		expectedErrString string
		inventoryReason   string
		inventoryStatus   string
		instancesPath     string
	}{
		"Nominal": {
			instanceID:        "70b7a72f4877d05880c487ef",
			inventoryReason:   "SyncOK",
			inventoryStatus:   "True",
			instancesPath:     "../../../test/e2e/data/atlasinventoryexpected.json",
			expectedErrString: "",
		},
		"InstanceIDNotFound": {
			instanceID:        "70b7a72f4877d05880c487efmissing",
			inventoryReason:   "SyncOK",
			inventoryStatus:   "True",
			instancesPath:     "../../../test/e2e/data/atlasinventoryexpected.json",
			expectedErrString: "",
		},
		"InventoryNotFound": {
			instanceID:        "70b7a72f4877d05880c487ef",
			inventoryReason:   "SyncOK",
			inventoryStatus:   "True",
			instancesPath:     "../../../test/e2e/data/atlasinventoryexpected.json",
			expectedErrString: "",
		},
		"AtlasDBUserDeleteFail": {
			instanceID:        "70b7a72f4877d05880c487ef",
			inventoryReason:   "SyncOK",
			inventoryStatus:   "True",
			instancesPath:     "../../../test/e2e/data/atlasinventoryexpected.json",
			expectedErrString: "failed to delete Atlas database user",
		},
	}

	for tcName, tc := range testCase {
		t.Run(tcName, func(t *testing.T) {
			instances := []dbaasv1alpha1.Instance{}
			if len(tc.instancesPath) > 0 {
				data, err := ioutil.ReadFile("../../../test/e2e/data/atlasinventoryexpected.json")
				assert.NoError(t, err)
				err = json.Unmarshal(data, &instances)
				assert.NoError(t, err)
			}
			secret := &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind: "Opaque",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("secret-%s", tcName),
					Namespace: "dbaas-operator",
				},
				Data: map[string][]byte{
					"orgId":         []byte("testorgid"),
					"privateApiKey": []byte("testprivatekey"),
					"publicApiKey":  []byte("testpublickey"),
				},
			}

			inventory := &dbaas.MongoDBAtlasInventory{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "dbaas.redhat.com/v1alpha1",
					Kind:       "MongoDBAtlasInventory",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("inventory-%s", tcName),
					Namespace: "dbaas-operator",
				},
				Spec: dbaasv1alpha1.DBaaSInventorySpec{
					CredentialsRef: &dbaasv1alpha1.NamespacedName{
						Name:      fmt.Sprintf("secret-%s", tcName),
						Namespace: "dbaas-operator",
					},
				},
				Status: dbaasv1alpha1.DBaaSInventoryStatus{
					Conditions: []metav1.Condition{
						{
							LastTransitionTime: metav1.Now(),
							Status:             metav1.ConditionStatus(tc.inventoryStatus),
							Reason:             tc.inventoryReason,
							Type:               dbaasv1alpha1.DBaaSInventoryProviderSyncType,
						},
					},
					Instances: instances,
				},
			}
			connection := &dbaas.MongoDBAtlasConnection{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "dbaas.redhat.com/v1alpha1",
					Kind:       "MongoDBAtlasConnection",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("connection-%s", tcName),
					Namespace: "dbaas-operator",
				},
				Spec: dbaasv1alpha1.DBaaSConnectionSpec{
					InventoryRef: dbaasv1alpha1.NamespacedName{
						Name:      fmt.Sprintf("inventory-%s", tcName),
						Namespace: "dbaas-operator",
					},
					InstanceID: tc.instanceID,
				},
			}
			objs := []runtime.Object{secret, connection}
			if tcName != "InventoryNotFound" {
				objs = append(objs, inventory)
			}

			if tcName == "AtlasDBUserDeleteFail" {
				simulateAtlasDBUserDeleteFailure = true
			} else {
				simulateAtlasDBUserDeleteFailure = false
			}
			// Create a fake client with the objects
			client := fake.NewClientBuilder().WithRuntimeObjects(objs...).WithScheme(s).Build()

			// Create a fake clientset
			clientSet := k8sfake.NewSimpleClientset()
			// Fake clientset does not generate resource names based on GenerateName,
			// so we add reactor to generate such names when a secret or configmap is created
			clientSet.PrependReactor("create", "secrets", GenerateNameReactor)
			clientSet.PrependReactor("create", "configmaps", GenerateNameReactor)

			r := &MongoDBAtlasConnectionReconciler{
				Client:          client,
				Clientset:       clientSet,
				AtlasClient:     atlasClient,
				Scheme:          s,
				Log:             logger.Sugar(),
				ResourceWatcher: watch.NewResourceWatcher(),
			}
			// Mock request
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      connection.Name,
					Namespace: connection.Namespace,
				},
			}

			log := r.Log.With("MongoDBAtlasConnection", kube.ObjectKeyFromObject(connection))
			_, err := r.Reconcile(context.Background(), req)
			assert.NoError(t, err)
			connectionUpdated := &dbaas.MongoDBAtlasConnection{}
			err = client.Get(context.Background(),
				types.NamespacedName{
					Name:      connection.Name,
					Namespace: connection.Namespace,
				}, connectionUpdated)
			assert.NoError(t, err)

			// Check db user deletion
			err = r.deleteDBUser(connectionUpdated, log)
			if err != nil {
				assert.Contains(t, err.Error(), tc.expectedErrString)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// GenerateNameReactor sets the metav1.Name of an object, if metav1.GenerateName was used.
// It returns "handled" == false, so the test client can continue to the next ReactionFunc.
func GenerateNameReactor(action ktesting.Action) (bool, runtime.Object, error) {
	obj := action.(ktesting.CreateAction).GetObject().(client.Object)
	if obj.GetName() == "" && obj.GetGenerateName() != "" {
		obj.SetName(fmt.Sprintf("%s%s", obj.GetGenerateName(), krand.String(8)))
	}
	return false, nil, nil
}
