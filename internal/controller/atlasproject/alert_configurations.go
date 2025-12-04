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

package atlasproject

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/paging"
)

func (r *AtlasProjectReconciler) ensureAlertConfigurations(service *workflow.Context, project *akov2.AtlasProject) workflow.DeprecatedResult {
	service.Log.Debug("starting alert configurations processing")
	defer service.Log.Debug("finished alert configurations processing")

	if project.Spec.AlertConfigurationSyncEnabled {
		specToSync := project.Spec.DeepCopy().AlertConfigurations

		alertConfigurationCondition := api.AlertConfigurationReadyType
		if len(specToSync) == 0 {
			service.UnsetCondition(alertConfigurationCondition)
			return workflow.OK()
		}
		err := r.readAlertConfigurationsSecretsData(project, service, specToSync)
		if err != nil {
			service.SetConditionFalseMsg(alertConfigurationCondition, err.Error())
			return workflow.Terminate(workflow.Internal, err)
		}
		result := syncAlertConfigurations(service, project.ID(), specToSync)
		if !result.IsOk() {
			service.SetConditionFromResult(alertConfigurationCondition, result)
			return result
		}
		service.SetConditionTrue(alertConfigurationCondition)
		return result
	}
	service.UnsetCondition(api.AlertConfigurationReadyType)
	service.Log.Debugf("Alert configuration sync is disabled for project %s", project.Name)
	return workflow.OK()
}

// This method reads secrets refs and fills the secret data for the related Notification
func (r *AtlasProjectReconciler) readAlertConfigurationsSecretsData(project *akov2.AtlasProject, service *workflow.Context, alertConfigs []akov2.AlertConfiguration) error {
	projectNs := project.Namespace

	for i := 0; i < len(alertConfigs); i++ {
		ac := &alertConfigs[i]
		for j := 0; j < len(ac.Notifications); j++ {
			nf := &ac.Notifications[j]
			switch {
			case nf.APITokenRef.Name != "":
				token, err := readNotificationSecret(service.Context, r.Client, nf.APITokenRef, projectNs, "APIToken")
				if err != nil {
					return err
				}
				nf.SetAPIToken(token)
			case nf.DatadogAPIKeyRef.Name != "":
				token, err := readNotificationSecret(service.Context, r.Client, nf.DatadogAPIKeyRef, projectNs, "DatadogAPIKey")
				if err != nil {
					return err
				}
				nf.SetDatadogAPIKey(token)
			case nf.FlowdockAPITokenRef.Name != "":
				token, err := readNotificationSecret(service.Context, r.Client, nf.FlowdockAPITokenRef, projectNs, "FlowdockAPIToken")
				if err != nil {
					return err
				}
				nf.SetFlowdockAPIToken(token)
			case nf.OpsGenieAPIKeyRef.Name != "":
				token, err := readNotificationSecret(service.Context, r.Client, nf.OpsGenieAPIKeyRef, projectNs, "OpsGenieAPIKey")
				if err != nil {
					return err
				}
				nf.SetOpsGenieAPIKey(token)
			case nf.ServiceKeyRef.Name != "":
				token, err := readNotificationSecret(service.Context, r.Client, nf.ServiceKeyRef, projectNs, "ServiceKey")
				if err != nil {
					return err
				}
				nf.SetServiceKey(token)
			case nf.VictorOpsSecretRef.Name != "":
				token, err := readNotificationSecret(service.Context, r.Client, nf.VictorOpsSecretRef, projectNs, "VictorOpsAPIKey")
				if err != nil {
					return err
				}
				nf.SetVictorOpsAPIKey(token)
				token, err = readNotificationSecret(service.Context, r.Client, nf.VictorOpsSecretRef, projectNs, "VictorOpsRoutingKey")
				if err != nil {
					return err
				}
				nf.SetVictorOpsRoutingKey(token)
			}
		}
	}
	return nil
}

func readNotificationSecret(ctx context.Context, kubeClient client.Client, res common.ResourceRefNamespaced, parentNamespace string, fieldName string) (string, error) {
	secret := &v1.Secret{}
	var ns string
	if res.Namespace == "" {
		ns = parentNamespace
	} else {
		ns = res.Namespace
	}

	secretObj := client.ObjectKey{Name: res.Name, Namespace: ns}

	if err := kubeClient.Get(ctx, secretObj, secret); err != nil {
		return "", err
	}
	val, exists := secret.Data[fieldName]
	switch {
	case !exists:
		return "", fmt.Errorf("secret '%s/%s' doesn't contain '%s' parameter", ns, res.Name, fieldName)
	case len(val) == 0:
		return "", fmt.Errorf("secret '%s/%s' contains an empty value for '%s' parameter", ns, res.Name, fieldName)
	}
	return string(val), nil
}

