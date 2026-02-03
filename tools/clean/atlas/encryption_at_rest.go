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

package atlas

import (
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
)

func (c *Cleaner) getEncryptionAtRest(ctx context.Context, projectID string) *admin.EncryptionAtRest {
	ear, _, err := c.client.EncryptionAtRestUsingCustomerKeyManagementApi.
		GetEncryptionAtRest(ctx, projectID).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to get encryption at rest for project %s: %s", projectID, err))

		return nil
	}

	if !ear.HasAwsKms() && !ear.HasGoogleCloudKms() && !ear.HasAzureKeyVault() {
		return nil
	}

	return ear
}

func (c *Cleaner) deleteEncryptionAtRest(ctx context.Context, projectID string, ear *admin.EncryptionAtRest) {
	if config, ok := ear.GetAwsKmsOk(); ok && config.GetEnabled() {
		err := c.aws.DeleteKMS(ctx, config.GetCustomerMasterKeyID(), config.GetRegion())
		if err != nil {
			fmt.Println(text.FgRed.Sprintf("\tFailed to delete AWS KMS key %s: %s", config.GetCustomerMasterKeyID(), err))
		}
	}

	if config, ok := ear.GetGoogleCloudKmsOk(); ok && config.GetEnabled() {
		config := ear.GetGoogleCloudKms()
		err := c.gcp.DeleteCryptoKey(ctx, config.GetKeyVersionResourceID())
		if err != nil {
			fmt.Println(text.FgRed.Sprintf("\tFailed to delete GCP Crypto key %s: %s", config.GetKeyVersionResourceID(), err))
		}
	}

	disabled := false
	_, _, err := c.client.EncryptionAtRestUsingCustomerKeyManagementApi.
		UpdateEncryptionAtRest(
			ctx,
			projectID,
			&admin.EncryptionAtRest{
				AwsKms:         &admin.AWSKMSConfiguration{Enabled: &disabled},
				AzureKeyVault:  &admin.AzureKeyVault{Enabled: &disabled},
				GoogleCloudKms: &admin.GoogleCloudKMS{Enabled: &disabled},
			},
		).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to delete encryption at rest for project %s: %s", projectID, err))
	}
}
