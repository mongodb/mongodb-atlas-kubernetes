// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package indexer

import (
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	zap "go.uber.org/zap"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// nolint:dupl
const ThirdPartyIntegrationBySecretIndex = "thirdpartyintegration.readTokenSecretRef,routingKeySecretRef,secretSecretRef,writeTokenSecretRef,apiKeySecretRef,apiTokenSecretRef,passwordSecretRef,serviceKeySecretRef,urlSecretRef,licenseKeySecretRef,microsoftTeamsWebhookUrlSecretRef"

type ThirdPartyIntegrationBySecretIndexer struct {
	logger *zap.SugaredLogger
}

func NewThirdPartyIntegrationBySecretIndexer(logger *zap.Logger) *ThirdPartyIntegrationBySecretIndexer {
	return &ThirdPartyIntegrationBySecretIndexer{logger: logger.Named(ThirdPartyIntegrationBySecretIndex).Sugar()}
}
func (*ThirdPartyIntegrationBySecretIndexer) Object() client.Object {
	return &v1.ThirdPartyIntegration{}
}
func (*ThirdPartyIntegrationBySecretIndexer) Name() string {
	return ThirdPartyIntegrationBySecretIndex
}

// Keys extracts the index key(s) from the given object
func (i *ThirdPartyIntegrationBySecretIndexer) Keys(object client.Object) []string {
	resource, ok := object.(*v1.ThirdPartyIntegration)
	if !ok {
		i.logger.Errorf("expected *v1.ThirdPartyIntegration but got %T", object)
		return nil
	}
	var keys []string
	if resource.Spec.V20250312.Entry.ReadTokenSecretRef != nil && resource.Spec.V20250312.Entry.ReadTokenSecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.ReadTokenSecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.RoutingKeySecretRef != nil && resource.Spec.V20250312.Entry.RoutingKeySecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.RoutingKeySecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.SecretSecretRef != nil && resource.Spec.V20250312.Entry.SecretSecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.SecretSecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.WriteTokenSecretRef != nil && resource.Spec.V20250312.Entry.WriteTokenSecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.WriteTokenSecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.ApiKeySecretRef != nil && resource.Spec.V20250312.Entry.ApiKeySecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.ApiKeySecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.ApiTokenSecretRef != nil && resource.Spec.V20250312.Entry.ApiTokenSecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.ApiTokenSecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.PasswordSecretRef != nil && resource.Spec.V20250312.Entry.PasswordSecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.PasswordSecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.ServiceKeySecretRef != nil && resource.Spec.V20250312.Entry.ServiceKeySecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.ServiceKeySecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.UrlSecretRef != nil && resource.Spec.V20250312.Entry.UrlSecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.UrlSecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.LicenseKeySecretRef != nil && resource.Spec.V20250312.Entry.LicenseKeySecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.LicenseKeySecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.MicrosoftTeamsWebhookUrlSecretRef != nil && resource.Spec.V20250312.Entry.MicrosoftTeamsWebhookUrlSecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.MicrosoftTeamsWebhookUrlSecretRef.Name)
	}
	return keys
}
func ThirdPartyIntegrationRequests(list *v1.ThirdPartyIntegrationList) []reconcile.Request {
	requests := make([]reconcile.Request, 0, len(list.Items))
	for _, item := range list.Items {
		requests = append(requests, toRequest(&item))
	}
	return requests
}
