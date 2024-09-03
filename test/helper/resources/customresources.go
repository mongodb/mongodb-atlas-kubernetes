package resources

import (
	"context"
	"fmt"

	"github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
)

func CheckCondition(k8sClient client.Client, createdResource api.AtlasCustomResource, expectedCondition api.Condition, checksIfFail ...func(api.AtlasCustomResource)) bool {
	// This is only used from test code
	if ok := ReadAtlasResource(context.Background(), k8sClient, createdResource); !ok {
		return false
	}
	// Atlas Operator hasn't started working yet
	if createdResource.GetGeneration() != createdResource.GetStatus().GetObservedGeneration() {
		return false
	}

	match, err := gomega.ContainElement(conditions.MatchCondition(expectedCondition)).Match(createdResource.GetStatus().GetConditions())
	if err != nil || !match {
		if len(checksIfFail) > 0 {
			for _, f := range checksIfFail {
				f(createdResource)
			}
		}
		return false
	}
	return true
}

func ReadAtlasResource(ctx context.Context, k8sClient client.Client, createdResource api.AtlasCustomResource) bool {
	if err := k8sClient.Get(ctx, kube.ObjectKeyFromObject(createdResource), createdResource); err != nil {
		// The only error we tolerate is "not found"
		gomega.Expect(apierrors.IsNotFound(err)).To(gomega.BeTrue(), fmt.Sprintf("Unexpected error: %s", err))
		return false
	}
	return true
}
