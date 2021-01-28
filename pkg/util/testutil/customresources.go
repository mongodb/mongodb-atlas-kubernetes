package testutil

import (
	"context"
	"fmt"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/onsi/gomega"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WaitFor waits until the AO Custom Resource reaches some state - this is configured by 'expectedCondition'.
// It's possible to specify optional callbacks to check the state of the object if it hasn't reached the expected condition.
// This allows to validate the object in case it's in "pending" phase.
func WaitFor(k8sClient client.Client, createdResource mdbv1.AtlasCustomResource, expectedCondition status.Condition, check ...func(mdbv1.AtlasCustomResource)) func() bool {
	return func() bool {
		// Get requests may return the cached versions while the passed 'createdResource' is the correct one populated
		// with the updated object. So we take the expected generation from the passed parameter only. Also create the
		// safe copy
		notCachedGeneration := createdResource.GetGeneration()
		if ok := ReadAtlasResource(k8sClient, createdResource); !ok {
			return false
		}
		fmt.Printf("Generation: %+v, observed Generation: %+v, not cached generation: %+v\n", createdResource.GetGeneration(), createdResource.GetStatus().GetObservedGeneration(), notCachedGeneration)
		fmt.Printf("%+v\n", createdResource.GetStatus().GetConditions())
		// Atlas Operator hasn't started working yet
		if createdResource.GetGeneration() != createdResource.GetStatus().GetObservedGeneration() {
			return false
		}

		match, err := gomega.ContainElement(MatchCondition(expectedCondition)).Match(createdResource.GetStatus().GetConditions())
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

func ReadAtlasResource(k8sClient client.Client, createdResource mdbv1.AtlasCustomResource) bool {
	if err := k8sClient.Get(context.Background(), kube.ObjectKeyFromObject(createdResource), createdResource); err != nil {
		// The only error we tolerate is "not found"
		gomega.Expect(apiErrors.IsNotFound(err)).To(gomega.BeTrue(), fmt.Sprintf("Unexpected error: %s", err))
		return false
	}
	return true
}
