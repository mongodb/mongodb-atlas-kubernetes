package testutil

import (
	"context"

	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func EventExists(k8sClient client.Client, createdResource mdbv1.AtlasCustomResource, expectedType, expectedReason, expectedMessageRegexp string) {
	eventList := corev1.EventList{}
	gomega.Expect(k8sClient.List(context.Background(), &eventList, client.MatchingFieldsSelector{
		Selector: fields.AndSelectors(
			fields.OneTermEqualSelector("metadata.namespace", createdResource.GetNamespace()),
			fields.OneTermEqualSelector("involvedObject.name", createdResource.GetName()),
			fields.OneTermEqualSelector("reason", expectedReason),
			fields.OneTermEqualSelector("type", expectedType),
		),
	})).To(gomega.Succeed())

	gomega.Expect(eventList.Items).NotTo(gomega.BeEmpty())
	gomega.Expect(eventList.Items[0].Message).To(gomega.MatchRegexp(expectedMessageRegexp))
}
