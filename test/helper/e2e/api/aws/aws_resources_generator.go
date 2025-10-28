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

package helper

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/onsi/ginkgo/v2"
)

type AwsResourcesGenerator struct {
	t ginkgo.GinkgoTInterface

	iamClient *iam.Client
	s3Client  *s3.Client
}

const defaultRegion = "us-east-1"

type IAMPolicy *string

func NewAwsResourcesGenerator(t ginkgo.GinkgoTInterface, region *string) *AwsResourcesGenerator {
	t.Helper()

	if region == nil {
		region = aws.String(defaultRegion)
	}

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, func(lo *config.LoadOptions) error {
		lo.Region = *region
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	return &AwsResourcesGenerator{
		t: t,

		iamClient: iam.NewFromConfig(cfg),
		s3Client:  s3.NewFromConfig(cfg),
	}
}

func (g *AwsResourcesGenerator) GetIAMRole(name string) (string, error) {
	ctx := context.TODO()
	input := &iam.GetRoleInput{
		RoleName: aws.String(name),
	}

	role, err := g.iamClient.GetRole(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to get iam role: %w", err)
	}

	b, err := json.MarshalIndent(role, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal role output: %w", err)
	}

	return string(b), nil
}

func (g *AwsResourcesGenerator) CreateIAMRole(name string, policy func() IAMPolicy) (string, error) {
	ctx := context.TODO()
	input := &iam.CreateRoleInput{
		RoleName:                 aws.String(name),
		AssumeRolePolicyDocument: policy(),
	}

	role, err := g.iamClient.CreateRole(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to create iam role: %w", err)
	}

	return *role.Role.Arn, nil
}

func (g *AwsResourcesGenerator) DeleteIAMRole(name string) error {
	ctx := context.TODO()
	input := &iam.DeleteRoleInput{
		RoleName: aws.String(name),
	}

	_, err := g.iamClient.DeleteRole(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create iam role: %w", err)
	}

	return nil
}

func (g *AwsResourcesGenerator) CreatePolicy(name string, policy func() IAMPolicy) (string, error) {
	ctx := context.TODO()
	input := &iam.CreatePolicyInput{
		PolicyDocument: policy(),
		PolicyName:     aws.String(name),
		Tags: []iamtypes.Tag{
			{Key: aws.String(OwnerTag), Value: aws.String(AKOTeam)},
			{Key: aws.String(OwnerEmailTag), Value: aws.String(AKOEmail)},
			{Key: aws.String(CostCenterTag), Value: aws.String(AKOCostCenter)},
			{Key: aws.String(EnvironmentTag), Value: aws.String(AKOEnvTest)},
		},
	}

	r, err := g.iamClient.CreatePolicy(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to create iam policy: %w", err)
	}

	return *r.Policy.Arn, nil
}

func (g *AwsResourcesGenerator) DeletePolicy(arn string) error {
	ctx := context.TODO()
	input := &iam.DeletePolicyInput{
		PolicyArn: aws.String(arn),
	}

	_, err := g.iamClient.DeletePolicy(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete iam policy: %w", err)
	}

	return nil
}

func (g *AwsResourcesGenerator) AttachRolePolicy(roleName, policyArn string) error {
	ctx := context.TODO()
	input := &iam.AttachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
		RoleName:  aws.String(roleName),
	}

	_, err := g.iamClient.AttachRolePolicy(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to attach iam policy to role: %w", err)
	}

	return nil
}

func (g *AwsResourcesGenerator) DetachRolePolicy(roleName, policyArn string) error {
	ctx := context.TODO()
	input := &iam.DetachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
		RoleName:  aws.String(roleName),
	}

	_, err := g.iamClient.DetachRolePolicy(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to detach iam policy from role: %w", err)
	}

	return nil
}

func (g *AwsResourcesGenerator) ListAttachedRolePolicy(roleName string) (string, error) {
	ctx := context.TODO()
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	}

	r, err := g.iamClient.ListAttachedRolePolicies(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to list iam policies from role: %w", err)
	}

	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal role output: %w", err)
	}

	return string(b), nil
}

func (g *AwsResourcesGenerator) CreateBucket(name string) error {
	ctx := context.TODO()
	input := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}

	_, err := g.s3Client.CreateBucket(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create aws bucket: %w", err)
	}

	tagSet := &s3types.Tagging{
		TagSet: []s3types.Tag{
			{Key: aws.String(OwnerTag), Value: aws.String(AKOTeam)},
			{Key: aws.String(OwnerEmailTag), Value: aws.String(AKOEmail)},
			{Key: aws.String(CostCenterTag), Value: aws.String(AKOCostCenter)},
			{Key: aws.String(EnvironmentTag), Value: aws.String(AKOEnvTest)},
		},
	}

	taggingInput := &s3.PutBucketTaggingInput{
		Bucket:  aws.String(name),
		Tagging: tagSet,
	}

	if _, err := g.s3Client.PutBucketTagging(ctx, taggingInput); err != nil {
		return fmt.Errorf("failed to tag bucket %s: %w", name, err)
	}

	return nil
}

func (g *AwsResourcesGenerator) DeleteBucket(name string) error {
	ctx := context.TODO()
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}

	_, err := g.s3Client.DeleteBucket(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete aws bucket: %w", err)
	}

	return nil
}

func (g *AwsResourcesGenerator) EmptyBucket(name string) error {
	objs, err := g.ListObjects(name)
	if err != nil {
		return err
	}

	for _, obj := range objs {
		err = g.DeleteObject(name, *obj.Key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *AwsResourcesGenerator) ListObjects(name string) ([]s3types.Object, error) {
	ctx := context.TODO()
	input := &s3.ListObjectsInput{
		Bucket: aws.String(name),
	}

	objs, err := g.s3Client.ListObjects(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects of bucket %s: %w", name, err)
	}

	return objs.Contents, nil
}

func (g *AwsResourcesGenerator) DeleteObject(bucketName, objectKey string) error {
	ctx := context.TODO()
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	_, err := g.s3Client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete object %s in bucket %s: %w", bucketName, objectKey, err)
	}

	return nil
}

func (g *AwsResourcesGenerator) Cleanup(task func()) {
	g.t.Cleanup(task)
}

func CloudProviderAccessPolicy(atlasAWSAccountArn, atlasAssumedRoleExternalID string) IAMPolicy {
	policy := `{
   "Version":"2012-10-17",
   "Statement":[
      {
         "Effect":"Allow",
         "Principal":{
            "AWS":"` + atlasAWSAccountArn + `"
         },
         "Action":"sts:AssumeRole",
         "Condition":{
            "StringEquals":{
               "sts:ExternalId":"` + atlasAssumedRoleExternalID + `"
            }
         }
      }
   ]
}`

	return &policy
}

func BucketExportPolicy(bucketName string) IAMPolicy {
	policy := `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "s3:GetBucketLocation",
            "Resource": "arn:aws:s3:::` + bucketName + `"
        },
        {
            "Effect": "Allow",
            "Action": "s3:PutObject",
            "Resource": "arn:aws:s3:::` + bucketName + `/*"
        }
    ]
}`

	return &policy
}
