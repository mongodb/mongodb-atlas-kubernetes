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

package events

import (
	"context"

	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
)

func EventExists(k8sClient client.Client, createdResource api.AtlasCustomResource, expectedType, expectedReason, expectedMessageRegexp string) {
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
