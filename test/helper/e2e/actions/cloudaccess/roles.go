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

package cloudaccess

import (
	"context"
	"fmt"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
)

type Role struct {
	Name       string
	AccessRole akov2.CloudProviderIntegration
}

func CreateRoles(ctx context.Context, roles []Role) error {
	for i, role := range roles {
		switch role.AccessRole.ProviderName {
		case string(provider.ProviderAWS):
			arn, err := CreateAWSIAMRole(ctx, role.Name)
			if err != nil {
				return err
			}
			roles[i].AccessRole.IamAssumedRoleArn = arn
		default:
			return fmt.Errorf("unsupported provider %s", role.AccessRole.ProviderName)
		}
	}
	return nil
}

func AddAtlasStatementToRole(ctx context.Context, roles []Role, roleStatuses []status.CloudProviderIntegration) error {
	if len(roles) != len(roleStatuses) {
		return fmt.Errorf("number of roles %d does not match number of statuses %d", len(roles), len(roleStatuses))
	}
	for _, role := range roles {
		for _, roleStatus := range roleStatuses {
			if role.AccessRole.ProviderName == roleStatus.ProviderName && role.AccessRole.IamAssumedRoleArn == roleStatus.IamAssumedRoleArn {
				if err := AddAtlasStatementToAWSIAMRole(ctx, roleStatus.AtlasAWSAccountArn, roleStatus.AtlasAssumedRoleExternalID, role.Name); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func DeleteCloudProviderIntegrations(ctx context.Context, roles []akov2.CloudProviderIntegration) []error {
	var errorList []error
	for _, role := range roles {
		switch role.ProviderName {
		case string(provider.ProviderAWS):
			if err := DeleteAWSIAMRoleByArn(ctx, role.IamAssumedRoleArn); err != nil {
				errorList = append(errorList, err)
			}
		default:
			errorList = append(errorList, fmt.Errorf("unsupported provider %s", role.ProviderName))
		}
	}
	return errorList
}
