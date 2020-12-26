/*
Copyright 2020 MongoDB.

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

package atlascluster

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// AtlasClusterReconciler reconciles a AtlasCluster object
type AtlasClusterReconciler struct {
	client.Client
	Log    *zap.SugaredLogger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasclusters/status,verbs=get;update;patch

func (r *AtlasClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.With("atlascluster", req.NamespacedName)
	wctx := workflow.NewContext(log)

	cluster := &mdbv1.AtlasCluster{}
	if err := r.Client.Get(ctx, kube.ObjectKey(req.Namespace, req.Name), cluster); err != nil {
		// TODO make generic (update status, log message)
		log.Error(err, "Failed to read the AtlasCluster")
		return kube.ResultRetry, nil
	}

	log.Infow("-> Starting AtlasCluster reconciliation", "spec", cluster.Spec)

	connection, err := atlas.ReadConnection(wctx, r.Client, "TODO!", cluster.ConnectionSecretObjectKey())
	if err != nil {
		log.Errorf("Failed to read Atlas Connection details: %s", err)
		return kube.ResultRetry, nil
	}

	client, err := atlas.Client(connection, log)
	if err != nil {
		log.Errorf("Failed to read Atlas Connection details: %s", err)
		return kube.ResultRetry, nil
	}

	c, resp, err := client.Clusters.Get(ctx, cluster.Status.GroupID, cluster.Spec.Name)
	if err != nil && resp.StatusCode == http.StatusNotFound {
		c, _, err = client.Clusters.Create(ctx, cluster.Status.GroupID, cluster.Spec.Cluster())
		if err != nil {
			log.Errorf("Cannot get or create cluster %q: %w", cluster.Spec.Name, err)
			return kube.ResultRetry, nil
		}
	}

	switch c.StateName {
	case "IDLE":
		return kube.ResultSuccess, nil

	case "CREATING", "UPDATING", "REPAIRING":
		return kube.ResultRetry, nil

	default:
		log.Errorf("Unknown cluster state %q", c.StateName)
		return kube.ResultRetry, nil
	}
}

func (r *AtlasClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mdbv1.AtlasCluster{}).
		Complete(r)
}
