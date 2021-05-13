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

package dbaasprovider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/fgrosse/zaptest"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	dbaasoperator "github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"

	dbaas "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/dbaas"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

const mongDBAtlasOperatorLabel = "mongodb-atlas-kubernetes.v0.0.0"

func TestDBaaSProviderCreate(t *testing.T) {
	s := scheme.Scheme
	utilruntime.Must(scheme.AddToScheme(s))
	utilruntime.Must(mdbv1.AddToScheme(s))
	utilruntime.Must(dbaas.AddToScheme(s))
	logger := zaptest.Logger(t)

	testCase := map[string]struct {
		crdChecker        func(groupVersion, kind string) (bool, error)
		expectedRequeue   bool
		expectedErrString string
	}{
		"Nominal": {
			crdChecker:        cdrCheckerOK,
			expectedErrString: "",
			expectedRequeue:   false,
		},
		"CRDCheckFail": {
			crdChecker:        cdrCheckerFail,
			expectedErrString: "failed to check DBaaSProvider CRD",
			expectedRequeue:   false,
		},
		"CRDCheckNotFound": {
			crdChecker:        cdrCheckerNotFound,
			expectedErrString: "",
			expectedRequeue:   true,
		},
	}
	for tcName, tc := range testCase {
		t.Run(tcName, func(t *testing.T) {
			d := &v1.Deployment{}
			data, err := ioutil.ReadFile("../../../test/e2e/data/dummy_deployment.json")
			assert.NoError(t, err)
			err = json.Unmarshal(data, d)
			assert.NoError(t, err)

			gvk := schema.GroupVersionKind{
				Group:   "dbaas.redhat.com",
				Version: "v1alpha1",
				Kind:    "DBaaSProvider",
			}
			clusterrole := &rbac.ClusterRole{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "rbac.authorization.k8s.io/v1",
					Kind:       "ClusterRole",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("mongodb-atlas-kubernetes-%s", tcName),
					Namespace: "dbaas-operator",
					Labels: map[string]string{
						"olm.owner":      mongDBAtlasOperatorLabel,
						"olm.owner.kind": "ClusterServiceVersion",
					},
				},
			}

			// Register DBaaSProvider CRD with the scheme
			s.AddKnownTypeWithName(gvk, &dbaasoperator.DBaaSProvider{})
			// Create a fake client with the objects
			client := fake.NewClientBuilder().WithRuntimeObjects(d, clusterrole).WithScheme(s).Build()
			// Create a fake clientset
			clientSet := k8sfake.NewSimpleClientset()

			// Mock request to simulate Reconcile() being called on an event for a
			// watched resource.
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      d.Name,
					Namespace: d.Namespace,
				},
			}
			// Create a MongoDBAtlasInventoryReconciler object with the scheme and fake client.
			r := &DBaaSProviderReconciler{
				Client:              client,
				Clientset:           clientSet,
				Scheme:              s,
				Log:                 logger.Sugar(),
				cdrChecker:          tc.crdChecker,
				operatorNameVersion: mongDBAtlasOperatorLabel,
				providerFile:        "../../../config/dbaasprovider/dbaas_provider.yaml",
			}

			res, err := r.Reconcile(context.Background(), req)
			if err != nil {
				assert.Contains(t, err.Error(), tc.expectedErrString)
				// Move on to next test case
				return
			}
			assert.Equal(t, tc.expectedRequeue, res.Requeue)
			if res.Requeue {
				// Move on to next test case
				return
			}
			// Check the dbaasprovider CR has been created
			instance := &dbaasoperator.DBaaSProvider{}
			err = r.Client.Get(context.Background(), types.NamespacedName{Name: resourceName}, instance)
			assert.NoError(t, err)
		})
	}
}

// Mock functions to to check DBaaSProvider CRD
func cdrCheckerOK(groupVersion, kind string) (bool, error) {
	return true, nil
}
func cdrCheckerNotFound(groupVersion, kind string) (bool, error) {
	return false, nil
}
func cdrCheckerFail(groupVersion, kind string) (bool, error) {
	return false, errors.New("failed to check DBaaSProvider CRD:")
}
