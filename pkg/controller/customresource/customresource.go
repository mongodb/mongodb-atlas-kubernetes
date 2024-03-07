package customresource

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/Masterminds/semver"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/version"
)

const (
	ResourcePolicyAnnotation       = "mongodb.com/atlas-resource-policy"
	ReconciliationPolicyAnnotation = "mongodb.com/atlas-reconciliation-policy"
	ResourceVersion                = "mongodb.com/atlas-resource-version"
	ResourceVersionOverride        = "mongodb.com/atlas-resource-version-policy"
	ResourcePolicyKeep             = "keep"
	ResourcePolicyDelete           = "delete"
	ReconciliationPolicySkip       = "skip"
	ResourceVersionAllow           = "allow"
)

// PrepareResource queries the Custom Resource 'request.NamespacedName' and populates the 'resource' pointer.
func PrepareResource(ctx context.Context, client client.Client, request reconcile.Request, resource akov2.AtlasCustomResource, log *zap.SugaredLogger) workflow.Result {
	err := client.Get(ctx, request.NamespacedName, resource)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			log.Debugf("Object %s doesn't exist, was it deleted after reconcile request?", request.NamespacedName)
			return workflow.TerminateSilently().WithoutRetry()
		}
		// Error reading the object - requeue the request. Note, that we don't intend to update resource status
		// as most of all it will fail as well.
		log.Errorf("Failed to query object %s: %s", request.NamespacedName, err)
		return workflow.TerminateSilently()
	}

	return workflow.OK()
}

func ValidateResourceVersion(ctx *workflow.Context, resource akov2.AtlasCustomResource, log *zap.SugaredLogger) workflow.Result {
	valid, err := ResourceVersionIsValid(resource)
	if err != nil {
		log.Debugf("resource version for '%s' is invalid", resource.GetName())
		result := workflow.Terminate(workflow.AtlasResourceVersionIsInvalid, err.Error())
		ctx.SetConditionFromResult(status.ResourceVersionStatus, result)
		return result
	}

	if !valid {
		log.Debugf("resource '%s' version mismatch", resource.GetName())
		result := workflow.Terminate(workflow.AtlasResourceVersionMismatch,
			fmt.Sprintf("version of the resource '%s' is higher than the operator version '%s'. ",
				resource.GetName(),
				version.Version))
		ctx.SetConditionFromResult(status.ResourceVersionStatus, result)
		return result
	}

	log.Debugf("resource '%s' version is valid", resource.GetName())
	ctx.SetConditionTrue(status.ResourceVersionStatus)
	return workflow.OK()
}

// MarkReconciliationStarted updates the status of the Atlas Resource to indicate that the Operator has started working on it.
// Internally this will also update the 'observedGeneration' field that notify clients that the resource is being worked on
func MarkReconciliationStarted(client client.Client, resource akov2.AtlasCustomResource, log *zap.SugaredLogger, context context.Context) *workflow.Context {
	updatedConditions := status.EnsureConditionExists(status.FalseCondition(status.ReadyType), resource.GetStatus().GetConditions())

	ctx := workflow.NewContext(log, updatedConditions, context)
	statushandler.Update(ctx, client, nil, resource)

	return ctx
}

func IsResourcePolicyKeepOrDefault(resource akov2.AtlasCustomResource, protectionFlag bool) bool {
	if policy, ok := resource.GetAnnotations()[ResourcePolicyAnnotation]; ok {
		return policy == ResourcePolicyKeep
	}

	return protectionFlag
}

// IsResourcePolicyKeep returns 'true' if the resource should not be removed from Atlas on K8s resource removal.
func IsResourcePolicyKeep(resource akov2.AtlasCustomResource) bool {
	if v, ok := resource.GetAnnotations()[ResourcePolicyAnnotation]; ok {
		return v == ResourcePolicyKeep
	}
	return false
}

// ResourceVersionIsValid returns 'true' if current version of resource is <= current version of the operator.
func ResourceVersionIsValid(resource akov2.AtlasCustomResource) (bool, error) {
	// proceed if label is not present
	resourceVersion, ok := resource.GetLabels()[ResourceVersion]
	if !ok {
		return true, nil
	}

	// error for an invalid resource version (non-semver)
	rv, err := semver.NewVersion(resourceVersion)
	if err != nil {
		return false, fmt.Errorf("%s is not a valid semver version for label %s", resourceVersion, ResourceVersion)
	}

	// no errors for invalid operator version
	ov, err := semver.NewVersion(version.Version)
	if err != nil {
		return true, nil
	}

	// proceed if resource version <= operator version
	if rv.Compare(ov) <= 0 {
		return true, nil
	}

	// proceed if ResourceOverride annotation is present
	if v, ok := resource.GetAnnotations()[ResourceVersionOverride]; ok {
		if v == ResourceVersionAllow {
			return true, nil
		}
	}

	return false, nil
}

// ReconciliationShouldBeSkipped returns 'true' if reconciliation should be skipped for this resource.
func ReconciliationShouldBeSkipped(resource akov2.AtlasCustomResource) bool {
	if v, ok := resource.GetAnnotations()[ReconciliationPolicyAnnotation]; ok {
		return v == ReconciliationPolicySkip
	}
	return false
}

// SetAnnotation sets an annotation in resource while respecting the rest of annotations.
func SetAnnotation(resource akov2.AtlasCustomResource, key, value string) {
	annot := resource.GetAnnotations()
	if annot == nil {
		annot = map[string]string{}
	}
	annot[key] = value
	resource.SetAnnotations(annot)
}
