package cloudaccess

import (
	"fmt"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
)

type Role struct {
	Name       string
	AccessRole v1.CloudProviderAccessRole
}

func CreateRoles(roles []Role) error {
	for i, role := range roles {
		switch role.AccessRole.ProviderName {
		case string(provider.ProviderAWS):
			arn, err := CreateAWSIAMRole(role.Name)
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

func AddAtlasStatementToRole(roles []Role, roleStatuses []status.CloudProviderAccessRole) error {
	if len(roles) != len(roleStatuses) {
		return fmt.Errorf("number of roles %d does not match number of statuses %d", len(roles), len(roleStatuses))
	}
	for _, role := range roles {
		for _, roleStatus := range roleStatuses {
			if role.AccessRole.ProviderName == roleStatus.ProviderName && role.AccessRole.IamAssumedRoleArn == roleStatus.IamAssumedRoleArn {
				if err := AddAtlasStatementToAWSIAMRole(roleStatus.AtlasAWSAccountArn, roleStatus.AtlasAssumedRoleExternalID, role.Name); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func DeleteRoles(roles []v1.CloudProviderAccessRole) []error {
	var errorList []error
	for _, role := range roles {
		switch role.ProviderName {
		case string(provider.ProviderAWS):
			if err := DeleteAWSIAMRoleByArn(role.IamAssumedRoleArn); err != nil {
				errorList = append(errorList, err)
			}
		default:
			errorList = append(errorList, fmt.Errorf("unsupported provider %s", role.ProviderName))
		}
	}
	return errorList
}
