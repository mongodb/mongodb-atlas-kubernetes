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

package validate

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
)

func Project(project *akov2.AtlasProject, isGov bool) error {
	if !isGov && project.Spec.RegionUsageRestrictions != "" && project.Spec.RegionUsageRestrictions != "NONE" {
		return errors.New("regionUsageRestriction can be used only with Atlas for government")
	}

	if isGov {
		if err := projectForGov(project); err != nil {
			return err
		}
	}

	if err := projectIPAccessList(project.Spec.ProjectIPAccessList); err != nil {
		return err
	}

	if err := projectCustomRoles(project.Spec.CustomRoles); err != nil {
		return err
	}

	if project.Spec.AlertConfigurationSyncEnabled {
		if err := alertConfigs(project.Spec.AlertConfigurations); err != nil {
			return err
		}
	}

	return nil
}

func projectForGov(project *akov2.AtlasProject) error {
	var err error

	if len(project.Spec.NetworkPeers) > 0 {
		for _, peer := range project.Spec.NetworkPeers {
			if peer.ProviderName != "AWS" {
				err = errors.Join(err, errors.New("atlas for government only supports AWS provider. one or more network peers are not set to AWS"))
			}

			regionErr := validCloudGovRegion(project.Spec.RegionUsageRestrictions, peer.AccepterRegionName)
			if regionErr != nil {
				err = errors.Join(err, fmt.Errorf("network peering in atlas for government support a restricted set of regions: %w", regionErr))
			}
		}
	}

	if project.Spec.EncryptionAtRest != nil {
		if project.Spec.EncryptionAtRest.AzureKeyVault.Enabled != nil && *project.Spec.EncryptionAtRest.AzureKeyVault.Enabled {
			err = errors.Join(err, errors.New("atlas for government only supports AWS provider. disable encryption at rest for Azure"))
		}

		if project.Spec.EncryptionAtRest.GoogleCloudKms.Enabled != nil && *project.Spec.EncryptionAtRest.GoogleCloudKms.Enabled {
			err = errors.Join(err, errors.New("atlas for government only supports AWS provider. disable encryption at rest for Google Cloud"))
		}

		if project.Spec.EncryptionAtRest.AwsKms.Enabled != nil && *project.Spec.EncryptionAtRest.AwsKms.Enabled {
			regionErr := validCloudGovRegion(project.Spec.RegionUsageRestrictions, project.Spec.EncryptionAtRest.AwsKms.Region)
			if regionErr != nil {
				err = errors.Join(err, fmt.Errorf("encryption at rest in atlas for government support a restricted set of regions: %w", regionErr))
			}
		}
	}

	if len(project.Spec.PrivateEndpoints) > 0 {
		for _, pe := range project.Spec.PrivateEndpoints {
			if pe.Provider != "AWS" {
				err = errors.Join(err, errors.New("atlas for government only supports AWS provider. one or more private endpoints are not set to AWS"))
			}

			regionErr := validCloudGovRegion(project.Spec.RegionUsageRestrictions, pe.Region)
			if regionErr != nil {
				err = errors.Join(err, fmt.Errorf("private endpoint in atlas for government support a restricted set of regions: %w", regionErr))
			}
		}
	}

	return err
}

func validCloudGovRegion(restriction, region string) error {
	fedRampRegions := map[string]struct{}{
		"US_EAST_1": {},
		"US_EAST_2": {},
		"US_WEST_1": {},
		"US_WEST_2": {},
		"us-east-1": {},
		"us-east-2": {},
		"us-west-1": {},
		"us-west-2": {},
	}
	govRegions := map[string]struct{}{
		"US_GOV_EAST_1": {},
		"US_GOV_WEST_1": {},
		"us-gov-east-1": {},
		"us-gov-west-1": {},
	}

	switch restriction {
	case "GOV_REGIONS_ONLY":
		if _, ok := govRegions[region]; !ok {
			return fmt.Errorf("%s is not part of AWS for government regions", region)
		}
	default:
		if _, ok := fedRampRegions[region]; !ok {
			return fmt.Errorf("%s is not part of AWS FedRAMP regions", region)
		}
	}

	return nil
}

func DatabaseUser(_ *akov2.AtlasDatabaseUser) error {
	return nil
}