func syncAlertConfigurations(service *workflow.Context, groupID string, alertSpec []akov2.AlertConfiguration) workflow.DeprecatedResult {
	logger := service.Log
	existedAlertConfigs, err := paging.ListAll(service.Context, func(ctx context.Context, pageNum int) (paging.Response[admin.GroupAlertsConfig], *http.Response, error) {
		return service.SdkClientSet.SdkClient20250312009.AlertConfigurationsApi.
			ListAlertConfigs(ctx, groupID).
			PageNum(pageNum).
			Execute()
	})
	if err != nil {
		logger.Errorf("failed to list alert configurations: %v", err)
		return workflow.Terminate(workflow.ProjectAlertConfigurationIsNotReadyInAtlas, fmt.Errorf("failed to list alert configurations: %w", err))
	}

	diff := sortAlertConfigs(logger, alertSpec, existedAlertConfigs)
	logger.Debugf("to create %v, to create statuses %v, to delete %v", len(diff.Create), len(diff.CreateStatus), len(diff.Delete))

	newStatuses := createAlertConfigs(service, groupID, diff.Create)

	for _, existedAlertConfig := range diff.CreateStatus {
		newStatuses = append(newStatuses, status.ParseAlertConfiguration(existedAlertConfig, service.Log))
	}

	service.EnsureStatusOption(status.AtlasProjectSetAlertConfigOption(&newStatuses))

	err = deleteAlertConfigs(service, groupID, diff.Delete)
	if err != nil {
		return workflow.Terminate(workflow.ProjectAlertConfigurationIsNotReadyInAtlas, fmt.Errorf("failed to delete alert configurations: %w", err))
	}

	return checkAlertConfigurationStatuses(newStatuses)
}

func checkAlertConfigurationStatuses(statuses []status.AlertConfiguration) workflow.DeprecatedResult {
	for _, alertConfigurationStatus := range statuses {
		if alertConfigurationStatus.ErrorMessage != "" {
			return workflow.Terminate(workflow.ProjectAlertConfigurationIsNotReadyInAtlas,
				fmt.Errorf("failed to create alert configuration: %s", alertConfigurationStatus.ErrorMessage))
		}
	}
	return workflow.OK()
}

func deleteAlertConfigs(workflowCtx *workflow.Context, groupID string, alertConfigIDs []string) error {
	logger := workflowCtx.Log
	for _, alertConfigID := range alertConfigIDs {
		_, err := workflowCtx.SdkClientSet.SdkClient20250312009.AlertConfigurationsApi.
			DeleteAlertConfig(workflowCtx.Context, groupID, alertConfigID).
			Execute()
		if err != nil {
			logger.Errorf("failed to delete alert configuration: %v", err)
			return err
		}
		logger.Infof("Alert configuration %s deleted.", alertConfigID)
	}

	return nil
}

func createAlertConfigs(workflowCtx *workflow.Context, groupID string, alertSpec []akov2.AlertConfiguration) []status.AlertConfiguration {
	logger := workflowCtx.Log
	var result []status.AlertConfiguration
	for _, alert := range alertSpec {
		atlasAlert, err := alert.ToAtlas()
		if err != nil {
			logger.Errorf("failed to convert spec to atlas alert configuration: %v", err)
			raw, errMarshal := json.Marshal(alert)
			if errMarshal != nil {
				logger.Errorf("failed to marshal alert configuration: %v", errMarshal)
				continue
			}
			result = append(result, status.NewFailedParseAlertConfigStatus(fmt.Sprintf("failed to parse atlas alert configuration: %v", err), string(raw)))
			continue
		}

		alertConfiguration, _, err := workflowCtx.SdkClientSet.SdkClient20250312009.AlertConfigurationsApi.
			CreateAlertConfig(workflowCtx.Context, groupID, atlasAlert).
			Execute()
		if err != nil {
			logger.Errorf("failed to create alert configuration: %v", err)
			result = append(result, status.NewIncorrectAlertConfigStatus(fmt.Sprintf("failed to create atlas alert configuration: %v", err), atlasAlert, workflowCtx.Log))
		} else {
			if alertConfiguration == nil {
				logger.Errorf("failed to create alert configuration: %v", err)
				result = append(result, status.NewIncorrectAlertConfigStatus(fmt.Sprintf("failed to create atlas alert configuration: %v", err), atlasAlert, workflowCtx.Log))
			} else {
				result = append(result, status.ParseAlertConfiguration(*alertConfiguration, workflowCtx.Log))
			}
		}
	}
	return result
}

