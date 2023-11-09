package atlas

import (
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"go.mongodb.org/atlas-sdk/v20231001002/admin"
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
	if ear.HasAwsKms() {
		config := ear.GetAwsKms()
		err := c.aws.DeleteKMS(config.GetCustomerMasterKeyID(), config.GetRegion())
		if err != nil {
			fmt.Println(text.FgRed.Sprintf("\tFailed to delete AWS KMS key %s: %s", config.GetCustomerMasterKeyID(), err))
		}
	}

	if ear.HasGoogleCloudKms() {
		config := ear.GetGoogleCloudKms()
		err := c.gcp.DeleteCryptoKey(ctx, config.GetKeyVersionResourceID())
		if err != nil {
			fmt.Println(text.FgRed.Sprintf("\tFailed to delete GCP Crypto key %s: %s", config.GetKeyVersionResourceID(), err))
		}
	}

	_, _, err := c.client.EncryptionAtRestUsingCustomerKeyManagementApi.
		UpdateEncryptionAtRest(ctx, projectID, &admin.EncryptionAtRest{}).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to delete encryption at rest for project %s: %s", projectID, err))
	}
}
