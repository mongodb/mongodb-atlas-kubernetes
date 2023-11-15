package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/toptr"
)

var (
	excludedClusterFieldsOurs   = map[string]bool{}
	excludedClusterFieldsTheirs = map[string]bool{}
)

func init() {
	excludedClusterFieldsOurs["projectRef"] = true

	// Global deployment fields
	excludedClusterFieldsOurs["customZoneMapping"] = true
	excludedClusterFieldsOurs["managedNamespaces"] = true

	excludedClusterFieldsTheirs["backupEnabled"] = true
	excludedClusterFieldsTheirs["id"] = true
	excludedClusterFieldsTheirs["groupId"] = true
	excludedClusterFieldsTheirs["createDate"] = true
	excludedClusterFieldsTheirs["links"] = true
	excludedClusterFieldsTheirs["versionReleaseSystem"] = true

	// Deprecated
	excludedClusterFieldsTheirs["replicationSpec"] = true
	excludedClusterFieldsTheirs["replicationFactor"] = true

	// Termination protection
	excludedClusterFieldsTheirs["terminationProtectionEnabled"] = true

	// Root cert type
	excludedClusterFieldsTheirs["rootCertType"] = true

	// These fields are shown in the status
	excludedClusterFieldsTheirs["mongoDBVersion"] = true
	excludedClusterFieldsTheirs["mongoURI"] = true
	excludedClusterFieldsTheirs["mongoURIUpdated"] = true
	excludedClusterFieldsTheirs["mongoURIWithOptions"] = true
	excludedClusterFieldsTheirs["connectionStrings"] = true
	excludedClusterFieldsTheirs["srvAddress"] = true
	excludedClusterFieldsTheirs["stateName"] = true
	excludedClusterFieldsTheirs["links"] = true
	excludedClusterFieldsTheirs["createDate"] = true
	excludedClusterFieldsTheirs["versionReleaseSystem"] = true
	excludedClusterFieldsTheirs["serverlessBackupOptions"] = true
}

func TestIsEqual(t *testing.T) {
	operatorArgs := ProcessArgs{
		JavascriptEnabled: toptr.MakePtr(true),
	}

	atlasArgs := mongodbatlas.ProcessArgs{
		JavascriptEnabled: toptr.MakePtr(false),
	}

	areTheyEqual := operatorArgs.IsEqual(atlasArgs)
	assert.False(t, areTheyEqual, "should NOT be equal if pointer values are different")

	atlasArgs.JavascriptEnabled = toptr.MakePtr(true)

	areTheyEqual = operatorArgs.IsEqual(atlasArgs)
	assert.True(t, areTheyEqual, "should be equal if all pointer values are the same")

	areTheyEqual = operatorArgs.IsEqual(&atlasArgs)
	assert.True(t, areTheyEqual, "should be equal if Atlas args is a pointer")

	atlasArgs.DefaultReadConcern = "available"

	areTheyEqual = operatorArgs.IsEqual(atlasArgs)
	assert.True(t, areTheyEqual, "should be equal if Atlas args have more values")

	operatorArgs.DefaultReadConcern = "available"

	areTheyEqual = operatorArgs.IsEqual(atlasArgs)
	assert.True(t, areTheyEqual, "should work for non-pointer values")

	operatorArgs.OplogSizeMB = toptr.MakePtr[int64](8)

	areTheyEqual = operatorArgs.IsEqual(atlasArgs)
	assert.False(t, areTheyEqual, "should NOT be equal if Operator has more args")

	atlasArgs.OplogSizeMB = toptr.MakePtr[int64](8)

	areTheyEqual = operatorArgs.IsEqual(atlasArgs)
	assert.True(t, areTheyEqual, "should become equal")

	operatorArgs.OplogMinRetentionHours = "2.0"
	atlasArgs.OplogMinRetentionHours = toptr.MakePtr[float64](2)

	areTheyEqual = operatorArgs.IsEqual(atlasArgs)
	assert.True(t, areTheyEqual, "should be equal when OplogMinRetentionHours field is the same")
}

func TestToAtlas(t *testing.T) {
	operatorArgs := ProcessArgs{
		JavascriptEnabled:      toptr.MakePtr(true),
		OplogMinRetentionHours: "2.0",
	}

	atlasArgs, err := operatorArgs.ToAtlas()
	assert.NoError(t, err, "no errors should occur")

	areTheyEqual := operatorArgs.IsEqual(atlasArgs)
	assert.True(t, areTheyEqual, "should be equal after conversion")
}
