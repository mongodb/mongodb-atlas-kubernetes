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

package builder

import (
	controllerruntime "sigs.k8s.io/controller-runtime"
	ctrlrtbuilder "sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	mckpredicate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/predicate"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/ratelimit"
)

func NewDefaultSetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, obj client.Object) error {
	return controllerruntime.NewControllerManagedBy(mgr).
		For(
			obj,
			ctrlrtbuilder.WithPredicates(
				predicate.Or(
					mckpredicate.AnnotationChanged("mongodb.com/reapply-period"),
					predicate.GenerationChangedPredicate{},
				),
				mckpredicate.IgnoreDeletedPredicate[client.Object](),
			),
		).
		WithOptions(controller.Options{
			RateLimiter: ratelimit.NewRateLimiter[reconcile.Request](),
		}).Complete(rec)
}
