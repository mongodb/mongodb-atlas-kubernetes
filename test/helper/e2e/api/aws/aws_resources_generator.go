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

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/tags"
)

type AwsResourcesGenerator struct {
	iamClient *iam.Client
	s3Client  *s3.Client
}

const defaultRegion = "us-east-1"

type IAMPolicy *string

func NewAwsResourcesGenerator(ctx context.Context, t ginkgo.GinkgoTInterface, region *string) *AwsResourcesGenerator {
	t.Helper()

	if region == nil {
		region = aws.String(defaultRegion)
	}

	cfg, err := config.LoadDefaultConfig(ctx, func(lo *config.LoadOptions) error {
		lo.Region = *region
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	return &AwsResourcesGenerator{
		iamClient: iam.NewFromConfig(cfg),
		s3Client:  s3.NewFromConfig(cfg),
	}
}

func (g *AwsResourcesGenerator) GetIAMRole(ctx context.Context, name string) (string, error) {
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

func (g *AwsResourcesGenerator) CreateIAMRole(ctx context.Context, name string, policy func() IAMPolicy) (string, error) {
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

func (g *AwsResourcesGenerator) DeleteIAMRole(ctx context.Context, name string) error {
	input := &iam.DeleteRoleInput{
		RoleName: aws.String(name),
	}

	_, err := g.iamClient.DeleteRole(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create iam role: %w", err)
	}

	return nil
}

func (g *AwsResourcesGenerator) CreatePolicy(ctx context.Context, name string, policy func() IAMPolicy) (string, error) {
	input := &iam.CreatePolicyInput{
		PolicyDocument: policy(),
		PolicyName:     aws.String(name),
		Tags: []iamtypes.Tag{
			{Key: aws.String(tags.OwnerTag), Value: aws.String(tags.AKOTeam)},
			{Key: aws.String(tags.OwnerEmailTag), Value: aws.String(tags.AKOEmail)},
			{Key: aws.String(tags.CostCenterTag), Value: aws.String(tags.AKOCostCenter)},
			{Key: aws.String(tags.EnvironmentTag), Value: aws.String(tags.AKOEnvTest)},
		},
	}

	r, err := g.iamClient.CreatePolicy(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to create iam policy: %w", err)
	}

	return *r.Policy.Arn, nil
}

func (g *AwsResourcesGenerator) DeletePolicy(ctx context.Context, arn string) error {
	input := &iam.DeletePolicyInput{
		PolicyArn: aws.String(arn),
	}

	_, err := g.iamClient.DeletePolicy(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete iam policy: %w", err)
	}

	return nil
}

func (g *AwsResourcesGenerator) AttachRolePolicy(ctx context.Context, roleName, policyArn string) error {
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

func (g *AwsResourcesGenerator) DetachRolePolicy(ctx context.Context, roleName, policyArn string) error {
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

func (g *AwsResourcesGenerator) ListAttachedRolePolicy(ctx context.Context, roleName string) (string, error) {
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

func (g *AwsResourcesGenerator) CreateBucket(ctx context.Context, name string) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}

	_, err := g.s3Client.CreateBucket(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create aws bucket: %w", err)
	}

	tagSet := &s3types.Tagging{
		TagSet: []s3types.Tag{
			{Key: aws.String(tags.OwnerTag), Value: aws.String(tags.AKOTeam)},
			{Key: aws.String(tags.OwnerEmailTag), Value: aws.String(tags.AKOEmail)},
			{Key: aws.String(tags.CostCenterTag), Value: aws.String(tags.AKOCostCenter)},
			{Key: aws.String(tags.EnvironmentTag), Value: aws.String(tags.AKOEnvTest)},
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

func (g *AwsResourcesGenerator) DeleteBucket(ctx context.Context, name string) error {
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}

	_, err := g.s3Client.DeleteBucket(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete aws bucket: %w", err)
	}

	return nil
}

func (g *AwsResourcesGenerator) EmptyBucket(ctx context.Context, name string) error {
	objs, err := g.ListObjects(ctx, name)
	if err != nil {
		return err
	}

	for _, obj := range objs {
		err = g.DeleteObject(ctx, name, *obj.Key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *AwsResourcesGenerator) ListObjects(ctx context.Context, name string) ([]s3types.Object, error) {
	input := &s3.ListObjectsInput{
		Bucket: aws.String(name),
	}

	objs, err := g.s3Client.ListObjects(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects of bucket %s: %w", name, err)
	}

	return objs.Contents, nil
}

func (g *AwsResourcesGenerator) DeleteObject(ctx context.Context, bucketName, objectKey string) error {
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
