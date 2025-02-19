/*
Copyright (C) MongoDB, Inc. 2025-present.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
*/

package atlasnetworkcontainer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
)

const (
	testProjectID = "project-id"
)

func TestReconcile(t *testing.T) {
	ctx := context.Background()

	testScheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(testScheme))

	tests := map[string]struct {
		request        reconcile.Request
		expectedResult reconcile.Result
		expectedLogs   []string
	}{
		"failed to prepare resource": {
			request:        reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "nc0"}},
			expectedResult: reconcile.Result{},
			expectedLogs: []string{
				"-> Starting AtlasNetworkContainer reconciliation",
				"Object default/nc0 doesn't exist, was it deleted after reconcile request?",
			},
		},
		"prepare resource for reconciliation": {
			request:        reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "nc1"}},
			expectedResult: reconcile.Result{},
			expectedLogs: []string{
				"-> Starting AtlasNetworkContainer reconciliation",
				"-> Skipping AtlasNetworkContainer reconciliation as annotation mongodb.com/atlas-reconciliation-policy=skip",
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			core, logs := observer.New(zap.DebugLevel)
			fakeClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(testNetworkContainer()).
				Build()
			r := &AtlasNetworkContainerReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: fakeClient,
					Log:    zap.New(core).Sugar(),
				},
			}
			result, _ := r.Reconcile(ctx, tc.request)
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, len(tc.expectedLogs), logs.Len())
			for i, log := range logs.All() {
				assert.Equal(t, tc.expectedLogs[i], log.Message)
			}
		})
	}
}

func testNetworkContainer() *akov2.AtlasNetworkContainer {
	return &akov2.AtlasNetworkContainer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nc1",
			Namespace: "default",
			Annotations: map[string]string{
				customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
			},
		},
		Spec: akov2.AtlasNetworkContainerSpec{
			ProjectDualReference: akov2.ProjectDualReference{
				ExternalProjectRef: &akov2.ExternalProjectReference{
					ID: testProjectID,
				},
				ConnectionSecret: &api.LocalObjectReference{},
			},
			Provider: "AWS",
		},
		Status: status.AtlasNetworkContainerStatus{
			ID:          "container-id",
			Provisioned: true,
		},
	}
}
