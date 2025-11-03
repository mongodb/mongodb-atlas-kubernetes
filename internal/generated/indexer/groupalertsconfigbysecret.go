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
const GroupAlertsConfigBySecretIndex = "groupalertsconfig.apiTokenSecretRef,microsoftTeamsWebhookUrlSecretRef,serviceKeySecretRef,victorOpsRoutingKeySecretRef,webhookUrlSecretRef,datadogApiKeySecretRef,notificationTokenSecretRef,opsGenieApiKeySecretRef,victorOpsApiKeySecretRef,webhookSecretSecretRef"

type GroupAlertsConfigBySecretIndexer struct {
	logger *zap.SugaredLogger
}

func NewGroupAlertsConfigBySecretIndexer(logger *zap.Logger) *GroupAlertsConfigBySecretIndexer {
	return &GroupAlertsConfigBySecretIndexer{logger: logger.Named(GroupAlertsConfigBySecretIndex).Sugar()}
}
func (*GroupAlertsConfigBySecretIndexer) Object() client.Object {
	return &v1.GroupAlertsConfig{}
}
func (*GroupAlertsConfigBySecretIndexer) Name() string {
	return GroupAlertsConfigBySecretIndex
}

// Keys extracts the index key(s) from the given object
func (i *GroupAlertsConfigBySecretIndexer) Keys(object client.Object) []string {
	resource, ok := object.(*v1.GroupAlertsConfig)
	if !ok {
		i.logger.Errorf("expected *v1.GroupAlertsConfig but got %T", object)
		return nil
	}
	var keys []string
	if resource.Spec.V20250312.Entry.Notifications.Items.ApiTokenSecretRef != nil && resource.Spec.V20250312.Entry.Notifications.Items.ApiTokenSecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.Notifications.Items.ApiTokenSecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.Notifications.Items.MicrosoftTeamsWebhookUrlSecretRef != nil && resource.Spec.V20250312.Entry.Notifications.Items.MicrosoftTeamsWebhookUrlSecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.Notifications.Items.MicrosoftTeamsWebhookUrlSecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.Notifications.Items.ServiceKeySecretRef != nil && resource.Spec.V20250312.Entry.Notifications.Items.ServiceKeySecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.Notifications.Items.ServiceKeySecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.Notifications.Items.VictorOpsRoutingKeySecretRef != nil && resource.Spec.V20250312.Entry.Notifications.Items.VictorOpsRoutingKeySecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.Notifications.Items.VictorOpsRoutingKeySecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.Notifications.Items.WebhookUrlSecretRef != nil && resource.Spec.V20250312.Entry.Notifications.Items.WebhookUrlSecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.Notifications.Items.WebhookUrlSecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.Notifications.Items.DatadogApiKeySecretRef != nil && resource.Spec.V20250312.Entry.Notifications.Items.DatadogApiKeySecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.Notifications.Items.DatadogApiKeySecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.Notifications.Items.NotificationTokenSecretRef != nil && resource.Spec.V20250312.Entry.Notifications.Items.NotificationTokenSecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.Notifications.Items.NotificationTokenSecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.Notifications.Items.OpsGenieApiKeySecretRef != nil && resource.Spec.V20250312.Entry.Notifications.Items.OpsGenieApiKeySecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.Notifications.Items.OpsGenieApiKeySecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.Notifications.Items.VictorOpsApiKeySecretRef != nil && resource.Spec.V20250312.Entry.Notifications.Items.VictorOpsApiKeySecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.Notifications.Items.VictorOpsApiKeySecretRef.Name)
	}
	if resource.Spec.V20250312.Entry.Notifications.Items.WebhookSecretSecretRef != nil && resource.Spec.V20250312.Entry.Notifications.Items.WebhookSecretSecretRef.Name != "" {
		keys = append(keys, resource.Namespace+"/"+resource.Spec.V20250312.Entry.Notifications.Items.WebhookSecretSecretRef.Name)
	}
	return keys
}
func GroupAlertsConfigRequests(list *v1.GroupAlertsConfigList) []reconcile.Request {
	requests := make([]reconcile.Request, 0, len(list.Items))
	for _, item := range list.Items {
		requests = append(requests, toRequest(&item))
	}
	return requests
}
