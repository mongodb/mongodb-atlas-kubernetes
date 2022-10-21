package atlasproject

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func ensureAlertConfigurations(service *workflow.Context, project *mdbv1.AtlasProject, groupID string) workflow.Result {
	if project.Spec.AlertConfigurationSyncEnabled {
		specToSync := project.Spec.DeepCopy().AlertConfigurations

		alertConfigurationCondition := status.AlertConfigurationReadyType
		ctx := context.Background()
		if len(specToSync) == 0 {
			service.UnsetCondition(alertConfigurationCondition)
			return workflow.OK()
		}
		result := syncAlertConfigurations(ctx, service, groupID, specToSync)
		if !result.IsOk() {
			service.SetConditionFromResult(alertConfigurationCondition, result)
			return result
		}
		service.SetConditionTrue(alertConfigurationCondition)
		return result
	}
	service.UnsetCondition(status.AlertConfigurationReadyType)
	service.Log.Debugf("Alert configuration sync is disabled for project %s", project.Name)
	return workflow.OK()
}

func syncAlertConfigurations(context context.Context, service *workflow.Context, groupID string, alertSpec []mdbv1.AlertConfiguration) workflow.Result {
	logger := service.Log
	existedAlertConfigs, _, err := service.Client.AlertConfigurations.List(context, groupID, nil)
	if err != nil {
		logger.Errorf("failed to list alert configurations: %v", err)
		return workflow.Terminate(workflow.ProjectAlertConfigurationIsNotReadyInAtlas, fmt.Sprintf("failed to list alert configurations: %v", err))
	}

	diff := sortAlertConfigs(logger, alertSpec, existedAlertConfigs)
	logger.Debugf("to create %v, to create statuses %v, to delete %v", len(diff.Create), len(diff.CreateStatus), len(diff.Delete))

	newStatuses := createAlertConfigs(context, service, groupID, diff.Create)

	for _, existedAlertConfig := range diff.CreateStatus {
		newStatuses = append(newStatuses, status.ParseAlertConfiguration(existedAlertConfig))
	}

	service.EnsureStatusOption(status.AtlasProjectSetAlertConfigOption(&newStatuses))

	err = deleteAlertConfigs(context, service, groupID, diff.Delete)
	if err != nil {
		return workflow.Terminate(workflow.ProjectAlertConfigurationIsNotReadyInAtlas, fmt.Sprintf("failed to delete alert configurations: %v", err))
	}

	return checkAlertConfigurationStatuses(newStatuses)
}

func checkAlertConfigurationStatuses(statuses []status.AlertConfiguration) workflow.Result {
	for _, alertConfigurationStatus := range statuses {
		if alertConfigurationStatus.ErrorMessage != "" {
			return workflow.Terminate(workflow.ProjectAlertConfigurationIsNotReadyInAtlas,
				fmt.Sprintf("failed to create alert configuration: %s", alertConfigurationStatus.ErrorMessage))
		}
	}
	return workflow.OK()
}

func deleteAlertConfigs(context context.Context, ctx *workflow.Context, groupID string, alertConfigIDs []string) error {
	logger := ctx.Log
	for _, alertConfigID := range alertConfigIDs {
		_, err := ctx.Client.AlertConfigurations.Delete(context, groupID, alertConfigID)
		if err != nil {
			logger.Errorf("failed to delete alert configuration: %v", err)
			return err
		}
		logger.Infof("Alert configuration %s deleted.", alertConfigID)
	}
	return nil
}

func createAlertConfigs(context context.Context, ctx *workflow.Context, groupID string, alertSpec []mdbv1.AlertConfiguration) []status.AlertConfiguration {
	logger := ctx.Log
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

		alertConfiguration, _, err := ctx.Client.AlertConfigurations.Create(context, groupID, atlasAlert)
		if err != nil {
			logger.Errorf("failed to create alert configuration: %v", err)
			result = append(result, status.NewIncorrectAlertConfigStatus(fmt.Sprintf("failed to create atlas alert configuration: %v", err), atlasAlert))
		} else {
			if alertConfiguration == nil {
				logger.Errorf("failed to create alert configuration: %v", err)
				result = append(result, status.NewIncorrectAlertConfigStatus(fmt.Sprintf("failed to create atlas alert configuration: %v", err), atlasAlert))
			} else {
				result = append(result, status.ParseAlertConfiguration(*alertConfiguration))
			}
		}
	}
	return result
}

func sortAlertConfigs(logger *zap.SugaredLogger, alertConfigSpecs []mdbv1.AlertConfiguration, atlasAlertConfigs []mongodbatlas.AlertConfiguration) alertConfigurationDiff {
	var result alertConfigurationDiff
	for _, alertConfigSpec := range alertConfigSpecs {
		found := false
		for _, atlasAlertConfig := range atlasAlertConfigs {
			if isAlertConfigSpecEqualToAtlas(logger, alertConfigSpec, atlasAlertConfig) {
				found = true
				logger.Debugf("Alert configuration %s already exists.", atlasAlertConfig.ID)
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
			if atlasAlertConfig.ID == alertConfigSpec.ID {
				found = true
			}
		}
		if !found {
			result.Delete = append(result.Delete, atlasAlertConfig.ID)
		}
	}

	return result
}

type alertConfigurationDiff struct {
	Create       []mdbv1.AlertConfiguration
	Delete       []string
	CreateStatus []mongodbatlas.AlertConfiguration
}

func isAlertConfigSpecEqualToAtlas(logger *zap.SugaredLogger, alertConfigSpec mdbv1.AlertConfiguration, atlasAlertConfig mongodbatlas.AlertConfiguration) bool {
	if alertConfigSpec.EventTypeName != atlasAlertConfig.EventTypeName {
		return false
	}
	if atlasAlertConfig.Enabled == nil {
		logger.Debugf("Alert configuration %s is not nil", atlasAlertConfig.ID)
		return false
	}
	if alertConfigSpec.Enabled != *atlasAlertConfig.Enabled {
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
	if len(alertConfigSpec.Notifications) != len(atlasAlertConfig.Notifications) {
		logger.Debugf("len(alertConfigSpec.NotificationTokenNames) %v != len(atlasAlertConfig.NotificationTokenNames) %v", len(alertConfigSpec.Notifications), len(atlasAlertConfig.Notifications))
		return false
	}
	for _, notification := range alertConfigSpec.Notifications {
		found := false
		for _, atlasNotification := range atlasAlertConfig.Notifications {
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
	if len(alertConfigSpec.Matchers) != len(atlasAlertConfig.Matchers) {
		logger.Debugf("len(alertConfigSpec.Matchers) %v != len(atlasAlertConfig.Matchers) %v", len(alertConfigSpec.Matchers), len(atlasAlertConfig.Matchers))
		return false
	}
	for _, matcher := range alertConfigSpec.Matchers {
		found := false
		for _, atlasMatcher := range atlasAlertConfig.Matchers {
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
