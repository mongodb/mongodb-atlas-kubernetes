package status

import (
	"go.mongodb.org/atlas/mongodbatlas"
)

type CloudProviderAccessRole struct {
	AtlasAWSAccountArn         string         `json:"atlasAWSAccountArn,omitempty"`
	AtlasAssumedRoleExternalID string         `json:"atlasAssumedRoleExternalId"`
	AuthorizedDate             string         `json:"authorizedDate,omitempty"`
	CreatedDate                string         `json:"createdDate,omitempty"`
	FeatureUsages              []FeatureUsage `json:"featureUsages,omitempty"`
	IamAssumedRoleArn          string         `json:"iamAssumedRoleArn,omitempty"`
	ProviderName               string         `json:"providerName"`
	RoleID                     string         `json:"roleId,omitempty"`
	Status                     string         `json:"status,omitempty"`
	ErrorMessage               string         `json:"errorMessage,omitempty"`
}

type FeatureUsage struct {
	FeatureType string `json:"featureType,omitempty"`
	FeatureID   string `json:"featureId,omitempty"`
}

const (
	StatusFailed  = "FAILED"
	StatusCreated = "CREATED"
	StatusReady   = "READY"
)

func NewCloudProviderAccessRole(providerName, assumedRoleArn string) CloudProviderAccessRole {
	return CloudProviderAccessRole{
		ProviderName:      providerName,
		IamAssumedRoleArn: assumedRoleArn,
		Status:            StatusCreated,
	}
}

func (c *CloudProviderAccessRole) Failed(errorMessage string) {
	c.Status = StatusFailed
	c.ErrorMessage = errorMessage
}

func (c *CloudProviderAccessRole) FailedToAuthorise(errorMessage string) {
	c.ErrorMessage = errorMessage
}

func (c *CloudProviderAccessRole) Update(role mongodbatlas.AWSIAMRole) {
	c.RoleID = role.RoleID
	c.AtlasAssumedRoleExternalID = role.AtlasAssumedRoleExternalID
	c.AtlasAWSAccountArn = role.AtlasAWSAccountARN
	c.AuthorizedDate = role.AuthorizedDate
	c.CreatedDate = role.CreatedDate
	for _, featureUsage := range role.FeatureUsages {
		if featureUsage != nil {
			featureUsageID, ok := featureUsage.FeatureID.(string)
			if ok {
				c.FeatureUsages = append(c.FeatureUsages, FeatureUsage{
					FeatureType: featureUsage.FeatureType,
					FeatureID:   featureUsageID,
				})
			}
		}
	}

	switch role.IAMAssumedRoleARN {
	case "":
		c.Status = StatusCreated
	case c.IamAssumedRoleArn:
		c.Status = StatusReady
		c.ErrorMessage = ""
	default:
		c.Status = StatusFailed
		c.ErrorMessage = "IAMAssumedRoleARN is different from the previous one"
	}
}
