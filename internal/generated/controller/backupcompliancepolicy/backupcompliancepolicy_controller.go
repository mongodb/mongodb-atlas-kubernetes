// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backupcompliancepolicy

import (
	zap "go.uber.org/zap"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	cluster "sigs.k8s.io/controller-runtime/pkg/cluster"

	atlas "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	reconciler "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
)

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=backupcompliancepolicys,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=backupcompliancepolicys/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=backupcompliancepolicys/finalizers,verbs=update
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=backupcompliancepolicys,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=backupcompliancepolicys/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=backupcompliancepolicys/finalizers,verbs=update
type BackupCompliancePolicyHandler struct {
	ctrlstate.StateHandler[v1.BackupCompliancePolicy]
	reconciler.AtlasReconciler
	handlerv20250312 *BackupCompliancePolicyHandlerv20250312
}

func NewBackupCompliancePolicyReconciler(c cluster.Cluster, atlasProvider atlas.Provider, logger *zap.Logger, globalSecretRef client.ObjectKey, reapplySupport bool) *ctrlstate.Reconciler[v1.BackupCompliancePolicy] {
	// Create version-specific handlers

	handlerv20250312 := NewBackupCompliancePolicyHandlerv20250312(atlasProvider, c.GetClient(), logger.Named("controllers").Named("BackupCompliancePolicy-v20250312").Sugar(), globalSecretRef)

	// Create main handler dispatcher
	backupcompliancepolicyHandler := &BackupCompliancePolicyHandler{
		AtlasReconciler: reconciler.AtlasReconciler{
			AtlasProvider:   atlasProvider,
			Client:          c.GetClient(),
			GlobalSecretRef: globalSecretRef,
			Log:             logger.Named("controllers").Named("AtlasBackupCompliancePolicy").Sugar(),
		},
		handlerv20250312: handlerv20250312,
	}

	return ctrlstate.NewStateReconciler(backupcompliancepolicyHandler, ctrlstate.WithCluster[v1.BackupCompliancePolicy](c), ctrlstate.WithReapplySupport[v1.BackupCompliancePolicy](reapplySupport))
}
