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

package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/ratelimit"
)

func TestCtrlStateReconciler_SetupWithManager_NoGomock(t *testing.T) {
	fakeMgr := &fakeManager{}
	mock := mockStateReconciler{}
	fakeReconciler := ctrlstate.NewStateReconciler(&mock)
	skipNameValidation := true

	r := newCtrlStateReconciler(fakeReconciler)
	require.NoError(t, r.SetupWithManager(fakeMgr, skipNameValidation))
	require.Equal(t, fakeMgr, mock.ReceivedMgr)
	wantOpts := controller.TypedOptions[reconcile.Request]{
		SkipNameValidation: pointer.MakePtr(skipNameValidation),
		RateLimiter:        ratelimit.NewRateLimiter[reconcile.Request](),
	}
	assert.Equal(t, wantOpts, mock.ReceivedOpts)
}

type fakeManager struct {
	ctrl.Manager
}

type mockStateReconciler struct {
	ctrlstate.StateHandler[mockStateReconciler]
	ReceivedMgr  ctrl.Manager
	ReceivedOpts controller.TypedOptions[reconcile.Request]
}

func (m *mockStateReconciler) SetupWithManager(mgr ctrl.Manager, reconciler reconcile.Reconciler, opts controller.Options) error {
	m.ReceivedMgr = mgr
	m.ReceivedOpts = opts
	return nil
}
