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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/onsi/ginkgo/v2"
)

type AwsResourcesGenerator struct {
	t ginkgo.GinkgoTInterface

	iamClient *iam.IAM
	s3Client  *s3.S3
}

const defaultRegion = "us-east-1"

type IAMPolicy *string

func NewAwsResourcesGenerator(t ginkgo.GinkgoTInterface, region *string) *AwsResourcesGenerator {
	t.Helper()

	if region == nil {
		region = aws.String(defaultRegion)
	}

	awsSession, err := session.NewSession(
		&aws.Config{
			Region: region,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return &AwsResourcesGenerator{
		t: t,

		iamClient: iam.New(awsSession),
		s3Client:  s3.New(awsSession),
	}
}

func (g *AwsResourcesGenerator) GetIAMRole(name string) (string, error) {
	input := &iam.GetRoleInput{
		RoleName: aws.String(name),
	}

	role, err := g.iamClient.GetRole(input)
	if err != nil {
		return "", fmt.Errorf("failed to get iam role: %w", err)
	}

	return role.GoString(), nil
}

func (g *AwsResourcesGenerator) CreateIAMRole(name string, policy func() IAMPolicy) (string, error) {
	input := &iam.CreateRoleInput{
		RoleName:                 aws.String(name),
		AssumeRolePolicyDocument: policy(),
	}

	role, err := g.iamClient.CreateRole(input)
	if err != nil {
		return "", fmt.Errorf("failed to create iam role: %w", err)
	}

	return *role.Role.Arn, nil
}

func (g *AwsResourcesGenerator) DeleteIAMRole(name string) error {
	input := &iam.DeleteRoleInput{
		RoleName: aws.String(name),
	}

	_, err := g.iamClient.DeleteRole(input)
	if err != nil {
		return fmt.Errorf("failed to create iam role: %w", err)
	}

	return nil
}

func (g *AwsResourcesGenerator) CreatePolicy(name string, policy func() IAMPolicy) (string, error) {
	input := &iam.CreatePolicyInput{
		PolicyDocument: policy(),
		PolicyName:     aws.String(name),
	}

	r, err := g.iamClient.CreatePolicy(input)
	if err != nil {
		return "", fmt.Errorf("failed to create iam policy: %w", err)
	}

	return *r.Policy.Arn, nil
}

func (g *AwsResourcesGenerator) DeletePolicy(arn string) error {
	input := &iam.DeletePolicyInput{
		PolicyArn: aws.String(arn),
	}

	_, err := g.iamClient.DeletePolicy(input)
	if err != nil {
		return fmt.Errorf("failed to delete iam policy: %w", err)
	}

	return nil
}

func (g *AwsResourcesGenerator) AttachRolePolicy(roleName, policyArn string) error {
	input := &iam.AttachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
		RoleName:  aws.String(roleName),
	}

	_, err := g.iamClient.AttachRolePolicy(input)
	if err != nil {
		return fmt.Errorf("failed to attach iam policy to role: %w", err)
	}

	return nil
}

func (g *AwsResourcesGenerator) DetachRolePolicy(roleName, policyArn string) error {
	input := &iam.DetachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
		RoleName:  aws.String(roleName),
	}

	_, err := g.iamClient.DetachRolePolicy(input)
	if err != nil {
		return fmt.Errorf("failed to detach iam policy from role: %w", err)
	}

	return nil
}

func (g *AwsResourcesGenerator) ListAttachedRolePolicy(roleName string) (string, error) {
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	}

	r, err := g.iamClient.ListAttachedRolePolicies(input)
	if err != nil {
		return "", fmt.Errorf("failed to list iam policies from role: %w", err)
	}

	return r.GoString(), nil
}

func (g *AwsResourcesGenerator) CreateBucket(name string) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}

	_, err := g.s3Client.CreateBucket(input)
	if err != nil {
		return fmt.Errorf("failed to create aws bucket: %w", err)
	}

	return nil
}

func (g *AwsResourcesGenerator) DeleteBucket(name string) error {
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}

	_, err := g.s3Client.DeleteBucket(input)
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

func (g *AwsResourcesGenerator) ListObjects(name string) ([]*s3.Object, error) {
	input := &s3.ListObjectsInput{
		Bucket: aws.String(name),
	}

	objs, err := g.s3Client.ListObjects(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects of bucket %s: %w", name, err)
	}

	return objs.Contents, nil
}

func (g *AwsResourcesGenerator) DeleteObject(bucketName, objectKey string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	_, err := g.s3Client.DeleteObject(input)
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
