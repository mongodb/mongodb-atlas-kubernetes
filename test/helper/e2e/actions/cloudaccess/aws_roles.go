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
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
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
	IAMClient := iam.New(session.Must(session.NewSession()))
	roleInput := iam.CreateRoleInput{}
	roleInput.SetRoleName(roleName)
	roleInput.SetAssumeRolePolicyDocument(policy)
	//roleInput.SetTags([]*iam.Tag{
	//	{
	//		Key:   aws.String(config.TagForTestKey),
	//		Value: aws.String(config.TagForTestValue),
	//	},
	//})
	role, err := IAMClient.CreateRole(&roleInput)
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
	IAMClient := iam.New(session.Must(session.NewSession()))
	roleUpdate := iam.UpdateAssumeRolePolicyInput{}
	roleUpdate.SetPolicyDocument(updatedPolicy)
	roleUpdate.SetRoleName(roleName)
	req, _ := IAMClient.UpdateAssumeRolePolicyRequest(&roleUpdate)
	err = req.Send()
	if err != nil {
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
	IAMClient := iam.New(session.Must(session.NewSession()))
	_, err := IAMClient.DeleteRole(&iam.DeleteRoleInput{
		RoleName: aws.String(NameFromArn(arn)),
	})
	if err != nil {
		return err
	}
	return nil
}
