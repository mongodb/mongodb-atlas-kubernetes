// Copyright 2026 MongoDB Inc
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
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/cluster"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret"
	akov2generatedcluster "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/experimental/controller/cluster"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/experimental/controller/databaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/experimental/controller/flexcluster"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/experimental/controller/group"
)

func (r *Registry) experimentalReconcilers(c cluster.Cluster, ap atlas.Provider) ([]Reconciler, error) {
	var reconcilers []Reconciler
	// Add experimental controllers here
	reconcilers = append(reconcilers, connectionsecret.NewConnectionSecretReconciler(c, r.defaultPredicates(), ap, r.logger, r.globalSecretRef))

	groupReconciler, err := group.NewGroupReconciler(c, ap, r.logger, r.globalSecretRef, r.deletionProtection, true, r.defaultPredicates())
	if err != nil {
		return nil, fmt.Errorf("error creating experimental group reconciler: %w", err)
	}

	clusterController, err := akov2generatedcluster.NewClusterReconciler(c, ap, r.logger, r.globalSecretRef, r.deletionProtection, true, r.defaultPredicates())
	if err != nil {
		return nil, fmt.Errorf("error creating experimental cluster reconciler: %w", err)
	}

	flexController, err := flexcluster.NewFlexClusterReconciler(c, ap, r.logger, r.globalSecretRef, r.deletionProtection, true, r.defaultPredicates())
	if err != nil {
		return nil, fmt.Errorf("error creating experimental flex cluster reconciler: %w", err)
	}

	databaseUserReconciler, err := databaseuser.NewDatabaseUserReconciler(c, ap, r.logger, r.globalSecretRef, r.deletionProtection, true, r.defaultPredicates())
	if err != nil {
		return nil, fmt.Errorf("error creating experimental database user reconciler: %w", err)
	}

	reconcilers = append(reconcilers,
		newCtrlStateReconciler(groupReconciler, r.maxConcurrentReconciles),
		newCtrlStateReconciler(clusterController, r.maxConcurrentReconciles),
		newCtrlStateReconciler(flexController, r.maxConcurrentReconciles),
		newCtrlStateReconciler(databaseUserReconciler, r.maxConcurrentReconciles),
	)
	return reconcilers, nil
}