func sortAlertConfigs(logger *zap.SugaredLogger, alertConfigSpecs []akov2.AlertConfiguration, atlasAlertConfigs []admin.GroupAlertsConfig) alertConfigurationDiff {
	var result alertConfigurationDiff
	for _, alertConfigSpec := range alertConfigSpecs {
		found := false
		for _, atlasAlertConfig := range atlasAlertConfigs {
			if isAlertConfigSpecEqualToAtlas(logger, alertConfigSpec, atlasAlertConfig) {
				found = true
				logger.Debugf("Alert configuration %s already exists.", atlasAlertConfig.GetId())
				result.CreateStatus = append(result.CreateStatus, atlasAlertConfig)
				break
			}
		}
		if !found {
			result.Create = append(result.Create, alertConfigSpec)
		}
	}

	for _, atlasAlertConfig := range atlasAlertConfigs {
		found := false
		for _, alertConfigSpec := range result.CreateStatus {
			if atlasAlertConfig.GetId() == alertConfigSpec.GetId() {
				found = true
			}
		}
		if !found {
			result.Delete = append(result.Delete, atlasAlertConfig.GetId())
		}
	}

	return result
}

type alertConfigurationDiff struct {
	Create       []akov2.AlertConfiguration
	Delete       []string
	CreateStatus []admin.GroupAlertsConfig
}

func isAlertConfigSpecEqualToAtlas(logger *zap.SugaredLogger, alertConfigSpec akov2.AlertConfiguration, atlasAlertConfig admin.GroupAlertsConfig) bool {
	if alertConfigSpec.EventTypeName != atlasAlertConfig.GetEventTypeName() {
		return false
	}
	if alertConfigSpec.SeverityOverride != atlasAlertConfig.GetSeverityOverride() {
		return false
	}
	if atlasAlertConfig.Enabled == nil {
		logger.Debugf("Alert configuration %s is not nil", atlasAlertConfig.GetId())
		return false
	}
	if alertConfigSpec.Enabled != atlasAlertConfig.GetEnabled() {
		logger.Debugf("alertConfigSpec.Enabled %v != *atlasAlertConfig.Enabled %v", alertConfigSpec.Enabled, *atlasAlertConfig.Enabled)
		return false
	}

	if !alertConfigSpec.Threshold.IsEqual(atlasAlertConfig.Threshold) {
		logger.Debugf("alertConfigSpec.Threshold %v != atlasAlertConfig.Threshold %v", alertConfigSpec.Threshold, atlasAlertConfig.Threshold)
		return false
	}

	if !alertConfigSpec.MetricThreshold.IsEqual(atlasAlertConfig.MetricThreshold) {
		logger.Debugf("alertConfigSpec.MetricThreshold %v != atlasAlertConfig.MetricThreshold %v", alertConfigSpec.MetricThreshold, atlasAlertConfig.MetricThreshold)
		return false
	}

	// Notifications
	if len(alertConfigSpec.Notifications) != len(atlasAlertConfig.GetNotifications()) {
		logger.Debugf("len(alertConfigSpec.NotificationTokenNames) %v != len(atlasAlertConfig.NotificationTokenNames) %v", len(alertConfigSpec.Notifications), len(atlasAlertConfig.GetNotifications()))
		return false
	}
	for _, notification := range alertConfigSpec.Notifications {
		found := false
		for _, atlasNotification := range atlasAlertConfig.GetNotifications() {
			if notification.IsEqual(atlasNotification) {
				found = true
			}
		}
		if !found {
			logger.Debugf("notification %v not found in atlasAlertConfig.Notifications %v", notification, atlasAlertConfig.Notifications)
			return false
		}
	}

	// Matchers
	if len(alertConfigSpec.Matchers) != len(atlasAlertConfig.GetMatchers()) {
		logger.Debugf("len(alertConfigSpec.Matchers) %v != len(atlasAlertConfig.Matchers) %v", len(alertConfigSpec.Matchers), len(atlasAlertConfig.GetMatchers()))
		return false
	}

	atlasMatchers := []akov2.Matcher{}
	err := compat.JSONCopy(atlasMatchers, atlasAlertConfig.GetMatchers())
	if err != nil {
		logger.Errorf("unable to convert matchers to structured type: %s", err)
		return false
	}
	for _, matcher := range alertConfigSpec.Matchers {
		found := false
		for _, atlasMatcher := range atlasMatchers {
			if matcher.IsEqual(atlasMatcher) {
				found = true
			}
		}
		if !found {
			logger.Debugf("matcher %v not found in atlasAlertConfig.Matchers %v", matcher, atlasAlertConfig.Matchers)
			return false
		}
	}

	return true
}
