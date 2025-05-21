package builder

import (
	controllerruntime "sigs.k8s.io/controller-runtime"
	builder2 "sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	mckpredicate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/predicate"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/ratelimit"
)

func NewDefaultControllerBuilder(mgr controllerruntime.Manager, obj client.Object) *builder2.Builder {
	return controllerruntime.NewControllerManagedBy(mgr).
		For(
			obj,
			builder2.WithPredicates(
				predicate.Or(
					mckpredicate.AnnotationChanged("mongodb.com/reapply-period"),
					predicate.GenerationChangedPredicate{},
				),
				mckpredicate.IgnoreDeletedPredicate[client.Object](),
			),
		).
		WithOptions(controller.Options{
			RateLimiter: ratelimit.NewRateLimiter[reconcile.Request](),
		})
}
