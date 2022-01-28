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

package atlasinventory

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/fgrosse/zaptest"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	dbaasv1alpha1 "github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
	dbaas "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/dbaas"
	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints.
	baseURLPath = "/api-v1"
)

var isSingleProjectWithoutCluster = false

// setupMockAltasServer sets up a test HTTP server along with a mongodbatlas.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setupMockAltasServer(t *testing.T) (client *mongodbatlas.Client, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	router := mux.NewRouter()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, router))
	router.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		var data []byte
		var err error
		if isSingleProjectWithoutCluster {
			data, err = ioutil.ReadFile("../../../test/e2e/data/atlasprojectlistresp_single.json")
		} else {
			data, err = ioutil.ReadFile("../../../test/e2e/data/atlasprojectlistresp.json")
		}
		assert.NoError(t, err)
		if err == nil {
			fmt.Fprint(w, string(data))
		}
	}).Methods(http.MethodGet)

	router.HandleFunc("/groups/{group-id}/clusters", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		groupID, ok := vars["group-id"]
		if !ok {
			fmt.Fprint(w, "group-id is missing in parameters")
		}
		data, err := ioutil.ReadFile(fmt.Sprintf("../../../test/e2e/data/atlasclusterlistresp_%s.json", groupID))
		assert.NoError(t, err)
		if err == nil {
			fmt.Fprint(w, string(data))
		}
	}).Methods(http.MethodGet)

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the Atlas client being tested and is
	// configured to use test server.
	client = mongodbatlas.NewClient(nil)
	u, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = u

	return client, server.Close
}

// Test special case for instance discovery: nominal
func TestDiscoverInstancesNominal(t *testing.T) {
	isSingleProjectWithoutCluster = false
	client, teardown := setupMockAltasServer(t)
	defer teardown()

	instances, res := discoverInstances(client)
	assert.True(t, res.IsOk())
	instancesExpected := []dbaasv1alpha1.Instance{}
	dataExpected, err := ioutil.ReadFile("../../../test/e2e/data/atlasinventoryexpected.json")
	assert.NoError(t, err)
	err = json.Unmarshal(dataExpected, &instancesExpected)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(instancesExpected, instances))
}

// Test special case for instance discovery: no instances found
func TestDiscoverInstancesEmpty(t *testing.T) {
	isSingleProjectWithoutCluster = true
	client, teardown := setupMockAltasServer(t)
	defer teardown()

	instances, res := discoverInstances(client)
	assert.True(t, res.IsOk())
	assert.True(t, len(instances) == 0)
}

func TestAtlasInventoryReconcile(t *testing.T) {
	isSingleProjectWithoutCluster = false
	atlasClient, teardown := setupMockAltasServer(t)
	defer teardown()

	s := scheme.Scheme
	utilruntime.Must(scheme.AddToScheme(s))
	utilruntime.Must(mdbv1.AddToScheme(s))
	utilruntime.Must(dbaas.AddToScheme(s))

	client := fake.NewClientBuilder().WithScheme(s).Build()

	logger := zaptest.Logger(t)

	// Create a MongoDBAtlasInventoryReconciler object with the scheme and fake client.
	r := &MongoDBAtlasInventoryReconciler{
		Client:          client,
		AtlasClient:     atlasClient,
		Scheme:          s,
		Log:             logger.Sugar(),
		ResourceWatcher: watch.NewResourceWatcher(),
	}

	testCase := map[string]struct {
		createInventory        bool
		createSecret           bool
		expectedRequeue        bool
		expectedErrString      string
		expectedReadyCondition string
		expectedReasonString   string
		expectedInstancesPath  string
	}{
		"Nominal": {
			createInventory:        true,
			createSecret:           true,
			expectedReadyCondition: "True",
			expectedReasonString:   "SyncOK",
			expectedInstancesPath:  "../../../test/e2e/data/atlasinventoryexpected.json",
			expectedErrString:      "",
			expectedRequeue:        false,
		},
		"InventoryCRNotFound": {
			createInventory:        false,
			createSecret:           false,
			expectedReadyCondition: "",
			expectedReasonString:   "",
			expectedInstancesPath:  "",
			expectedErrString:      "",
			expectedRequeue:        false,
		},
		"SecretNotFound": {
			createInventory:        true,
			createSecret:           false,
			expectedReadyCondition: "False",
			expectedReasonString:   "InputError",
			expectedInstancesPath:  "",
			expectedErrString:      "",
			expectedRequeue:        false,
		},
		"NoCredentialRef": {
			createInventory:        true,
			createSecret:           false,
			expectedReadyCondition: "False",
			expectedReasonString:   "InputError",
			expectedInstancesPath:  "",
			expectedErrString:      "",
			expectedRequeue:        false,
		},
	}
	for tcName, tc := range testCase {
		t.Run(tcName, func(t *testing.T) {
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
			}
			if tcName == "NoCredentialRef" {
				inventory.Spec.CredentialsRef = nil
			}
			secret := &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind: "Opaque",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("secret-%s", tcName),
					Namespace: "dbaas-operator",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"orgId":         []byte("testorgid"),
					"privateApiKey": []byte("testprivatekey"),
					"publicApiKey":  []byte("testpublickey"),
				},
			}
			if tc.createInventory {
				err := client.Create(context.Background(), inventory)
				assert.NoError(t, err)
			}
			if tc.createSecret {
				err := client.Create(context.Background(), secret)
				assert.NoError(t, err)
			}
			// Mock request to simulate Reconcile() being called on an event for a
			// watched resource .
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      inventory.Name,
					Namespace: inventory.Namespace,
				},
			}

			res, err := r.Reconcile(context.Background(), req)
			if err != nil {
				assert.Contains(t, err, tc.expectedErrString)
				// Move on to next test case
				return
			}
			assert.Equal(t, tc.expectedRequeue, res.Requeue)
			if tcName == "InventoryCRNotFound" {
				// Special case: the CR does not exist
				// No further checking is needed
				return
			}
			inventoryUpdated := &dbaas.MongoDBAtlasInventory{}
			err = client.Get(context.Background(),
				types.NamespacedName{
					Name:      inventory.Name,
					Namespace: inventory.Namespace,
				}, inventoryUpdated)
			assert.NoError(t, err)

			if len(tc.expectedReadyCondition) > 0 {
				assert.Equal(t, tc.expectedReadyCondition, string(inventoryUpdated.Status.Conditions[0].Status))
				assert.Equal(t, tc.expectedReasonString, inventoryUpdated.Status.Conditions[0].Reason)
			}

			if len(tc.expectedInstancesPath) == 0 {
				// No need to check discovered instances
				// Move on to next test case
				return
			}
			instancesExpected := []dbaasv1alpha1.Instance{}
			dataExpected, err := ioutil.ReadFile("../../../test/e2e/data/atlasinventoryexpected.json")
			assert.NoError(t, err)
			err = json.Unmarshal(dataExpected, &instancesExpected)
			assert.NoError(t, err)
			assert.True(t, reflect.DeepEqual(instancesExpected, inventoryUpdated.Status.Instances))
		})
	}
}
