package testutil

import (
	"context"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	. "github.com/onsi/gomega"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WaitFor waits until the AO Custom Resource reaches some state - this is configured by 'expectedCondition'.
// It's possible to specify optional callbacks to check the state of the object if it hasn't reached the expected condition.
// This allows to validate the object in case it's in "pending" phase.
func WaitFor(k8sClient client.Client, resource, createdResource mdbv1.AtlasCustomResource, expectedCondition status.Condition, check ...func(mdbv1.AtlasCustomResource)) func() bool {
	return func() bool {
		if ok := ReadAtlasResource(k8sClient, resource, createdResource); !ok {
			return false
		}
		// Atlas Operator hasn't started working yet
		if createdResource.GetGeneration() != createdResource.GetStatus().GetObservedGeneration() {
			return false
		}

		match, err := ContainElement(MatchCondition(expectedCondition)).Match(createdResource.GetStatus().GetConditions())
		if err != nil || !match {
			if len(check) > 0 {
				for _, f := range check {
					f(createdResource)
				}
			}
			return false
		}
		return true
	}
}

func ReadAtlasResource(k8sClient client.Client, resource, createdResource mdbv1.AtlasCustomResource) bool {
	if err := k8sClient.Get(context.Background(), kube.ObjectKeyFromObject(resource), createdResource); err != nil {
		// The only error we tolerate is "not found"
		Expect(apiErrors.IsNotFound(err)).To(BeTrue())
		return false
	}
	return true
}
