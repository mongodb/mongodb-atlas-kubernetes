package atlasproject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

func TestFillStructFields(t *testing.T) {
	myStruct := struct {
		StrPublic      string
		strPrivate     string
		WrongFieldType int32
	}{
		StrPublic:      "",
		strPrivate:     "",
		WrongFieldType: 0,
	}

	data := map[string]string{
		"StrPublic":  "test-str",
		"strPrivate": "test",
	}

	fillStructFields(data, &myStruct)

	assert.Equal(t, data["StrPublic"], myStruct.StrPublic)
	assert.Equal(t, "", myStruct.strPrivate)
	assert.Equal(t, int32(0), myStruct.WrongFieldType)
}

func TestIsEncryptionAtlasEmpty(t *testing.T) {
	spec := &v1.EncryptionAtRest{}
	isEmpty := IsEncryptionSpecEmpty(spec)
	assert.True(t, isEmpty, "Empty spec should be empty")

	spec.AwsKms.Enabled = toptr.MakePtr(true)
	isEmpty = IsEncryptionSpecEmpty(spec)
	assert.False(t, isEmpty, "Non-empty spec")

	spec.AwsKms.Enabled = toptr.MakePtr(false)
	isEmpty = IsEncryptionSpecEmpty(spec)
	assert.True(t, isEmpty, "Enabled flag set to false is same as empty")
}

func TestAtlasInSync(t *testing.T) {
	areInSync, err := AtlasInSync(nil, nil)
	assert.NoError(t, err)
	assert.True(t, areInSync, "Both atlas and spec are nil")

	groupID := "0"
	atlas := mongodbatlas.EncryptionAtRest{
		GroupID: groupID,
		AwsKms: mongodbatlas.AwsKms{
			Enabled: toptr.MakePtr(true),
		},
	}
	spec := v1.EncryptionAtRest{
		AwsKms: v1.AwsKms{
			Enabled: toptr.MakePtr(true),
		},
	}

	areInSync, err = AtlasInSync(nil, &spec)
	assert.NoError(t, err)
	assert.False(t, areInSync, "Nil atlas")

	areInSync, err = AtlasInSync(&atlas, nil)
	assert.NoError(t, err)
	assert.False(t, areInSync, "Nil spec")

	areInSync, err = AtlasInSync(&atlas, &spec)
	assert.NoError(t, err)
	assert.True(t, areInSync, "Both are the same")

	spec.AwsKms.Enabled = toptr.MakePtr(false)
	areInSync, err = AtlasInSync(&atlas, &spec)
	assert.NoError(t, err)
	assert.False(t, areInSync, "Atlas is disabled")

	atlas.AwsKms.Enabled = toptr.MakePtr(false)
	areInSync, err = AtlasInSync(&atlas, &spec)
	assert.NoError(t, err)
	assert.True(t, areInSync, "Both are disabled")

	atlas.AwsKms.RoleID = "example"
	areInSync, err = AtlasInSync(&atlas, &spec)
	assert.NoError(t, err)
	assert.True(t, areInSync, "Both are disabled but atlas RoleID field")

	spec.AwsKms.Enabled = toptr.MakePtr(true)
	areInSync, err = AtlasInSync(&atlas, &spec)
	assert.NoError(t, err)
	assert.False(t, areInSync, "Spec is re-enabled")

	atlas.AwsKms.Enabled = toptr.MakePtr(true)
	areInSync, err = AtlasInSync(&atlas, &spec)
	assert.NoError(t, err)
	assert.True(t, areInSync, "Both are re-enabled and only RoleID is different")

	atlas = mongodbatlas.EncryptionAtRest{
		AwsKms: mongodbatlas.AwsKms{
			Enabled:             toptr.MakePtr(true),
			CustomerMasterKeyID: "example",
			Region:              "US_EAST_1",
			RoleID:              "example",
			Valid:               toptr.MakePtr(true),
		},
		AzureKeyVault: mongodbatlas.AzureKeyVault{
			Enabled: toptr.MakePtr(false),
		},
		GoogleCloudKms: mongodbatlas.GoogleCloudKms{
			Enabled: toptr.MakePtr(false),
		},
	}
	spec = v1.EncryptionAtRest{
		AwsKms: v1.AwsKms{
			Enabled:             toptr.MakePtr(true),
			CustomerMasterKeyID: "example",
			Region:              "US_EAST_1",
			Valid:               toptr.MakePtr(true),
		},
		AzureKeyVault:  v1.AzureKeyVault{},
		GoogleCloudKms: v1.GoogleCloudKms{},
	}

	areInSync, err = AtlasInSync(&atlas, &spec)
	assert.NoError(t, err)
	assert.True(t, areInSync, "Realistic exampel. should be equal")
}