func BackupSchedule(bSchedule *akov2.AtlasBackupSchedule, deployment *akov2.AtlasDeployment) error {
	var err error

	if bSchedule.Spec.Export == nil && bSchedule.Spec.AutoExportEnabled {
		err = errors.Join(err, errors.New("you must specify export policy when auto export is enabled"))
	}

	replicaSets := map[string]struct{}{}
	if deployment.Status.ReplicaSets != nil {
		for _, replicaSet := range deployment.Status.ReplicaSets {
			replicaSets[replicaSet.ID] = struct{}{}
		}
	}

	if len(bSchedule.Spec.CopySettings) > 0 && len(deployment.Status.ReplicaSets) == 0 {
		err = errors.Join(err, fmt.Errorf("deployment %s doesn't have replication status available", deployment.GetDeploymentName()))
	}

	for position, copySetting := range bSchedule.Spec.CopySettings {
		if copySetting.RegionName == nil {
			err = errors.Join(err, fmt.Errorf("copy setting at position %d: you must set a region name", position))
		}

		if copySetting.ShouldCopyOplogs != nil && *copySetting.ShouldCopyOplogs {
			if deployment.Spec.DeploymentSpec != nil &&
				(deployment.Spec.DeploymentSpec.PitEnabled == nil ||
					!*deployment.Spec.DeploymentSpec.PitEnabled) {
				err = errors.Join(err, fmt.Errorf("copy setting at position %d: you must enable pit before enable copyOplogs", position))
			}
		}
	}

	return err
}

func projectIPAccessList(ipAccessList []project.IPAccessList) error {
	if len(ipAccessList) == 0 {
		return nil
	}

	var err error
	for _, item := range ipAccessList {
		if item.IPAddress == "" && item.CIDRBlock == "" && item.AwsSecurityGroup == "" {
			err = errors.Join(err, errors.New("invalid config! one of option must be configured"))
		}

		if item.CIDRBlock != "" {
			if item.AwsSecurityGroup != "" || item.IPAddress != "" {
				err = errors.Join(err, errors.New("don't set ipAddress or awsSecurityGroup when configuring cidrBlock"))
			}

			_, _, cidrErr := net.ParseCIDR(item.CIDRBlock)
			if cidrErr != nil {
				err = errors.Join(err, fmt.Errorf("invalid cidrBlock: %s", item.CIDRBlock))
			}
		}

		if item.IPAddress != "" {
			if item.AwsSecurityGroup != "" || item.CIDRBlock != "" {
				err = errors.Join(err, errors.New("don't set cidrBlock or awsSecurityGroup when configuring ipAddress"))
			}

			ip := net.ParseIP(item.IPAddress)
			if ip == nil {
				err = errors.Join(err, fmt.Errorf("invalid ipAddress: %s", item.IPAddress))
			}
		}

		if item.AwsSecurityGroup != "" {
			if item.IPAddress != "" || item.CIDRBlock != "" {
				err = errors.Join(err, errors.New("don't set cidrBlock or ipAddress when configuring awsSecurityGroup"))
			}

			reg := regexp.MustCompile("^([0-9]*/)?sg-([0-9]*)")
			if !reg.MatchString(item.AwsSecurityGroup) {
				err = errors.Join(err, fmt.Errorf("invalid awsSecurityGroup: %s", item.AwsSecurityGroup))
			}
		}

		if item.DeleteAfterDate != "" {
			_, delErr := timeutil.ParseISO8601(item.DeleteAfterDate)
			if delErr != nil {
				err = errors.Join(err, fmt.Errorf("invalid deleteAfterDate: %s. value should follow ISO8601 format", item.DeleteAfterDate))
			}
		}
	}

	return err
}

func projectCustomRoles(customRoles []akov2.CustomRole) error {
	if len(customRoles) == 0 {
		return nil
	}

	var err error
	customRolesMap := map[string]struct{}{}

	for _, customRole := range customRoles {
		if _, ok := customRolesMap[customRole.Name]; ok {
			err = errors.Join(err, fmt.Errorf("the custom role \"%s\" is duplicate. custom role name must be unique", customRole.Name))
		}

		customRolesMap[customRole.Name] = struct{}{}
	}

	return err
}

func alertConfigs(alertConfigs []akov2.AlertConfiguration) error {
	seenConfigs := []akov2.AlertConfiguration{}
	for j, cfg := range alertConfigs {
		for i, seenCfg := range seenConfigs {
			if reflect.DeepEqual(seenCfg, cfg) {
				return fmt.Errorf("alert config at position %d is a duplicate of "+
					"alert config at position %d: %v", j, i, cfg)
			}
		}
		seenConfigs = append(seenConfigs, cfg)
	}
	return nil
}
