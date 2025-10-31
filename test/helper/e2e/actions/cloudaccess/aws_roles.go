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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	awshelper "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/aws"
)

type AssumeRolePolicyDocument struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

type Statement struct {
	Effect    string    `json:"Effect"`
	Principal Principal `json:"Principal"`
	Action    string    `json:"Action"`
	Condition Condition `json:"Condition,omitempty"`
}

type Principal struct {
	AWS     string `json:"AWS,omitempty"`
	Service string `json:"Service,omitempty"`
}

type Condition struct {
	StringEquals StringEquals `json:"StringEquals,omitempty"`
}

type StringEquals struct {
	StsExternalId string `json:"sts:ExternalId,omitempty"`
}

func defaultPolicy() AssumeRolePolicyDocument {
	return AssumeRolePolicyDocument{
		Version: "2012-10-17",
		Statement: []Statement{
			{
				Effect: "Allow",
				Principal: Principal{
					Service: "ec2.amazonaws.com",
				},
				Action: "sts:AssumeRole",
			},
		},
	}
}

func EC2RolePolicyString() (string, error) {
	policy := defaultPolicy()
	byteStr, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}
	return string(byteStr), nil
}

func PolicyWithAtlasArn(atlasAWSAccountArn, atlasAssumedRoleExternalId string) (string, error) {
	policy := defaultPolicy()
	policy.Statement = append(policy.Statement, Statement{
		Effect: "Allow",
		Principal: Principal{
			AWS: atlasAWSAccountArn,
		},
		Action: "sts:AssumeRole",
		Condition: Condition{
			StringEquals: StringEquals{
				StsExternalId: atlasAssumedRoleExternalId,
			},
		},
	})
	byteStr, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}
	return string(byteStr), nil
}

func CreateAWSIAMRole(roleName string) (string, error) {
	policy, err := EC2RolePolicyString()
	if err != nil {
		return "", err
	}
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create an AWS config: %w", err)
	}
	IAMClient := iam.NewFromConfig(cfg)
	roleInput := iam.CreateRoleInput{
		RoleName:                 &roleName,
		AssumeRolePolicyDocument: &policy,
	}
	roleInput.Tags = []types.Tag{
		{Key: aws.String(awshelper.OwnerTag), Value: aws.String(awshelper.AKOTeam)},
		{Key: aws.String(awshelper.OwnerEmailTag), Value: aws.String(awshelper.AKOEmail)},
		{Key: aws.String(awshelper.CostCenterTag), Value: aws.String(awshelper.AKOCostCenter)},
		{Key: aws.String(awshelper.EnvironmentTag), Value: aws.String(awshelper.AKOEnvTest)},
	}
	role, err := IAMClient.CreateRole(ctx, &roleInput)
	if err != nil {
		return "", err
	}
	return *role.Role.Arn, nil
}

func AddAtlasStatementToAWSIAMRole(atlasAWSAccountArn, atlasAssumedRoleExternalId, roleName string) error {
	updatedPolicy, err := PolicyWithAtlasArn(atlasAWSAccountArn, atlasAssumedRoleExternalId)
	if err != nil {
		return err
	}
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to create an AWS config: %w", err)
	}
	IAMClient := iam.NewFromConfig(cfg)
	roleUpdate := iam.UpdateAssumeRolePolicyInput{
		RoleName:       &roleName,
		PolicyDocument: &updatedPolicy,
	}
	if _, err := IAMClient.UpdateAssumeRolePolicy(ctx, &roleUpdate); err != nil {
		return err
	}
	return nil
}

func NameFromArn(arn string) string {
	// It's a little hacky, but it works. AWS doesn't have an API for finding role by arn
	// arn format is arn:aws:iam::<account_id>:role/<role_name>
	return arn[strings.LastIndex(arn, "/")+1:]
}

func DeleteAWSIAMRoleByArn(arn string) error {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to create an AWS config: %w", err)
	}
	IAMClient := iam.NewFromConfig(cfg)
	_, err = IAMClient.DeleteRole(ctx, &iam.DeleteRoleInput{
		RoleName: aws.String(NameFromArn(arn)),
	})
	if err != nil {
		return err
	}
	return nil
}
