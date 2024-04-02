package launcher_test

import (
	"testing"
	"time"

	_ "embed"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/launcher"
)

//go:embed test.yml
var testYml string

const (
	testVersion = "2.1.0"
)

func TestLaunch(t *testing.T) {
	if !control.Enabled("AKO_LAUNCHER_TEST") {
		t.Skip("Skipping int tests, AKO_LAUNCHER_TEST is not set")
	}

	l := launcher.NewFromEnv(testVersion)
	assert.NoError(t, l.Launch(
		testYml,
		launcher.WaitReady("atlasprojects/my-project", 30*time.Second)))
	// retest should also work and faster (no need to re-create resources)
	assert.NoError(t, l.Launch(
		testYml,
		launcher.WaitReady("atlasprojects/my-project", 30*time.Second)))

	require.NoError(t, l.Cleanup())
}
